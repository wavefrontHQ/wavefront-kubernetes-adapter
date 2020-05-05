// Copyright 2018-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

// RESTMapper
type fakeRESTMapper struct {
	kindForInput schema.GroupVersionResource
}

func (f *fakeRESTMapper) KindFor(resource schema.GroupVersionResource) (schema.GroupVersionKind, error) {
	f.kindForInput = resource
	return schema.GroupVersionKind{
		Version: "v1",
		Kind:    "pod",
	}, nil
}

func (f *fakeRESTMapper) KindsFor(resource schema.GroupVersionResource) ([]schema.GroupVersionKind, error) {
	result := []schema.GroupVersionKind{}
	result = append(result, schema.GroupVersionKind{
		Version: "v1",
		Kind:    "pod",
	})
	return result, nil
}

func (f *fakeRESTMapper) ResourceFor(input schema.GroupVersionResource) (schema.GroupVersionResource, error) {
	resources, err := f.ResourcesFor(input)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}
	return resources[0], nil
}

func (f *fakeRESTMapper) ResourcesFor(input schema.GroupVersionResource) ([]schema.GroupVersionResource, error) {
	result := []schema.GroupVersionResource{}
	result = append(result, getGroupVersionResource("test-deployment-7f54684694-2cg5v"))
	result = append(result, getGroupVersionResource("test-deployment-7f54684694-cbts9"))
	result = append(result, getGroupVersionResource("test-deployment-7f54684694-mm49g"))
	result = append(result, getGroupVersionResource("test-deployment-7f54684694-t57tb"))
	result = append(result, getGroupVersionResource("test-deployment-7f54684694-xnxfp"))
	return result, nil
}

func (f *fakeRESTMapper) RESTMapping(gk schema.GroupKind, versions ...string) (*meta.RESTMapping, error) {
	return nil, nil
}

func (f *fakeRESTMapper) RESTMappings(gk schema.GroupKind, versions ...string) ([]*meta.RESTMapping, error) {
	return nil, nil
}

func (f *fakeRESTMapper) ResourceSingularizer(resource string) (singular string, err error) {
	return "", nil
}

func getGroupVersionResource(name string) schema.GroupVersionResource {
	return schema.GroupVersionResource{Resource: name}
}

// ExternalMetricsDriver

type fakeExternalDriver struct{}

func (d *fakeExternalDriver) loadConfig() {}

func (d *fakeExternalDriver) registerListener(listener ExternalConfigListener) {}

func (d *fakeExternalDriver) getMetricNames() []string {
	result := make([]string, 0)
	result = append(result, "externalMetric1")
	result = append(result, "externalMetric2")
	result = append(result, "externalMetric3")
	result = append(result, "externalMetric4")
	result = append(result, "externalMetric5")
	return result
}

func (d *fakeExternalDriver) getQuery(metric string) string {
	if strings.HasPrefix(metric, "external") {
		return "ts(cpu.usage.idle)"
	}
	return ""
}
