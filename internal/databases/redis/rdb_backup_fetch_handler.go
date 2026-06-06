package redis

import (
	"path/filepath"
	"os/exec"
	"context"

	"github.com/lateos-ai/wal-g/internal/databases/redis/archive"
	"github.com/lateos-ai/wal-g/pkg/storages/storage"
	"github.com/lateos-ai/wal-g/internal"
	conf "github.com/lateos-ai/wal-g/internal/config"
)

func HandleBackupFetch(ctx context.Context, folder storage.Folder, backupName string, restoreCmd *exec.Cmd, skipClean bool) error {
	backup, err := archive.SentinelWithExistenceCheck(folder, backupName)
	if err != nil {
		return err
	}

	if !skipClean {
		dataFolder, _ := conf.GetSetting(conf.RedisDataPath)
		aofFolder, _ := conf.GetSetting(conf.RedisAppendonlyFolder)
		aofPath := filepath.Join(dataFolder, aofFolder)
		aofFolderInfo := archive.CreateAofFolderInfo(aofPath)

		err = aofFolderInfo.CleanPathAndParent()
		if err != nil {
			return err
		}
	}

	return internal.StreamBackupToCommandStdin(restoreCmd, backup.ToInternal(folder))
}
