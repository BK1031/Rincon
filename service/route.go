package service

import (
	"fmt"
	"rincon/config"
	"rincon/database"
	"rincon/model"
	"rincon/utils"
	"strings"
	"time"
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
	route.CreatedAt = time.Now()
	defer BuildRouteGraph()

	if GetRouteByID(route.Route).Route != "" {
		if route.ServiceName != GetRouteByID(route.Route).ServiceName {
			if config.OverwriteRoutes == "true" {
				DeleteRoute(route.Route)
			} else {
				return fmt.Errorf("route with id %s already exists", route.Route)
			}
		} else {
			utils.SugarLogger.Debugf("route with id %s for service %s already exists", route.Route, route.ServiceName)
			return nil
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
	BuildRouteGraph()
}

func BuildRouteGraph() {
	children := make(map[string][]model.RouteNode)
	for _, r := range database.Local.Routes {
		slugs := strings.Split(r.Route, "/")
		parent := ""
		for i := 0; i < len(slugs); i++ {
			if slugs[i] != "" {
				if parent == "" {
					parent = "/"
				}
				println("i: ", i)
				println("len: ", len(slugs))
				println("Parent: ", parent)
				if _, exists := children[parent]; !exists {
					children[parent] = make([]model.RouteNode, 0)
				}
				println("Slug: ", slugs[i])
				// delete existing node
				for j, n := range children[parent] {
					if n.Path == slugs[i] {
						children[parent] = append(children[parent][:j], children[parent][j+1:]...)
						break
					}
				}
				name := ""
				if i == len(slugs)-1 {
					name = r.ServiceName
				}
				children[parent] = append(children[parent], model.RouteNode{
					ID:          parent + "/" + slugs[i],
					Path:        slugs[i],
					ServiceName: name,
					CreatedAt:   time.Now(),
				})
				if parent == "/" {
					parent += slugs[i]
				} else {
					parent += "/" + slugs[i]
				}
			}
		}
	}
	println("=====================================")
	println("Route Graph")
	println("Nodes:", len(children))
	for parent, nodes := range children {
		println(parent)
		for _, n := range nodes {
			println(" -> " + n.Path + " (" + n.ServiceName + ")")
		}
	}
}
