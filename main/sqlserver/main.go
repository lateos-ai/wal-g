package main

import (
	"github.com/lateos-ai/wal-g/cmd/sqlserver"
	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	sqlserver.Execute()
}
