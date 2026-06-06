package mongo

import (
	"io"
	"encoding/json"

	"github.com/lateos-ai/wal-g/pkg/storages/storage"
	"github.com/lateos-ai/wal-g/internal/databases/mongo/common"
)

// HandleBackupShow prints sentinel contents.
func HandleBackupShow(backupFolder storage.Folder, backupName string, output io.Writer, pretty bool) (err error) {
	sentinel, err := common.DownloadSentinel(backupFolder, backupName)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(output)
	if pretty {
		encoder.SetIndent("", "    ")
	}
	return encoder.Encode(sentinel)
}
