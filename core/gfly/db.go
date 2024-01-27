package gfly

import (
	"app/core/log"
	"app/core/utils"
)

// ===========================================================================================================
// 										DB
// ===========================================================================================================

func setupDB() {
	if utils.Getenv("DB_DRIVER", "") == "" {
		log.Trace("Disable Database")

		return
	}

	log.Trace("Connect Database")
}
