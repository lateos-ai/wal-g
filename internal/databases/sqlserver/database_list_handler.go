package sqlserver

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/lateos-ai/wal-g/internal"
	"github.com/lateos-ai/wal-g/utility"
	"github.com/wal-g/tracelog"
)

func HandleDatabaseList(backupName string) {
	ctx, cancel := context.WithCancel(context.Background())
	signalHandler := utility.NewSignalHandler(ctx, cancel, []os.Signal{syscall.SIGINT, syscall.SIGTERM})
	defer func() { _ = signalHandler.Close() }()
	storage, err := internal.ConfigureStorage()
	tracelog.ErrorLogger.FatalOnError(err)
	backup, err := internal.GetBackupByName(backupName, utility.BaseBackupPath, storage.RootFolder())
	if err != nil {
		tracelog.ErrorLogger.Fatalf("can't find backup %s: %v", backupName, err)
	}
	sentinel := new(SentinelDto)
	err = backup.FetchSentinel(sentinel)
	tracelog.ErrorLogger.FatalOnError(err)
	for _, name := range sentinel.Databases {
		fmt.Println(name)
	}
}
