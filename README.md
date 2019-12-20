# wavefront-kubernetes-adapter [![build status][ci-img]][ci] [![Go Report Card][go-report-img]][go-report] [![Docker Pulls][docker-pull-img]][docker-img]

This repository provides a Wavefront Kubernetes HPA (Horizontal Pod Autoscaler) adapter that implements the custom metrics (`custom.metrics.k8s.io/v1beta1`) and external metrics (`external.metrics.k8s.io/v1beta1`) APIs. The adapter can be used with the `autoscaling/v2` HPA in Kubernetes 1.9+.

## Prerequisites

- Kubernetes 1.9+
- The [aggregation layer](https://kubernetes.io/docs/tasks/access-kubernetes-api/configure-aggregation-layer/) needs to be enabled in your Kubernetes cluster prior to deploying the Wavefront adapter.

## Introduction
See the [introduction](https://github.com/wavefrontHQ/wavefront-kubernetes-adapter/blob/master/docs/introduction.md) for an overview of the functionality provided by this adapter.

## Configuration

The adapter takes the standard Kubernetes generic API server arguments.

Additionally, the following arguments are specific to this adapter:

- `wavefront-url`: Wavefront URL of the form *https://INSTANCE.wavefront.com*
- `wavefront-token`: Wavefront API token with permissions to query for points
- `wavefront-metric-prefix`: Metrics under the prefix are exposed in the custom metrics API. Defaults to `kubernetes`.
- `metrics-relist-interval`: The interval at which to fetch the list of metrics from Wavefront. Defaults to 10 minutes.
- `external-metrics-config`: Optional configuration file driving the external metrics API. If omitted, the external metrics API will not be supported.

## Installation

1. Clone this repo.
2. Edit the `wavefront-url` and `wavefront-token` properties in `deploy/manifests/custom-metrics-apiserver-deployment.yaml`.
3. Optionally, edit the `deploy/manifests/custom-metrics-config-map.yaml` and modify the external metrics you wish to export.
4. Finally run `kubectl apply -f deploy/manifests` to deploy the adapter in your Kubernetes cluster.

To verify the installation, run `kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1" | jq .`. You should get the list of supported metrics similar to:

```json
{
  "kind": "APIResourceList",
  "apiVersion": "v1",
  "groupVersion": "custom.metrics.k8s.io/v1beta1",
  "resources": [
    {
      "name": "nodes/cpu.node_reservation",
      "singularName": "",
      "namespaced": false,
      "kind": "MetricValueList",
      "verbs": [
        "get"
      ]
    },
    {
      "name": "pods/network.rx_errors_rate",
      "singularName": "",
      "namespaced": true,
      "kind": "MetricValueList",
      "verbs": [
        "get"
      ]
    },
    {
      "name": "pods/network.tx_rate",
      "singularName": "",
      "namespaced": true,
      "kind": "MetricValueList",
      "verbs": [
        "get"
      ]
    },
    {
      "name": "pods/cpu.request",
      "singularName": "",
      "namespaced": true,
      "kind": "MetricValueList",
      "verbs": [
        "get"
      ]
    }
  ]
}    
```

You can similarly run `kubectl get --raw "/apis/external.metrics.k8s.io/v1beta1" | jq .` to verify the external metrics API.

## HPA Autoscaling

See the sample configurations under [hpa-examples](https://github.com/wavefrontHQ/wavefront-kubernetes-adapter/tree/master/deploy/hpa-examples) on how to autoscale deployments based on metrics under the custom or external metrics APIs.

Run `kubectl describe hpa example-hpa-custom-metrics` to verify the autoscaling works.

[ci-img]: https://travis-ci.com/wavefrontHQ/wavefront-kubernetes-adapter.svg?branch=master
[ci]: https://travis-ci.com/wavefrontHQ/wavefront-kubernetes-adapter
[go-report-img]: https://goreportcard.com/badge/github.com/wavefronthq/wavefront-kubernetes-adapter
[go-report]: https://goreportcard.com/report/github.com/wavefronthq/wavefront-kubernetes-adapter
[docker-pull-img]: https://img.shields.io/docker/pulls/wavefronthq/wavefront-hpa-adapter.svg?logo=docker
[docker-img]: https://hub.docker.com/r/wavefronthq/wavefront-hpa-adapter/
