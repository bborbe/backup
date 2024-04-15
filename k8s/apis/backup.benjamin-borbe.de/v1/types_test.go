// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package v1_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	v1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
)

var _ = Describe("BackupSpec", func() {
	var a, b v1.BackupSpec
	BeforeEach(func() {
		a = v1.BackupSpec{
			Name: "my-backup",
			Annotations: map[string]string{
				"a": "b",
			},
			Expression: "true",
			For:        "1h",
			Labels: map[string]string{
				"c": "d",
			},
		}
		b = *a.DeepCopy()
	})
	Context("Equal", func() {
		var result bool
		JustBeforeEach(func() {
			result = a.Equal(b)
		})
		Context("everything is equal", func() {
			It("is equal", func() {
				Expect(result).To(BeTrue())
			})
		})
		Context("Name not equal", func() {
			BeforeEach(func() {
				b.Name = "banana"
			})
			It("is not equal", func() {
				Expect(result).To(BeFalse())
			})
		})
		Context("For not equal", func() {
			BeforeEach(func() {
				b.For = "2h"
			})
			It("is not equal", func() {
				Expect(result).To(BeFalse())
			})
		})
		Context("Expression not equal", func() {
			BeforeEach(func() {
				b.Expression = "false"
			})
			It("is not equal", func() {
				Expect(result).To(BeFalse())
			})
		})
		Context("Labels not equal", func() {
			BeforeEach(func() {
				b.Labels["newlabel"] = "banana"
			})
			It("is not equal", func() {
				Expect(result).To(BeFalse())
			})
		})
		Context("Annotations not equal", func() {
			BeforeEach(func() {
				b.Annotations["newanno"] = "banana"
			})
			It("is not equal", func() {
				Expect(result).To(BeFalse())
			})
		})
	})
})
