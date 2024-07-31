package model

import (
	"rincon/config"
	"slices"
	"strings"
	"time"
)

var validMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD", "*"}

type Route struct {
	Route       string    `json:"route" gorm:"primaryKey"`
	Method      string    `json:"method" gorm:"primaryKey"`
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
	ID          string    `json:"id" gorm:"primaryKey"`
	Path        string    `json:"path"`
	ServiceName string    `json:"service_name"`
	Method      string    `json:"method"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;precision:6"`
}
