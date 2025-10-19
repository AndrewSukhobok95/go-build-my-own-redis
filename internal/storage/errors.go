package storage

import (
	"errors"
)

var (
	ErrNotInteger = errors.New("value is not an integer or out of range")
	ErrOverflow   = errors.New("increment or decrement would overflow")
	ErrWrongType  = errors.New("wrong type")
)
