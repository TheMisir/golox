package main

import "errors"

var catch = errors.New("catch")

func tryCatch(cb func()) (err error) {
	defer func() {
		// recover from panic if one occurred. Set err to nil otherwise.
		if recover() != nil {
			err = catch
		}
	}()

	cb()
	return
}
