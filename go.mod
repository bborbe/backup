module github.com/bborbe/backup

go 1.22.3

exclude (
	k8s.io/api v0.29.0
	k8s.io/api v0.29.1
	k8s.io/api v0.29.2
	k8s.io/api v0.29.3
	k8s.io/api v0.29.4
	k8s.io/api v0.30.0
	k8s.io/apiextensions-apiserver v0.29.0
	k8s.io/apiextensions-apiserver v0.29.1
	k8s.io/apiextensions-apiserver v0.29.2
	k8s.io/apiextensions-apiserver v0.29.3
	k8s.io/apiextensions-apiserver v0.29.4
	k8s.io/apiextensions-apiserver v0.30.0
	k8s.io/apimachinery v0.29.0
	k8s.io/apimachinery v0.29.1
	k8s.io/apimachinery v0.29.2
	k8s.io/apimachinery v0.29.3
	k8s.io/apimachinery v0.29.4
	k8s.io/apimachinery v0.30.0
	k8s.io/client-go v0.29.0
	k8s.io/client-go v0.29.1
	k8s.io/client-go v0.29.2
	k8s.io/client-go v0.29.3
	k8s.io/client-go v0.29.4
	k8s.io/client-go v0.30.0
	k8s.io/code-generator v0.29.0
	k8s.io/code-generator v0.29.1
	k8s.io/code-generator v0.29.2
	k8s.io/code-generator v0.29.3
	k8s.io/code-generator v0.29.4
	k8s.io/code-generator v0.30.0
)

replace github.com/imdario/mergo => github.com/imdario/mergo v0.3.16

replace github.com/antlr/antlr4/runtime/Go/antlr/v4 => github.com/antlr4-go/antlr/v4 v4.13.0

require (
	github.com/bborbe/collection v1.4.0
	github.com/bborbe/cron v1.1.0
	github.com/bborbe/errors v1.2.0
	github.com/bborbe/http v1.1.0
	github.com/bborbe/k8s v1.0.0
	github.com/bborbe/log v1.0.0
	github.com/bborbe/run v1.5.2
	github.com/bborbe/sentry v1.6.0
	github.com/bborbe/service v1.3.0
	github.com/bborbe/time v1.2.0
	github.com/bborbe/validation v1.0.0
	github.com/getsentry/sentry-go v0.27.0
	github.com/golang/glog v1.2.1
	github.com/google/addlicense v1.1.1
	github.com/gorilla/mux v1.8.1
	github.com/incu6us/goimports-reviser v0.1.6
	github.com/kisielk/errcheck v1.7.0
	github.com/maxbrunsfeld/counterfeiter/v6 v6.8.1
	github.com/onsi/ginkgo/v2 v2.17.3
	github.com/onsi/gomega v1.33.1
	github.com/prometheus/client_golang v1.19.1
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616
	golang.org/x/vuln v1.1.0
	k8s.io/apiextensions-apiserver v0.28.9
	k8s.io/apimachinery v0.28.9
	k8s.io/client-go v0.28.9
	k8s.io/code-generator v0.28.9
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1
)

require (
	github.com/bborbe/argument/v2 v2.0.4 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bmatcuk/doublestar/v4 v4.6.1 // indirect
	github.com/certifi/gocertifi v0.0.0-20210507211836-431795d63e8d // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emicklei/go-restful/v3 v3.12.0 // indirect
	github.com/evanphx/json-patch v5.9.0+incompatible // indirect
	github.com/getsentry/raven-go v0.2.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/pprof v0.0.0-20240508145209-1db217f89380 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.53.0 // indirect
	github.com/prometheus/procfs v0.14.0 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/oauth2 v0.20.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/term v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	golang.org/x/tools v0.21.0 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.28.9 // indirect
	k8s.io/gengo v0.0.0-20240404160639-a0386bf69313 // indirect
	k8s.io/gengo/v2 v2.0.0-20240404160639-a0386bf69313 // indirect
	k8s.io/klog/v2 v2.120.1 // indirect
	k8s.io/kube-openapi v0.0.0-20240430033511-f0e62f92d13f // indirect
	k8s.io/utils v0.0.0-20240502163921-fe8a2dddb1d0 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)
