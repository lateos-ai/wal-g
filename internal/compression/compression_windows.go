//go:build windows
// +build windows

package compression

import (
	"github.com/lateos-ai/wal-g/internal/compression/lz4"
	"github.com/lateos-ai/wal-g/internal/compression/lzma"
)

var CompressingAlgorithms = []string{lz4.AlgorithmName, lzma.AlgorithmName}

var Compressors = map[string]Compressor{
	lz4.AlgorithmName:  lz4.Compressor{},
	lzma.AlgorithmName: lzma.Compressor{},
}

var Decompressors = []Decompressor{
	lz4.Decompressor{},
	lzma.Decompressor{},
}
