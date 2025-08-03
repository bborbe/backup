// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg_test

import (
	"context"
	"errors"

	bberrors "github.com/bborbe/errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
	"github.com/bborbe/backup/mocks"
	"github.com/bborbe/backup/pkg"
)

var _ = Describe("TargetFinder", func() {
	var ctx context.Context
	var err error
	var mockK8sConnector *mocks.K8sConnector
	var targets v1.Targets

	BeforeEach(func() {
		ctx = context.Background()
		mockK8sConnector = &mocks.K8sConnector{}

		targets = v1.Targets{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "target-name-1",
				},
				Spec: v1.BackupSpec{
					Host: "host1.example.com",
					Port: 22,
					User: "root",
					Dirs: v1.BackupDirs{"/data"},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "target-name-2",
				},
				Spec: v1.BackupSpec{
					Host: "host2.example.com",
					Port: 22,
					User: "root",
					Dirs: v1.BackupDirs{"/backup"},
				},
			},
		}
	})

	Describe("TargetFinderByHostname", func() {
		var targetFinder pkg.TargetFinder
		var result *v1.Target

		BeforeEach(func() {
			targetFinder = pkg.NewTargetFinderByHostname(mockK8sConnector)
		})

		JustBeforeEach(func() {
			result, err = targetFinder.Target(ctx, "host1.example.com")
		})

		Context("when targets exist", func() {
			BeforeEach(func() {
				mockK8sConnector.TargetsReturns(targets, nil)
			})

			It("returns nil error", func() {
				Expect(err).To(BeNil())
			})

			It("calls Targets", func() {
				Expect(mockK8sConnector.TargetsCallCount()).To(Equal(1))
			})

			It("returns the target with matching hostname", func() {
				Expect(result).ToNot(BeNil())
				Expect(result.Name).To(Equal("target-name-1"))
				Expect(result.Spec.Host.String()).To(Equal("host1.example.com"))
			})
		})

		Context("when hostname doesn't exist", func() {
			BeforeEach(func() {
				mockK8sConnector.TargetsReturns(targets, nil)
			})

			JustBeforeEach(func() {
				result, err = targetFinder.Target(ctx, "nonexistent.example.com")
			})

			It("returns TargetNotFoundError", func() {
				Expect(err).ToNot(BeNil())
				Expect(bberrors.Is(err, pkg.TargetNotFoundError)).To(BeTrue())
			})

			It("returns nil result", func() {
				Expect(result).To(BeNil())
			})
		})

		Context("when k8s connector fails", func() {
			BeforeEach(func() {
				mockK8sConnector.TargetsReturns(nil, errors.New("k8s error"))
			})

			It("returns error", func() {
				Expect(err).ToNot(BeNil())
			})

			It("returns nil result", func() {
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("TargetFinderList", func() {
		var targetFinderList pkg.TargetFinderList
		var mockTargetFinder1 *mocks.TargetFinder
		var mockTargetFinder2 *mocks.TargetFinder
		var result *v1.Target

		BeforeEach(func() {
			mockTargetFinder1 = &mocks.TargetFinder{}
			mockTargetFinder2 = &mocks.TargetFinder{}
			targetFinderList = pkg.TargetFinderList{
				mockTargetFinder1,
				mockTargetFinder2,
			}
		})

		JustBeforeEach(func() {
			result, err = targetFinderList.Target(ctx, "test-target")
		})

		Context("when first finder succeeds", func() {
			BeforeEach(func() {
				mockTargetFinder1.TargetReturns(&targets[0], nil)
			})

			It("returns nil error", func() {
				Expect(err).To(BeNil())
			})

			It("returns target from first finder", func() {
				Expect(result).ToNot(BeNil())
				Expect(result.Name).To(Equal("target-name-1"))
			})

			It("calls first finder", func() {
				Expect(mockTargetFinder1.TargetCallCount()).To(Equal(1))
			})

			It("doesn't call second finder", func() {
				Expect(mockTargetFinder2.TargetCallCount()).To(Equal(0))
			})
		})

		Context("when first finder fails but second succeeds", func() {
			BeforeEach(func() {
				mockTargetFinder1.TargetReturns(nil, errors.New("not found"))
				mockTargetFinder2.TargetReturns(&targets[1], nil)
			})

			It("returns nil error", func() {
				Expect(err).To(BeNil())
			})

			It("returns target from second finder", func() {
				Expect(result).ToNot(BeNil())
				Expect(result.Name).To(Equal("target-name-2"))
			})

			It("calls both finders", func() {
				Expect(mockTargetFinder1.TargetCallCount()).To(Equal(1))
				Expect(mockTargetFinder2.TargetCallCount()).To(Equal(1))
			})
		})

		Context("when all finders fail", func() {
			BeforeEach(func() {
				mockTargetFinder1.TargetReturns(nil, errors.New("not found"))
				mockTargetFinder2.TargetReturns(nil, errors.New("not found"))
			})

			It("returns TargetNotFoundError", func() {
				Expect(err).ToNot(BeNil())
				Expect(bberrors.Is(err, pkg.TargetNotFoundError)).To(BeTrue())
			})

			It("returns nil result", func() {
				Expect(result).To(BeNil())
			})

			It("calls all finders", func() {
				Expect(mockTargetFinder1.TargetCallCount()).To(Equal(1))
				Expect(mockTargetFinder2.TargetCallCount()).To(Equal(1))
			})
		})
	})

	Describe("NewCombinedTargetFinder", func() {
		var combinedTargetFinder pkg.TargetFinder
		var result *v1.Target

		BeforeEach(func() {
			combinedTargetFinder = pkg.NewCombinedTargetFinder(mockK8sConnector)
		})

		Context("when target found by name", func() {
			BeforeEach(func() {
				mockK8sConnector.TargetReturns(&targets[0], nil)
			})

			JustBeforeEach(func() {
				result, err = combinedTargetFinder.Target(ctx, "target-name-1")
			})

			It("returns nil error", func() {
				Expect(err).To(BeNil())
			})

			It("returns the target", func() {
				Expect(result).ToNot(BeNil())
				Expect(result.Name).To(Equal("target-name-1"))
			})

			It("calls Target (name-based lookup)", func() {
				Expect(mockK8sConnector.TargetCallCount()).To(Equal(1))
			})

			It("doesn't call Targets (hostname-based lookup)", func() {
				Expect(mockK8sConnector.TargetsCallCount()).To(Equal(0))
			})
		})

		Context("when target not found by name but found by hostname", func() {
			BeforeEach(func() {
				mockK8sConnector.TargetReturns(nil, errors.New("not found by name"))
				mockK8sConnector.TargetsReturns(targets, nil)
			})

			JustBeforeEach(func() {
				result, err = combinedTargetFinder.Target(ctx, "host1.example.com")
			})

			It("returns nil error", func() {
				Expect(err).To(BeNil())
			})

			It("returns the target found by hostname", func() {
				Expect(result).ToNot(BeNil())
				Expect(result.Name).To(Equal("target-name-1"))
				Expect(result.Spec.Host.String()).To(Equal("host1.example.com"))
			})

			It("calls both Target and Targets", func() {
				Expect(mockK8sConnector.TargetCallCount()).To(Equal(1))
				Expect(mockK8sConnector.TargetsCallCount()).To(Equal(1))
			})
		})

		Context("when target not found by name or hostname", func() {
			BeforeEach(func() {
				mockK8sConnector.TargetReturns(nil, errors.New("not found by name"))
				mockK8sConnector.TargetsReturns(targets, nil)
			})

			JustBeforeEach(func() {
				result, err = combinedTargetFinder.Target(ctx, "nonexistent")
			})

			It("returns TargetNotFoundError", func() {
				Expect(err).ToNot(BeNil())
				Expect(bberrors.Is(err, pkg.TargetNotFoundError)).To(BeTrue())
			})

			It("returns nil result", func() {
				Expect(result).To(BeNil())
			})
		})
	})
})
