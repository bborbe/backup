// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"

	bberrors "github.com/bborbe/errors"
	libhttp "github.com/bborbe/http"
	muxlib "github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
	"github.com/bborbe/backup/mocks"
	"github.com/bborbe/backup/pkg"
	"github.com/bborbe/backup/pkg/handler"
)

var _ = Describe("CleanupHandler", func() {
	var ctx context.Context
	var router *muxlib.Router
	var mockTargetFinder *mocks.TargetFinder
	var mockBackupCleaner *mocks.BackupCleaner
	var target *v1.Target
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		ctx = context.Background()
		mockTargetFinder = &mocks.TargetFinder{}
		mockBackupCleaner = &mocks.BackupCleaner{}

		target = &v1.Target{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-target",
			},
			Spec: v1.BackupSpec{
				Host:     "host1.example.com",
				Port:     v1.BackupPort(22),
				User:     v1.BackupUser("testuser"),
				Dirs:     v1.BackupDirs{"/var/www"},
				Excludes: v1.ParseBackupExcludes([]string{"*.log"}),
			},
		}

		h := handler.NewCleanupHandler(mockTargetFinder, mockBackupCleaner)
		router = muxlib.NewRouter()
		router.Path("/cleanup/{name}").Handler(libhttp.NewJSONErrorHandler(h))
		recorder = httptest.NewRecorder()
	})

	Context("CleanupAlreadyRunning error", func() {
		BeforeEach(func() {
			mockTargetFinder.TargetReturns(target, nil)
		})

		Context("direct sentinel error", func() {
			BeforeEach(func() {
				mockBackupCleaner.CleanReturns(pkg.CleanupAlreadyRunningError)
			})

			It("returns 409 Conflict", func() {
				req := httptest.NewRequest(http.MethodPost, "/cleanup/test-target", nil)
				router.ServeHTTP(recorder, req)
				Expect(recorder.Code).To(Equal(http.StatusConflict))
			})

			It("returns application/json content type", func() {
				req := httptest.NewRequest(http.MethodPost, "/cleanup/test-target", nil)
				router.ServeHTTP(recorder, req)
				Expect(
					recorder.Header().Get("Content-Type"),
				).To(ContainSubstring("application/json"))
			})

			It("returns CLEANUP_ALREADY_RUNNING error code", func() {
				req := httptest.NewRequest(http.MethodPost, "/cleanup/test-target", nil)
				router.ServeHTTP(recorder, req)
				Expect(
					recorder.Body.String(),
				).To(ContainSubstring(`"code":"CLEANUP_ALREADY_RUNNING"`))
			})

			It("returns non-empty message containing host and already running", func() {
				req := httptest.NewRequest(http.MethodPost, "/cleanup/test-target", nil)
				router.ServeHTTP(recorder, req)
				body := recorder.Body.String()
				Expect(body).To(ContainSubstring("host1.example.com"))
				Expect(body).To(ContainSubstring("already running"))
			})

			It("returns empty details", func() {
				req := httptest.NewRequest(http.MethodPost, "/cleanup/test-target", nil)
				router.ServeHTTP(recorder, req)
				body := recorder.Body.String()
				Expect(body).To(ContainSubstring(`"code":"CLEANUP_ALREADY_RUNNING"`))
				Expect(
					body,
				).To(ContainSubstring("cleanup for host1.example.com is already running"))
				Expect(body).NotTo(ContainSubstring(`"details"`))
			})
		})

		Context("wrapped sentinel error", func() {
			BeforeEach(func() {
				mockBackupCleaner.CleanReturns(
					bberrors.Wrapf(ctx, pkg.CleanupAlreadyRunningError, "wrapped"),
				)
			})

			It("returns 409 Conflict", func() {
				req := httptest.NewRequest(http.MethodPost, "/cleanup/test-target", nil)
				router.ServeHTTP(recorder, req)
				Expect(recorder.Code).To(Equal(http.StatusConflict))
			})

			It("returns application/json content type", func() {
				req := httptest.NewRequest(http.MethodPost, "/cleanup/test-target", nil)
				router.ServeHTTP(recorder, req)
				Expect(
					recorder.Header().Get("Content-Type"),
				).To(ContainSubstring("application/json"))
			})

			It("returns CLEANUP_ALREADY_RUNNING error code", func() {
				req := httptest.NewRequest(http.MethodPost, "/cleanup/test-target", nil)
				router.ServeHTTP(recorder, req)
				Expect(
					recorder.Body.String(),
				).To(ContainSubstring(`"code":"CLEANUP_ALREADY_RUNNING"`))
			})

			It("returns non-empty message containing host and already running", func() {
				req := httptest.NewRequest(http.MethodPost, "/cleanup/test-target", nil)
				router.ServeHTTP(recorder, req)
				body := recorder.Body.String()
				Expect(body).To(ContainSubstring("host1.example.com"))
				Expect(body).To(ContainSubstring("already running"))
			})

			It("returns empty details", func() {
				req := httptest.NewRequest(http.MethodPost, "/cleanup/test-target", nil)
				router.ServeHTTP(recorder, req)
				body := recorder.Body.String()
				Expect(body).To(ContainSubstring(`"code":"CLEANUP_ALREADY_RUNNING"`))
				Expect(
					body,
				).To(ContainSubstring("cleanup for host1.example.com is already running"))
				Expect(body).NotTo(ContainSubstring(`"details"`))
			})
		})
	})

	Context("generic error", func() {
		BeforeEach(func() {
			mockTargetFinder.TargetReturns(target, nil)
			mockBackupCleaner.CleanReturns(errors.New("disk on fire"))
		})

		It("returns 500 Internal Server Error", func() {
			req := httptest.NewRequest(http.MethodPost, "/cleanup/test-target", nil)
			router.ServeHTTP(recorder, req)
			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
		})

		It("returns INTERNAL_ERROR code", func() {
			req := httptest.NewRequest(http.MethodPost, "/cleanup/test-target", nil)
			router.ServeHTTP(recorder, req)
			Expect(recorder.Body.String()).To(ContainSubstring(`"code":"INTERNAL_ERROR"`))
		})
	})

	Context("target not found", func() {
		BeforeEach(func() {
			mockTargetFinder.TargetReturns(nil, pkg.TargetNotFoundError)
		})

		It("returns 500 Internal Server Error", func() {
			req := httptest.NewRequest(http.MethodPost, "/cleanup/test-target", nil)
			router.ServeHTTP(recorder, req)
			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
		})

		It("does not call backup cleaner", func() {
			req := httptest.NewRequest(http.MethodPost, "/cleanup/test-target", nil)
			router.ServeHTTP(recorder, req)
			Expect(mockBackupCleaner.CleanCallCount()).To(Equal(0))
		})
	})

	Context("successful cleanup", func() {
		BeforeEach(func() {
			mockTargetFinder.TargetReturns(target, nil)
			mockBackupCleaner.CleanReturns(nil)
		})

		It("returns 200 OK", func() {
			req := httptest.NewRequest(http.MethodPost, "/cleanup/test-target", nil)
			router.ServeHTTP(recorder, req)
			Expect(recorder.Code).To(Equal(http.StatusOK))
		})

		It("returns completion message", func() {
			req := httptest.NewRequest(http.MethodPost, "/cleanup/test-target", nil)
			router.ServeHTTP(recorder, req)
			Expect(recorder.Body.String()).To(ContainSubstring("cleanup test-target completed"))
		})
	})
})
