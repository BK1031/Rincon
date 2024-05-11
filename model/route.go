package model

import "time"

type Route struct {
	Route       string    `json:"route" gorm:"primaryKey"`
	ServiceName string    `json:"service_name"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (Route) TableName() string {
	return "route"
}

type RouteNode struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Path        string    `json:"path"`
	ServiceName string    `json:"service_name"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (RouteNode) TableName() string {
	return "route_node"
}
