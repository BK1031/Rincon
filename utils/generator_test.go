package utils

import (
	"rincon/config"
	"testing"
)

func TestGenerateID(t *testing.T) {
	InitializeLogger()
	t.Run("Test Generate ID", func(t *testing.T) {
		InitializeLogger()
		config.ServiceIDLength = "6"
		id := GenerateID(0)
		println(id)
		if id < 100000 || id > 999999 {
			t.Errorf("Service ID Length is not 6")
		}
	})
	t.Run("Test Generate ID with Length", func(t *testing.T) {
		InitializeLogger()
		id := GenerateID(10)
		println(id)
		if id < 1000000000 || id > 9999999999 {
			t.Errorf("Service ID Length is not 10")
		}
	})
}
