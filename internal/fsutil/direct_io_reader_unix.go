//go:build !windows
// +build !windows

package fsutil

import (
	"os"
	"syscall"
)

func isEOFError(err error) bool {
	return errors.Is(err, io.ErrUnexpectedEOF)
}
