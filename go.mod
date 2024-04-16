module github.com/bborbe/backup

go 1.22.2

exclude (
	k8s.io/api v0.29.0
	k8s.io/api v0.29.1
	k8s.io/api v0.29.2
	k8s.io/api v0.29.3
	k8s.io/apiextensions-apiserver v0.29.0
	k8s.io/apiextensions-apiserver v0.29.1
	k8s.io/apiextensions-apiserver v0.29.2
	k8s.io/apiextensions-apiserver v0.29.3
	k8s.io/apimachinery v0.29.0
	k8s.io/apimachinery v0.29.1
	k8s.io/apimachinery v0.29.2
	k8s.io/apimachinery v0.29.3
	k8s.io/client-go v0.29.0
	k8s.io/client-go v0.29.1
	k8s.io/client-go v0.29.2
	k8s.io/client-go v0.29.3
	k8s.io/code-generator v0.29.0
	k8s.io/code-generator v0.29.1
	k8s.io/code-generator v0.29.2
	k8s.io/code-generator v0.29.3
)

require (
	github.com/bborbe/collection v1.3.1
	github.com/bborbe/errors v1.2.0
	github.com/bborbe/http v1.1.0
	github.com/bborbe/log v1.0.0
	github.com/bborbe/run v1.5.0
	github.com/bborbe/sentry v1.0.0
	github.com/bborbe/service v1.0.0
	github.com/bborbe/time v1.1.1
	github.com/bborbe/validation v1.0.0
	github.com/golang/glog v1.2.1
	github.com/google/addlicense v1.1.1
	github.com/gorilla/mux v1.8.1
	github.com/incu6us/goimports-reviser v0.1.6
	github.com/kisielk/errcheck v1.7.0
	github.com/maxbrunsfeld/counterfeiter/v6 v6.8.1
	github.com/onsi/ginkgo/v2 v2.17.1
	github.com/onsi/gomega v1.32.0
	github.com/prometheus/client_golang v1.19.0
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616
	golang.org/x/vuln v1.0.4
	k8s.io/apiextensions-apiserver v0.28.8
	k8s.io/apimachinery v0.28.8
	k8s.io/client-go v0.28.8
	k8s.io/code-generator v0.28.8
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3
)

require (
	github.com/bborbe/argument/v2 v2.0.4 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bmatcuk/doublestar/v4 v4.6.1 // indirect
	github.com/certifi/gocertifi v0.0.0-20210507211836-431795d63e8d // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emicklei/go-restful/v3 v3.9.0 // indirect
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/getsentry/raven-go v0.2.0 // indirect
	github.com/getsentry/sentry-go v0.27.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/pprof v0.0.0-20240320155624-b11c3daa6f07 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_model v0.6.0 // indirect
	github.com/prometheus/common v0.51.1 // indirect
	github.com/prometheus/procfs v0.13.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/oauth2 v0.18.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/term v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/tools v0.20.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.28.8 // indirect
	k8s.io/gengo v0.0.0-20220902162205-c0856e24416d // indirect
	k8s.io/klog/v2 v2.100.1 // indirect
	k8s.io/kube-openapi v0.0.0-20230717233707-2695361300d9 // indirect
	k8s.io/utils v0.0.0-20230406110748-d93618cff8a2 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)