// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package k8s

import (
	"context"
	"regexp"
	"strings"

	"github.com/bborbe/errors"
	"github.com/bborbe/validation"
)

var replaceLeading = regexp.MustCompile(`^[0-9-]+`)
var replaceFolling = regexp.MustCompile(`[0-9-]+$`)
var replaceIllegalCharacters = regexp.MustCompile(`[^a-z0-9-]+`)
var replaceMultiDash = regexp.MustCompile(`-+`)

// BuildName from the given string. Replace all illegal characters with underscore
func BuildName(names ...string) Name {
	name := strings.Join(names, "-")
	name = strings.ToLower(name)
	name = replaceIllegalCharacters.ReplaceAllString(name, "-")
	name = replaceLeading.ReplaceAllString(name, "")
	name = replaceFolling.ReplaceAllString(name, "")
	name = replaceMultiDash.ReplaceAllString(name, "-")
	return Name(name)
}

type Name string

func (n Name) Validate(ctx context.Context) error {
	if len(n) == 0 {
		return errors.Wrapf(ctx, validation.Error, "name empty")
	}
	if len(n) > 253 {
		return errors.Wrapf(ctx, validation.Error, "name longer than 253")
	}
	return nil
}

func (n Name) String() string {
	return string(n)
}

func (n Name) Bytes() []byte {
	return []byte(n)
}

func (n Name) Ptr() *Name {
	return &n
}

func (n Name) Add(name string) Name {
	return BuildName(n.String(), name)
}
