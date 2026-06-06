package etcd

import (
	"github.com/wal-g/tracelog"
	"github.com/spf13/cobra"

	"github.com/lateos-ai/wal-g/utility"
	"github.com/lateos-ai/wal-g/internal"
)

const backupListShortDescription = "Prints available backups"

// backupListCmd represents the backupList command
var backupListCmd = &cobra.Command{
	Use:   "backup-list",
	Short: backupListShortDescription,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := internal.ConfigureStorage()
		tracelog.ErrorLogger.FatalOnError(err)
		internal.HandleDefaultBackupList(storage.RootFolder().GetSubFolder(utility.BaseBackupPath), false, false)
	},
}

func init() {
	cmd.AddCommand(backupListCmd)
}
