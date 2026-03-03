package storage

import "errors"

var (
	errURLNotFound = errors.New("url not found")
	errURLExists   = errors.New("url aready exists")
)
