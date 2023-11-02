package main

import (
	"github.com/hudayberdipolatov/go-auth-with-session/models"
	"github.com/hudayberdipolatov/go-auth-with-session/routes"
	"log"
	"net/http"
)

const portNumber = ":8080"

func main() {
	models.ConnectDataBase()
	log.Println("Server run Port", portNumber)
	http.ListenAndServe(portNumber, routes.Routes())
}
