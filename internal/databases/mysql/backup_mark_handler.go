package mysql

import (
	"github.com/lateos-ai/wal-g/internal"
	"github.com/wal-g/tracelog"
)

// MarkBackup marks a backup as permanent or impermanent
func MarkBackup(uploader internal.Uploader, backupName string, toPermanent bool) {
	tracelog.InfoLogger.Printf("Retrieving previous related backups to be marked: toPermanent=%t", toPermanent)
	internal.HandleBackupMark(uploader, backupName, toPermanent, NewGenericMetaInteractor())
}
