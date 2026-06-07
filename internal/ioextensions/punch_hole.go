//go:build !linux

package ioextensions

import (
	"syscall"
	"os"
)

func PunchHole(f *os.File, offset int64, size int64) error {
	// do nothing on non-linux platforms
	return syscall.EOPNOTSUPP
}
