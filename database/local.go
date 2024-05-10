package database

import "rincon/model"

var Local *LocalStore

type LocalStore struct {
	Services            []model.Service
	ServiceDependencies []model.ServiceDependency
	Routes              []model.Route
	RouteNodes          []model.RouteNode
	RouteEdges          []model.RouteEdge
}

func InitializeLocal() {
	Local = &LocalStore{
		Services:            make([]model.Service, 0),
		ServiceDependencies: make([]model.ServiceDependency, 0),
		Routes:              make([]model.Route, 0),
		RouteNodes:          make([]model.RouteNode, 0),
		RouteEdges:          make([]model.RouteEdge, 0),
	}
}
