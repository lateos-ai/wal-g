package etcd

import (
	"github.com/wal-g/tracelog"
	"github.com/spf13/cobra"

	"github.com/lateos-ai/wal-g/internal/databases/etcd"
	"github.com/lateos-ai/wal-g/internal"
)

const fetchSinceFlagShortDescr = "backup name starting from which you want to fetch wals"

var fetchBackupName string

var WalFetchCmd = &cobra.Command{
	Use:   "wal-fetch dest-dir",
	Short: "Fetch wal from storage and save it to the specified dir",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := internal.ConfigureStorage()
		tracelog.ErrorLogger.FatalOnError(err)
		folderReader := internal.NewFolderReader(storage.RootFolder())
		etcd.HandleWalFetch(storage.RootFolder(), fetchBackupName, args[0], folderReader)
	},
}

func init() {
	WalFetchCmd.PersistentFlags().StringVar(&fetchBackupName, "since", "LATEST", fetchSinceFlagShortDescr)
	cmd.AddCommand(WalFetchCmd)
}
