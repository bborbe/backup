package status_server_handler

import (
	"errors"
	"net/http"
	"testing"

	. "github.com/bborbe/assert"
	backup_dto "github.com/bborbe/backup/dto"
	backup_status_checker "github.com/bborbe/backup/status_checker"
	server_mock "github.com/bborbe/server/mock"
)

func TestImplementsStatusHandler(t *testing.T) {
	var statusChecker backup_status_checker.StatusChecker
	object := NewStatusHandler(statusChecker)
	var expected *http.Handler
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestStatusCheckerFailure(t *testing.T) {
	logger.Debug("TestStatusCheckerFailure")
	var status []*backup_dto.Status
	err := errors.New("baem!")
	statusChecker := backup_status_checker.NewStatusCheckerMock(status, err)
	handler := NewStatusHandler(statusChecker)
	response := server_mock.NewHttpResponseWriterMock()
	request, err := server_mock.NewHttpRequestMock("http://www.example.com")
	if err != nil {
		t.Error(err)
	}
	handler.ServeHTTP(response, request)
	err = AssertThat(response.Status(), Is(http.StatusInternalServerError))
	if err != nil {
		t.Error(err)
	}
	err = AssertThat(response.String(), Is(`{"status":500,"message":"baem!"}`))
	if err != nil {
		t.Error(err)
	}
}

func TestStatusCheckerNil(t *testing.T) {
	logger.Debug("TestStatusCheckerNil")
	var status []*backup_dto.Status
	var err error
	statusChecker := backup_status_checker.NewStatusCheckerMock(status, err)
	handler := NewStatusHandler(statusChecker)
	response := server_mock.NewHttpResponseWriterMock()
	request, err := server_mock.NewHttpRequestMock("http://www.example.com")
	if err != nil {
		t.Error(err)
	}
	handler.ServeHTTP(response, request)
	err = AssertThat(response.Status(), Is(http.StatusOK))
	if err != nil {
		t.Error(err)
	}
	err = AssertThat(response.String(), Is("[]"))
	if err != nil {
		t.Error(err)
	}
}

func TestStatusCheckerOne(t *testing.T) {
	logger.Debug("TestStatusCheckerNil")
	var status []*backup_dto.Status
	var err error
	status = []*backup_dto.Status{
		createStatus(true, "fire.example.com"),
	}
	statusChecker := backup_status_checker.NewStatusCheckerMock(status, err)
	handler := NewStatusHandler(statusChecker)
	response := server_mock.NewHttpResponseWriterMock()
	request, err := server_mock.NewHttpRequestMock("http://www.example.com")
	if err != nil {
		t.Error(err)
	}
	handler.ServeHTTP(response, request)
	err = AssertThat(response.Status(), Is(http.StatusOK))
	if err != nil {
		t.Error(err)
	}
	err = AssertThat(response.String(), Is(`[{"host":"fire.example.com","status":true,"latestBackup":""}]`))
	if err != nil {
		t.Error(err)
	}
}

func TestStatusCheckerTwo(t *testing.T) {
	logger.Debug("TestStatusCheckerNil")
	var status []*backup_dto.Status
	var err error
	status = []*backup_dto.Status{
		createStatus(true, "fire.example.com"),
		createStatus(false, "burn.example.com"),
	}
	statusChecker := backup_status_checker.NewStatusCheckerMock(status, err)
	handler := NewStatusHandler(statusChecker)
	response := server_mock.NewHttpResponseWriterMock()
	request, err := server_mock.NewHttpRequestMock("http://www.example.com")
	if err != nil {
		t.Error(err)
	}
	handler.ServeHTTP(response, request)
	err = AssertThat(response.Status(), Is(http.StatusOK))
	if err != nil {
		t.Error(err)
	}
	err = AssertThat(response.String(), Is(`[{"host":"fire.example.com","status":true,"latestBackup":""},{"host":"burn.example.com","status":false,"latestBackup":""}]`))
	if err != nil {
		t.Error(err)
	}
}

func createStatus(status bool, host string) *backup_dto.Status {
	s := new(backup_dto.Status)
	s.Status = status
	s.Host = host
	return s
}
