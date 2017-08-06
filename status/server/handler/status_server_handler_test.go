package handler

import (
	"errors"
	"net/http"
	"testing"

	"os"

	. "github.com/bborbe/assert"
	backup_dto "github.com/bborbe/backup/dto"
	backup_status_checker "github.com/bborbe/backup/status/server/checker"
	server_mock "github.com/bborbe/http/mock"
	"github.com/golang/glog"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	glog.Flush()
	os.Exit(exit)
}

func TestImplementsStatusHandler(t *testing.T) {
	var statusChecker backup_status_checker.StatusChecker
	object := New(statusChecker)
	var expected *http.Handler
	if err := AssertThat(object, Implements(expected)); err != nil {
		t.Fatal(err)
	}
}

func TestStatusCheckerFailure(t *testing.T) {
	glog.V(2).Info("TestStatusCheckerFailure")
	var status []backup_dto.Status
	err := errors.New("baem!")
	statusChecker := backup_status_checker.NewStatusCheckerMock(status, err)
	handler := New(statusChecker)
	response := server_mock.NewHttpResponseWriterMock()
	request, err := server_mock.NewHttpRequestMock("http://www.example.com")
	if err != nil {
		t.Error(err)
	}
	handler.ServeHTTP(response, request)
	if err := AssertThat(response.Status(), Is(http.StatusInternalServerError)); err != nil {
		t.Error(err)
	}
	if err := AssertThat(response.String(), Is("{\"status\":500,\"message\":\"baem!\"}\n")); err != nil {
		t.Error(err)
	}
}

func TestStatusCheckerNil(t *testing.T) {
	glog.V(2).Info("TestStatusCheckerNil")
	var status []backup_dto.Status
	var err error
	statusChecker := backup_status_checker.NewStatusCheckerMock(status, err)
	handler := New(statusChecker)
	response := server_mock.NewHttpResponseWriterMock()
	request, err := server_mock.NewHttpRequestMock("http://www.example.com")
	if err != nil {
		t.Error(err)
	}
	handler.ServeHTTP(response, request)
	if err := AssertThat(response.Status(), Is(http.StatusOK)); err != nil {
		t.Error(err)
	}
	if err := AssertThat(response.String(), Is("[]")); err != nil {
		t.Error(err)
	}
}

func TestStatusCheckerOne(t *testing.T) {
	glog.V(2).Info("TestStatusCheckerNil")
	var status []backup_dto.Status
	var err error
	status = []backup_dto.Status{
		createStatus(true, "fire.example.com"),
	}
	statusChecker := backup_status_checker.NewStatusCheckerMock(status, err)
	handler := New(statusChecker)
	response := server_mock.NewHttpResponseWriterMock()
	request, err := server_mock.NewHttpRequestMock("http://www.example.com")
	if err != nil {
		t.Error(err)
	}
	handler.ServeHTTP(response, request)
	if err := AssertThat(response.Status(), Is(http.StatusOK)); err != nil {
		t.Error(err)
	}
	if err := AssertThat(response.String(), Is(`[{"host":"fire.example.com","status":true,"latestBackup":""}]`)); err != nil {
		t.Error(err)
	}
}

func TestStatusCheckerTwo(t *testing.T) {
	glog.V(2).Info("TestStatusCheckerNil")
	var status []backup_dto.Status
	var err error
	status = []backup_dto.Status{
		createStatus(true, "fire.example.com"),
		createStatus(false, "burn.example.com"),
	}
	statusChecker := backup_status_checker.NewStatusCheckerMock(status, err)
	handler := New(statusChecker)
	response := server_mock.NewHttpResponseWriterMock()
	request, err := server_mock.NewHttpRequestMock("http://www.example.com")
	if err != nil {
		t.Error(err)
	}
	handler.ServeHTTP(response, request)
	if err := AssertThat(response.Status(), Is(http.StatusOK)); err != nil {
		t.Error(err)
	}
	if err := AssertThat(response.String(), Is(`[{"host":"burn.example.com","status":false,"latestBackup":""},{"host":"fire.example.com","status":true,"latestBackup":""}]`)); err != nil {
		t.Error(err)
	}
}

func createStatus(status bool, host string) backup_dto.Status {
	return backup_dto.Status{
		Status: status,
		Host:   host,
	}
}

func TestStatusCheckerSort(t *testing.T) {
	glog.V(2).Info("TestStatusCheckerNil")
	var status []backup_dto.Status
	var err error
	status = []backup_dto.Status{
		backup_dto.Status{Status: true, Host: "example.com", LatestBackup: "2017-08-06T15:04:05"},
		backup_dto.Status{Status: true, Host: "example.com", LatestBackup: "2017-08-06T15:04:06"},
	}
	statusChecker := backup_status_checker.NewStatusCheckerMock(status, err)
	handler := New(statusChecker)
	response := server_mock.NewHttpResponseWriterMock()
	request, err := server_mock.NewHttpRequestMock("http://www.example.com")
	if err != nil {
		t.Error(err)
	}
	handler.ServeHTTP(response, request)
	if err := AssertThat(response.Status(), Is(http.StatusOK)); err != nil {
		t.Error(err)
	}
	if err := AssertThat(response.String(), Is(`[{"host":"example.com","status":true,"latestBackup":"2017-08-06T15:04:05"},{"host":"example.com","status":true,"latestBackup":"2017-08-06T15:04:06"}]`)); err != nil {
		t.Error(err)
	}
}

func TestStatusCheckerSortInverse(t *testing.T) {
	glog.V(2).Info("TestStatusCheckerNil")
	var status []backup_dto.Status
	var err error
	status = []backup_dto.Status{
		backup_dto.Status{Status: true, Host: "example.com", LatestBackup: "2017-08-06T15:04:06"},
		backup_dto.Status{Status: true, Host: "example.com", LatestBackup: "2017-08-06T15:04:05"},
	}
	statusChecker := backup_status_checker.NewStatusCheckerMock(status, err)
	handler := New(statusChecker)
	response := server_mock.NewHttpResponseWriterMock()
	request, err := server_mock.NewHttpRequestMock("http://www.example.com")
	if err != nil {
		t.Error(err)
	}
	handler.ServeHTTP(response, request)
	if err := AssertThat(response.Status(), Is(http.StatusOK)); err != nil {
		t.Error(err)
	}
	if err := AssertThat(response.String(), Is(`[{"host":"example.com","status":true,"latestBackup":"2017-08-06T15:04:05"},{"host":"example.com","status":true,"latestBackup":"2017-08-06T15:04:06"}]`)); err != nil {
		t.Error(err)
	}
}

func TestStatusCheckerSortMissing(t *testing.T) {
	glog.V(2).Info("TestStatusCheckerNil")
	var status []backup_dto.Status
	var err error
	status = []backup_dto.Status{
		{Status: true, Host: "example.com", LatestBackup: "2017-08-06T15:04:05"},
		{Status: true, Host: "example.com", LatestBackup: ""},
		{Status: true, Host: "example.com", LatestBackup: "2017-08-06T15:04:06"},
	}
	statusChecker := backup_status_checker.NewStatusCheckerMock(status, err)
	handler := New(statusChecker)
	response := server_mock.NewHttpResponseWriterMock()
	request, err := server_mock.NewHttpRequestMock("http://www.example.com")
	if err != nil {
		t.Error(err)
	}
	handler.ServeHTTP(response, request)
	if err := AssertThat(response.Status(), Is(http.StatusOK)); err != nil {
		t.Error(err)
	}
	if err := AssertThat(response.String(), Is(`[{"host":"example.com","status":true,"latestBackup":""},{"host":"example.com","status":true,"latestBackup":"2017-08-06T15:04:05"},{"host":"example.com","status":true,"latestBackup":"2017-08-06T15:04:06"}]`)); err != nil {
		t.Error(err)
	}
}

func TestStatusCheckerSortMissingInverse(t *testing.T) {
	glog.V(2).Info("TestStatusCheckerNil")
	var status []backup_dto.Status
	var err error
	status = []backup_dto.Status{
		{Status: true, Host: "example.com", LatestBackup: "2017-08-06T15:04:06"},
		{Status: true, Host: "example.com", LatestBackup: ""},
		{Status: true, Host: "example.com", LatestBackup: "2017-08-06T15:04:05"},
	}
	statusChecker := backup_status_checker.NewStatusCheckerMock(status, err)
	handler := New(statusChecker)
	response := server_mock.NewHttpResponseWriterMock()
	request, err := server_mock.NewHttpRequestMock("http://www.example.com")
	if err != nil {
		t.Error(err)
	}
	handler.ServeHTTP(response, request)
	if err := AssertThat(response.Status(), Is(http.StatusOK)); err != nil {
		t.Error(err)
	}
	if err := AssertThat(response.String(), Is(`[{"host":"example.com","status":true,"latestBackup":""},{"host":"example.com","status":true,"latestBackup":"2017-08-06T15:04:05"},{"host":"example.com","status":true,"latestBackup":"2017-08-06T15:04:06"}]`)); err != nil {
		t.Error(err)
	}
}
