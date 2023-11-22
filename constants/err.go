package constants

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrNotCategory  = errors.New("not found")
	ErrNotResume    = errors.New("not found")
	ErrDeleteResume = errors.New("not found")
	ErrLastResume   = errors.New("last resume")
	ErrOpenFile     = errors.New("open file")
)
