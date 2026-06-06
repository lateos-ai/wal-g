package memory

import (
	"testing"

	"github.com/lateos-ai/wal-g/pkg/storages/storage"
)

func TestMemoryFolder(t *testing.T) {
	storage.RunFolderTest(NewFolder("in_memory/", NewKVS()), t)
}
