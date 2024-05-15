package api

import (
	"os"
	"rincon/database"
	"rincon/utils"
	"testing"
)

func TestMain(m *testing.M) {
	utils.InitializeLogger()
	utils.VerifyConfig()
	database.InitializeLocal()
	router := setupRouter()
	exitVal := m.Run()
	os.Exit(exitVal)
}