package authController

import (
	"html/template"
	"log"
	"net/http"
)

type AuthController struct{}

func (auth AuthController) RegisterPage(w http.ResponseWriter, r *http.Request) {
	view, err := template.ParseFiles("views/auth/register.html")
	if err != nil {
		log.Println(err)
		return
	}
	_ = view.Execute(w, nil)
}
func (auth AuthController) Register(w http.ResponseWriter, r *http.Request) {

}

func (auth AuthController) LoginPage(w http.ResponseWriter, r *http.Request) {
	view, err := template.ParseFiles("views/auth/login.html")
	if err != nil {
		log.Println(err)
		return
	}
	_ = view.Execute(w, nil)
}

func (auth AuthController) Login(w http.ResponseWriter, r *http.Request) {

}
