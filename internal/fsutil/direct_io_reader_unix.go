//go:build !windows
// +build !windows

package fsutil

import (
	"io"
	"errors"
)

func isEOFError(err error) bool {
	return errors.Is(err, io.ErrUnexpectedEOF)
}
