package status_handler

import (
	"errors"
	. "github.com/bborbe/assert"
	"github.com/bborbe/backup/dto"
	"github.com/bborbe/backup/status_checker"
	"github.com/bborbe/server/mock"
	"net/http"
	"testing"
)

func TestImplementsStatusHandler(t *testing.T) {
	var statusChecker status_checker.StatusChecker
	object := NewStatusHandler(statusChecker)
	var expected *http.Handler
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestStatusCheckerFailure(t *testing.T) {
	logger.Debug("TestStatusCheckerFailure")
	var status []dto.Status
	err := errors.New("baem!")
	statusChecker := status_checker.NewStatusCheckerMock(status, err)
	handler := NewStatusHandler(statusChecker)
	response := mock.NewHttpResponseWriterMock()
	request, err := mock.NewHttpRequestMock("http://www.example.com")
	if err != nil {
		t.Error(err)
	}
	handler.ServeHTTP(response, request)
	err = AssertThat(response.Status(), Is(http.StatusInternalServerError))
	if err != nil {
		t.Error(err)
	}
	err = AssertThat(string(response.Content()), Is(`{"status":500,"message":"baem!"}`))
	if err != nil {
		t.Error(err)
	}
}

func TestStatusCheckerNil(t *testing.T) {
	logger.Debug("TestStatusCheckerNil")
	var status []dto.Status
	var err error
	statusChecker := status_checker.NewStatusCheckerMock(status, err)
	handler := NewStatusHandler(statusChecker)
	response := mock.NewHttpResponseWriterMock()
	request, err := mock.NewHttpRequestMock("http://www.example.com")
	if err != nil {
		t.Error(err)
	}
	handler.ServeHTTP(response, request)
	err = AssertThat(response.Status(), Is(http.StatusOK))
	if err != nil {
		t.Error(err)
	}
	err = AssertThat(string(response.Content()), Is("null"))
	if err != nil {
		t.Error(err)
	}
}

func TestStatusCheckerOne(t *testing.T) {
	logger.Debug("TestStatusCheckerNil")
	var status []dto.Status
	var err error
	status = []dto.Status{
		createStatus(true, "fire.example.com"),
	}
	statusChecker := status_checker.NewStatusCheckerMock(status, err)
	handler := NewStatusHandler(statusChecker)
	response := mock.NewHttpResponseWriterMock()
	request, err := mock.NewHttpRequestMock("http://www.example.com")
	if err != nil {
		t.Error(err)
	}
	handler.ServeHTTP(response, request)
	err = AssertThat(response.Status(), Is(http.StatusOK))
	if err != nil {
		t.Error(err)
	}
	err = AssertThat(string(response.Content()), Is(`[{"host":"fire.example.com","status":true}]`))
	if err != nil {
		t.Error(err)
	}
}

func TestStatusCheckerTwo(t *testing.T) {
	logger.Debug("TestStatusCheckerNil")
	var status []dto.Status
	var err error
	status = []dto.Status{
		createStatus(true, "fire.example.com"),
		createStatus(false, "burn.example.com"),
	}
	statusChecker := status_checker.NewStatusCheckerMock(status, err)
	handler := NewStatusHandler(statusChecker)
	response := mock.NewHttpResponseWriterMock()
	request, err := mock.NewHttpRequestMock("http://www.example.com")
	if err != nil {
		t.Error(err)
	}
	handler.ServeHTTP(response, request)
	err = AssertThat(response.Status(), Is(http.StatusOK))
	if err != nil {
		t.Error(err)
	}
	err = AssertThat(string(response.Content()), Is(`[{"host":"fire.example.com","status":true},{"host":"burn.example.com","status":false}]`))
	if err != nil {
		t.Error(err)
	}
}

func createStatus(status bool, host string) dto.Status {
	s := dto.NewStatus()
	s.SetStatus(status)
	s.SetHost(host)
	return s
}
