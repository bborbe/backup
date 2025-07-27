// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg_test

import (
	"context"
	"errors"

	"github.com/bborbe/run"
	libsentrymocks "github.com/bborbe/sentry/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
	"github.com/bborbe/backup/mocks"
	"github.com/bborbe/backup/pkg"
)

var _ = Describe("BackupAction", func() {
	var ctx context.Context
	var err error
	var backupAction run.Runnable
	var mockSentryClient *libsentrymocks.SentryClient
	var mockK8sConnector *mocks.K8sConnector
	var mockBackupExecutor *mocks.BackupExecutor
	var targets v1.Targets

	BeforeEach(func() {
		ctx = context.Background()

		// Set up mocks
		mockSentryClient = &libsentrymocks.SentryClient{}
		mockK8sConnector = &mocks.K8sConnector{}
		mockBackupExecutor = &mocks.BackupExecutor{}

		// Create backup action
		backupAction = pkg.NewBackupAction(
			mockSentryClient,
			mockK8sConnector,
			mockBackupExecutor,
		)

		// Set up default successful targets
		targets = v1.Targets{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-target-1",
				},
				Spec: v1.BackupSpec{
					Host:     "host1.example.com",
					Port:     v1.BackupPort(22),
					User:     v1.BackupUser("testuser"),
					Dirs:     v1.BackupDirs{"/var/www"},
					Excludes: v1.ParseBackupExcludes([]string{"*.log"}),
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-target-2",
				},
				Spec: v1.BackupSpec{
					Host:     "host2.example.com",
					Port:     v1.BackupPort(22),
					User:     v1.BackupUser("testuser"),
					Dirs:     v1.BackupDirs{"/etc"},
					Excludes: v1.ParseBackupExcludes([]string{"*.tmp"}),
				},
			},
		}
	})

	Context("Run", func() {
		JustBeforeEach(func() {
			err = backupAction.Run(ctx)
		})

		Context("with successful backup execution", func() {
			BeforeEach(func() {
				mockK8sConnector.TargetsReturns(targets, nil)
				mockBackupExecutor.BackupReturns(nil)
			})

			It("returns no error", func() {
				Expect(err).To(BeNil())
			})

			It("calls k8s connector to get targets", func() {
				Expect(mockK8sConnector.TargetsCallCount()).To(Equal(1))
				actualCtx := mockK8sConnector.TargetsArgsForCall(0)
				Expect(actualCtx).To(Equal(ctx))
			})

			It("backs up all targets", func() {
				Expect(mockBackupExecutor.BackupCallCount()).To(Equal(2))

				// Verify first target backup
				actualCtx1, actualSpec1 := mockBackupExecutor.BackupArgsForCall(0)
				Expect(actualCtx1).To(Equal(ctx))
				Expect(actualSpec1.Host).To(Equal(v1.BackupHost("host1.example.com")))
				Expect(actualSpec1.User).To(Equal(v1.BackupUser("testuser")))

				// Verify second target backup
				actualCtx2, actualSpec2 := mockBackupExecutor.BackupArgsForCall(1)
				Expect(actualCtx2).To(Equal(ctx))
				Expect(actualSpec2.Host).To(Equal(v1.BackupHost("host2.example.com")))
				Expect(actualSpec2.User).To(Equal(v1.BackupUser("testuser")))
			})

			It("does not capture any exceptions", func() {
				Expect(mockSentryClient.CaptureExceptionCallCount()).To(Equal(0))
			})
		})

		Context("when k8s connector fails to get targets", func() {
			BeforeEach(func() {
				mockK8sConnector.TargetsReturns(nil, errors.New("k8s connection failed"))
			})

			It("returns error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("get target failed"))
				Expect(err.Error()).To(ContainSubstring("k8s connection failed"))
			})

			It("calls k8s connector", func() {
				Expect(mockK8sConnector.TargetsCallCount()).To(Equal(1))
			})

			It("does not attempt any backups", func() {
				Expect(mockBackupExecutor.BackupCallCount()).To(Equal(0))
			})

			It("does not capture exceptions", func() {
				Expect(mockSentryClient.CaptureExceptionCallCount()).To(Equal(0))
			})
		})

		Context("when backup executor fails for one target", func() {
			BeforeEach(func() {
				mockK8sConnector.TargetsReturns(targets, nil)

				// First backup succeeds, second fails
				mockBackupExecutor.BackupReturnsOnCall(0, nil)
				mockBackupExecutor.BackupReturnsOnCall(1, errors.New("backup execution failed"))
			})

			It("returns no error (continues processing)", func() {
				Expect(err).To(BeNil())
			})

			It("attempts backup for all targets", func() {
				Expect(mockBackupExecutor.BackupCallCount()).To(Equal(2))
			})

			It("captures exception for failed backup", func() {
				Expect(mockSentryClient.CaptureExceptionCallCount()).To(Equal(1))

				actualErr, actualHint, actualScope := mockSentryClient.CaptureExceptionArgsForCall(0)
				Expect(actualErr.Error()).To(ContainSubstring("backup execution failed"))
				Expect(actualHint).NotTo(BeNil())
				Expect(actualHint.Context).To(Equal(ctx))
				Expect(actualScope).To(BeNil())

				// Verify exception data includes target information
				data := actualHint.Data.(map[string]interface{})
				Expect(data).To(HaveKey("name"))
				Expect(data["name"]).To(Equal("test-target-2"))
				Expect(data).To(HaveKey("host"))
				Expect(data["host"]).To(Equal(v1.BackupHost("host2.example.com")))
				Expect(data).To(HaveKey("user"))
				Expect(data["user"]).To(Equal(v1.BackupUser("testuser")))
			})
		})

		Context("when backup executor fails for all targets", func() {
			BeforeEach(func() {
				mockK8sConnector.TargetsReturns(targets, nil)
				mockBackupExecutor.BackupReturns(errors.New("backup execution failed"))
			})

			It("returns no error (continues processing)", func() {
				Expect(err).To(BeNil())
			})

			It("attempts backup for all targets", func() {
				Expect(mockBackupExecutor.BackupCallCount()).To(Equal(2))
			})

			It("captures exceptions for all failed backups", func() {
				Expect(mockSentryClient.CaptureExceptionCallCount()).To(Equal(2))
			})
		})

		Context("with no targets", func() {
			BeforeEach(func() {
				mockK8sConnector.TargetsReturns(v1.Targets{}, nil)
			})

			It("returns no error", func() {
				Expect(err).To(BeNil())
			})

			It("calls k8s connector", func() {
				Expect(mockK8sConnector.TargetsCallCount()).To(Equal(1))
			})

			It("does not attempt any backups", func() {
				Expect(mockBackupExecutor.BackupCallCount()).To(Equal(0))
			})

			It("does not capture any exceptions", func() {
				Expect(mockSentryClient.CaptureExceptionCallCount()).To(Equal(0))
			})
		})

		Context("with context cancellation", func() {
			var cancelCtx context.Context
			var cancel context.CancelFunc

			BeforeEach(func() {
				cancelCtx, cancel = context.WithCancel(ctx)

				// Set up targets
				mockK8sConnector.TargetsReturns(targets, nil)

				// Set up backup executor to succeed but we'll cancel context
				mockBackupExecutor.BackupStub = func(ctx context.Context, spec v1.BackupSpec) error {
					// Cancel context during first backup
					if spec.Host == "host1.example.com" {
						cancel()
					}
					select {
					case <-ctx.Done():
						return ctx.Err()
					default:
						return nil
					}
				}
			})

			JustBeforeEach(func() {
				err = backupAction.Run(cancelCtx)
			})

			It("returns context cancellation error", func() {
				Expect(err).To(Equal(context.Canceled))
			})

			It("stops processing when context is cancelled", func() {
				// Should attempt first backup, get cancelled, then return
				Expect(mockBackupExecutor.BackupCallCount()).To(BeNumerically(">=", 1))
			})
		})

		Context("with complex target configurations", func() {
			BeforeEach(func() {
				complexTargets := v1.Targets{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "complex-target",
						},
						Spec: v1.BackupSpec{
							Host:     "complex.example.com",
							Port:     v1.BackupPort(2222),
							User:     v1.BackupUser("backupuser"),
							Dirs:     v1.BackupDirs{"/var/www", "/etc", "/opt/data"},
							Excludes: v1.ParseBackupExcludes([]string{"*.log", "*.tmp", "*.cache"}),
						},
					},
				}
				mockK8sConnector.TargetsReturns(complexTargets, nil)
				mockBackupExecutor.BackupReturns(nil)
			})

			It("handles complex configurations correctly", func() {
				Expect(err).To(BeNil())
				Expect(mockBackupExecutor.BackupCallCount()).To(Equal(1))

				_, actualSpec := mockBackupExecutor.BackupArgsForCall(0)
				Expect(actualSpec.Host).To(Equal(v1.BackupHost("complex.example.com")))
				Expect(actualSpec.Port).To(Equal(v1.BackupPort(2222)))
				Expect(actualSpec.User).To(Equal(v1.BackupUser("backupuser")))
				Expect(actualSpec.Dirs).To(HaveLen(3))
				Expect(actualSpec.Dirs).To(ContainElement(v1.BackupDir("/var/www")))
				Expect(actualSpec.Dirs).To(ContainElement(v1.BackupDir("/etc")))
				Expect(actualSpec.Dirs).To(ContainElement(v1.BackupDir("/opt/data")))
			})
		})

		Context("exception data verification", func() {
			BeforeEach(func() {
				mockK8sConnector.TargetsReturns(targets, nil)
				mockBackupExecutor.BackupReturnsOnCall(0, errors.New("test error"))
			})

			It("includes all relevant target data in exception", func() {
				Expect(err).To(BeNil())
				Expect(mockSentryClient.CaptureExceptionCallCount()).To(Equal(1))

				_, actualHint, _ := mockSentryClient.CaptureExceptionArgsForCall(0)

				// Verify all expected fields are present
				data := actualHint.Data.(map[string]interface{})
				Expect(data).To(HaveKey("name"))
				Expect(data).To(HaveKey("host"))
				Expect(data).To(HaveKey("port"))
				Expect(data).To(HaveKey("user"))
				Expect(data).To(HaveKey("dirs"))
				Expect(data).To(HaveKey("excludes"))

				// Verify field values
				Expect(data["name"]).To(Equal("test-target-1"))
				Expect(data["host"]).To(Equal(v1.BackupHost("host1.example.com")))
				Expect(data["port"]).To(Equal(v1.BackupPort(22)))
				Expect(data["user"]).To(Equal(v1.BackupUser("testuser")))
				Expect(data["dirs"]).To(Equal(v1.BackupDirs{"/var/www"}))
				Expect(data["excludes"]).To(Equal(v1.ParseBackupExcludes([]string{"*.log"})))
			})
		})

		Context("error handling patterns", func() {
			Context("with different error types", func() {
				BeforeEach(func() {
					mockK8sConnector.TargetsReturns(targets, nil)
				})

				Context("timeout error", func() {
					BeforeEach(func() {
						timeoutErr := context.DeadlineExceeded
						mockBackupExecutor.BackupReturns(timeoutErr)
					})

					It("captures timeout error", func() {
						Expect(err).To(BeNil())
						Expect(mockSentryClient.CaptureExceptionCallCount()).To(Equal(2))

						actualErr, _, _ := mockSentryClient.CaptureExceptionArgsForCall(0)
						Expect(actualErr).To(Equal(context.DeadlineExceeded))
					})
				})

				Context("permission error", func() {
					BeforeEach(func() {
						permissionErr := errors.New("permission denied")
						mockBackupExecutor.BackupReturns(permissionErr)
					})

					It("captures permission error", func() {
						Expect(err).To(BeNil())
						Expect(mockSentryClient.CaptureExceptionCallCount()).To(Equal(2))

						actualErr, _, _ := mockSentryClient.CaptureExceptionArgsForCall(0)
						Expect(actualErr.Error()).To(ContainSubstring("permission denied"))
					})
				})
			})
		})

	})
})
