package mysql

import (
	"github.com/wal-g/tracelog"
	gomysql "github.com/go-mysql-org/go-mysql/mysql"

	"github.com/lateos-ai/wal-g/utility"
	"github.com/lateos-ai/wal-g/pkg/storages/storage"
)

func HandleBinlogFind(folder storage.Folder, gtid string) {
	db, err := getMySQLConnection()
	tracelog.ErrorLogger.FatalOnError(err)
	defer utility.LoggedClose(db, "")
	flavor, err := getMySQLFlavor(db)
	tracelog.ErrorLogger.FatalOnError(err)
	var gtidSet gomysql.GTIDSet
	if gtid == "" {
		gtidSet, err = getMySQLGTIDExecuted(db, flavor)
		tracelog.ErrorLogger.FatalOnError(err)
	} else {
		gtidSet, err = gomysql.ParseGTIDSet(flavor, gtid)
		tracelog.ErrorLogger.FatalOnError(err)
	}
	name, err := getLastUploadedBinlogBeforeGTID(folder, gtidSet, flavor)
	tracelog.ErrorLogger.FatalOnError(err)
	tracelog.InfoLogger.Println(name)
}
