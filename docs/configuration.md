# Configuration

The Wavefront HPA Adapter is configured via command-line flags and an optional configuration file (for external metrics).

## Flags
The adapter takes the standard Kubernetes generic API server arguments.

Additionally, the following arguments are specific to this adapter:
```
Usage:
      --wavefront-url string                                    Wavefront url of the form https://INSTANCE.wavefront.com
      --wavefront-token string                                  Wavefront API token with permissions to query for points
      --wavefront-metric-prefix string                          Metrics under this prefix are exposed in the custom metrics API. (default "kubernetes.")
      --metrics-relist-interval duration                        interval at which to fetch the list of custom metric names from Wavefront (default 10m0s)
      --api-client-timeout duration                             API client timeout (default 10s)
      --external-metrics-config string                          Configuration file for driving external metrics API
      --log-level string                                        one of info, debug, warn or trace (default "info")      
```

## External Metrics Configuration file

Source: [config.go](https://github.com/wavefrontHQ/wavefront-kubernetes-adapter/blob/master/pkg/config/config.go)

The configuration file is written in YAML and provided using the `--external-metrics-config` flag. The adapter can reload configuration changes at runtime.

A reference example is provided [here](https://github.com/wavefrontHQ/wavefront-kubernetes-adapter/blob/master/deploy/manifests/04-custom-metrics-config-map.yaml).
