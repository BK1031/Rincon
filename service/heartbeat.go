package service

import (
	"github.com/go-co-op/gocron/v2"
	"rincon/config"
	"rincon/utils"
	"strconv"
	"time"
)

func InitializeHeartbeat() {
	s, err := gocron.NewScheduler()
	if err != nil {
		utils.SugarLogger.Fatalln("Error creating scheduler")
	}

	interval, err := strconv.Atoi(config.HeartbeatInterval)
	if err != nil {
		utils.SugarLogger.Debugln("HEARTBEAT_INTERVAL is invalid, defaulting to 10")
		interval = 10
	}
	if config.HeartbeatType == "client" {
		j, err := s.NewJob(
			gocron.DurationJob(time.Duration(interval)*time.Second),
			gocron.NewTask(ClientHeartbeat, interval),
		)
		if err != nil {
			utils.SugarLogger.Errorf("Error creating job: %v", err)
		}
		utils.SugarLogger.Infof("Job ID: %d", j.ID)
		s.Start()
	} else {
		j, err := s.NewJob(
			gocron.DurationJob(time.Duration(interval)*time.Second),
			gocron.NewTask(ServerHeartbeat, interval),
		)
		if err != nil {
			utils.SugarLogger.Errorf("Error creating job: %v", err)
		}
		utils.SugarLogger.Infof("Job ID: %d", j.ID)
		s.Start()
	}
}

func ServerHeartbeat(interval int) {
	utils.SugarLogger.Infoln("Running server heartbeats...")
}

func ClientHeartbeat(interval int) {
	utils.SugarLogger.Infoln("Checking client heartbeats...")
	for _, s := range GetAllServices() {
		delta := time.Now().Sub(s.UpdatedAt).Milliseconds()
		utils.SugarLogger.Infof("Last %s (%d) ping was %d milliseconds ago", s.Name, s.ID, delta)
		if delta > int64((interval+1)*1000) {
			utils.SugarLogger.Errorf("Service %s (%d) registration expired!", s.Name, s.ID)
			go RemoveService(s.ID)
		}
	}
}
