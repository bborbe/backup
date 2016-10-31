package checker

import "github.com/bborbe/backup/dto"

type statusCheckerMock struct {
	status []*dto.Status
	err    error
}

func NewStatusCheckerMock(status []*dto.Status, err error) StatusChecker {
	s := new(statusCheckerMock)
	s.status = status
	s.err = err
	return s
}

func (s *statusCheckerMock) Check() ([]*dto.Status, error) {
	return s.status, s.err
}
