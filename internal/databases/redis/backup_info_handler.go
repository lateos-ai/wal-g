package redis

import (
	"reflect"
	"io"
	"fmt"

	"github.com/wal-g/tracelog"

	"github.com/lateos-ai/wal-g/pkg/storages/storage"
	"github.com/lateos-ai/wal-g/internal/printlist"
	"github.com/lateos-ai/wal-g/internal/databases/redis/archive"
)

func HandleBackupInfo(folder storage.Folder, backupName string, output io.Writer, tag string) {
	backupDetails, err := archive.SentinelWithExistenceCheck(folder, backupName)
	tracelog.ErrorLogger.FatalOnError(err)

	if tag != "" {
		v, err := getField(folder, &backupDetails, tag)
		tracelog.ErrorLogger.FatalOnError(err)
		_, err = fmt.Fprintln(output, v)
		tracelog.ErrorLogger.FatalOnError(err)
		return
	}

	pretty := false
	json := true
	err = printlist.List([]printlist.Entity{backupDetails}, output, pretty, json)
	tracelog.ErrorLogger.FatalfOnError("Print backup info: %v", err)
}

func getField(folder storage.Folder, v *archive.Backup, field string) (string, error) {
	if field == "Slots" {
		return archive.FetchSlotsDataFromStorage(folder, v)
	}

	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	if f.IsValid() {
		return f.String(), nil
	}
	return "", fmt.Errorf("no %s field in struct %v", field, &v)
}
