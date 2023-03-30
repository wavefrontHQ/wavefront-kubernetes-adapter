# Operations for Applications HPA Adapter for Kubernetes
[![build status][ci-img]][ci] [![Go Report Card][go-report-img]][go-report] [![Docker Pulls][docker-pull-img]][docker-img]

[VMware Aria Operations for Applications](https://docs.wavefront.com) is a high-performance streaming analytics platform for monitoring and optimizing your environment and applications.

The Operations for Applications HPA (Horizontal Pod Autoscaler) adapter for Kubernetes implements the Kubernetes `custom.metrics.k8s.io/v1beta1` and `external.metrics.k8s.io/v1beta1` APIs allowing you to autoscale pods based on metrics available within Operations for Applications.

## Prerequisites

- Kubernetes 1.9+
- The [aggregation layer](https://kubernetes.io/docs/tasks/access-kubernetes-api/configure-aggregation-layer/) needs to be enabled in your Kubernetes cluster prior to deploying the Operations for Applications HPA adapter.

## Getting Started
Refer to the [Getting Started](/docs/introduction.md) guide for an overview of the functionality provided by this adapter.

## Configuration

Refer to the [Configuration](/docs/configuration.md) documentation for detailed configuration options.

## Installation

### Helm install
Refer to the [helm chart](https://github.com/wavefrontHQ/helm#installation) to install the adapter using Helm.

### Manual install
1. Clone this repo.
2. Edit the `wavefront-url` and `wavefront-token` properties in `deploy/manifests/05-custom-metrics-apiserver-deployment.yaml`.
3. Optionally, edit the `deploy/manifests/04-custom-metrics-config-map.yaml` and modify the external metrics you wish to export.
4. Finally run `kubectl apply -f deploy/manifests` to deploy the adapter in your Kubernetes cluster.

## Debugging

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

See the sample configurations under [hpa-examples](/deploy/hpa-examples/) on how to autoscale deployments based on metrics under the custom or external metrics APIs.

Run `kubectl describe hpa example-hpa-custom-metrics` to verify the autoscaling works.

[ci-img]: https://github.com/wavefrontHQ/wavefront-kubernetes-adapter/actions/workflows/go.yml/badge.svg
[ci]: https://github.com/wavefrontHQ/wavefront-kubernetes-adapter/actions/workflows/go.yml
[go-report-img]: https://goreportcard.com/badge/github.com/wavefronthq/wavefront-kubernetes-adapter
[go-report]: https://goreportcard.com/report/github.com/wavefronthq/wavefront-kubernetes-adapter
[docker-pull-img]: https://img.shields.io/docker/pulls/wavefronthq/wavefront-hpa-adapter.svg?logo=docker
[docker-img]: https://hub.docker.com/r/wavefronthq/wavefront-hpa-adapter/
