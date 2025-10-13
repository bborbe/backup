// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	libtime "github.com/bborbe/time"
	libtimetest "github.com/bborbe/time/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	v1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
	"github.com/bborbe/backup/mocks"
	"github.com/bborbe/backup/pkg"
)

var _ = Describe("BackupExecutor", func() {
	var ctx context.Context
	var err error
	var backupExecutor pkg.BackupExectuor
	var mockRsyncExecutor *mocks.RsyncExecutor
	var currentTime libtime.CurrentTime
	var backupRootDir pkg.Path
	var sshPrivateKey pkg.SSHPrivateKey
	var backupSpec v1.BackupSpec
	var tempDir string
	var fixedTime time.Time

	BeforeEach(func() {
		ctx = context.Background()

		// Create temporary directory for testing
		tempDir, err = os.MkdirTemp("", "backup-test-*")
		Expect(err).To(BeNil())

		// Set up fixed time for deterministic testing
		fixedTime = libtimetest.ParseDateTime("2023-12-25T12:00:00Z").Time()
		currentTime = libtime.NewCurrentTime()
		currentTime.SetNow(fixedTime)

		// Set up mocks and test data
		mockRsyncExecutor = &mocks.RsyncExecutor{}
		backupRootDir = pkg.Path(tempDir)
		sshPrivateKey = pkg.SSHPrivateKey("/path/to/ssh/key")

		backupSpec = v1.BackupSpec{
			Host:     "test-host.example.com",
			Port:     v1.BackupPort(22),
			User:     "testuser",
			Dirs:     v1.BackupDirs{"/var/www", "/etc"},
			Excludes: v1.ParseBackupExcludesFromString("*.log\n*.tmp\n"),
		}

		// Create backup executor
		backupExecutor = pkg.NewBackupExectuor(
			currentTime,
			mockRsyncExecutor,
			backupRootDir,
			sshPrivateKey,
		)
	})

	AfterEach(func() {
		// Clean up temporary directory
		_ = os.RemoveAll(tempDir)
	})

	Context("Backup", func() {
		JustBeforeEach(func() {
			err = backupExecutor.Backup(ctx, backupSpec)
		})

		Context("with valid backup spec", func() {
			BeforeEach(func() {
				// Mock successful rsync execution
				mockRsyncExecutor.RsyncReturns(nil)
			})

			Context("when no backup exists yet", func() {
				It("returns no error", func() {
					Expect(err).To(BeNil())
				})

				It("calls rsync with correct parameters", func() {
					Expect(mockRsyncExecutor.RsyncCallCount()).To(Equal(1))

					actualCtx, actualArgs := mockRsyncExecutor.RsyncArgsForCall(0)
					Expect(actualCtx).To(Equal(ctx))

					// Verify rsync arguments contain expected flags
					Expect(actualArgs).To(ContainElement("-a"))
					Expect(actualArgs).To(ContainElement("-m"))
					Expect(actualArgs).To(ContainElement("--progress"))
					Expect(actualArgs).To(ContainElement("--whole-file"))
					Expect(actualArgs).To(ContainElement("--numeric-ids"))
					Expect(actualArgs).To(ContainElement("--delete"))
					Expect(actualArgs).To(ContainElement("--delete-excluded"))

					// Verify SSH parameters
					sshFlag := false
					for _, arg := range actualArgs {
						if arg == "-e" {
							sshFlag = true
						} else if sshFlag {
							Expect(arg).To(ContainSubstring("ssh -T -x -o StrictHostKeyChecking=no"))
							Expect(arg).To(ContainSubstring("-p 22"))
							Expect(arg).To(ContainSubstring("-i /path/to/ssh/key"))
							sshFlag = false
						}
					}

					// Verify source paths
					Expect(actualArgs).To(ContainElement("testuser@test-host.example.com:/var/www"))
					Expect(actualArgs).To(ContainElement("testuser@test-host.example.com:/etc"))

					// Verify destination path
					expectedIncompletePath := filepath.Join(
						tempDir,
						"test-host.example.com",
						"incomplete",
					)
					Expect(actualArgs).To(ContainElement(expectedIncompletePath))
				})

				It("creates the backup directory structure", func() {
					hostDir := filepath.Join(tempDir, "test-host.example.com")

					// Check that the dated backup directory exists
					backupDir := filepath.Join(hostDir, fixedTime.Format(time.DateOnly))
					Expect(backupDir).To(BeADirectory())

					// Check that current symlink points to the backup
					currentLink := filepath.Join(hostDir, "current")
					Expect(currentLink).To(BeAnExistingFile())

					linkTarget, err := os.Readlink(currentLink)
					Expect(err).To(BeNil())
					Expect(linkTarget).To(Equal(fixedTime.Format(time.DateOnly)))
				})

				It("creates exclude file with correct content", func() {
					excludePath := filepath.Join("/tmp", "test-host.example.com.excludes")
					content, err := os.ReadFile(excludePath)
					Expect(err).To(BeNil())
					Expect(string(content)).To(Equal("*.log\n*.tmp\n\n"))
				})

				It("removes empty directory if it exists", func() {
					hostDir := filepath.Join(tempDir, "test-host.example.com")
					emptyDir := filepath.Join(hostDir, "empty")
					Expect(emptyDir).NotTo(BeADirectory())
				})
			})

			Context("when backup already exists", func() {
				BeforeEach(func() {
					// Create existing backup directory
					hostDir := filepath.Join(tempDir, "test-host.example.com")
					backupDir := filepath.Join(hostDir, fixedTime.Format(time.DateOnly))
					err := os.MkdirAll(backupDir, 0755)
					Expect(err).To(BeNil())
				})

				It("returns no error", func() {
					Expect(err).To(BeNil())
				})

				It("does not call rsync", func() {
					Expect(mockRsyncExecutor.RsyncCallCount()).To(Equal(0))
				})
			})

			Context("when current symlink does not exist", func() {
				It("creates empty directory and current symlink", func() {
					hostDir := filepath.Join(tempDir, "test-host.example.com")
					currentLink := filepath.Join(hostDir, "current")

					// Initially, current should point to empty
					Expect(currentLink).To(BeAnExistingFile())

					linkTarget, err := os.Readlink(currentLink)
					Expect(err).To(BeNil())
					Expect(linkTarget).To(Equal(fixedTime.Format(time.DateOnly)))
				})
			})
		})

		Context("with invalid backup spec", func() {
			Context("missing host", func() {
				BeforeEach(func() {
					backupSpec.Host = ""
				})

				It("returns validation error", func() {
					Expect(err).NotTo(BeNil())
					Expect(err.Error()).To(ContainSubstring("valid backup faild"))
				})

				It("does not call rsync", func() {
					Expect(mockRsyncExecutor.RsyncCallCount()).To(Equal(0))
				})
			})

			Context("missing user", func() {
				BeforeEach(func() {
					backupSpec.User = ""
				})

				It("returns validation error", func() {
					Expect(err).NotTo(BeNil())
					Expect(err.Error()).To(ContainSubstring("valid backup faild"))
				})
			})

			Context("missing directories", func() {
				BeforeEach(func() {
					backupSpec.Dirs = nil
				})

				It("returns validation error", func() {
					Expect(err).NotTo(BeNil())
					Expect(err.Error()).To(ContainSubstring("valid backup faild"))
				})
			})
		})

		Context("when rsync fails", func() {
			BeforeEach(func() {
				mockRsyncExecutor.RsyncReturns(errors.New("rsync execution failed"))
			})

			It("returns error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("rsync failed"))
				Expect(err.Error()).To(ContainSubstring("rsync execution failed"))
			})

			It("calls rsync", func() {
				Expect(mockRsyncExecutor.RsyncCallCount()).To(Equal(1))
			})
		})

		Context("time-based functionality", func() {
			Context("when time advances", func() {
				BeforeEach(func() {
					// Set up initial backup
					mockRsyncExecutor.RsyncReturns(nil)

					// Execute first backup
					err := backupExecutor.Backup(ctx, backupSpec)
					Expect(err).To(BeNil())

					// Advance time by one day
					newTime := fixedTime.Add(24 * time.Hour)
					currentTime.SetNow(newTime)

					// Reset mock for second backup
					mockRsyncExecutor.RsyncReturns(nil)
				})

				It("creates new backup directory with new date", func() {
					// Execute second backup
					err := backupExecutor.Backup(ctx, backupSpec)
					Expect(err).To(BeNil())

					hostDir := filepath.Join(tempDir, "test-host.example.com")

					// Check both backup directories exist
					firstBackupDir := filepath.Join(hostDir, fixedTime.Format(time.DateOnly))
					secondBackupDir := filepath.Join(
						hostDir,
						fixedTime.Add(24*time.Hour).Format(time.DateOnly),
					)

					Expect(firstBackupDir).To(BeADirectory())
					Expect(secondBackupDir).To(BeADirectory())

					// Check current symlink points to latest backup
					currentLink := filepath.Join(hostDir, "current")
					linkTarget, err := os.Readlink(currentLink)
					Expect(err).To(BeNil())
					Expect(
						linkTarget,
					).To(Equal(fixedTime.Add(24 * time.Hour).Format(time.DateOnly)))
				})

				It("uses link-dest pointing to previous backup", func() {
					// Execute second backup
					err := backupExecutor.Backup(ctx, backupSpec)
					Expect(err).To(BeNil())

					// Verify rsync was called twice (first and second backup)
					Expect(mockRsyncExecutor.RsyncCallCount()).To(Equal(2))

					// Check the second call uses link-dest
					_, secondCallArgs := mockRsyncExecutor.RsyncArgsForCall(1)
					linkDestFound := false
					for _, arg := range secondCallArgs {
						if len(arg) > 12 && arg[:12] == "--link-dest=" {
							linkDestFound = true
							expectedCurrentPath := filepath.Join(
								tempDir,
								"test-host.example.com",
								"current",
							)
							Expect(arg).To(Equal("--link-dest=" + expectedCurrentPath))
						}
					}
					Expect(linkDestFound).To(BeTrue())
				})
			})
		})

		Context("SSH key handling", func() {
			Context("with different SSH key path", func() {
				BeforeEach(func() {
					sshPrivateKey = pkg.SSHPrivateKey("/custom/ssh/key/path")

					backupExecutor = pkg.NewBackupExectuor(
						currentTime,
						mockRsyncExecutor,
						backupRootDir,
						sshPrivateKey,
					)

					mockRsyncExecutor.RsyncReturns(nil)
				})

				It("uses correct SSH key in rsync command", func() {
					_, args := mockRsyncExecutor.RsyncArgsForCall(0)

					sshFound := false
					for i, arg := range args {
						if arg == "-e" && i+1 < len(args) {
							Expect(args[i+1]).To(ContainSubstring("-i /custom/ssh/key/path"))
							sshFound = true
						}
					}
					Expect(sshFound).To(BeTrue())
				})
			})
		})

		Context("port handling", func() {
			Context("with custom port", func() {
				BeforeEach(func() {
					backupSpec.Port = v1.BackupPort(2222)
					mockRsyncExecutor.RsyncReturns(nil)
				})

				It("uses correct port in SSH and rsync commands", func() {
					_, args := mockRsyncExecutor.RsyncArgsForCall(0)

					// Check SSH port
					sshFound := false
					portFound := false
					for i, arg := range args {
						if arg == "-e" && i+1 < len(args) {
							Expect(args[i+1]).To(ContainSubstring("-p 2222"))
							sshFound = true
						}
						if arg == "--port=2222" {
							portFound = true
						}
					}
					Expect(sshFound).To(BeTrue())
					Expect(portFound).To(BeTrue())
				})
			})
		})
	})
})
