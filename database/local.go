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
		Services:            []model.Service{},
		ServiceDependencies: []model.ServiceDependency{},
		Routes:              []model.Route{},
		RouteNodes:          []model.RouteNode{},
		RouteEdges:          []model.RouteEdge{},
	}
}
