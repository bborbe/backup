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

var _ = Describe("BackupCleanerOnlyOnce", func() {
	var ctx context.Context
	var err error
	var mockBackupCleaner *mocks.BackupCleaner
	var backupCleanerOnlyOnce pkg.BackupCleaner
	var backupHost v1.BackupHost

	BeforeEach(func() {
		ctx = context.Background()
		mockBackupCleaner = &mocks.BackupCleaner{}
		backupCleanerOnlyOnce = pkg.NewBackupCleanerOnlyOnce(mockBackupCleaner)
		backupHost = "test-host.example.com"
	})

	JustBeforeEach(func() {
		err = backupCleanerOnlyOnce.Clean(ctx, backupHost)
	})

	Context("when no cleanup is running", func() {
		BeforeEach(func() {
			mockBackupCleaner.CleanReturns(nil)
		})

		It("returns nil error", func() {
			Expect(err).To(BeNil())
		})

		It("calls underlying backup cleaner", func() {
			Expect(mockBackupCleaner.CleanCallCount()).To(Equal(1))
			actualCtx, actualHost := mockBackupCleaner.CleanArgsForCall(0)
			Expect(actualCtx).To(Equal(ctx))
			Expect(actualHost).To(Equal(backupHost))
		})
	})

	Context("when underlying backup cleaner fails", func() {
		BeforeEach(func() {
			mockBackupCleaner.CleanReturns(errors.New("cleanup failed"))
		})

		It("returns the error", func() {
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("cleanup failed"))
		})

		It("calls underlying backup cleaner", func() {
			Expect(mockBackupCleaner.CleanCallCount()).To(Equal(1))
		})
	})

	Context("when cleanup is already running", func() {
		var secondErr error
		var wg sync.WaitGroup
		var freshMockCleaner *mocks.BackupCleaner
		var freshCleanerOnlyOnce pkg.BackupCleaner

		BeforeEach(func() {
			// Create fresh instances
			freshMockCleaner = &mocks.BackupCleaner{}
			freshCleanerOnlyOnce = pkg.NewBackupCleanerOnlyOnce(freshMockCleaner)

			// Make the first cleanup block
			freshMockCleaner.CleanStub = func(ctx context.Context, host v1.BackupHost) error {
				time.Sleep(100 * time.Millisecond)
				return nil
			}
		})

		JustBeforeEach(func() {
			// Start first cleanup in goroutine
			wg.Add(1)
			go func() {
				defer wg.Done()
				err = freshCleanerOnlyOnce.Clean(ctx, backupHost)
			}()

			// Wait a bit then try second cleanup
			time.Sleep(10 * time.Millisecond)
			secondErr = freshCleanerOnlyOnce.Clean(ctx, backupHost)
			wg.Wait()
		})

		It("first cleanup succeeds", func() {
			Expect(err).To(BeNil())
		})

		It("second cleanup returns CleanupAlreadyRunningError", func() {
			Expect(secondErr).ToNot(BeNil())
			Expect(bberrors.Is(secondErr, pkg.CleanupAlreadyRunningError)).To(BeTrue())
		})

		It("only calls underlying backup cleaner once", func() {
			Expect(freshMockCleaner.CleanCallCount()).To(Equal(1))
		})
	})

	Context("when multiple cleanups are attempted concurrently", func() {
		var results []error
		var wg sync.WaitGroup
		var freshMockCleaner *mocks.BackupCleaner
		var freshCleanerOnlyOnce pkg.BackupCleaner

		BeforeEach(func() {
			// Create fresh instances
			freshMockCleaner = &mocks.BackupCleaner{}
			freshCleanerOnlyOnce = pkg.NewBackupCleanerOnlyOnce(freshMockCleaner)

			// Make cleanup take some time
			freshMockCleaner.CleanStub = func(ctx context.Context, host v1.BackupHost) error {
				time.Sleep(50 * time.Millisecond)
				return nil
			}
		})

		JustBeforeEach(func() {
			results = make([]error, 5)

			// Start 5 concurrent cleanups
			for i := 0; i < 5; i++ {
				wg.Add(1)
				go func(index int) {
					defer wg.Done()
					results[index] = freshCleanerOnlyOnce.Clean(ctx, backupHost)
				}(i)
			}
			wg.Wait()
		})

		It("only one cleanup succeeds", func() {
			successCount := 0
			alreadyRunningCount := 0

			for _, result := range results {
				if result == nil {
					successCount++
				} else if bberrors.Is(result, pkg.CleanupAlreadyRunningError) {
					alreadyRunningCount++
				}
			}

			Expect(successCount).To(Equal(1))
			Expect(alreadyRunningCount).To(Equal(4))
		})

		It("only calls underlying backup cleaner once", func() {
			Expect(freshMockCleaner.CleanCallCount()).To(Equal(1))
		})
	})

	Context("when cleanup completes and another is started", func() {
		var secondErr error
		var freshMockCleaner *mocks.BackupCleaner
		var freshCleanerOnlyOnce pkg.BackupCleaner

		BeforeEach(func() {
			// Create fresh instances
			freshMockCleaner = &mocks.BackupCleaner{}
			freshCleanerOnlyOnce = pkg.NewBackupCleanerOnlyOnce(freshMockCleaner)
			freshMockCleaner.CleanReturns(nil)
		})

		JustBeforeEach(func() {
			// First cleanup
			err = freshCleanerOnlyOnce.Clean(ctx, backupHost)
			// Second cleanup after first completes
			secondErr = freshCleanerOnlyOnce.Clean(ctx, backupHost)
		})

		It("both cleanups succeed", func() {
			Expect(err).To(BeNil())
			Expect(secondErr).To(BeNil())
		})

		It("calls underlying backup cleaner twice", func() {
			Expect(freshMockCleaner.CleanCallCount()).To(Equal(2))
		})
	})
})
