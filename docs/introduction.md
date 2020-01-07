# Introduction

The Wavefront Kubernetes HPA Adapter can be used to horizontally autoscale your Kubernetes Pods based on metrics available within Wavefront.

By default the Kubernetes HorizontalPodAutoscaler (HPA) controller fetches metrics from a series of [APIs](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis). This adapter implements the APIs detailed below.

## custom.metrics.k8s.io
For the custom metrics API, this adapter provides all Kubernetes metrics collected using [Heapster](https://docs.wavefront.com/kubernetes.html) or the [Wavefront Kubernetes Collector](https://github.com/wavefrontHQ/wavefront-kubernetes-collector).

This can be configured using the `wavefront-metric-prefix` adapter property. See [metrics.md](https://github.com/wavefrontHQ/wavefront-kubernetes-adapter/blob/master/docs/metrics.md) for the list of metrics provided through this API.

Use the external metrics API detailed below if you wish to use non Kubernetes metrics or Kubernetes metrics collected using a different mechanism other than Heapster or the Wavefront Collector.

## external.metrics.k8s.io

For the external metrics API, you have complete control over what metrics are provided by this adapter:

1. Deploy a Kubernetes [config map](https://github.com/wavefrontHQ/wavefront-kubernetes-adapter/blob/master/deploy/manifests/04-custom-metrics-config-map.yaml) listing the metrics you wish to autoscale on.
2. Configure the [`external-metrics-config`](https://github.com/wavefrontHQ/wavefront-kubernetes-adapter/blob/06b75dcf6fd9813a2b8a5a5762be1ae922d35ce7/deploy/manifests/custom-metrics-apiserver-deployment.yaml#L31) adapter property based on the config map and redeploy the adapter.
3. Deploy a [HPA](https://github.com/wavefrontHQ/wavefront-kubernetes-adapter/blob/master/deploy/hpa-examples/hpa-external.yaml) based on an external metric.


For example, you can autoscale a Pod based on AWS SQS queue sizes or custom business metrics specific to your application services.

## Kubernetes HPA Spec
Refer to the [autoscaling spec](https://godoc.org/k8s.io/api/autoscaling/v2beta1#MetricSpec) for more details on configuring HPAs based on the custom or external metrics APIs.
