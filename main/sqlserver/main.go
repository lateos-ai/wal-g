package main

import (
	_ "github.com/microsoft/go-mssqldb"

	"github.com/lateos-ai/wal-g/cmd/sqlserver"
)

func main() {
	sqlserver.Execute()
}
