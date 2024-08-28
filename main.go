package main

import (
	"rincon/api"
	"rincon/config"
	"rincon/database"
	"rincon/model"
	"rincon/service"
	"rincon/utils"
)

func main() {
	config.PrintStartupBanner()
	utils.InitializeLogger()
	defer utils.Logger.Sync()

	utils.VerifyConfig()
	database.InitializeDB()
	service.RegisterSelf()
	service.InitializeHeartbeat()

	service.CreateService(model.Service{
		Name:        "New York",
		Endpoint:    "http://localhost:3000",
		HealthCheck: "http://localhost:3000/health",
		Version:     "1.0.0",
	})
	service.CreateService(model.Service{
		Name:        "San Francisco",
		Endpoint:    "http://localhost:4000",
		HealthCheck: "http://localhost:4000/health",
		Version:     "1.0.0",
	})

	router := api.SetupRouter()
	api.InitializeRoutes(router)
	err := router.Run(":" + config.Port)
	if err != nil {
		utils.SugarLogger.Fatalln(err)
	}
}
