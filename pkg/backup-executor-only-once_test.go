// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg_test

import (
	"context"
	"errors"
	"sync"
	"time"

	bberrors "github.com/bborbe/errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	v1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
	"github.com/bborbe/backup/mocks"
	"github.com/bborbe/backup/pkg"
)

var _ = Describe("BackupExecutorOnlyOnce", func() {
	var ctx context.Context
	var err error
	var mockBackupExecutor *mocks.BackupExecutor
	var backupExecutorOnlyOnce pkg.BackupExectuor
	var backupSpec v1.BackupSpec

	BeforeEach(func() {
		ctx = context.Background()
		mockBackupExecutor = &mocks.BackupExecutor{}
		backupExecutorOnlyOnce = pkg.NewBackupExectuorOnlyOnce(mockBackupExecutor)

		backupSpec = v1.BackupSpec{
			Host: "test-host.example.com",
			Port: 22,
			User: "root",
			Dirs: v1.BackupDirs{"/data"},
		}
	})

	JustBeforeEach(func() {
		err = backupExecutorOnlyOnce.Backup(ctx, backupSpec)
	})

	Context("when no backup is running", func() {
		BeforeEach(func() {
			mockBackupExecutor.BackupReturns(nil)
		})

		It("returns nil error", func() {
			Expect(err).To(BeNil())
		})

		It("calls underlying backup executor", func() {
			Expect(mockBackupExecutor.BackupCallCount()).To(Equal(1))
			actualCtx, actualSpec := mockBackupExecutor.BackupArgsForCall(0)
			Expect(actualCtx).To(Equal(ctx))
			Expect(actualSpec).To(Equal(backupSpec))
		})
	})

	Context("when underlying backup executor fails", func() {
		BeforeEach(func() {
			mockBackupExecutor.BackupReturns(errors.New("backup failed"))
		})

		It("returns the error", func() {
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("backup failed"))
		})

		It("calls underlying backup executor", func() {
			Expect(mockBackupExecutor.BackupCallCount()).To(Equal(1))
		})
	})

	Context("when backup is already running", func() {
		var secondErr error
		var wg sync.WaitGroup
		var freshMockExecutor *mocks.BackupExecutor
		var freshExecutorOnlyOnce pkg.BackupExectuor

		BeforeEach(func() {
			// Create fresh instances
			freshMockExecutor = &mocks.BackupExecutor{}
			freshExecutorOnlyOnce = pkg.NewBackupExectuorOnlyOnce(freshMockExecutor)

			// Make the first backup block
			freshMockExecutor.BackupStub = func(ctx context.Context, spec v1.BackupSpec) error {
				time.Sleep(100 * time.Millisecond)
				return nil
			}
		})

		JustBeforeEach(func() {
			// Start first backup in goroutine
			wg.Add(1)
			go func() {
				defer wg.Done()
				err = freshExecutorOnlyOnce.Backup(ctx, backupSpec)
			}()

			// Wait a bit then try second backup
			time.Sleep(10 * time.Millisecond)
			secondErr = freshExecutorOnlyOnce.Backup(ctx, backupSpec)
			wg.Wait()
		})

		It("first backup succeeds", func() {
			Expect(err).To(BeNil())
		})

		It("second backup returns BackupAlreadyRunningError", func() {
			Expect(secondErr).ToNot(BeNil())
			Expect(bberrors.Is(secondErr, pkg.BackupAlreadyRunningError)).To(BeTrue())
		})

		It("only calls underlying backup executor once", func() {
			Expect(freshMockExecutor.BackupCallCount()).To(Equal(1))
		})
	})

	Context("when multiple backups are attempted concurrently", func() {
		var results []error
		var wg sync.WaitGroup
		var freshMockExecutor *mocks.BackupExecutor
		var freshExecutorOnlyOnce pkg.BackupExectuor

		BeforeEach(func() {
			// Create fresh instances
			freshMockExecutor = &mocks.BackupExecutor{}
			freshExecutorOnlyOnce = pkg.NewBackupExectuorOnlyOnce(freshMockExecutor)

			// Make backup take some time
			freshMockExecutor.BackupStub = func(ctx context.Context, spec v1.BackupSpec) error {
				time.Sleep(50 * time.Millisecond)
				return nil
			}
		})

		JustBeforeEach(func() {
			results = make([]error, 5)

			// Start 5 concurrent backups
			for i := 0; i < 5; i++ {
				wg.Add(1)
				go func(index int) {
					defer wg.Done()
					results[index] = freshExecutorOnlyOnce.Backup(ctx, backupSpec)
				}(i)
			}
			wg.Wait()
		})

		It("only one backup succeeds", func() {
			successCount := 0
			alreadyRunningCount := 0

			for _, result := range results {
				if result == nil {
					successCount++
				} else if bberrors.Is(result, pkg.BackupAlreadyRunningError) {
					alreadyRunningCount++
				}
			}

			Expect(successCount).To(Equal(1))
			Expect(alreadyRunningCount).To(Equal(4))
		})

		It("only calls underlying backup executor once", func() {
			Expect(freshMockExecutor.BackupCallCount()).To(Equal(1))
		})
	})

	Context("when backup completes and another is started", func() {
		var secondErr error
		var freshMockExecutor *mocks.BackupExecutor
		var freshExecutorOnlyOnce pkg.BackupExectuor

		BeforeEach(func() {
			// Create fresh instances to avoid interference from previous tests
			freshMockExecutor = &mocks.BackupExecutor{}
			freshExecutorOnlyOnce = pkg.NewBackupExectuorOnlyOnce(freshMockExecutor)
			freshMockExecutor.BackupReturns(nil)
		})

		JustBeforeEach(func() {
			// First backup
			err = freshExecutorOnlyOnce.Backup(ctx, backupSpec)
			// Second backup after first completes
			secondErr = freshExecutorOnlyOnce.Backup(ctx, backupSpec)
		})

		It("both backups succeed", func() {
			Expect(err).To(BeNil())
			Expect(secondErr).To(BeNil())
		})

		It("calls underlying backup executor twice", func() {
			Expect(freshMockExecutor.BackupCallCount()).To(Equal(2))
		})
	})
})
