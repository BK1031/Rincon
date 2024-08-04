package service

import (
	"rincon/config"
	"rincon/model"
	"testing"
)

func TestGetAllRoutesLocal(t *testing.T) {
	ResetLocalDB()
	t.Run("Test No Routes", func(t *testing.T) {
		routes := GetAllRoutes()
		if len(routes) != 0 {
			t.Errorf("Expected length to be 0")
		}
	})
	t.Run("Test One Route", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		routes := GetAllRoutes()
		if len(routes) != 1 {
			t.Errorf("Expected length to be 1")
		}
	})
	t.Run("Test Two Routes", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/test2",
			ServiceName: "Service 2",
			Method:      "GET",
		})
		routes := GetAllRoutes()
		if len(routes) != 2 {
			t.Errorf("Expected length to be 2")
		}
	})
}

func TestGetNumRoutesLocal(t *testing.T) {
	ResetLocalDB()
	t.Run("Test No Routes", func(t *testing.T) {
		num := GetNumRoutes()
		if num != 0 {
			t.Errorf("Expected length to be 0")
		}
	})
	t.Run("Test One Route", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		num := GetNumRoutes()
		if num != 1 {
			t.Errorf("Expected length to be 1")
		}
	})
}

func TestGetRoutesByRouteLocal(t *testing.T) {
	ResetLocalDB()
	t.Run("Test No Routes", func(t *testing.T) {
		routes := GetRoutesByRoute("/test")
		if len(routes) != 0 {
			t.Errorf("Expected length to be 0")
		}
	})
	t.Run("Test One Route", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		routes := GetRoutesByRoute("/test")
		if len(routes) != 1 {
			t.Errorf("Expected length to be 1")
		}
	})
	t.Run("Test Two Routes", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/test2",
			ServiceName: "Service 2",
			Method:      "GET",
		})
		CreateRoute(model.Route{
			Route:       "/test2",
			ServiceName: "Service 1",
			Method:      "POST",
		})
		routes := GetRoutesByRoute("/test2")
		if len(routes) != 2 {
			t.Errorf("Expected length to be 2")
		}
	})
}

func TestGetRoutesByServiceNameLocal(t *testing.T) {
	ResetLocalDB()
	t.Run("Test No Routes", func(t *testing.T) {
		routes := GetRoutesByServiceName("Service 1")
		if len(routes) != 0 {
			t.Errorf("Expected length to be 0")
		}
	})
	t.Run("Test One Route", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		CreateRoute(model.Route{
			Route:       "/test/example",
			ServiceName: "Service 2",
			Method:      "GET",
		})
		routes := GetRoutesByServiceName("Service 1")
		if len(routes) != 1 {
			t.Errorf("Expected length to be 1")
		}
	})
	t.Run("Test Two Routes", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/another/route",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		routes := GetRoutesByServiceName("Service 1")
		if len(routes) != 2 {
			t.Errorf("Expected length to be 2")
		}
	})
}

func TestGetRouteByRouteAndMethodLocal(t *testing.T) {
	ResetLocalDB()
	t.Run("Test No Routes", func(t *testing.T) {
		route := GetRouteByRouteAndMethod("/test", "GET")
		if route.ID != "" {
			t.Errorf("Expected route id to be empty")
		}
	})
	t.Run("Test One Route", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		route := GetRouteByRouteAndMethod("/test", "GET")
		if route.ID != "/test-[GET]" {
			t.Errorf("Expected route id to be set")
		}
	})
	t.Run("Test Two Routes", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 2",
			Method:      "POST",
		})
		route := GetRouteByRouteAndMethod("/test", "GET")
		if route.ID != "/test-[GET]" {
			t.Errorf("Expected route id to be set")
		}
	})
	t.Run("Test Wildcard Method", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/amazing",
			ServiceName: "Service 3",
			Method:      "*",
		})
		route := GetRouteByRouteAndMethod("/amazing", "GET")
		if route.ID != "/amazing-[*]" {
			t.Errorf("Expected route id to be set")
		}
	})
}

func TestGetRouteByRouteAndServiceLocal(t *testing.T) {
	ResetLocalDB()
	t.Run("Test No Routes", func(t *testing.T) {
		route := GetRouteByRouteAndService("/test", "Service 1")
		if route.ID != "" {
			t.Errorf("Expected route id to be empty")
		}
	})
	t.Run("Test One Route", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		route := GetRouteByRouteAndService("/test", "Service 1")
		if route.ID != "/test-[GET]" {
			t.Errorf("Expected route id to be set")
		}
	})
	t.Run("Test Two Routes", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 2",
			Method:      "POST",
		})
		route := GetRouteByRouteAndService("/test", "Service 2")
		if route.ID != "/test-[POST]" {
			t.Errorf("Expected route id to be set")
		}
	})
}

func TestCreateRouteLocal(t *testing.T) {
	ResetLocalDB()
	config.OverwriteRoutes = "false"
	t.Run("Test No Route", func(t *testing.T) {
		route := model.Route{
			ServiceName: "Service 1",
		}
		err := CreateRoute(route)
		if err == nil {
			t.Errorf("No error when creating route: %v", err)
		}
	})
	t.Run("Test No Service Name", func(t *testing.T) {
		route := model.Route{
			Route: "/test",
		}
		err := CreateRoute(route)
		if err == nil {
			t.Errorf("No error when creating route: %v", err)
		}
	})
	t.Run("Test No Method Name", func(t *testing.T) {
		route := model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
		}
		err := CreateRoute(route)
		if err == nil {
			t.Errorf("No error when creating route: %v", err)
		}
	})
	t.Run("Test Route Ends With Slash", func(t *testing.T) {
		route := model.Route{
			Route:       "/test/",
			ServiceName: "Service 1",
			Method:      "GET",
		}
		err := CreateRoute(route)
		if err == nil {
			t.Errorf("No error when creating route: %v", err)
		}
	})
	t.Run("Test Create Route", func(t *testing.T) {
		route := model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		}
		err := CreateRoute(route)
		if err != nil {
			t.Errorf("Error when creating route: %v", err)
		}
	})
	t.Run("Test Route Exists Upgrade Method", func(t *testing.T) {
		route := model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET,POST",
		}
		err := CreateRoute(route)
		if err != nil {
			t.Errorf("Error when creating route: %v", err)
		}
		route = GetRouteByRouteAndService("/test", "Service 1")
		if route.Method != "GET,POST" {
			t.Errorf("Route method not updated, found %s", route.Method)
		}
	})
	t.Run("Test Route Exists No Overlap", func(t *testing.T) {
		route := model.Route{
			Route:       "/test",
			ServiceName: "Service 2",
			Method:      "DELETE",
		}
		err := CreateRoute(route)
		if err != nil {
			t.Errorf("Error when creating route: %v", err)
		}
	})
	t.Run("Test Route Exists Overlap Deny Overwrite", func(t *testing.T) {
		route := model.Route{
			Route:       "/test",
			ServiceName: "Service 2",
			Method:      "GET,DELETE",
		}
		config.OverwriteRoutes = "false"
		err := CreateRoute(route)
		if err == nil {
			t.Errorf("Expected error when creating route")
		}
	})
	t.Run("Test Route Exists Overlap Allow Overwrite", func(t *testing.T) {
		route := model.Route{
			Route:       "/test",
			ServiceName: "Service 2",
			Method:      "GET,DELETE",
		}
		config.OverwriteRoutes = "true"
		err := CreateRoute(route)
		if err != nil {
			t.Errorf("Error when creating route: %v", err)
		}
		route = GetRouteByRouteAndService("/test", "Service 2")
		if route.Method != "GET,DELETE" {
			t.Errorf("Service 2 route method not updated, found %s", route.Method)
		}
		route = GetRouteByRouteAndService("/test", "Service 1")
		if route.ID != "" {
			t.Errorf("Found old route for Service 1")
		}
	})
}

func TestGetOverlappingRoutes(t *testing.T) {
	t.Run("Test No Overlap Same Service", func(t *testing.T) {
		ResetLocalDB()
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		overlap := GetOverlappingRoutes(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "POST",
		})
		if len(overlap) != 0 {
			t.Errorf("Expected no overlap")
		}
	})
	t.Run("Test No Overlap Different Service", func(t *testing.T) {
		ResetLocalDB()
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		overlap := GetOverlappingRoutes(model.Route{
			Route:       "/test",
			ServiceName: "Service 2",
			Method:      "POST",
		})
		if len(overlap) != 0 {
			t.Errorf("Expected no overlap")
		}
	})
	t.Run("Test Overlap Same Service", func(t *testing.T) {
		ResetLocalDB()
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		overlap := GetOverlappingRoutes(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET,POST",
		})
		if len(overlap) != 1 {
			t.Errorf("Expected overlap to be 1")
		}
		if overlap[0].ID != "/test-[GET]" {
			t.Errorf("Expected overlap to be /test-[GET]")
		}
	})
	t.Run("Test Overlap Different Service", func(t *testing.T) {
		ResetLocalDB()
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		overlap := GetOverlappingRoutes(model.Route{
			Route:       "/test",
			ServiceName: "Service 2",
			Method:      "GET,POST",
		})
		if len(overlap) != 1 {
			t.Errorf("Expected overlap to be 1")
		}
		if overlap[0].ID != "/test-[GET]" {
			t.Errorf("Expected overlap to be /test-[GET]")
		}
	})
	t.Run("Test Wildcard Existing Service", func(t *testing.T) {
		ResetLocalDB()
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "*",
		})
		overlap := GetOverlappingRoutes(model.Route{
			Route:       "/test",
			ServiceName: "Service 2",
			Method:      "GET,POST",
		})
		if len(overlap) != 1 {
			t.Errorf("Expected overlap to be 1")
		}
		if overlap[0].ID != "/test-[*]" {
			t.Errorf("Expected overlap to be /test-[*]")
		}
	})
	t.Run("Test Wildcard New Service", func(t *testing.T) {
		ResetLocalDB()
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 2",
			Method:      "PUT",
		})
		overlap := GetOverlappingRoutes(model.Route{
			Route:       "/test",
			ServiceName: "Service 2",
			Method:      "*",
		})
		if len(overlap) != 2 {
			t.Errorf("Expected overlap to be 2")
		}
		if overlap[0].ID != "/test-[GET]" {
			t.Errorf("Expected overlap 1 to be /test-[GET]")
		}
		if overlap[1].ID != "/test-[PUT]" {
			t.Errorf("Expected overlap 2 to be /test-[PUT]")
		}
	})
}

func TestGetAllRoutesSQL(t *testing.T) {
	ResetSQLDB()
	t.Run("Test No Routes", func(t *testing.T) {
		routes := GetAllRoutes()
		if len(routes) != 0 {
			t.Errorf("Expected length to be 0")
		}
	})
	t.Run("Test One Route", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		routes := GetAllRoutes()
		if len(routes) != 1 {
			t.Errorf("Expected length to be 1")
		}
	})
	t.Run("Test Two Routes", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/test2",
			ServiceName: "Service 2",
			Method:      "GET",
		})
		routes := GetAllRoutes()
		if len(routes) != 2 {
			t.Errorf("Expected length to be 2")
		}
	})
}

func TestGetNumRoutesSQL(t *testing.T) {
	ResetSQLDB()
	t.Run("Test No Routes", func(t *testing.T) {
		num := GetNumRoutes()
		if num != 0 {
			t.Errorf("Expected length to be 0")
		}
	})
	t.Run("Test One Route", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		num := GetNumRoutes()
		if num != 1 {
			t.Errorf("Expected length to be 1")
		}
	})
}

func TestGetRoutesByRouteSQL(t *testing.T) {
	ResetSQLDB()
	t.Run("Test No Routes", func(t *testing.T) {
		routes := GetRoutesByRoute("/test")
		if len(routes) != 0 {
			t.Errorf("Expected length to be 0")
		}
	})
	t.Run("Test One Route", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		routes := GetRoutesByRoute("/test")
		if len(routes) != 1 {
			t.Errorf("Expected length to be 1")
		}
	})
	t.Run("Test Two Routes", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/test2",
			ServiceName: "Service 2",
			Method:      "GET",
		})
		CreateRoute(model.Route{
			Route:       "/test2",
			ServiceName: "Service 1",
			Method:      "POST",
		})
		routes := GetRoutesByRoute("/test2")
		if len(routes) != 2 {
			t.Errorf("Expected length to be 2")
		}
	})
}

func TestGetRoutesByServiceNameSQL(t *testing.T) {
	ResetSQLDB()
	t.Run("Test No Routes", func(t *testing.T) {
		routes := GetRoutesByServiceName("Service 1")
		if len(routes) != 0 {
			t.Errorf("Expected length to be 0")
		}
	})
	t.Run("Test One Route", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		CreateRoute(model.Route{
			Route:       "/test/example",
			ServiceName: "Service 2",
			Method:      "GET",
		})
		routes := GetRoutesByServiceName("Service 1")
		if len(routes) != 1 {
			t.Errorf("Expected length to be 1")
		}
	})
	t.Run("Test Two Routes", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/another/route",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		routes := GetRoutesByServiceName("Service 1")
		if len(routes) != 2 {
			t.Errorf("Expected length to be 2")
		}
	})
}

func TestGetRouteByRouteAndMethodSQL(t *testing.T) {
	ResetSQLDB()
	t.Run("Test No Routes", func(t *testing.T) {
		route := GetRouteByRouteAndMethod("/test", "GET")
		if route.ID != "" {
			t.Errorf("Expected route id to be empty")
		}
	})
	t.Run("Test One Route", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		route := GetRouteByRouteAndMethod("/test", "GET")
		if route.ID != "/test-[GET]" {
			t.Errorf("Expected route id to be set")
		}
	})
	t.Run("Test Two Routes", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 2",
			Method:      "POST",
		})
		route := GetRouteByRouteAndMethod("/test", "GET")
		if route.ID != "/test-[GET]" {
			t.Errorf("Expected route id to be set")
		}
	})
	t.Run("Test Wildcard Method", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/amazing",
			ServiceName: "Service 3",
			Method:      "*",
		})
		route := GetRouteByRouteAndMethod("/amazing", "GET")
		if route.ID != "/amazing-[*]" {
			t.Errorf("Expected route id to be set")
		}
	})
}

func TestGetRouteByRouteAndServiceSQL(t *testing.T) {
	ResetSQLDB()
	t.Run("Test No Routes", func(t *testing.T) {
		route := GetRouteByRouteAndService("/test", "Service 1")
		if route.ID != "" {
			t.Errorf("Expected route id to be empty")
		}
	})
	t.Run("Test One Route", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		route := GetRouteByRouteAndService("/test", "Service 1")
		if route.ID != "/test-[GET]" {
			t.Errorf("Expected route id to be set")
		}
	})
	t.Run("Test Two Routes", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 2",
			Method:      "POST",
		})
		route := GetRouteByRouteAndService("/test", "Service 2")
		if route.ID != "/test-[POST]" {
			t.Errorf("Expected route id to be set")
		}
	})
}

func TestCreateRouteSQL(t *testing.T) {
	ResetSQLDB()
	config.OverwriteRoutes = "false"
	t.Run("Test No Route", func(t *testing.T) {
		route := model.Route{
			ServiceName: "Service 1",
		}
		err := CreateRoute(route)
		if err == nil {
			t.Errorf("No error when creating route: %v", err)
		}
	})
	t.Run("Test No Service Name", func(t *testing.T) {
		route := model.Route{
			Route: "/test",
		}
		err := CreateRoute(route)
		if err == nil {
			t.Errorf("No error when creating route: %v", err)
		}
	})
	t.Run("Test No Method Name", func(t *testing.T) {
		route := model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
		}
		err := CreateRoute(route)
		if err == nil {
			t.Errorf("No error when creating route: %v", err)
		}
	})
	t.Run("Test Route Ends With Slash", func(t *testing.T) {
		route := model.Route{
			Route:       "/test/",
			ServiceName: "Service 1",
			Method:      "GET",
		}
		err := CreateRoute(route)
		if err == nil {
			t.Errorf("No error when creating route: %v", err)
		}
	})
	t.Run("Test Create Route", func(t *testing.T) {
		route := model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		}
		err := CreateRoute(route)
		if err != nil {
			t.Errorf("Error when creating route: %v", err)
		}
	})
	t.Run("Test Route Exists Upgrade Method", func(t *testing.T) {
		route := model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET,POST",
		}
		err := CreateRoute(route)
		if err != nil {
			t.Errorf("Error when creating route: %v", err)
		}
		route = GetRouteByRouteAndService("/test", "Service 1")
		if route.Method != "GET,POST" {
			t.Errorf("Route method not updated, found %s", route.Method)
		}
	})
	t.Run("Test Route Exists No Overlap", func(t *testing.T) {
		route := model.Route{
			Route:       "/test",
			ServiceName: "Service 2",
			Method:      "DELETE",
		}
		err := CreateRoute(route)
		if err != nil {
			t.Errorf("Error when creating route: %v", err)
		}
	})
	t.Run("Test Route Exists Overlap Deny Overwrite", func(t *testing.T) {
		route := model.Route{
			Route:       "/test",
			ServiceName: "Service 2",
			Method:      "GET,DELETE",
		}
		config.OverwriteRoutes = "false"
		err := CreateRoute(route)
		if err == nil {
			t.Errorf("Expected error when creating route")
		}
	})
	t.Run("Test Route Exists Overlap Allow Overwrite", func(t *testing.T) {
		route := model.Route{
			Route:       "/test",
			ServiceName: "Service 2",
			Method:      "GET,DELETE",
		}
		config.OverwriteRoutes = "true"
		err := CreateRoute(route)
		if err != nil {
			t.Errorf("Error when creating route: %v", err)
		}
		route = GetRouteByRouteAndService("/test", "Service 2")
		if route.Method != "GET,DELETE" {
			t.Errorf("Service 2 route method not updated, found %s", route.Method)
		}
		route = GetRouteByRouteAndService("/test", "Service 1")
		if route.ID != "" {
			t.Errorf("Found old route for Service 1")
		}
	})
}

func TestGetOverlappingRoutesSQL(t *testing.T) {
	t.Run("Test No Overlap Same Service", func(t *testing.T) {
		ResetSQLDB()
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		overlap := GetOverlappingRoutes(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "POST",
		})
		if len(overlap) != 0 {
			t.Errorf("Expected no overlap")
		}
	})
	t.Run("Test No Overlap Different Service", func(t *testing.T) {
		ResetSQLDB()
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		overlap := GetOverlappingRoutes(model.Route{
			Route:       "/test",
			ServiceName: "Service 2",
			Method:      "POST",
		})
		if len(overlap) != 0 {
			t.Errorf("Expected no overlap")
		}
	})
	t.Run("Test Overlap Same Service", func(t *testing.T) {
		ResetSQLDB()
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		overlap := GetOverlappingRoutes(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET,POST",
		})
		if len(overlap) != 1 {
			t.Errorf("Expected overlap to be 1")
		}
		if overlap[0].ID != "/test-[GET]" {
			t.Errorf("Expected overlap to be /test-[GET]")
		}
	})
	t.Run("Test Overlap Different Service", func(t *testing.T) {
		ResetSQLDB()
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		overlap := GetOverlappingRoutes(model.Route{
			Route:       "/test",
			ServiceName: "Service 2",
			Method:      "GET,POST",
		})
		if len(overlap) != 1 {
			t.Errorf("Expected overlap to be 1")
		}
		if overlap[0].ID != "/test-[GET]" {
			t.Errorf("Expected overlap to be /test-[GET]")
		}
	})
	t.Run("Test Wildcard Existing Service", func(t *testing.T) {
		ResetSQLDB()
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "*",
		})
		overlap := GetOverlappingRoutes(model.Route{
			Route:       "/test",
			ServiceName: "Service 2",
			Method:      "GET,POST",
		})
		if len(overlap) != 1 {
			t.Errorf("Expected overlap to be 1")
		}
		if overlap[0].ID != "/test-[*]" {
			t.Errorf("Expected overlap to be /test-[*]")
		}
	})
	t.Run("Test Wildcard New Service", func(t *testing.T) {
		ResetSQLDB()
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
			Method:      "GET",
		})
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 2",
			Method:      "PUT",
		})
		overlap := GetOverlappingRoutes(model.Route{
			Route:       "/test",
			ServiceName: "Service 2",
			Method:      "*",
		})
		if len(overlap) != 2 {
			t.Errorf("Expected overlap to be 2")
		}
		if overlap[0].ID != "/test-[GET]" {
			t.Errorf("Expected overlap 1 to be /test-[GET]")
		}
		if overlap[1].ID != "/test-[PUT]" {
			t.Errorf("Expected overlap 2 to be /test-[PUT]")
		}
	})
}

func TestMatchRoute(t *testing.T) {
	ResetLocalDB()
	CreateService(model.Service{
		Name:        "Montecito",
		Version:     "1.4.2",
		Endpoint:    "http://localhost:10312",
		HealthCheck: "http://localhost:10312/health",
	})
	CreateService(model.Service{
		Name:        "Lacumbre",
		Version:     "2.7.9",
		Endpoint:    "http://localhost:10313",
		HealthCheck: "http://localhost:10313/health",
	})
	CreateRoute(model.Route{
		Route:       "/service/ping",
		ServiceName: "Montecito",
		Method:      "*",
	})
	CreateRoute(model.Route{
		Route:       "/service/*/awesome",
		ServiceName: "Montecito",
		Method:      "*",
	})
	CreateRoute(model.Route{
		Route:       "/service/**",
		ServiceName: "Lacumbre",
		Method:      "*",
	})
	CreateRoute(model.Route{
		Route:       "/awesome/test",
		ServiceName: "Service DNE",
		Method:      "GET",
	})
	CreateRoute(model.Route{
		Route:       "/awesome/test",
		ServiceName: "Montecito",
		Method:      "POST",
	})
	CreateRoute(model.Route{
		Route:       "/no/service",
		ServiceName: "No Service",
		Method:      "*",
	})
	t.Run("Test Match Route", func(t *testing.T) {
		route := MatchRoute("service/ping", "GET")
		if route.Name != "montecito" {
			t.Errorf("MatchRoute returned wrong service, got %s", route.Name)
		}
	})
	t.Run("Test Match Route 2", func(t *testing.T) {
		route := MatchRoute("service/1/awesome", "GET")
		if route.Name != "montecito" {
			t.Errorf("MatchRoute returned wrong service, got %s", route.Name)
		}
	})
	t.Run("Test Match Route 3", func(t *testing.T) {
		route := MatchRoute("service/1/awesome/2", "GET")
		if route.Name != "lacumbre" {
			t.Errorf("MatchRoute returned wrong service, got %s", route.Name)
		}
	})
	t.Run("Test Match Route 4", func(t *testing.T) {
		route := MatchRoute("epic/awesome", "GET")
		if route.Name != "" {
			t.Errorf("MatchRoute returned wrong service, got %s", route.Name)
		}
	})
	t.Run("Test Match Route 5", func(t *testing.T) {
		route := MatchRoute("no/service", "GET")
		if route.Name != "" {
			t.Errorf("MatchRoute returned wrong service, got %s", route.Name)
		}
	})
	t.Run("Test Match Route 6", func(t *testing.T) {
		route := MatchRoute("awesome/test", "GET")
		if route.Name != "" {
			t.Errorf("MatchRoute expected no service, got %s", route.Name)
		}
	})
	t.Run("Test Match Route 7", func(t *testing.T) {
		route := MatchRoute("awesome/test", "POST")
		if route.Name != "montecito" {
			t.Errorf("MatchRoute returned wrong service, got %s", route.Name)
		}
	})
	t.Run("Test Match Route 8", func(t *testing.T) {
		route := MatchRoute("awesome/test", "DELETE")
		if route.Name != "" {
			t.Errorf("MatchRoute expected no service, got %s", route.Name)
		}
	})
}

func TestPrintRouteGraph(t *testing.T) {
	ResetLocalDB()
	CreateService(model.Service{
		Name:        "Montecito",
		Version:     "1.4.2",
		Endpoint:    "http://localhost:10312",
		HealthCheck: "http://localhost:10312/health",
	})
	CreateService(model.Service{
		Name:        "Lacumbre",
		Version:     "2.7.9",
		Endpoint:    "http://localhost:10313",
		HealthCheck: "http://localhost:10313/health",
	})
	CreateRoute(model.Route{
		Route:       "/service/ping",
		ServiceName: "Montecito",
	})
	CreateRoute(model.Route{
		Route:       "/service/*/awesome",
		ServiceName: "Montecito",
	})
	CreateRoute(model.Route{
		Route:       "/service/**",
		ServiceName: "Lacumbre",
	})
	t.Run("Test Print Route Graph", func(t *testing.T) {
		PrintRouteGraph()
	})
}
