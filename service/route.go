package service

import (
	"fmt"
	"rincon/config"
	"rincon/database"
	"rincon/model"
	"rincon/utils"
	"slices"
	"strings"
	"time"
)

func GetAllRoutes() []model.Route {
	routes := make([]model.Route, 0)
	if config.StorageMode == "sql" {
		database.DB.Find(&routes)
	} else {
		routes = database.Local.Routes
	}
	return routes
}

func GetNumRoutes() int {
	if config.StorageMode == "sql" {
		var count int64
		database.DB.Model(&model.Route{}).Count(&count)
		return int(count)

	} else {
		return len(database.Local.Routes)
	}
}

func GetRoutesByRoute(route string) []model.Route {
	route = strings.TrimPrefix(route, "/")
	route = strings.TrimSuffix(route, "/")
	route = "/" + route
	routes := make([]model.Route, 0)
	if config.StorageMode == "sql" {
		database.DB.Where("route = ?", route).Find(&routes)
	} else {
		for _, r := range database.Local.Routes {
			if r.Route == route {
				routes = append(routes, r)
			}
		}
	}
	return routes
}

func GetRoutesByServiceName(name string) []model.Route {
	name = utils.NormalizeName(name)
	routes := make([]model.Route, 0)
	if config.StorageMode == "sql" {
		database.DB.Where("service_name = ?", name).Find(&routes)
	} else {
		for _, r := range database.Local.Routes {
			if r.ServiceName == name {
				routes = append(routes, r)
			}
		}
	}
	return routes
}

func GetRouteByRouteAndMethod(route string, method string) model.Route {
	route = strings.TrimPrefix(route, "/")
	route = strings.TrimSuffix(route, "/")
	route = "/" + route
	routes := GetRoutesByRoute(route)
	for _, r := range routes {
		if strings.Contains(r.Method, method) || strings.Contains(r.Method, "*") {
			return r
		}
	}
	return model.Route{}
}

func GetRouteByRouteAndService(route string, service string) model.Route {
	route = strings.TrimPrefix(route, "/")
	route = strings.TrimSuffix(route, "/")
	route = "/" + route
	service = utils.NormalizeName(service)
	routes := GetRoutesByRoute(route)
	for _, r := range routes {
		if r.ServiceName == service {
			return r
		}
	}
	return model.Route{}
}

func CreateRoute(route model.Route) error {
	if route.Route == "" {
		return fmt.Errorf("route cannot be empty")
	} else if route.ServiceName == "" {
		return fmt.Errorf("service name cannot be empty")
	} else if !strings.HasPrefix(route.Route, "/") {
		return fmt.Errorf("route must start with a slash")
	} else if strings.HasSuffix(route.Route, "/") {
		return fmt.Errorf("route cannot end with a slash")
	} else if !route.IsMethodValid() {
		return fmt.Errorf("invalid method %s", route.Method)
	}
	route.Method = strings.ToUpper(route.Method)
	route.ServiceName = utils.NormalizeName(route.ServiceName)
	route.ID = fmt.Sprintf("%s-[%s]", route.Route, route.Method)
	route.CreatedAt = time.Now()

	overlapRoutes := GetOverlappingRoutes(route)
	if len(overlapRoutes) == 1 && overlapRoutes[0].ServiceName == route.ServiceName {
		utils.SugarLogger.Debugf("stacking existing route %s for service %s", route.Route, route.ServiceName)
		route.Method = StackMethods(route.Method, overlapRoutes[0].Method)
		route.ID = fmt.Sprintf("%s-[%s]", route.Route, route.Method)
		DeleteRoute(overlapRoutes[0].ID)
	} else if len(overlapRoutes) > 0 {
		if config.OverwriteRoutes == "true" {
			for _, r := range overlapRoutes {
				if r.ServiceName == route.ServiceName {
					route.Method = StackMethods(route.Method, r.Method)
					route.ID = fmt.Sprintf("%s-[%s]", route.Route, route.Method)
				}
				DeleteRoute(r.ID)
			}
		} else {
			return fmt.Errorf("route with id %s overlaps with existing routes [%s]", route.ID, PrintRouteArray(overlapRoutes))
		}
	}
	existingRoute := GetRouteByRouteAndService(route.Route, route.ServiceName)
	if existingRoute.ID != "" {
		route.Method = StackMethods(route.Method, existingRoute.Method)
		route.ID = fmt.Sprintf("%s-[%s]", route.Route, route.Method)
		DeleteRoute(existingRoute.ID)
	}

	if config.StorageMode == "sql" {
		database.DB.Create(&route)
	} else {
		database.Local.Routes = append(database.Local.Routes, route)
	}
	utils.SugarLogger.Infof("route with id %s registered for service %s", route.ID, route.ServiceName)
	return nil
}

func PrintRouteArray(routes []model.Route) string {
	s := ""
	for i, r := range routes {
		s += fmt.Sprintf("[%s] %s (%s)", r.Method, r.Route, r.ServiceName)
		if i != len(routes)-1 {
			s += ", "
		}
	}
	return s
}

func GetOverlappingRoutes(route model.Route) []model.Route {
	methodMap := make(map[string]string)
	for _, method := range model.ValidMethods {
		methodMap[method] = ""
	}

	route.Method = strings.ToUpper(route.Method)
	route.ServiceName = utils.NormalizeName(route.ServiceName)
	overlapRoutes := make([]model.Route, 0)
	overlapServices := make([]string, 0)
	existingRoutes := GetRoutesByRoute(route.Route)

	for _, r := range existingRoutes {
		for _, m := range strings.Split(r.Method, ",") {
			methodMap[m] = r.ServiceName
			if m == "*" {
				return existingRoutes
			}
		}
	}
	for _, m := range strings.Split(route.Method, ",") {
		if m == "*" {
			return existingRoutes
		}
		if methodMap[m] != "" {
			if !slices.Contains(overlapServices, methodMap[m]) {
				overlapServices = append(overlapServices, methodMap[m])
				overlapRoutes = append(overlapRoutes, GetRouteByRouteAndService(route.Route, methodMap[m]))
			}
		}
	}
	return overlapRoutes
}

func StackMethods(m1 string, m2 string) string {
	methodMap := make(map[string]bool)
	for _, method := range model.ValidMethods {
		methodMap[method] = false
	}
	for _, method := range strings.Split(m1, ",") {
		methodMap[method] = true
	}
	for _, method := range strings.Split(m2, ",") {
		methodMap[method] = true
	}
	if methodMap["*"] {
		return "*"
	}
	methods := make([]string, 0)
	for method, enabled := range methodMap {
		if enabled {
			methods = append(methods, method)
		}
	}
	return strings.Join(methods, ",")
}

func DeleteRoute(id string) {
	if config.StorageMode == "sql" {
		database.DB.Where("id = ?", id).Delete(&model.Route{})
	} else {
		for i, r := range database.Local.Routes {
			if r.ID == id {
				database.Local.Routes = append(database.Local.Routes[:i], database.Local.Routes[i+1:]...)
				break
			}
		}
	}
	utils.SugarLogger.Infof("route with id %s deleted", id)
}

func MatchRoute(route string, method string) model.Service {
	if utils.SugarLogger.Level().String() == "debug" {
		PrintRouteGraph()
	}
	var service model.Service
	graph := GetRouteGraph()
	utils.SugarLogger.Debugf("Matching route  /" + route)
	matchedRoute := TraverseGraph("", route, method, graph)
	if matchedRoute == "" {
		utils.SugarLogger.Errorf("No route found for /%s", route)
		return service
	}
	utils.SugarLogger.Debugf("Matched to " + matchedRoute)
	for _, r := range GetAllRoutes() {
		if r.Route == matchedRoute && (strings.Contains(r.Method, method) || strings.Contains(r.Method, "*")) {
			service.Name = r.ServiceName
			break
		}
	}
	service = LoadBalance(service.Name, "random")
	if service.ID == 0 {
		utils.SugarLogger.Infoln("No eligible service instance found for route /" + route)
		go DeleteRoute(matchedRoute)
	} else {
		utils.SugarLogger.Infof("Matched route /%s to %s for service %s (%d)", route, matchedRoute, service.Name, service.ID)
	}
	return service
}

func TraverseGraph(path string, route string, method string, graph map[string][]model.RouteNode) string {
	currPathCount := strings.Count(path, "/")
	routeSlugCount := strings.Count("/"+route, "/")
	lastSlug := strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
	pathWithoutLastSlug := strings.TrimSuffix(path, "/"+lastSlug)

	utils.SugarLogger.Debugf("Traversing graph with path \"%s\" and route \"/%s\"", path, route)

	if pathWithoutLastSlug == "" {
		pathWithoutLastSlug = "/"
	}
	cindex := HasChildPath(lastSlug, graph[pathWithoutLastSlug])
	if lastSlug != "" && cindex == -1 {
		utils.SugarLogger.Debugf("Child path \"%s\" does not exist", lastSlug)
		return ""
	}
	if lastSlug == "**" && CanRouteHandleMethod(graph[pathWithoutLastSlug][cindex], method) {
		utils.SugarLogger.Debugf("Found all path wildcard (**)")
		return path
	}
	utils.SugarLogger.Debugf("Child path \"%s\" exists", lastSlug)

	if currPathCount == routeSlugCount {
		utils.SugarLogger.Debugf("Reached end of route")
		if CanRouteHandleMethod(graph[pathWithoutLastSlug][cindex], method) {
			utils.SugarLogger.Debugf("Route can handle method %s", method)
			return path
		} else {
			utils.SugarLogger.Debugf("Route cannot handle method %s", method)
			return ""
		}
	}

	nextSlug := strings.Split("/"+route, "/")[currPathCount+1]
	slugBranch := TraverseGraph(path+"/"+nextSlug, route, method, graph)
	if slugBranch != "" {
		return slugBranch
	}
	anyBranch := TraverseGraph(path+"/*", route, method, graph)
	if anyBranch != "" {
		return anyBranch
	}
	allBranch := TraverseGraph(path+"/**", route, method, graph)
	if allBranch != "" {
		return allBranch
	}

	return ""
}

func HasChildPath(path string, children []model.RouteNode) int {
	for i, c := range children {
		if c.Path == path {
			return i
		}
	}
	return -1
}

func CanRouteHandleMethod(route model.RouteNode, method string) bool {
	for _, s := range route.Services {
		if strings.Contains(s.Method, method) || strings.Contains(s.Method, "*") {
			return true
		}
	}
	return false
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
				endpoint := false
				if i == len(slugs)-1 {
					endpoint = true
				}
				if index := HasChildPath(slugs[i], children[parent]); index != -1 {
					if endpoint {
						children[parent][index].Services = append(children[parent][index].Services, model.RouteService{
							ServiceName: r.ServiceName,
							Method:      r.Method,
						})
					}
				} else {
					if endpoint {
						children[parent] = append(children[parent], model.RouteNode{
							ID:        parent + "/" + slugs[i],
							Path:      slugs[i],
							Services:  []model.RouteService{{ServiceName: r.ServiceName, Method: r.Method}},
							CreatedAt: time.Now(),
						})
					} else {
						children[parent] = append(children[parent], model.RouteNode{
							ID:        parent + "/" + slugs[i],
							Path:      slugs[i],
							Services:  []model.RouteService{},
							CreatedAt: time.Now(),
						})
					}
				}
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
			println("   -> " + n.Path + " (" + PrintRouteServices(n.Services) + ")")
		}
	}
	println("===========================")
}

func PrintRouteServices(rs []model.RouteService) string {
	s := ""
	for i, r := range rs {
		s += fmt.Sprintf("[%s] %s", r.Method, r.ServiceName)
		if i != len(rs)-1 {
			s += ", "
		}
	}
	return s
}
