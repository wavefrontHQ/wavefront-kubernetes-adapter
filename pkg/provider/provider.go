package provider

import (
	"fmt"
	"time"

	apierr "k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/apis/custom_metrics"
	"k8s.io/metrics/pkg/apis/external_metrics"

	"github.com/golang/glog"
	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider/helpers"

	wave "github.com/wavefronthq/wavefront-kubernetes-adapter/pkg/client"
)

type wavefrontProvider struct {
	mapper         apimeta.RESTMapper
	dynClient      dynamic.Interface
	waveClient     wave.WavefrontClient
	lister         MetricsLister
	externalDriver ExternalMetricsDriver

	Translator
}

type WavefrontProviderConfig struct {
	DynClient    dynamic.Interface
	KubeClient   kubernetes.Interface
	Mapper       apimeta.RESTMapper
	WaveClient   wave.WavefrontClient
	Prefix       string
	ListInterval time.Duration
	ExternalCfg  string
}

func NewWavefrontProvider(cfg WavefrontProviderConfig) (provider.MetricsProvider, MetricsLister) {
	glog.Infof("wavefrontProvider Prefix: %s, ListInterval: %d", cfg.Prefix, cfg.ListInterval)

	translator := NewWavefrontTranslator(cfg.Prefix)
	externalDriver := NewExternalMetricsDriver(cfg.KubeClient, cfg.ExternalCfg)

	lister := &WavefrontMetricsLister{
		Prefix:         cfg.Prefix,
		UpdateInterval: cfg.ListInterval,
		waveClient:     cfg.WaveClient,
		externalDriver: externalDriver,
		Translator:     translator,
	}

	return &wavefrontProvider{
		dynClient:      cfg.DynClient,
		mapper:         cfg.Mapper,
		waveClient:     cfg.WaveClient,
		lister:         lister,
		externalDriver: externalDriver,
		Translator:     translator,
	}, lister
}

func (p *wavefrontProvider) query(info provider.CustomMetricInfo, namespace string, names ...string) (wave.QueryResult, error) {
	query, found := p.QueryFor(info, namespace, names...)
	if !found {
		return wave.QueryResult{}, provider.NewMetricNotFoundError(info.GroupResource, info.Metric)
	}
	return p.doQuery(query)
}

func (p *wavefrontProvider) doQuery(query string) (wave.QueryResult, error) {
	now := time.Now()
	start := now.Add(time.Duration(-30) * time.Second)
	queryResult, err := p.waveClient.Query(start.Unix(), query)
	if err != nil {
		glog.Errorf("unable to fetch metrics from wavefront: %v", err)
		// don't leak implementation details to the user
		return wave.QueryResult{}, apierr.NewInternalError(fmt.Errorf("unable to fetch metrics"))
	}
	return queryResult, nil
}

func (p *wavefrontProvider) metricFor(value float64, name types.NamespacedName, info provider.CustomMetricInfo) (*custom_metrics.MetricValue, error) {

	objRef, err := helpers.ReferenceFor(p.mapper, name, info)
	if err != nil {
		return nil, err
	}

	return &custom_metrics.MetricValue{
		DescribedObject: objRef,
		MetricName:      info.Metric,
		Timestamp:       metav1.Time{time.Now()},
		Value:           *resource.NewMilliQuantity(int64(value*1000.0), resource.DecimalSI),
	}, nil
}

func (p *wavefrontProvider) metricsFor(queryResult wave.QueryResult, namespace string, info provider.CustomMetricInfo, names []string) (*custom_metrics.MetricValueList, error) {

	values, found := p.MatchValuesToNames(queryResult, info.GroupResource)
	if !found {
		return nil, provider.NewMetricNotFoundError(info.GroupResource, info.Metric)
	}
	glog.V(5).Info("DEBUG:---metricsFor values", values)

	res := make([]custom_metrics.MetricValue, len(names))
	for i, name := range names {
		value, err := p.metricFor(values[name], types.NamespacedName{Namespace: namespace, Name: name}, info)
		if err != nil {
			return nil, err
		}
		res[i] = *value
	}

	return &custom_metrics.MetricValueList{
		Items: res,
	}, nil
}

func (p *wavefrontProvider) getSingle(info provider.CustomMetricInfo, name types.NamespacedName) (*custom_metrics.MetricValue, error) {
	queryResult, err := p.query(info, name.Namespace, name.Name)
	if err != nil {
		return nil, err
	}

	if len(queryResult.Timeseries) < 1 {
		return nil, provider.NewMetricNotFoundForError(info.GroupResource, info.Metric, name.Name)
	}

	namedValues, found := p.MatchValuesToNames(queryResult, info.GroupResource)
	if !found {
		return nil, provider.NewMetricNotFoundError(info.GroupResource, info.Metric)
	}

	if len(namedValues) > 1 {
		glog.V(2).Infof("Got more than one result (%v results) when fetching metric %s for %q, using the first one with a matching name...",
			len(queryResult.Timeseries), info.String(), name)
	}

	resultValue, nameFound := namedValues[name.Name]
	if !nameFound {
		glog.Errorf("None of the results returned by when fetching metric %s for %q matched the resource name", info.String(), name)
		return nil, provider.NewMetricNotFoundForError(info.GroupResource, info.Metric, name.Name)
	}
	return p.metricFor(resultValue, name, info)
}

func (p *wavefrontProvider) getMultiple(info provider.CustomMetricInfo, namespace string, selector labels.Selector) (*custom_metrics.MetricValueList, error) {
	resourceNames, err := helpers.ListObjectNames(p.mapper, p.dynClient, namespace, selector, info)
	if err != nil {
		return nil, err
	}
	glog.V(5).Infof("DEBUG:---resourceNames: %s", resourceNames)

	// query Wavefront for points
	queryResult, err := p.query(info, namespace, resourceNames...)
	if err != nil {
		return nil, err
	}
	return p.metricsFor(queryResult, namespace, info, resourceNames)
}

func (p *wavefrontProvider) GetMetricByName(name types.NamespacedName, info provider.CustomMetricInfo) (*custom_metrics.MetricValue, error) {
	glog.V(5).Info("DEBUG:---GetMetricByName", name, info)
	return p.getSingle(info, name)
}

func (p *wavefrontProvider) GetMetricBySelector(namespace string, selector labels.Selector, info provider.CustomMetricInfo) (*custom_metrics.MetricValueList, error) {
	glog.V(5).Info("DEBUG:---GetMetricBySelector", namespace, selector, info)
	return p.getMultiple(info, namespace, selector)
}

// Provides a list of all available metrics at the current time.
// Note that we cache and periodically update this list, instead of querying every time.
func (p *wavefrontProvider) ListAllMetrics() []provider.CustomMetricInfo {
	return p.lister.ListCustomMetrics()
}

func (p *wavefrontProvider) GetExternalMetric(namespace string, metricSelector labels.Selector, info provider.ExternalMetricInfo) (*external_metrics.ExternalMetricValueList, error) {
	glog.V(5).Info("DEBUG:---GetExternalMetric", namespace, metricSelector, info)

	if p.externalDriver == nil {
		return nil, apierr.NewInternalError(fmt.Errorf("missing external driver for external metric: %s", info.Metric))
	}

	query := p.externalDriver.getQuery(info.Metric)
	if query == "" {
		return nil, apierr.NewInternalError(fmt.Errorf("missing query for external metric: %s", info.Metric))
	}

	queryResult, err := p.doQuery(query)
	if err != nil {
		return nil, apierr.NewInternalError(fmt.Errorf("error fetching metrics for external metric: %s error=%v", info.Metric, err))
	}
	return p.ExternalValuesFor(queryResult, info.Metric)
}

func (p *wavefrontProvider) ListAllExternalMetrics() []provider.ExternalMetricInfo {
	return p.lister.ListExternalMetrics()
}
