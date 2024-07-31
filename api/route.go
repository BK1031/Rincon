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
	route := strings.ReplaceAll(c.Param("id"), "<->", "/")
	route = strings.TrimPrefix(route, "/")
	route = strings.TrimSuffix(route, "/")
	result := service.GetRouteByID("/" + route)
	if result.Route == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": "No route /" + route + " found"})
		return
	}
	c.JSON(http.StatusOK, result)
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
	c.JSON(http.StatusOK, service.GetRouteByID(input.Route))
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

	result := service.MatchRoute(route)
	if result.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("No route [%s] /%s found", method, route)})
		return
	}
	c.JSON(http.StatusOK, result)
}
