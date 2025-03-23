package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not fount")
	ErrURLExists   = errors.New("url exists")
)
