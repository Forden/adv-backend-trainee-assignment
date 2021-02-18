package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"adv-backend-trainee-assignment/config"
	"adv-backend-trainee-assignment/src/db"
	"adv-backend-trainee-assignment/src/routes"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.RawQuery = strings.ToLower(r.URL.RawQuery)
		h.ServeHTTP(w, r)
	})
}

func newRouter(cfg config.MyConfig) *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	var server routes.APIServer
	var err error
	switch cfg.UsedDB {
	case "postgresql":
		server.DBManager, err = db.NewPostgreSQLManager(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.PostgreSQL.Username, cfg.PostgreSQL.Password, cfg.PostgreSQL.Host, cfg.PostgreSQL.Port, cfg.PostgreSQL.DBName, cfg.PostgreSQL.SSLMode))
	}
	if err != nil {
		log.Fatalf("couldn't connect to db: %s", err)
	}
	for _, route := range routes.GenerateRoutes(server) {
		s.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(Middleware(route.HandlerFunc))
	}
	return r
}

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalf("empty path to config file")
	}
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("couldn't load config. error: [%s] path to config: [%s]", err, configPath)
	} else {
		router := newRouter(cfg)
		address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
		log.Printf("Starting server on: %s", address)
		log.Fatal(http.ListenAndServe(address, router))
	}
}
