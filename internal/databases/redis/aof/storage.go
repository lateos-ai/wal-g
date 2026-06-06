package aof

import (
	"fmt"

	"github.com/lateos-ai/wal-g/pkg/storages/storage"
	"github.com/lateos-ai/wal-g/internal/databases/redis/archive"
)

func SentinelWithExistenceCheck(folder storage.Folder, backupName string) (archive.Backup, error) {
	sentinel, err := archive.SentinelWithExistenceCheck(folder, backupName)
	if err != nil {
		return archive.Backup{}, err
	}
	if sentinel.Version == "" {
		return archive.Backup{}, fmt.Errorf("expecting sentinel file for aof backup with always filled version: %+v", sentinel)
	}
	return sentinel, nil
}
