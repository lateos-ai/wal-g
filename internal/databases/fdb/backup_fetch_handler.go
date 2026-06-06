package fdb

import (
	"os/exec"
	"context"

	"github.com/lateos-ai/wal-g/pkg/storages/storage"
	"github.com/lateos-ai/wal-g/internal"
)

func HandleBackupFetch(ctx context.Context,
	folder storage.Folder,
	targetBackupSelector internal.BackupSelector,
	restoreCmd *exec.Cmd) {
	internal.HandleBackupFetch(folder, targetBackupSelector, internal.GetBackupToCommandFetcher(restoreCmd))
}
