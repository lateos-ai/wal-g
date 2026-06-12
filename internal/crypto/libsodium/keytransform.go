//go:build libsodium
// +build libsodium

package libsodium

// NOTE: #cgo uses ${SRCDIR} (absolute path to this package's directory) to
// reference the tmp/libsodium/ tree in the repository root. The resulting
// absolute -I/-L paths are trusted by Go 1.21+'s -mod=vendor security checks
// (which only sanitize relative paths). The Makefile also exports
// CGO_ENABLED=1 CGO_CFLAGS/CGO_LDFLAGS for manual builds outside of make.
// Before building, run link_libsodium.sh (or equivalent) to populate
// tmp/libsodium/ with headers and the static library.

// #cgo CFLAGS: -I${SRCDIR}/../../../tmp/libsodium/include -I${SRCDIR}/../../../tmp/libsodium/include/sodium
// #cgo LDFLAGS: -L${SRCDIR}/../../../tmp/libsodium/lib -lsodium
// #include <sodium.h>

import "C"

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

const KeyTransformBase64 = "base64"

const KeyTransformHex = "hex"

const KeyTransformNone = "none"

type keyTransformRegEntry struct {
	typ string

	fun func(userInput string) ([]byte, error)
}

var keyTransformReg = []keyTransformRegEntry{
	{typ: KeyTransformBase64, fun: keyTransformBase64},

	{typ: KeyTransformHex, fun: keyTransformHex},

	{typ: KeyTransformNone, fun: keyTransformNone},
}

func keyTransform(userInput string, transformType string, expectedLen int) ([]byte, error) {
	for _, entry := range keyTransformReg {
		if entry.typ == transformType {
			decoded, err := entry.fun(userInput)

			if err != nil {
				return nil, err
			}

			if len(decoded) != expectedLen {
				return nil, fmt.Errorf("key must be exactly %d bytes (got %d bytes)", expectedLen, len(decoded))
			}

			return decoded, nil
		}
	}

	// unknown transform

	var builder strings.Builder

	for idx, entry := range keyTransformReg {
		if idx > 0 {
			if idx+1 == len(keyTransformReg) {
				builder.WriteString(" or ")
			} else {
				builder.WriteString(", ")
			}
		}

		builder.WriteString(entry.typ)
	}

	return nil, fmt.Errorf("unknown key transform '%s' (must be %s)", transformType, builder.String())
}

func keyTransformBase64(userInput string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(userInput)

	if err != nil {
		return nil, fmt.Errorf("while base64 decoding key: %v", err)
	}

	return decoded, nil
}

func keyTransformHex(userInput string) ([]byte, error) {
	decoded, err := hex.DecodeString(userInput)

	if err != nil {
		return nil, fmt.Errorf("while hex decoding key: %v", err)
	}

	return decoded, nil
}

// Mimics the behavior of older versions of wal-g.

func keyTransformNone(userInput string) ([]byte, error) {
	if len(userInput) < minimalKeyLength {
		return nil, newErrShortKey(len(userInput))
	}

	if len(userInput) > libsodiumKeybytes {
		return []byte(userInput[:libsodiumKeybytes]), nil
	}

	if len(userInput) < libsodiumKeybytes {
		buf := make([]byte, libsodiumKeybytes)

		copy(buf[:libsodiumKeybytes], userInput)

		return buf, nil
	}

	return []byte(userInput), nil
}
