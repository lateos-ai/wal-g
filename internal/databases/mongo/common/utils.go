package common

import (
	"strings"

	"github.com/lateos-ai/wal-g/internal/config"
)

func SystemDBs() *map[string]struct{} {
	res := map[string]struct{}{
		"admin":  {},
		"local":  {},
		"config": {},
	}

	extraSystemDBs, ok := config.GetSetting(config.MongoDBExtraInternalDatabases)
	if ok {
		for _, systemDB := range strings.Split(extraSystemDBs, ",") {
			res[systemDB] = struct{}{}
		}
	}

	return &res
}
