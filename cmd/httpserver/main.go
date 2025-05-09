// @title           IP Detector API
// @version         1.0
// @description     Register with IP detecting
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
package main

import (
	"database/sql"
	"ip_detector/internal/logger"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"

	_ "github.com/lib/pq"

	_ "ip_detector/docs"
	"ip_detector/internal/adapter/db/postgres"
	"ip_detector/internal/adapter/external/geoip"
	"ip_detector/internal/adapter/http/router"
	"ip_detector/internal/app/service"
	"ip_detector/internal/config"
)

func main() {
	logger.Init()
	cfg := config.LoadConfig()

	applyMigrations(cfg.DSN())

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	userRepo := postgres.NewPostgresUserRepo(db)
	geoIP := geoip.NewIPAPIService("http://ip-api.com/json")

	serviceConfig := &service.Config{
		JWTSecret:     cfg.JWTSecret,
		JWTExpiration: cfg.JWTExpiration,
	}

	userService := service.NewUserService(userRepo, geoIP, serviceConfig)

	r := router.SetupRouter(userService, cfg.JWTSecret).(*mux.Router)

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func applyMigrations(dsn string) {
	m, err := migrate.New(
		"file://./migrations",
		dsn,
	)

	if err != nil {
		log.Fatalf("Migrations initialisation failed: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Applying migration failed: %v", err)
	}

	log.Println("Migrations applied")
}
