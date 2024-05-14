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
	} else if strings.HasSuffix(route.Route, "/") {
		return fmt.Errorf("route cannot end with a slash")
	}
	route.ServiceName = utils.NormalizeName(route.ServiceName)
	route.CreatedAt = time.Now()

	if GetRouteByID(route.Route).Route != "" {
		if route.ServiceName != GetRouteByID(route.Route).ServiceName {
			if config.OverwriteRoutes == "true" {
				DeleteRoute(route.Route)
			} else {
				utils.SugarLogger.Errorf("route with id %s already exists", route.Route)
				return fmt.Errorf("route with id %s already exists", route.Route)
			}
		} else {
			utils.SugarLogger.Debugf("route with id %s for service %s already exists", route.Route, route.ServiceName)
			return nil
		}
	}
	database.Local.Routes = append(database.Local.Routes, route)
	utils.SugarLogger.Infof("route with id %s registered for service %s", route.Route, route.ServiceName)
	return nil
}

func DeleteRoute(id string) {
	for i, r := range database.Local.Routes {
		if r.Route == id {
			database.Local.Routes = append(database.Local.Routes[:i], database.Local.Routes[i+1:]...)
			break
		}
	}
	utils.SugarLogger.Infof("route with id %s deleted", id)
}

func MatchRoute(route string) model.Service {
	PrintRouteGraph()
	var service model.Service
	graph := GetRouteGraph()
	slugs := strings.Split(route, "/")
	println("MATCHING ROUTE FOR  /" + route)
	println("slugs are", len(slugs))
	matchedRoute := TraverseGraph("", route, graph)
	if matchedRoute == "" {
		println("no route found")
	}
	println("matched route is", matchedRoute)
	return service
}

func TraverseGraph(path string, route string, graph map[string][]model.RouteNode) string {
	// assume path starts as "/"
	// assume route looks like "/rincon/services/awesome/routes"
	println(path, "path has ", strings.Count(path, "/"), "slugs")
	println(route, "route has ", strings.Count("/"+route, "/"), "slugs")

	currPathIndex := strings.Count(path, "/")
	routeSlugCount := strings.Count("/"+route, "/")

	if currPathIndex == routeSlugCount {
		return path
	}

	if path == "" {
		path = "/"
	}
	// get current slug
	slugs := strings.Split(route, "/")
	currentSlug := slugs[currPathIndex]
	println("current slug is", currentSlug)

	if CheckChildren(currentSlug, graph[path]) != "" {
		if path == "/" {
			path += currentSlug
		} else {
			path += "/" + currentSlug
		}
		return TraverseGraph(path, route, graph)
	} else if CheckChildren("*", graph[path]) != "" {
		if path == "/" {
			path += "*"
		} else {
			path += "/*"
		}
		return TraverseGraph(path, route, graph)
	} else if CheckChildren("**", graph[path]) != "" {
		if path == "/" {
			path += "**"
		} else {
			path += "/**" + currentSlug
		}
		return path
	}
	return ""
}

func CheckChildren(path string, children []model.RouteNode) string {
	for _, c := range children {
		println("checking", c.Path, "against", path)
		if c.Path == path {
			println("found match!")
			return c.Path
		}
	}
	return ""
}

func GetRouteGraph() map[string][]model.RouteNode {
	children := make(map[string][]model.RouteNode)
	routes := GetAllRoutes()
	for _, r := range routes {
		slugs := strings.Split(r.Route, "/")
		parent := ""
		for i := 0; i < len(slugs); i++ {
			if slugs[i] != "" {
				if parent == "" {
					parent = "/"
				}
				if _, exists := children[parent]; !exists {
					children[parent] = make([]model.RouteNode, 0)
				}
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
	return children
}

func PrintRouteGraph() {
	println("======= ROUTE GRAPH =======")
	graph := GetRouteGraph()
	for k, v := range graph {
		println(k)
		for _, n := range v {
			println("   -> " + n.Path + " (" + n.ServiceName + ")")
		}
	}
	println("===========================")
}
