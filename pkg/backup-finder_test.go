// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg_test

import (
	"context"
	"os"
	"path/filepath"
	"time"

	libtime "github.com/bborbe/time"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
	"github.com/bborbe/backup/pkg"
)

var _ = Describe("BackupFinder", func() {
	var ctx context.Context
	var err error
	var backupFinder pkg.BackupFinder
	var backupRootDir pkg.Path
	var backupHost v1.BackupHost
	var result []libtime.Date
	var tempDir string

	BeforeEach(func() {
		ctx = context.Background()

		// Create temporary directory for testing
		tempDir, err = os.MkdirTemp("", "backup-finder-test-*")
		Expect(err).To(BeNil())

		backupRootDir = pkg.Path(tempDir)
		backupHost = v1.BackupHost("test-host.example.com")

		// Create backup finder
		backupFinder = pkg.NewBackupFinder(backupRootDir)
	})

	AfterEach(func() {
		// Clean up temporary directory
		_ = os.RemoveAll(tempDir)
	})

	Context("List", func() {
		JustBeforeEach(func() {
			result, err = backupFinder.List(ctx, backupHost)
		})

		Context("when backup directory does not exist", func() {
			It("returns error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("list failed"))
			})

			It("returns no results", func() {
				Expect(result).To(BeNil())
			})
		})

		Context("when backup directory exists but is empty", func() {
			BeforeEach(func() {
				hostDir := filepath.Join(tempDir, string(backupHost))
				err := os.MkdirAll(hostDir, 0755)
				Expect(err).To(BeNil())
			})

			It("returns no error", func() {
				Expect(err).To(BeNil())
			})

			It("returns empty result", func() {
				Expect(result).To(HaveLen(0))
			})
		})

		Context("when backup directory contains valid date directories", func() {
			BeforeEach(func() {
				hostDir := filepath.Join(tempDir, string(backupHost))
				err := os.MkdirAll(hostDir, 0755)
				Expect(err).To(BeNil())

				// Create directories with valid date names
				validDates := []string{
					"2023-12-25",
					"2023-12-24",
					"2023-12-23",
				}

				for _, date := range validDates {
					dateDir := filepath.Join(hostDir, date)
					err := os.MkdirAll(dateDir, 0755)
					Expect(err).To(BeNil())
				}
			})

			It("returns no error", func() {
				Expect(err).To(BeNil())
			})

			It("returns correct number of dates", func() {
				Expect(result).To(HaveLen(3))
			})

			It("returns correct dates", func() {
				dateStrings := make([]string, len(result))
				for i, date := range result {
					dateStrings[i] = date.String()
				}

				Expect(dateStrings).To(ContainElement("2023-12-25"))
				Expect(dateStrings).To(ContainElement("2023-12-24"))
				Expect(dateStrings).To(ContainElement("2023-12-23"))
			})

			It("returns dates as libtime.Date objects", func() {
				for _, date := range result {
					// Verify we can convert back to time.Time
					parsedTime := date.Time()
					Expect(parsedTime).NotTo(BeZero())

					// Verify the date is in the expected range
					expectedTime := time.Date(2023, 12, 23, 0, 0, 0, 0, time.UTC)
					Expect(parsedTime).To(BeTemporally(">=", expectedTime))

					expectedTime = time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
					Expect(parsedTime).To(BeTemporally("<=", expectedTime))
				}
			})
		})

		Context("when backup directory contains invalid entries", func() {
			BeforeEach(func() {
				hostDir := filepath.Join(tempDir, string(backupHost))
				err := os.MkdirAll(hostDir, 0755)
				Expect(err).To(BeNil())

				// Create mix of valid and invalid entries
				validDates := []string{"2023-12-25", "2023-12-24"}
				invalidEntries := []string{"current", "incomplete", "empty", "invalid-date", "backup.txt"}

				for _, date := range validDates {
					dateDir := filepath.Join(hostDir, date)
					err := os.MkdirAll(dateDir, 0755)
					Expect(err).To(BeNil())
				}

				for _, entry := range invalidEntries {
					entryPath := filepath.Join(hostDir, entry)
					if entry == "backup.txt" {
						// Create a file instead of directory
						err := os.WriteFile(entryPath, []byte("test"), 0644)
						Expect(err).To(BeNil())
					} else {
						err := os.MkdirAll(entryPath, 0755)
						Expect(err).To(BeNil())
					}
				}
			})

			It("returns no error", func() {
				Expect(err).To(BeNil())
			})

			It("returns only valid dates", func() {
				Expect(result).To(HaveLen(2))

				dateStrings := make([]string, len(result))
				for i, date := range result {
					dateStrings[i] = date.String()
				}

				Expect(dateStrings).To(ContainElement("2023-12-25"))
				Expect(dateStrings).To(ContainElement("2023-12-24"))
			})

			It("ignores invalid entries", func() {
				// Verify that invalid entries are not included
				dateStrings := make([]string, len(result))
				for i, date := range result {
					dateStrings[i] = date.String()
				}

				Expect(dateStrings).NotTo(ContainElement("current"))
				Expect(dateStrings).NotTo(ContainElement("incomplete"))
				Expect(dateStrings).NotTo(ContainElement("empty"))
				Expect(dateStrings).NotTo(ContainElement("invalid-date"))
				Expect(dateStrings).NotTo(ContainElement("backup.txt"))
			})
		})

		Context("with different date formats", func() {
			BeforeEach(func() {
				hostDir := filepath.Join(tempDir, string(backupHost))
				err := os.MkdirAll(hostDir, 0755)
				Expect(err).To(BeNil())

				// Test various date formats to ensure only valid ISO dates are accepted
				dateEntries := []string{
					"2023-12-25", // Valid ISO date
					"2023-01-01", // Valid ISO date
					"2023-2-5",   // Invalid (wrong format)
					"23-12-25",   // Invalid (wrong format)
					"2023/12/25", // Invalid (wrong separator)
					"2023.12.25", // Invalid (wrong separator)
					"20231225",   // Invalid (no separators)
				}

				for _, entry := range dateEntries {
					entryPath := filepath.Join(hostDir, entry)
					err := os.MkdirAll(entryPath, 0755)
					Expect(err).To(BeNil())
				}
			})

			It("returns no error", func() {
				Expect(err).To(BeNil())
			})

			It("returns only valid ISO date format entries", func() {
				Expect(result).To(HaveLen(2))

				dateStrings := make([]string, len(result))
				for i, date := range result {
					dateStrings[i] = date.String()
				}

				Expect(dateStrings).To(ContainElement("2023-12-25"))
				Expect(dateStrings).To(ContainElement("2023-01-01"))
			})
		})

		Context("with different host names", func() {
			Context("with simple host name", func() {
				BeforeEach(func() {
					backupHost = v1.BackupHost("simple-host")

					hostDir := filepath.Join(tempDir, string(backupHost))
					err := os.MkdirAll(hostDir, 0755)
					Expect(err).To(BeNil())

					dateDir := filepath.Join(hostDir, "2023-12-25")
					err = os.MkdirAll(dateDir, 0755)
					Expect(err).To(BeNil())
				})

				It("finds backups correctly", func() {
					Expect(err).To(BeNil())
					Expect(result).To(HaveLen(1))
					Expect(result[0].String()).To(Equal("2023-12-25"))
				})
			})

			Context("with FQDN host name", func() {
				BeforeEach(func() {
					backupHost = v1.BackupHost("server.domain.com")

					hostDir := filepath.Join(tempDir, string(backupHost))
					err := os.MkdirAll(hostDir, 0755)
					Expect(err).To(BeNil())

					dateDir := filepath.Join(hostDir, "2023-12-25")
					err = os.MkdirAll(dateDir, 0755)
					Expect(err).To(BeNil())
				})

				It("finds backups correctly", func() {
					Expect(err).To(BeNil())
					Expect(result).To(HaveLen(1))
					Expect(result[0].String()).To(Equal("2023-12-25"))
				})
			})

			Context("with host name containing special characters", func() {
				BeforeEach(func() {
					backupHost = v1.BackupHost("server-01.example-domain.com")

					hostDir := filepath.Join(tempDir, string(backupHost))
					err := os.MkdirAll(hostDir, 0755)
					Expect(err).To(BeNil())

					dateDir := filepath.Join(hostDir, "2023-12-25")
					err = os.MkdirAll(dateDir, 0755)
					Expect(err).To(BeNil())
				})

				It("finds backups correctly", func() {
					Expect(err).To(BeNil())
					Expect(result).To(HaveLen(1))
					Expect(result[0].String()).To(Equal("2023-12-25"))
				})
			})
		})

		Context("with large number of backup directories", func() {
			BeforeEach(func() {
				hostDir := filepath.Join(tempDir, string(backupHost))
				err := os.MkdirAll(hostDir, 0755)
				Expect(err).To(BeNil())

				// Create many backup directories to test performance
				baseDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
				for i := 0; i < 100; i++ {
					date := baseDate.AddDate(0, 0, i)
					dateDir := filepath.Join(hostDir, date.Format(time.DateOnly))
					err := os.MkdirAll(dateDir, 0755)
					Expect(err).To(BeNil())
				}
			})

			It("returns no error", func() {
				Expect(err).To(BeNil())
			})

			It("returns all valid backup dates", func() {
				Expect(result).To(HaveLen(100))
			})

			It("returns dates in valid format", func() {
				for _, date := range result {
					// Verify each date is valid
					parsedTime := date.Time()
					Expect(parsedTime).NotTo(BeZero())

					// Verify date is in expected range
					expectedStart := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
					expectedEnd := time.Date(2023, 4, 10, 0, 0, 0, 0, time.UTC)
					Expect(parsedTime).To(BeTemporally(">=", expectedStart))
					Expect(parsedTime).To(BeTemporally("<=", expectedEnd))
				}
			})
		})

		Context("edge cases", func() {
			Context("with empty host name", func() {
				BeforeEach(func() {
					backupHost = v1.BackupHost("")

					// Create directory for empty host name
					hostDir := filepath.Join(tempDir, string(backupHost))
					err := os.MkdirAll(hostDir, 0755)
					Expect(err).To(BeNil())
				})

				It("handles empty host gracefully", func() {
					Expect(err).To(BeNil())
					Expect(result).To(HaveLen(0))
				})
			})

			Context("when backup root directory does not exist", func() {
				BeforeEach(func() {
					backupRootDir = pkg.Path("/nonexistent/directory")
					backupFinder = pkg.NewBackupFinder(backupRootDir)
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
					Expect(err.Error()).To(ContainSubstring("list failed"))
				})
			})

			Context("with read permission denied", func() {
				BeforeEach(func() {
					hostDir := filepath.Join(tempDir, string(backupHost))
					err := os.MkdirAll(hostDir, 0755)
					Expect(err).To(BeNil())

					// Remove read permissions (if running as non-root)
					err = os.Chmod(hostDir, 0000)
					if err != nil {
						Skip("Cannot change permissions - likely running as root")
					}
				})

				AfterEach(func() {
					// Restore permissions for cleanup
					hostDir := filepath.Join(tempDir, string(backupHost))
					_ = os.Chmod(hostDir, 0755)
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
					Expect(err.Error()).To(ContainSubstring("list failed"))
				})
			})
		})

		Context("date boundary testing", func() {
			BeforeEach(func() {
				hostDir := filepath.Join(tempDir, string(backupHost))
				err := os.MkdirAll(hostDir, 0755)
				Expect(err).To(BeNil())

				// Test various date boundaries
				boundaryDates := []string{
					"2023-01-01", // Year start
					"2023-12-31", // Year end
					"2023-02-28", // Non-leap year February end
					"2024-02-29", // Leap year February end
					"2023-04-30", // Month with 30 days
					"2023-05-31", // Month with 31 days
				}

				for _, date := range boundaryDates {
					dateDir := filepath.Join(hostDir, date)
					err := os.MkdirAll(dateDir, 0755)
					Expect(err).To(BeNil())
				}
			})

			It("handles date boundaries correctly", func() {
				Expect(err).To(BeNil())
				Expect(result).To(HaveLen(6))

				// Verify all boundary dates are included
				dateStrings := make([]string, len(result))
				for i, date := range result {
					dateStrings[i] = date.String()
				}

				Expect(dateStrings).To(ContainElement("2023-01-01"))
				Expect(dateStrings).To(ContainElement("2023-12-31"))
				Expect(dateStrings).To(ContainElement("2023-02-28"))
				Expect(dateStrings).To(ContainElement("2024-02-29"))
				Expect(dateStrings).To(ContainElement("2023-04-30"))
				Expect(dateStrings).To(ContainElement("2023-05-31"))
			})
		})
	})
})
