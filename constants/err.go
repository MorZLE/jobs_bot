package constants

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrNotCategory  = errors.New("not found category")
	ErrNotResume    = errors.New("not found")
	ErrDeleteResume = errors.New("resume delete")
	ErrLastResume   = errors.New("last resume")
	ErrOpenFile     = errors.New("open file")
)
