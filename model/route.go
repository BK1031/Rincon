package model

type Route struct {
	Route       string `json:"route" gorm:"primaryKey"`
	ServiceName string `json:"service_name"`
	CreatedAt   string `json:"created_at" gorm:"autoCreateTime"`
}

func (Route) TableName() string {
	return "route"
}

type RouteNode struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Path        string `json:"path"`
	ServiceName string `json:"service_name"`
	CreatedAt   string `json:"created_at" gorm:"autoCreateTime"`
}

func (RouteNode) TableName() string {
	return "route_node"
}

type RouteEdge struct {
	ParentID string `json:"parent_id" gorm:"primaryKey"`
	ChildID  string `json:"child_id" gorm:"primaryKey"`
}

func (RouteEdge) TableName() string {
	return "route_edge"
}
