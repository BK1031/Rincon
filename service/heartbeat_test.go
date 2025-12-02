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
		Name:        "Montecito",
		Version:     "2.7.9",
		Endpoint:    "http://localhost:10313",
		HealthCheck: "https://bk1031.dev",
	})
	t.Run("Test Invalid Heartbeat Interval", func(t *testing.T) {
		config.HeartbeatInterval = "bruh"
		InitializeHeartbeat()
	})
	t.Run("Test Initialize Server Heartbeat", func(t *testing.T) {
		config.HeartbeatType = "server"
		InitializeHeartbeat()
	})
	t.Run("Test Initialize Client Heartbeat", func(t *testing.T) {
		config.HeartbeatType = "client"
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
			HealthCheck: "https://bk1031.dev",
		})
		ClientHeartbeat(1)
	})
	t.Run("Test Server Heartbeat", func(t *testing.T) {
		config.HeartbeatType = "server"
		CreateService(model.Service{
			Name:        "Lacumbre",
			Version:     "2.7.9",
			Endpoint:    "http://localhost:10313",
			HealthCheck: "https://bk1031.dev",
		})
		CreateService(model.Service{
			Name:        "Montecito",
			Version:     "1.4.2",
			Endpoint:    "http://localhost:10312",
			HealthCheck: "http://localhost:10312/health",
		})
		CreateService(model.Service{
			Name:        "Montecito",
			Version:     "1.4.2",
			Endpoint:    "http://localhost:10314",
			HealthCheck: "https://bk1031.dev/health",
		})
		ServerHeartbeat(1)
	})
	t.Run("Test Server Heartbeat With Retry", func(t *testing.T) {
		config.HeartbeatType = "server"
		config.HeartbeatRetryCount = "2"
		config.HeartbeatRetryBackoff = "100"
		CreateService(model.Service{
			Name:        "RetryTest",
			Version:     "1.0.0",
			Endpoint:    "http://localhost:19999",
			HealthCheck: "http://localhost:19999/health",
		})
		ServerHeartbeat(1)
	})
	t.Run("Test Server Heartbeat With Invalid Retry Config", func(t *testing.T) {
		config.HeartbeatType = "server"
		config.HeartbeatRetryCount = "invalid"
		config.HeartbeatRetryBackoff = "invalid"
		CreateService(model.Service{
			Name:        "InvalidRetryTest",
			Version:     "1.0.0",
			Endpoint:    "http://localhost:19998",
			HealthCheck: "http://localhost:19998/health",
		})
		ServerHeartbeat(1)
		// Reset to valid values
		config.HeartbeatRetryCount = "3"
		config.HeartbeatRetryBackoff = "1000"
	})
	t.Run("Test Server Heartbeat With Zero Retries", func(t *testing.T) {
		config.HeartbeatType = "server"
		config.HeartbeatRetryCount = "0"
		config.HeartbeatRetryBackoff = "100"
		CreateService(model.Service{
			Name:        "ZeroRetryTest",
			Version:     "1.0.0",
			Endpoint:    "http://localhost:19997",
			HealthCheck: "http://localhost:19997/health",
		})
		ServerHeartbeat(1)
	})
}
