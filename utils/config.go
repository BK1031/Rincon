package utils

import "Rincon/config"

func VerifyConfig() {
	if config.Env == "" {
		config.Env = "PROD"
		SugarLogger.Debugln("ENV is not set, defaulting to PROD")
	}
	if config.Port == "" {
		config.Port = "10311"
		SugarLogger.Debugln("PORT is not set, defaulting to 10311")
	}
	if config.AuthUser == "" {
		config.AuthUser = "admin"
		SugarLogger.Debugln("AUTH_USER is not set, defaulting to admin")
	}
	if config.AuthPassword == "" {
		config.AuthPassword = "admin"
		SugarLogger.Debugln("AUTH_PASSWORD is not set, defaulting to admin")
	}
	if config.StorageMode != "local" && config.StorageMode != "sql" && config.StorageMode != "redis" && config.StorageMode != "redis+sql" {
		if config.DatabaseDriver != "" && config.DatabaseHost != "" {
			SugarLogger.Infoln("STORAGE_MODE is not set, defaulting to sql")
			config.StorageMode = "sql"
		}
		config.StorageMode = "local"
		SugarLogger.Infoln("STORAGE_MODE is not set, defaulting to local")
	}
}
