# Configuration

The Operations for Applications HPA adapter is configured via command-line flags and an optional configuration file (for external metrics).

## Flags
The adapter takes the standard Kubernetes generic API server arguments.

Additionally, the following arguments are specific to this adapter:
```
Usage:
  --wavefront-url string                   Wavefront URL in the format https://YOUR_INSTANCE.wavefront.com.
  --wavefront-token string                 Wavefront API token with permissions to query for points.
  --wavefront-metric-prefix string         Metrics under this prefix are exposed in the custom metrics API. (default "kubernetes")
  --metrics-relist-interval duration       Interval at which to fetch the list of custom metric names from Operations for Applications. (default 10m0s)
  --api-client-timeout duration            Client timeout to Operations for Applications. (default 10s)
  --external-metrics-config string         Configuration file for driving external metrics API.
  --log-level string                       One of info, debug or trace. (default "info")
```

## External Metrics Configuration File

Source: [config.go](/pkg/config/config.go)

The configuration file is written in YAML and provided using the `--external-metrics-config` flag. The adapter can reload configuration changes at runtime.

A reference example is provided [here](/deploy/manifests/04-custom-metrics-config-map.yaml).
