package linkedlist

import "errors"

var (
	ErrEmptyList = errors.New("list is empty")
	ErrNotFound  = errors.New("value not found in list")
)
