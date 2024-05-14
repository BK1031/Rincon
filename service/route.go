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
	utils.SugarLogger.Debugf("Matching route  /" + route)
	//utils.SugarLogger.Debugf("input route has %d slugs", len(slugs))
	matchedRoute := TraverseGraph("", route, graph)
	utils.SugarLogger.Debugf("Matched to " + matchedRoute)
	for _, r := range GetAllRoutes() {
		if r.Route == matchedRoute {
			service.Name = r.ServiceName
			break
		}
	}
	if service.Name == "" {
		utils.SugarLogger.Debugf("No service found for route /" + route)
	}
	service = LoadBalance(service.Name, "random")
	if service.ID == 0 {
		utils.SugarLogger.Debugf("No eligible service instance found for" + service.Name)
	} else {
		utils.SugarLogger.Infof("Matched route /%s to %s for service %s (%d)", route, matchedRoute, service.Name, service.ID)
	}
	return service
}

func TraverseGraph(path string, route string, graph map[string][]model.RouteNode) string {
	currPathCount := strings.Count(path, "/")
	routeSlugCount := strings.Count("/"+route, "/")
	lastSlug := strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
	pathWithoutLastSlug := strings.TrimSuffix(path, "/"+lastSlug)

	utils.SugarLogger.Debugf("Traversing graph with path %s and route /%s", path, route)

	if pathWithoutLastSlug == "" {
		pathWithoutLastSlug = "/"
	}
	if lastSlug != "" && HasChildPath(lastSlug, graph[pathWithoutLastSlug]) == "" {
		utils.SugarLogger.Debugf("Child path %s does not exist", lastSlug)
		return ""
	}
	if lastSlug == "**" {
		utils.SugarLogger.Debugf("Found all path wildcard (**)")
		return path
	}
	println("child path found for " + lastSlug)

	if currPathCount == routeSlugCount {
		println("path and route have same slug count")
		return path
	}

	nextSlug := strings.Split("/"+route, "/")[currPathCount+1]
	println("nextSlug: " + nextSlug)

	slugBranch := TraverseGraph(path+"/"+nextSlug, route, graph)
	if slugBranch != "" {
		return slugBranch
	}
	anyBranch := TraverseGraph(path+"/*", route, graph)
	if anyBranch != "" {
		return anyBranch
	}
	allBranch := TraverseGraph(path+"/**", route, graph)
	if allBranch != "" {
		return allBranch
	}

	return ""
}

func HasChildPath(path string, children []model.RouteNode) string {
	for _, c := range children {
		if c.Path == path {
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
