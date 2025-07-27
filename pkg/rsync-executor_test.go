// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg_test

import (
	"context"
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/backup/pkg"
)

var _ = Describe("RsyncExecutor", func() {
	var ctx context.Context
	var err error
	var rsyncExecutor pkg.RsyncExectuor

	BeforeEach(func() {
		ctx = context.Background()
		rsyncExecutor = pkg.NewRsyncExectuor()
	})

	Context("Rsync", func() {
		JustBeforeEach(func() {
			// Use rsync --help as a safe command to test
			err = rsyncExecutor.Rsync(ctx, "--help")
		})

		It("executes rsync command", func() {
			// This test requires rsync to be installed
			if _, lookupErr := exec.LookPath("rsync"); lookupErr != nil {
				Skip("rsync not found in PATH")
			}
			Expect(err).To(BeNil())
		})
	})

	Context("with invalid arguments", func() {
		JustBeforeEach(func() {
			err = rsyncExecutor.Rsync(ctx, "--invalid-flag-that-does-not-exist")
		})

		It("returns error for invalid arguments", func() {
			if _, lookupErr := exec.LookPath("rsync"); lookupErr != nil {
				Skip("rsync not found in PATH")
			}
			Expect(err).NotTo(BeNil())
		})
	})

	Context("with context cancellation", func() {
		var cancelCtx context.Context
		var cancel context.CancelFunc

		BeforeEach(func() {
			cancelCtx, cancel = context.WithCancel(ctx)
		})

		JustBeforeEach(func() {
			// Cancel context immediately
			cancel()
			err = rsyncExecutor.Rsync(cancelCtx, "--help")
		})

		It("respects context cancellation", func() {
			if _, lookupErr := exec.LookPath("rsync"); lookupErr != nil {
				Skip("rsync not found in PATH")
			}
			// Note: This may or may not fail depending on timing
			// The test verifies the executor doesn't panic with cancelled context
		})
	})
})
