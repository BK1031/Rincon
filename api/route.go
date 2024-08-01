package api

import (
	"fmt"
	"net/http"
	"rincon/model"
	"rincon/service"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetAllRoutes(c *gin.Context) {
	result := service.GetAllRoutes()
	c.JSON(http.StatusOK, result)
}

func GetRoute(c *gin.Context) {
	r := c.Query("route")
	r = strings.TrimPrefix(r, "/")
	r = strings.TrimSuffix(r, "/")
	m := c.Query("method")
	m = strings.ToUpper(m)
	s := c.Query("service")

	if r == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Route query is required"})
		return
	}
	if s != "" {
		requestedRoute := service.GetRouteByRouteAndService(r, s)
		if requestedRoute.ID == "" {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("No route /%s for service %s found", r, s)})
			return
		}
		c.JSON(http.StatusOK, requestedRoute)
		return
	} else if m != "" {
		requestedRoute := service.GetRouteByRouteAndMethod(r, m)
		if requestedRoute.ID == "" {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("No route [%s] /%s found", m, r)})
			return
		}
		c.JSON(http.StatusOK, requestedRoute)
		return
	} else {
		c.JSON(http.StatusOK, service.GetRoutesByRoute(r))
		return
	}
}

func GetRoutesForService(c *gin.Context) {
	result := service.GetRoutesByServiceName(c.Param("name"))
	c.JSON(http.StatusOK, result)
}

func CreateRoute(c *gin.Context) {
	var input model.Route
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := service.CreateRoute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, service.GetRouteByRouteAndService(input.Route, input.ServiceName))
}

func MatchRoute(c *gin.Context) {
	route := c.Query("route")
	route = strings.TrimPrefix(route, "/")
	route = strings.TrimSuffix(route, "/")
	method := c.Query("method")
	method = strings.ToUpper(method)

	if route == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Route query is required"})
		return
	}
	if method == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Method query is required"})
		return
	}

	result := service.MatchRoute(route, method)
	if result.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("No route [%s] /%s found", method, route)})
		return
	}
	c.JSON(http.StatusOK, result)
}
