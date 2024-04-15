// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package v1

import (
	"context"
	"reflect"

	"github.com/bborbe/errors"
	"github.com/bborbe/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Targets []Target

func (a Targets) Specs() BackupSpecs {
	var result BackupSpecs
	for _, aa := range a {
		result = append(result, aa.Spec)
	}
	return result
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Target describes a database.
type Target struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec BackupSpec `json:"spec"`
}

type BackupSpecs []BackupSpec

// BackupSpec is the spec for a Foo resource
type BackupSpec struct {
	Name        string            `json:"name" yaml:"name"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Expression  string            `json:"expression,omitempty" yaml:"expression,omitempty"`
	For         string            `json:"for,omitempty" yaml:"for,omitempty"`
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
}

func (a BackupSpec) Equal(backup BackupSpec) bool {
	if a.Name != backup.Name {
		return false
	}
	if a.Expression != backup.Expression {
		return false
	}
	if a.For != backup.For {
		return false
	}
	if !reflect.DeepEqual(a.Annotations, backup.Annotations) {
		return false
	}
	if !reflect.DeepEqual(a.Labels, backup.Labels) {
		return false
	}
	return true
}

func (a BackupSpec) Validate(ctx context.Context) error {
	if a.Name == "" {
		return errors.Wrap(ctx, validation.Error, "Name is empty")
	}
	if a.Expression == "" {
		return errors.Wrap(ctx, validation.Error, "Expression is empty")
	}
	if len(a.Annotations) == 0 {
		return errors.Wrap(ctx, validation.Error, "Annotations is empty")
	}
	if len(a.Labels) == 0 {
		return errors.Wrap(ctx, validation.Error, "Labels is empty")
	}
	return nil
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TargetList is a list of Target resources
type TargetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Target `json:"items"`
}
