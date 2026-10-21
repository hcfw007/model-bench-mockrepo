package biz

import "errors"

var (
	ErrContactNotFound = errors.New("contact not found")
	ErrTagAlreadyExists = errors.New("tag already exists")
)
