package config

import (
	"os"
)

var Version = "0.1.0"
var Env = os.Getenv("ENV")
var Port = os.Getenv("PORT")

var AuthUser = os.Getenv("AUTH_USER")
var AuthPassword = os.Getenv("AUTH_PASSWORD")

// StorageMode is the mode of storage to use.
// It can be "local", "sql", "redis", 'redis+sql".
var StorageMode = os.Getenv("STORAGE_MODE")

var DatabaseDriver = os.Getenv("DB_DRIVER")
var DatabaseHost = os.Getenv("DB_HOST")
var DatabasePort = os.Getenv("DB_PORT")
var DatabaseName = os.Getenv("DB_NAME")
var DatabaseUser = os.Getenv("DB_USER")
var DatabasePassword = os.Getenv("DB_PASSWORD")
