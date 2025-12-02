package service

import (
	"rincon/config"
	"rincon/utils"
	"strconv"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/go-resty/resty/v2"
)

func InitializeHeartbeat() {
	s, _ := gocron.NewScheduler()
	interval, err := strconv.Atoi(config.HeartbeatInterval)
	if err != nil {
		utils.SugarLogger.Debugln("HEARTBEAT_INTERVAL is invalid, defaulting to 10")
		interval = 10
	}
	if config.HeartbeatType == "client" {
		j, _ := s.NewJob(
			gocron.DurationJob(time.Duration(interval)*time.Second),
			gocron.NewTask(ClientHeartbeat, interval),
		)
		utils.SugarLogger.Infof("Job ID: %d", j.ID)
		s.Start()
	} else {
		j, _ := s.NewJob(
			gocron.DurationJob(time.Duration(interval)*time.Second),
			gocron.NewTask(ServerHeartbeat, interval),
		)
		utils.SugarLogger.Infof("Job ID: %d", j.ID)
		s.Start()
	}
}

func ServerHeartbeat(interval int) {
	client := resty.New()
	retryCount, err := strconv.Atoi(config.HeartbeatRetryCount)
	if err != nil {
		utils.SugarLogger.Debugln("HEARTBEAT_RETRY_COUNT is invalid, defaulting to 3")
		retryCount = 3
	}
	retryBackoff, err := strconv.Atoi(config.HeartbeatRetryBackoff)
	if err != nil {
		utils.SugarLogger.Debugln("HEARTBEAT_RETRY_BACKOFF is invalid, defaulting to 5000ms")
		retryBackoff = 5000
	}

	for _, s := range GetAllServices() {
		success := false
		var lastErr error
		var lastResp *resty.Response

		for attempt := 0; attempt <= retryCount; attempt++ {
			if attempt > 0 {
				utils.SugarLogger.Debugf("Retrying heartbeat for %s (%d), attempt %d/%d", s.Name, s.ID, attempt, retryCount)
				time.Sleep(time.Duration(retryBackoff) * time.Millisecond)
			}

			resp, err := client.R().Get(s.HealthCheck)
			if err != nil {
				lastErr = err
				lastResp = nil
				utils.SugarLogger.Warnf("Error pinging %s (%d) on attempt %d: %v", s.Name, s.ID, attempt+1, err)
				continue
			}

			if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
				utils.SugarLogger.Infof("Pinged %s (%d) in %dms", s.Name, s.ID, resp.Time().Milliseconds())
				success = true
				break
			}

			lastErr = nil
			lastResp = resp
			utils.SugarLogger.Warnf("Error pinging %s (%d) on attempt %d: status %d", s.Name, s.ID, attempt+1, resp.StatusCode())
		}

		if !success {
			if lastErr != nil {
				utils.SugarLogger.Errorf("Failed to ping %s (%d) after %d retries: %v", s.Name, s.ID, retryCount, lastErr)
			} else if lastResp != nil {
				utils.SugarLogger.Errorf("Failed to ping %s (%d) after %d retries: status %d", s.Name, s.ID, retryCount, lastResp.StatusCode())
			}
			RemoveService(s.ID)
		}
	}
}

func ClientHeartbeat(interval int) {
	for _, s := range GetAllServices() {
		delta := time.Since(s.UpdatedAt).Milliseconds()
		utils.SugarLogger.Infof("Last %s (%d) ping was %dms ago", s.Name, s.ID, delta)
		if delta > int64((interval+1)*1000) && s.Name != "rincon" {
			utils.SugarLogger.Errorf("Service %s (%d) registration expired!", s.Name, s.ID)
			RemoveService(s.ID)
		}
	}
}
