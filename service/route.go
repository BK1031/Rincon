package service

import (
	"rincon/database"
	"rincon/model"
)

func GetAllRoutes() []model.Route {
	routes := make([]model.Route, 0)
	routes = database.Local.Routes
	return routes
}

func GetNumRoutes() int {
	return len(database.Local.Routes)
}
