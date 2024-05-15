package api

import (
	"os"
	"rincon/database"
	"rincon/service"
	"rincon/utils"
	"testing"
)

func TestMain(m *testing.M) {
	utils.InitializeLogger()
	utils.VerifyConfig()
	database.InitializeLocal()
	service.RegisterSelf()
	exitVal := m.Run()
	os.Exit(exitVal)
}
