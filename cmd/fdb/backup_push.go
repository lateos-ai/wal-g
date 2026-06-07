package fdb

import (
	"github.com/lateos-ai/wal-g/internal"
	conf "github.com/lateos-ai/wal-g/internal/config"
	"github.com/lateos-ai/wal-g/internal/databases/fdb"
	"github.com/lateos-ai/wal-g/utility"
	"github.com/spf13/cobra"
	"github.com/wal-g/tracelog"
)

const backupPushShortDescription = "Pushes backup to storage"

// backupPushCmd represents the backupPush command
var backupPushCmd = &cobra.Command{
	Use:   "backup-push",
	Short: backupPushShortDescription,
	Run: func(cmd *cobra.Command, args []string) {
		internal.ConfigureLimiters()

		uploader, err := internal.ConfigureUploader()
		tracelog.ErrorLogger.FatalOnError(err)
		uploader.ChangeDirectory(utility.BaseBackupPath)

		backupCmd, err := internal.GetCommandSetting(conf.NameStreamCreateCmd)
		tracelog.ErrorLogger.FatalOnError(err)
		fdb.HandleBackupPush(cmd.Context(), uploader, backupCmd)
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		conf.RequiredSettings[conf.NameStreamCreateCmd] = true
		err := internal.AssertRequiredSettingsSet()
		tracelog.ErrorLogger.FatalOnError(err)
	},
}

func init() {
	cmd.AddCommand(backupPushCmd)
}
