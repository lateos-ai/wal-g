package etcd

import (
	"github.com/lateos-ai/wal-g/internal"
	conf "github.com/lateos-ai/wal-g/internal/config"
	"github.com/lateos-ai/wal-g/internal/databases/etcd"
	"github.com/spf13/cobra"
	"github.com/wal-g/tracelog"
)

const backupFetchShortDescription = "Fetches desired backup from storage"

// backupFetchCmd represents the streamFetch command
var backupFetchCmd = &cobra.Command{
	Use:   "backup-fetch backup-name",
	Short: backupFetchShortDescription,
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		conf.RequiredSettings[conf.NameStreamRestoreCmd] = true
		err := internal.AssertRequiredSettingsSet()
		tracelog.ErrorLogger.FatalOnError(err)
	},
	Run: func(cmd *cobra.Command, args []string) {
		internal.ConfigureLimiters()
		ctx := cmd.Context()

		storage, err := internal.ConfigureStorage()
		tracelog.ErrorLogger.FatalOnError(err)

		restoreCmd, err := internal.GetCommandSettingContext(ctx, conf.NameStreamRestoreCmd)
		tracelog.ErrorLogger.FatalOnError(err)
		targetBackupSelector, err := internal.NewBackupNameSelector(args[0], true)
		tracelog.ErrorLogger.FatalOnError(err)
		etcd.HandleBackupFetch(ctx, storage.RootFolder(), targetBackupSelector, restoreCmd)
	},
}

func init() {
	cmd.AddCommand(backupFetchCmd)
}
