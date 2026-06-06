package mongo

import (
	"os"

	"github.com/wal-g/tracelog"
	"github.com/spf13/cobra"

	"github.com/lateos-ai/wal-g/internal/databases/mongo/common"
	"github.com/lateos-ai/wal-g/internal/databases/mongo"
)

const BackupShowShortDescription = "Prints information about backup"

// backupShowCmd represents the backupList command
var backupShowCmd = &cobra.Command{
	Use:   "backup-show <backup-name>",
	Short: BackupShowShortDescription,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		backupName := args[0]

		backupFolder, err := common.GetBackupFolder()
		tracelog.ErrorLogger.FatalOnError(err)

		err = mongo.HandleBackupShow(backupFolder, backupName, os.Stdout, true)
		tracelog.ErrorLogger.FatalOnError(err)
	},
}

func init() {
	cmd.AddCommand(backupShowCmd)
}
