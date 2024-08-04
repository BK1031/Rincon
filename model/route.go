package model

import (
	"rincon/config"
	"slices"
	"strings"
	"time"
)

var validMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD", "*"}

type Route struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Route       string    `json:"route"`
	Method      string    `json:"method"`
	ServiceName string    `json:"service_name"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;precision:6"`
}

func (Route) TableName() string {
	return config.DatabaseTablePrefix + "route"
}

func (r *Route) IsMethodValid() bool {
	methods := strings.Split(r.Method, ",")
	for _, method := range methods {
		if !slices.Contains(validMethods, method) {
			return false
		}
	}
	return true
}

type RouteNode struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	Path      string         `json:"path"`
	Services  []RouteService `json:"services"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime;precision:6"`
}

type RouteService struct {
	ServiceName string `json:"service_name"`
	Method      string `json:"method"`
}
