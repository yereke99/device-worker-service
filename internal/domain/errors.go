package domain

import "errors"

var (
	ErrBusyChannel = errors.New("error channel is busy")
)
