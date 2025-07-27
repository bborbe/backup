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

var _ = Describe("CleanAction", func() {
	var ctx context.Context
	var err error
	var cleanAction run.Runnable
	var mockSentryClient *libsentrymocks.SentryClient
	var mockK8sConnector *mocks.K8sConnector
	var mockBackupCleaner *mocks.BackupCleaner
	var targets v1.Targets

	BeforeEach(func() {
		ctx = context.Background()

		// Set up mocks
		mockSentryClient = &libsentrymocks.SentryClient{}
		mockK8sConnector = &mocks.K8sConnector{}
		mockBackupCleaner = &mocks.BackupCleaner{}

		// Create clean action
		cleanAction = pkg.NewCleanAction(
			mockSentryClient,
			mockK8sConnector,
			mockBackupCleaner,
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
					User:     "testuser",
					Dirs:     v1.BackupDirs{"/var/www"},
					Excludes: v1.ParseBackupExcludesFromString("*.log\n"),
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-target-2",
				},
				Spec: v1.BackupSpec{
					Host:     "host2.example.com",
					Port:     v1.BackupPort(22),
					User:     "testuser",
					Dirs:     v1.BackupDirs{"/etc"},
					Excludes: v1.ParseBackupExcludesFromString("*.tmp\n"),
				},
			},
		}
	})

	Context("Run", func() {
		JustBeforeEach(func() {
			err = cleanAction.Run(ctx)
		})

		Context("with successful cleanup execution", func() {
			BeforeEach(func() {
				mockK8sConnector.TargetsReturns(targets, nil)
				mockBackupCleaner.CleanReturns(nil)
			})

			It("returns no error", func() {
				Expect(err).To(BeNil())
			})

			It("calls k8s connector to get targets", func() {
				Expect(mockK8sConnector.TargetsCallCount()).To(Equal(1))
				actualCtx := mockK8sConnector.TargetsArgsForCall(0)
				Expect(actualCtx).To(Equal(ctx))
			})

			It("cleans up all targets", func() {
				Expect(mockBackupCleaner.CleanCallCount()).To(Equal(2))

				// Verify first target cleanup
				actualCtx1, actualHost1 := mockBackupCleaner.CleanArgsForCall(0)
				Expect(actualCtx1).To(Equal(ctx))
				Expect(actualHost1).To(Equal(v1.BackupHost("host1.example.com")))

				// Verify second target cleanup
				actualCtx2, actualHost2 := mockBackupCleaner.CleanArgsForCall(1)
				Expect(actualCtx2).To(Equal(ctx))
				Expect(actualHost2).To(Equal(v1.BackupHost("host2.example.com")))
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

			It("does not attempt any cleanup", func() {
				Expect(mockBackupCleaner.CleanCallCount()).To(Equal(0))
			})

			It("does not capture exceptions", func() {
				Expect(mockSentryClient.CaptureExceptionCallCount()).To(Equal(0))
			})
		})

		Context("when backup cleaner fails for one target", func() {
			BeforeEach(func() {
				mockK8sConnector.TargetsReturns(targets, nil)

				// First cleanup succeeds, second fails
				mockBackupCleaner.CleanReturnsOnCall(0, nil)
				mockBackupCleaner.CleanReturnsOnCall(1, errors.New("cleanup execution failed"))
			})

			It("returns no error (continues processing)", func() {
				Expect(err).To(BeNil())
			})

			It("attempts cleanup for all targets", func() {
				Expect(mockBackupCleaner.CleanCallCount()).To(Equal(2))
			})

			It("captures exception for failed cleanup", func() {
				Expect(mockSentryClient.CaptureExceptionCallCount()).To(Equal(1))

				actualErr, actualHint, actualScope := mockSentryClient.CaptureExceptionArgsForCall(0)
				Expect(actualErr.Error()).To(ContainSubstring("cleanup execution failed"))
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

		Context("when backup cleaner fails for all targets", func() {
			BeforeEach(func() {
				mockK8sConnector.TargetsReturns(targets, nil)
				mockBackupCleaner.CleanReturns(errors.New("cleanup execution failed"))
			})

			It("returns no error (continues processing)", func() {
				Expect(err).To(BeNil())
			})

			It("attempts cleanup for all targets", func() {
				Expect(mockBackupCleaner.CleanCallCount()).To(Equal(2))
			})

			It("captures exceptions for all failed cleanups", func() {
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

			It("does not attempt any cleanup", func() {
				Expect(mockBackupCleaner.CleanCallCount()).To(Equal(0))
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

				// Set up backup cleaner to succeed but we'll cancel context
				mockBackupCleaner.CleanStub = func(ctx context.Context, host v1.BackupHost) error {
					// Cancel context during first cleanup
					if host == "host1.example.com" {
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
				err = cleanAction.Run(cancelCtx)
			})

			It("returns context cancellation error", func() {
				Expect(err).To(Equal(context.Canceled))
			})

			It("stops processing when context is cancelled", func() {
				// Should attempt first cleanup, get cancelled, then return
				Expect(mockBackupCleaner.CleanCallCount()).To(BeNumerically(">=", 1))
			})
		})

		Context("with different host configurations", func() {
			BeforeEach(func() {
				complexTargets := v1.Targets{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "simple-host",
						},
						Spec: v1.BackupSpec{
							Host: "simple",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "fqdn-host",
						},
						Spec: v1.BackupSpec{
							Host: "server.example.com",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "complex-host",
						},
						Spec: v1.BackupSpec{
							Host: "backup-server-01.production.example.com",
						},
					},
				}
				mockK8sConnector.TargetsReturns(complexTargets, nil)
				mockBackupCleaner.CleanReturns(nil)
			})

			It("handles different host name formats correctly", func() {
				Expect(err).To(BeNil())
				Expect(mockBackupCleaner.CleanCallCount()).To(Equal(3))

				// Verify each host is passed correctly
				_, actualHost1 := mockBackupCleaner.CleanArgsForCall(0)
				_, actualHost2 := mockBackupCleaner.CleanArgsForCall(1)
				_, actualHost3 := mockBackupCleaner.CleanArgsForCall(2)

				Expect(actualHost1).To(Equal(v1.BackupHost("simple")))
				Expect(actualHost2).To(Equal(v1.BackupHost("server.example.com")))
				Expect(actualHost3).To(Equal(v1.BackupHost("backup-server-01.production.example.com")))
			})
		})

		Context("exception data verification", func() {
			BeforeEach(func() {
				mockK8sConnector.TargetsReturns(targets, nil)
				mockBackupCleaner.CleanReturnsOnCall(0, errors.New("test cleanup error"))
			})

			It("includes all relevant target data in exception", func() {
				Expect(err).To(BeNil())
				Expect(mockSentryClient.CaptureExceptionCallCount()).To(Equal(1))

				_, actualHint, _ := mockSentryClient.CaptureExceptionArgsForCall(0)

				// Verify all expected fields are present
				Expect(actualHint.Data).To(HaveKey("name"))
				Expect(actualHint.Data).To(HaveKey("host"))
				Expect(actualHint.Data).To(HaveKey("port"))
				Expect(actualHint.Data).To(HaveKey("user"))
				Expect(actualHint.Data).To(HaveKey("dirs"))
				Expect(actualHint.Data).To(HaveKey("excludes"))

				// Verify field values
				data := actualHint.Data.(map[string]interface{})
				Expect(data["name"]).To(Equal("test-target-1"))
				Expect(data["host"]).To(Equal(v1.BackupHost("host1.example.com")))
				Expect(data["port"]).To(Equal(v1.BackupPort(22)))
				Expect(data["user"]).To(Equal(v1.BackupUser("testuser")))
				Expect(data["dirs"]).To(Equal(v1.BackupDirs{"/var/www"}))
				Expect(data["excludes"]).To(Equal(v1.ParseBackupExcludesFromString("*.log\n")))
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
						mockBackupCleaner.CleanReturns(timeoutErr)
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
						mockBackupCleaner.CleanReturns(permissionErr)
					})

					It("captures permission error", func() {
						Expect(err).To(BeNil())
						Expect(mockSentryClient.CaptureExceptionCallCount()).To(Equal(2))

						actualErr, _, _ := mockSentryClient.CaptureExceptionArgsForCall(0)
						Expect(actualErr.Error()).To(ContainSubstring("permission denied"))
					})
				})

				Context("disk space error", func() {
					BeforeEach(func() {
						diskErr := errors.New("no space left on device")
						mockBackupCleaner.CleanReturns(diskErr)
					})

					It("captures disk space error", func() {
						Expect(err).To(BeNil())
						Expect(mockSentryClient.CaptureExceptionCallCount()).To(Equal(2))

						actualErr, _, _ := mockSentryClient.CaptureExceptionArgsForCall(0)
						Expect(actualErr.Error()).To(ContainSubstring("no space left on device"))
					})
				})
			})
		})

		Context("sentry client behavior", func() {
			Context("when sentry client is nil", func() {
				BeforeEach(func() {
					// Test with nil sentry client
					cleanAction = pkg.NewCleanAction(
						nil,
						mockK8sConnector,
						mockBackupCleaner,
					)

					mockK8sConnector.TargetsReturns(targets, nil)
					mockBackupCleaner.CleanReturns(errors.New("test error"))
				})

				It("handles nil sentry client gracefully", func() {
					// Should not panic with nil sentry client
					Expect(func() { _ = cleanAction.Run(ctx) }).NotTo(Panic())
				})
			})
		})

		Context("partial failure scenarios", func() {
			BeforeEach(func() {
				// Create multiple targets with mixed success/failure
				multipleTargets := v1.Targets{
					{
						ObjectMeta: metav1.ObjectMeta{Name: "target-1"},
						Spec:       v1.BackupSpec{Host: "host1.example.com"},
					},
					{
						ObjectMeta: metav1.ObjectMeta{Name: "target-2"},
						Spec:       v1.BackupSpec{Host: "host2.example.com"},
					},
					{
						ObjectMeta: metav1.ObjectMeta{Name: "target-3"},
						Spec:       v1.BackupSpec{Host: "host3.example.com"},
					},
				}
				mockK8sConnector.TargetsReturns(multipleTargets, nil)

				// First and third succeed, second fails
				mockBackupCleaner.CleanReturnsOnCall(0, nil)
				mockBackupCleaner.CleanReturnsOnCall(1, errors.New("cleanup failed for host2"))
				mockBackupCleaner.CleanReturnsOnCall(2, nil)
			})

			It("continues processing all targets despite failures", func() {
				Expect(err).To(BeNil())
				Expect(mockBackupCleaner.CleanCallCount()).To(Equal(3))
			})

			It("captures only the failed target exception", func() {
				Expect(mockSentryClient.CaptureExceptionCallCount()).To(Equal(1))

				actualErr, actualHint, _ := mockSentryClient.CaptureExceptionArgsForCall(0)
				Expect(actualErr.Error()).To(ContainSubstring("cleanup failed for host2"))
				data := actualHint.Data.(map[string]interface{})
				Expect(data["name"]).To(Equal("target-2"))
				Expect(data["host"]).To(Equal(v1.BackupHost("host2.example.com")))
			})
		})

		Context("empty target configurations", func() {
			Context("with target having empty host", func() {
				BeforeEach(func() {
					emptyHostTargets := v1.Targets{
						{
							ObjectMeta: metav1.ObjectMeta{Name: "empty-host-target"},
							Spec: v1.BackupSpec{
								Host: "",
								User: "testuser",
							},
						},
					}
					mockK8sConnector.TargetsReturns(emptyHostTargets, nil)
					mockBackupCleaner.CleanReturns(nil)
				})

				It("processes empty host correctly", func() {
					Expect(err).To(BeNil())
					Expect(mockBackupCleaner.CleanCallCount()).To(Equal(1))

					_, actualHost := mockBackupCleaner.CleanArgsForCall(0)
					Expect(actualHost).To(Equal(v1.BackupHost("")))
				})
			})
		})
	})
})
