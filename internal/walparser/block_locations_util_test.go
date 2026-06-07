package walparser_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/lateos-ai/wal-g/internal/walparser"
	"github.com/lateos-ai/wal-g/testtools"
	"github.com/stretchr/testify/assert"
)

func TestExtractBlockLocations(t *testing.T) {
	record, _ := testtools.GetXLogRecordData()
	expectedLocations := []walparser.BlockLocation{record.Blocks[0].Header.BlockLocation}
	actualLocations := walparser.ExtractBlockLocations([]walparser.XLogRecord{record})
	assert.Equal(t, expectedLocations, actualLocations)
}

func TestExtractLocationsFromWalFile(t *testing.T) {
	record, recordData := testtools.GetXLogRecordData()
	fileData := testtools.CreateWalPagesWithRecords(recordData)
	walFile := io.NopCloser(bytes.NewReader(fileData))
	expectedLocations := []walparser.BlockLocation{record.Blocks[0].Header.BlockLocation}
	actualLocations, err := walparser.ExtractLocationsFromWalFile(walparser.NewWalParser(), walFile)
	assert.NoError(t, err)
	assert.Equal(t, expectedLocations, actualLocations)
}

func TestExtractLocationsFromWalFile_MultipleRecords(t *testing.T) {
	record, recordData := testtools.GetXLogRecordData()
	fileData := testtools.CreateWalPagesWithRecords(recordData, recordData)
	walFile := io.NopCloser(bytes.NewReader(fileData))
	expectedLocations := []walparser.BlockLocation{
		record.Blocks[0].Header.BlockLocation, record.Blocks[0].Header.BlockLocation}
	actualLocations, err := walparser.ExtractLocationsFromWalFile(walparser.NewWalParser(), walFile)
	assert.NoError(t, err)
	assert.Equal(t, expectedLocations, actualLocations)
}
