package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/agroconnect/api-gateway/config"
	"github.com/agroconnect/api-gateway/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func newReverseProxy(targetURL string) http.Handler {
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("Invalid target URL for proxy %s: %v", targetURL, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = target.Host
	}
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Del("Access-Control-Allow-Origin")
		resp.Header.Del("Access-Control-Allow-Credentials")
		resp.Header.Del("Access-Control-Allow-Methods")
		resp.Header.Del("Access-Control-Allow-Headers")
		resp.Header.Del("Access-Control-Expose-Headers")
		return nil
	}

	return proxy
}

func main() {
	cfg := config.LoadConfig()

	catalogProxy := newReverseProxy(cfg.CatalogServiceURL)
	orderProxy := newReverseProxy(cfg.OrderServiceURL)
	weatherProxy := newReverseProxy(cfg.WeatherServiceURL)

	authMiddleware := middleware.JWTAuthMiddleware(cfg.JWTSecret)

	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"service": "agroconnect-api-gateway",
			"status":  "UP",
			"version": "1.0.0",
		})
	}).Methods("GET")

	router.Handle("/api/auth/register", orderProxy).Methods("POST")
	router.Handle("/api/auth/login", orderProxy).Methods("POST")

	router.Handle("/api/products", catalogProxy).Methods("GET")
	router.Handle("/api/products/{id}", catalogProxy).Methods("GET")
	router.Handle("/api/farmers", catalogProxy).Methods("GET")
	router.Handle("/api/farmers/{slug}", catalogProxy).Methods("GET")

	router.Handle("/api/weather", weatherProxy).Methods("GET")

	router.Handle("/api/orders/{code}", orderProxy).Methods("GET")

	router.Handle("/api/products", authMiddleware(catalogProxy)).Methods("POST")
	router.Handle("/api/products/{id}", authMiddleware(catalogProxy)).Methods("DELETE")
	router.Handle("/api/products/{id}/stock", authMiddleware(catalogProxy)).Methods("PATCH")

	router.Handle("/api/orders", authMiddleware(orderProxy)).Methods("POST")
	router.Handle("/api/orders/user", authMiddleware(orderProxy)).Methods("GET")

	router.PathPrefix("/api/products").Handler(catalogProxy)
	router.PathPrefix("/api/farmers").Handler(catalogProxy)
	router.PathPrefix("/api/auth").Handler(orderProxy)
	router.PathPrefix("/api/orders").Handler(orderProxy)
	router.PathPrefix("/api/weather").Handler(weatherProxy)

	corsHandler := cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler(router)

	log.Printf("AgroConnect API Gateway listening on port :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, corsHandler); err != nil {
		log.Fatalf("API Gateway terminated: %v", err)
	}
}
