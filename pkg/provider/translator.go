// Copyright 2018-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/metrics/pkg/apis/external_metrics"

	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
	wave "github.com/wavefronthq/wavefront-kubernetes-adapter/pkg/client"
)

type Translator interface {
	QueryFor(info provider.CustomMetricInfo, namespace string, names ...string) (string, bool)
	MatchValuesToNames(queryResult wave.QueryResult, groupResource schema.GroupResource) (map[string]float64, bool)
	CustomMetricsFor(metricNames []string) []provider.CustomMetricInfo
	ExternalMetricsFor(metricNames []string) []provider.ExternalMetricInfo
	ExternalValuesFor(queryResult wave.QueryResult, metric string) (*external_metrics.ExternalMetricValueList, error)
}

type wavefrontTranslator struct {
	prefix string
}

func NewWavefrontTranslator(prefix string) Translator {
	return &wavefrontTranslator{prefix: prefix}
}

// Translates given metric info into a Wavefront ts query
func (t wavefrontTranslator) QueryFor(info provider.CustomMetricInfo, namespace string, names ...string) (string, bool) {
	metric := info.Metric
	resType := resourceType(info.GroupResource.Resource)
	resourceFilter := filterFor(tagKey(resType), " or ", names...)
	namespaceFilter := ""
	if namespaced(resType) {
		namespaceFilter = filterFor("namespace_name", "", namespace)
	}
	filters := combine(resourceFilter, namespaceFilter)

	// if Prefix=kubernetes, metric='cpu.usage_rate', resType='pod', namespace='default' and names=['pod1', 'pod2']
	// ts(kubernetes.pod.cpu.usage_rate, (pod_name="pod1" or pod_name="pod2") and (namespace_name="default"))
	query := fmt.Sprintf("ts(%s.%s.%s%s)", t.prefix, resType, metric, filters)
	return query, true
}

func (t wavefrontTranslator) MatchValuesToNames(queryResult wave.QueryResult, groupResource schema.GroupResource) (map[string]float64, bool) {
	log.Debugf("MatchValuesToNames: %v", queryResult.Timeseries)

	if len(queryResult.Timeseries) == 0 {
		return nil, false
	}

	resType := resourceType(groupResource.Resource)
	tagKey := tagKey(resType)

	values := make(map[string]float64, len(queryResult.Timeseries))
	for _, timeseries := range queryResult.Timeseries {
		length := len(timeseries.Data)
		if length == 0 {
			return nil, false
		}
		key, found := timeseries.Tags[tagKey]
		if !found {
			return nil, false
		}
		value, err := trimFloat(timeseries.Data[length-1][1])
		if err != nil {
			return nil, false
		}
		values[key] = value
	}
	return values, true
}

func (t wavefrontTranslator) CustomMetricsFor(metricNames []string) []provider.CustomMetricInfo {
	var customMetrics []provider.CustomMetricInfo
	for _, metricName := range metricNames {
		resourceName, metric := splitMetric(t.prefix, metricName)
		if resourceName == "" || metric == "" {
			continue
		}
		customMetrics = append(customMetrics, provider.CustomMetricInfo{
			GroupResource: schema.GroupResource{Group: "", Resource: normalize(resourceName)},
			Metric:        metric,
			Namespaced:    namespaced(resourceName),
		})
	}
	return customMetrics
}

func (t wavefrontTranslator) ExternalMetricsFor(metricNames []string) []provider.ExternalMetricInfo {
	var externalMetrics []provider.ExternalMetricInfo
	for _, metricName := range metricNames {
		externalMetrics = append(externalMetrics, provider.ExternalMetricInfo{
			Metric: metricName,
		})
	}
	return externalMetrics
}

func (t wavefrontTranslator) ExternalValuesFor(queryResult wave.QueryResult, name string) (*external_metrics.ExternalMetricValueList, error) {
	var matchingMetrics []external_metrics.ExternalMetricValue
	for _, timeseries := range queryResult.Timeseries {
		length := len(timeseries.Data)
		if length == 0 {
			return nil, fmt.Errorf("no data for external metric: %s", name)
		}

		// use the last data point
		point := timeseries.Data[length-1]
		if len(point) != 2 {
			return nil, fmt.Errorf("invalid data point for external metric: %s", name)
		}
		value, err := trimFloat(point[1])
		if err != nil {
			log.Errorf("error converting external metric: %s value: %f", name, point[1])
			continue
		}
		metricValue := external_metrics.ExternalMetricValue{
			MetricName: name,
			Value:      *resource.NewMilliQuantity(int64(1000*value), resource.DecimalSI),
			Timestamp:  metav1.Now(),
		}
		matchingMetrics = append(matchingMetrics, metricValue)
	}
	return &external_metrics.ExternalMetricValueList{
		Items: matchingMetrics,
	}, nil
}

var (
	resourceMap = map[string]string{
		"cluster":       "clusters",
		"ns":            "namespaces",
		"pod":           "pods",
		"node":          "nodes",
		"pod_container": "pod_containers",
		"sys_container": "sys_containers",
	}
)

func resourceType(resource string) string {
	//TODO: maybe move into map as well? cleanup
	if resource == "pods" {
		return "pod"
	}
	if resource == "namespaces" {
		return "ns"
	}
	if resource == "nodes" {
		return "node"
	}
	return resource
}

// returns the tag key such as pod_name, namespace_name, nodename etc based on resourceType
func tagKey(resourceType string) string {
	// pod_name, namespace_name, nodename
	if resourceType == "node" {
		return "nodename"
	}
	return resourceType + "_name"
}

// given a list of tag values returns a ts filter string of the form:
// '(pod_name="pod1" or pod_name="pod2")'
func filterFor(tagKey, sep string, names ...string) string {
	filter := ""
	prefix := ""
	for _, name := range names {
		if name != "" {
			filter = fmt.Sprintf("%s%s%s=\"%s\"", filter, prefix, tagKey, name)
			prefix = sep
		}
	}
	if filter != "" {
		filter = fmt.Sprintf("(%s)", filter)
	}
	return filter
}

// combines multiple ts filters into a single string of the form:
// ',(pod_name="pod1" or pod_name="pod2") and (namespace_name="namespace1")'
func combine(filters ...string) string {
	and := ""
	result := ""
	for _, filter := range filters {
		if filter != "" {
			result = fmt.Sprintf("%s%s%s", result, and, filter)
			and = " and "
		}
	}
	if result != "" {
		result = fmt.Sprintf(", %s", result)
	}
	return result
}

func normalize(resourceName string) string {
	mappedName, found := resourceMap[resourceName]
	if !found {
		return resourceName
	}
	return mappedName
}

func namespaced(resourceName string) bool {
	// TODO: remove once https://github.com/kubernetes/kubernetes/issues/67777 is fixed
	return resourceName == "pod" || resourceName == "pod_container"
}

// splits a metric such as "kubernetes.pod.cpu.limit" into "pod" and "cpu.limit"
func splitMetric(prefix, metricName string) (string, string) {
	if strings.HasPrefix(metricName, prefix) {
		metricName = metricName[len(prefix)+1:]
	}
	parts := strings.SplitN(metricName, ".", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

// trims a float64 to 3 decimal digits
func trimFloat(value float64) (float64, error) {
	s := fmt.Sprintf("%.3f", value)
	return strconv.ParseFloat(s, 3)
}
