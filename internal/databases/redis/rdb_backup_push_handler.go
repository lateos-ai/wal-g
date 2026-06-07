package redis

import (
	"os/exec"

	"github.com/lateos-ai/wal-g/internal"
	"github.com/lateos-ai/wal-g/internal/databases/redis/rdb"
	"github.com/lateos-ai/wal-g/utility"
	"github.com/wal-g/tracelog"
)

type RDBBackupPushArgs struct {
	BackupCmd       *exec.Cmd
	Sharded         bool
	Uploader        internal.Uploader
	MetaConstructor internal.MetaConstructor
}

func HandleRDBBackupPush(args RDBBackupPushArgs) error {
	stdout, err := utility.StartCommandWithStdoutPipe(args.BackupCmd)
	tracelog.ErrorLogger.FatalfOnError("failed to start backup create command: %v", err)

	redisUploader := rdb.NewRedisStorageUploader(args.Uploader)
	uploadArgs := rdb.UploadBackupArgs{
		Cmd:             args.BackupCmd,
		MetaConstructor: args.MetaConstructor,
		Sharded:         args.Sharded,
		Stream:          stdout,
	}

	return redisUploader.UploadBackup(uploadArgs)
}
