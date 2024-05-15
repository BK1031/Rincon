package service

import (
	"rincon/config"
	"rincon/model"
	"testing"
	"time"
)

func TestInitializeHeartbeat(t *testing.T) {
	CreateService(model.Service{
		Name:        "Montecito",
		Version:     "1.4.2",
		Endpoint:    "http://localhost:10312",
		HealthCheck: "http://localhost:10312/health",
	})
	CreateService(model.Service{
		Name:        "Lacumbre",
		Version:     "2.7.9",
		Endpoint:    "http://localhost:10313",
		HealthCheck: "http://localhost:10313/health",
	})
	t.Run("Test Invalid Heartbeat Interval", func(t *testing.T) {
		config.HeartbeatInterval = "bruh"
		InitializeHeartbeat()
	})
	t.Run("Test Client Heartbeat", func(t *testing.T) {
		config.HeartbeatType = "client"
		config.HeartbeatInterval = "1"
		time.Sleep(3 * time.Second)
		CreateService(model.Service{
			Name:        "Lacumbre",
			Version:     "2.7.9",
			Endpoint:    "http://localhost:10313",
			HealthCheck: "http://localhost:10313/health",
		})
		InitializeHeartbeat()
		ClientHeartbeat(1)
	})
	t.Run("Test Server Heartbeat", func(t *testing.T) {
		config.HeartbeatType = "server"
		InitializeHeartbeat()
		ServerHeartbeat(1)
	})
}
