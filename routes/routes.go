package routes

import (
	"github.com/gorilla/mux"
	"github.com/hudayberdipolatov/go-auth-with-session/controllers/HomeController"
	"github.com/hudayberdipolatov/go-auth-with-session/controllers/authController"
	"net/http"
)

func Routes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", HomeController.HomeController{}.Index).Methods("GET")
	router.HandleFunc("/register", authController.AuthController{}.RegisterPage).Methods("GET")
	router.HandleFunc("/register", authController.AuthController{}.Register).Methods(http.MethodPost)
	router.HandleFunc("/login", authController.AuthController{}.LoginPage).Methods("GET")
	router.HandleFunc("/login", authController.AuthController{}.Login).Methods(http.MethodPost)
	router.HandleFunc("/logout", authController.AuthController{}.Logout).Methods("GET")
	return router
}
