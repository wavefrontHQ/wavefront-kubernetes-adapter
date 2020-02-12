# Getting Started

The Wavefront HPA Adapter can be used to horizontally autoscale your Kubernetes Pods based on metrics available within Wavefront.

By default the Kubernetes HorizontalPodAutoscaler (HPA) controller fetches metrics from a series of [APIs](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis). This adapter implements the APIs detailed below.

## custom.metrics.k8s.io
For the custom metrics API, this adapter provides all Kubernetes metrics collected using the [Wavefront Kubernetes Collector](https://github.com/wavefrontHQ/wavefront-kubernetes-collector).

This can be configured using the `wavefront-metric-prefix` adapter property. See [metrics.md](https://github.com/wavefrontHQ/wavefront-kubernetes-adapter/blob/master/docs/metrics.md) for the list of metrics provided through this API.

Use the external metrics API detailed below if you wish to use non Kubernetes metrics or Kubernetes metrics collected using a different mechanism other than the Wavefront Collector.

## external.metrics.k8s.io
The external metrics API allows you to autoscale on any arbitrary metric available in Wavefront.

Metrics can be specified via annotations on HPAs or via a static configuration file.

### Annotations
[Annotations](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/) are metadata you attach to Kubernetes objects.

The adapter can dynamically discover HPAs and source external metrics via annotations. The annotations should be of the form `wavefront.com.external.metric/<metric_name>: '<ts query>'`. For example:

```yaml
apiVersion: autoscaling/v2beta1
kind: HorizontalPodAutoscaler
metadata:
  name: example-app
  annotations:
    wavefront.com.external.metric/sqs_queue_size: 'align(5m, avg(ts("aws.sqs.approximatenumberofmessagesvisible", QueueName="app-queue")))'
spec:
  minReplicas: 1
  maxReplicas: 5
  metrics:
  - type: External
    external:
      metricName: sqs_queue_size
      targetAverageValue: 248000m
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: example-app
```
**Note:** The external metric names should be globally unique.

### Static configuration file
To specify external metrics via a configuration file:

1. Deploy a Kubernetes [config map](https://github.com/wavefrontHQ/wavefront-kubernetes-adapter/blob/master/deploy/manifests/04-custom-metrics-config-map.yaml) listing the metrics you wish to autoscale on.
2. Configure the [`external-metrics-config`](https://github.com/wavefrontHQ/wavefront-kubernetes-adapter/blob/06b75dcf6fd9813a2b8a5a5762be1ae922d35ce7/deploy/manifests/custom-metrics-apiserver-deployment.yaml#L31) adapter property based on the config map and redeploy the adapter.
3. Deploy a [HPA](https://github.com/wavefrontHQ/wavefront-kubernetes-adapter/blob/master/deploy/hpa-examples/hpa-external.yaml) based on an external metric.

## Kubernetes HPA Spec
Refer to the [autoscaling spec](https://godoc.org/k8s.io/api/autoscaling/v2beta1#MetricSpec) for more details on configuring HPAs based on the custom or external metrics APIs.
