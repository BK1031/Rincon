package api

import (
	"net/http"
	"net/http/httptest"
	"rincon/model"
	"strings"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func TestGetRoute(t *testing.T) {
	router := SetupRouter()
	InitializeRoutes(router)
	t.Run("Get All Routes", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/rincon/routes", nil)
		router.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Errorf("Expected status code 200, got %v", w.Code)
		}
	})
	t.Run("Get Routes By Route", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/rincon/routes?route=/rincon/ping", nil)
		router.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Errorf("Expected status code 200, got %v", w.Code)
		}
	})
	t.Run("Get Routes By Route And Service", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/rincon/routes?route=/rincon/ping&service=rincon", nil)
		router.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Errorf("Expected status code 200, got %v", w.Code)
		}
	})
	t.Run("Get Routes By Route And Service Not Found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/rincon/routes?route=/rincon/ping&service=montecito", nil)
		router.ServeHTTP(w, req)
		if w.Code != 404 {
			t.Errorf("Expected status code 404, got %v", w.Code)
		}
	})
	t.Run("Get Routes By Route And Method", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/rincon/routes?route=/rincon/ping&method=GET", nil)
		router.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Errorf("Expected status code 200, got %v", w.Code)
		}
	})
	t.Run("Get Routes By Route And Method Not Found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/rincon/routes?route=/rincon/ping/wow&method=POST", nil)
		router.ServeHTTP(w, req)
		if w.Code != 404 {
			t.Errorf("Expected status code 404, got %v", w.Code)
		}
	})
}

func TestGetRoutesForService(t *testing.T) {
	router := SetupRouter()
	InitializeRoutes(router)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/rincon/services/rincon/routes", nil)
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %v", w.Code)
	}
}

func TestCreateRoute(t *testing.T) {
	router := SetupRouter()
	InitializeRoutes(router)
	t.Run("Create Route", func(t *testing.T) {
		r := model.Route{
			Route:       "/rincon/epic",
			ServiceName: "rincon",
			Method:      "GET",
		}
		json, _ := jsoniter.MarshalToString(r)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/rincon/routes", strings.NewReader(json))
		req.SetBasicAuth("admin", "admin")
		router.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Errorf("Expected status code 200, got %v", w.Code)
		}
	})
	t.Run("Invalid Body", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/rincon/routes", strings.NewReader("{\"service_name\": 123}"))
		req.SetBasicAuth("admin", "admin")
		router.ServeHTTP(w, req)
		if w.Code != 400 {
			t.Errorf("Expected status code 400, got %v", w.Code)
		}
	})
	t.Run("Create Route Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/rincon/routes", strings.NewReader("{\"service_name\": \"rincon\"}"))
		req.SetBasicAuth("admin", "admin")
		router.ServeHTTP(w, req)
		if w.Code != 500 {
			t.Errorf("Expected status code 500, got %v", w.Code)
		}
	})
}

func TestMatchRoute(t *testing.T) {
	router := SetupRouter()
	InitializeRoutes(router)
	t.Run("Match Route", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/rincon/match?route=/rincon/ping&method=GET", nil)
		router.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Errorf("Expected status code 200, got %v", w.Code)
		}
	})
	t.Run("Match Route Not Found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/rincon/match?route=/rincon/bruh&method=GET", nil)
		router.ServeHTTP(w, req)
		if w.Code != 404 {
			t.Errorf("Expected status code 404, got %v", w.Code)
		}
	})
	t.Run("Match Route With No Route", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/rincon/match?method=GET", nil)
		router.ServeHTTP(w, req)
		if w.Code != 400 {
			t.Errorf("Expected status code 400, got %v", w.Code)
		}
	})
	t.Run("Match Route With No Method", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/rincon/match?route=/test", nil)
		router.ServeHTTP(w, req)
		if w.Code != 400 {
			t.Errorf("Expected status code 400, got %v", w.Code)
		}
	})
}
