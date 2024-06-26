//go:build tools
// +build tools

package tools

import (
	_ "github.com/google/addlicense"
	_ "github.com/incu6us/goimports-reviser"
	_ "github.com/kisielk/errcheck"
	_ "github.com/maxbrunsfeld/counterfeiter/v6"
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/vuln/cmd/govulncheck"
	_ "k8s.io/code-generator"
)
