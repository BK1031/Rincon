package service

import (
	"fmt"
	"rincon/config"
	"rincon/database"
	"rincon/model"
	"rincon/utils"
)

func GetAllRoutes() []model.Route {
	routes := make([]model.Route, 0)
	routes = database.Local.Routes
	return routes
}

func GetNumRoutes() int {
	return len(database.Local.Routes)
}

func GetRouteByID(id string) model.Route {
	var route model.Route
	for _, r := range database.Local.Routes {
		if r.Route == id {
			route = r
			break
		}
	}
	return route
}

func GetRoutesByServiceName(name string) []model.Route {
	name = utils.NormalizeName(name)
	routes := make([]model.Route, 0)
	for _, r := range database.Local.Routes {
		if r.ServiceName == name {
			routes = append(routes, r)
		}
	}
	return routes
}

func CreateRoute(route model.Route) error {
	if route.Route == "" {
		return fmt.Errorf("route cannot be empty")
	} else if route.ServiceName == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	route.ServiceName = utils.NormalizeName(route.ServiceName)

	if GetRouteByID(route.Route).Route != "" && route.ServiceName != GetRouteByID(route.Route).ServiceName {
		if config.OverwriteRoutes == "true" {
			DeleteRoute(route.Route)
		} else {
			return fmt.Errorf("route with id %s already exists", route.Route)
		}
	}
	database.Local.Routes = append(database.Local.Routes, route)
	return nil
}

func DeleteRoute(id string) {
	for i, r := range database.Local.Routes {
		if r.Route == id {
			database.Local.Routes = append(database.Local.Routes[:i], database.Local.Routes[i+1:]...)
			break
		}
	}
}
