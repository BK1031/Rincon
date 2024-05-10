package service

import (
	"rincon/config"
	"rincon/database"
	"rincon/model"
	"testing"
)

func TestCreateRouteLocal(t *testing.T) {
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
	t.Run("Test Create Route", func(t *testing.T) {
		route := model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
		}
		err := CreateRoute(route)
		if err != nil {
			t.Errorf("Error when creating route: %v", err)
		}
	})
	t.Run("Test Route Exists 1", func(t *testing.T) {
		route := model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
		}
		config.OverwriteRoutes = "true"
		err := CreateRoute(route)
		if err != nil {
			t.Errorf("Error when creating route: %v", err)
		}
	})
	t.Run("Test Route Exists 2", func(t *testing.T) {
		route := model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
		}
		config.OverwriteRoutes = "false"
		err := CreateRoute(route)
		if err != nil {
			t.Errorf("Error when creating route: %v", err)
		}
	})
	t.Run("Test Route Exists 3", func(t *testing.T) {
		route := model.Route{
			Route:       "/test",
			ServiceName: "Service 2",
		}
		config.OverwriteRoutes = "false"
		err := CreateRoute(route)
		if err == nil {
			t.Errorf("No error when creating route: %v", err)
		}
	})
	t.Run("Test Route Exists 4", func(t *testing.T) {
		route := model.Route{
			Route:       "/test",
			ServiceName: "Service 2",
		}
		config.OverwriteRoutes = "true"
		err := CreateRoute(route)
		if err != nil {
			t.Errorf("Error when creating route: %v", err)
		}
	})
}

func TestGetRoutesLocal(t *testing.T) {
	t.Run("Test Get All Routes", func(t *testing.T) {
		database.Local.Routes = make([]model.Route, 0)
		routes := GetAllRoutes()
		if len(routes) != 0 {
			t.Errorf("Expected length to be 0")
		}
	})
	t.Run("Test Get Num Routes", func(t *testing.T) {
		CreateRoute(model.Route{
			Route:       "/test",
			ServiceName: "Service 1",
		})
		num := GetNumRoutes()
		if num == 0 {
			t.Errorf("No routes found")
		}
	})
	t.Run("Test Get Route By ID", func(t *testing.T) {
		route := GetRouteByID("/test")
		if route.Route != "/test" {
			t.Errorf("No route found")
		}
	})
	t.Run("Test Get Routes By Service Name", func(t *testing.T) {
		routes := GetRoutesByServiceName("Service 1")
		if len(routes) == 0 {
			t.Errorf("No routes found")
		}
	})
}
