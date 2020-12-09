module github.com/wavefronthq/wavefront-kubernetes-adapter

go 1.13

require (
	github.com/grpc-ecosystem/grpc-gateway v1.12.1 // indirect
	github.com/kubernetes-sigs/custom-metrics-apiserver v0.0.0-20201110135240-8c12d6d92362
	github.com/onsi/gomega v1.8.1 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.6.1
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.19.3
	k8s.io/apimachinery v0.19.3
	k8s.io/client-go v0.19.3
	k8s.io/component-base v0.19.3
	k8s.io/metrics v0.19.3
)
