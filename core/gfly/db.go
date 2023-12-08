package gfly

import (
	"app/core/db"
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

	err := db.Connect()

	if err != nil {
		log.Fatal(err)
	}

	log.Trace("Connect Database")
}
