// Copyright 2018-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"fmt"
	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
	"github.com/wavefronthq/wavefront-kubernetes-adapter/pkg/client"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic/fake"
	"testing"
)

func TestListAllMetrics(t *testing.T) {
	waveProvider := fakeProvider()
	metrics := waveProvider.ListAllMetrics()

	if len(metrics) != 5 {
		t.Errorf("Invalid list of metrics, len: %d", len(metrics))
	}
}

func TestGetMetricByName(t *testing.T) {
	waveProvider := fakeProvider()
	_, err := waveProvider.GetMetricByName(namespacedName("test-deployment-7f54684694-2cg5v", "default"), fakeCustomMetricInfo())
	if err != nil {
		t.Error(err)
	}
}

func GetMetricBySelector(t *testing.T) {
	//TODO: there's a bug in the fake REST mapping that currently causes this test to fail
	waveProvider := fakeProvider()
	selector := labels.Everything()
	_, err := waveProvider.GetMetricBySelector("default", selector, fakeCustomMetricInfo())

	if err != nil {
		t.Error(err)
	}
}

func TestListAllExternalMetrics(t *testing.T) {
	waveProvider := fakeProvider()
	metrics := waveProvider.ListAllExternalMetrics()

	if len(metrics) != 5 {
		t.Errorf("Invalid list of metrics, len: %d", len(metrics))
	}
}

func TestGetExternalMetric(t *testing.T) {
	waveProvider := fakeProvider()

	_, err := waveProvider.GetExternalMetric("", nil, provider.ExternalMetricInfo{Metric: "failMetric"})
	if err == nil {
		t.Error("Expected error but no error returned")
	}

	values, err := waveProvider.GetExternalMetric("", nil, provider.ExternalMetricInfo{Metric: "externalMetric1"})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(values)
}

func fakeProvider() provider.MetricsProvider {
	restMapper := &fakeRESTMapper{}
	dynClient := &fake.FakeDynamicClient{}
	api := client.NewFakeWavefrontClient()
	prefix := "kubernetes"
	translator := &wavefrontTranslator{
		prefix: prefix,
	}

	lister := &WavefrontMetricsLister{
		Prefix:         prefix,
		waveClient:     api,
		Translator:     translator,
		externalDriver: &fakeExternalDriver{},
	}
	lister.updateMetrics()

	return &wavefrontProvider{
		dynClient:      dynClient,
		mapper:         restMapper,
		waveClient:     api,
		Translator:     translator,
		lister:         lister,
		externalDriver: &fakeExternalDriver{},
	}
}

func namespacedName(name, namespace string) types.NamespacedName {
	return types.NamespacedName{
		Name:      "test-deployment-7f54684694-2cg5v",
		Namespace: "default",
	}
}

func fakeCustomMetricInfo() provider.CustomMetricInfo {
	return provider.CustomMetricInfo{
		GroupResource: schema.GroupResource{Resource: "pods"},
		Namespaced:    true,
		Metric:        "cpu.usage_rate",
	}
}
