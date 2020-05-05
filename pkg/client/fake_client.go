// Copyright 2018-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"net/http"
	"net/url"
)

type FakeWavefrontClient struct{}

func NewFakeWavefrontClient() WavefrontClient {
	return FakeWavefrontClient{}
}

func (w FakeWavefrontClient) Do(verb, endpoint string, query url.Values) (*http.Response, error) {
	return &http.Response{}, nil
}

func (w FakeWavefrontClient) ListMetrics(prefix string) ([]string, error) {
	result := make([]string, 0)
	result = append(result, "kubernetes.node.cpu.node_reservation")
	result = append(result, "kubernetes.pod.network.rx_errors_rate")
	result = append(result, "kubernetes.pod.network.tx_rate")
	result = append(result, "kubernetes.pod.cpu.request")
	result = append(result, "kubernetes.pod.cpu.usage_rate")
	return result, nil
}

func (w FakeWavefrontClient) Query(ts int64, query string) (QueryResult, error) {
	return fakeQueryResult(), nil
}

func fakeQueryResult() QueryResult {
	timeseries := make([]Timeseries, 5)
	timeseries[0] = getTimeseries("pod_name", "test-deployment-7f54684694-2cg5v")
	timeseries[1] = getTimeseries("pod_name", "test-deployment-7f54684694-cbts9")
	timeseries[2] = getTimeseries("pod_name", "test-deployment-7f54684694-mm49g")
	timeseries[3] = getTimeseries("pod_name", "test-deployment-7f54684694-t57tb")
	timeseries[4] = getTimeseries("pod_name", "test-deployment-7f54684694-xnxfp")

	result := QueryResult{
		Timeseries: timeseries,
	}
	return result
}

func getTimeseries(key, name string) Timeseries {
	return Timeseries{
		Tags: getTags(key, name),
		Data: getData(),
	}
}

func getData() [][]float64 {
	data := make([][]float64, 1)
	data[0] = []float64{0.0, 2.3598}
	return data
}

func getTags(key, name string) map[string]string {
	tags := make(map[string]string, 2)
	tags[key] = name
	tags["nodename"] = "gke-cluster-default-pool-f63db08a-xrdh"
	return tags
}
