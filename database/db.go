package database

import (
	"rincon/config"

	"gorm.io/gorm"
)

var DB *gorm.DB

var dbRetries = 0

func InitializeDB() {
	if config.StorageMode == "sql" {
		InitializeSQL()
	} else {
		InitializeLocal()
	}
}
