package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"ip_detector/internal/adapter/http/handler"
	"ip_detector/internal/adapter/http/middleware"
	"ip_detector/internal/app/service"
)

func SetupRouter(userService *service.UserService, jwtSecret string) http.Handler {
	r := mux.NewRouter()

	userHandler := handler.NewUserHandler(userService)

	r.HandleFunc("/register", userHandler.RegisterUser).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")

	protected := r.NewRoute().Subrouter()
	protected.Use(middleware.JWTMiddleware(jwtSecret))
	protected.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	protected.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")

	return r
}
