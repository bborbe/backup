package status_checker

import "github.com/bborbe/backup/dto"

type StatusChecker interface {
	Check() ([]dto.Status, error)
}

type statusChecker struct {
}

func NewStatusChecker(rootdir string) StatusChecker {
	return new(statusChecker)
}

func (s *statusChecker) Check() ([]dto.Status, error) {
	return nil, nil
}
