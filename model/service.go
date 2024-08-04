package model

import (
	"rincon/config"
	"time"
)

type Service struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Endpoint    string    `json:"endpoint" gorm:"unique"`
	HealthCheck string    `json:"health_check"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime;precision:6"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;precision:6"`
}

func (Service) TableName() string {
	return config.DatabaseTablePrefix + "service"
}

type ServiceDependency struct {
	ParentID  string    `json:"parent_id" gorm:"primaryKey"`
	ChildID   string    `json:"child_id" gorm:"primaryKey"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;precision:6"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;precision:6"`
}

func (ServiceDependency) TableName() string {
	return config.DatabaseTablePrefix + "service_dependency"
}
