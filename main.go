package main

import (
	"rincon/config"
	"rincon/database"
	"rincon/utils"
)

func main() {
	config.PrintStartupBanner()
	utils.InitializeLogger()
	defer utils.Logger.Sync()

	utils.VerifyConfig()
	database.InitializeLocal()
	database.InitializeDB()
}
