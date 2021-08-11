module github.com/IIJ-Global-Solutions-Vietnam/namespace-admission-controller

go 1.16

require (
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/rancher/norman v0.0.0-20210225010917-c7fd1e24145b
	github.com/rancher/types v0.0.0-20210123000350-7cb436b3f0b0
	github.com/sirupsen/logrus v1.8.1
	github.com/slok/kubewebhook/v2 v2.0.1-0.20210427051517-8c5178e3a101
	github.com/stretchr/testify v1.7.0
	golang.org/x/net v0.0.0-20210415231046-e915ea6b2b7d // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	k8s.io/api v0.21.0
	k8s.io/apimachinery v0.21.0
)

replace (
	github.com/coreos/prometheus-operator => github.com/prometheus-operator/prometheus-operator v0.36.0
	k8s.io/api => k8s.io/api v0.21.0
	k8s.io/client-go => k8s.io/client-go v0.21.0
)
