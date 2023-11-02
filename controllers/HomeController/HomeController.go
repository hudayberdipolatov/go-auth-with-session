package HomeController

import (
	"html/template"
	"log"
	"net/http"
)

type HomeController struct{}

func (home HomeController) Index(w http.ResponseWriter, r *http.Request) {
	view, err := template.ParseFiles("views/Home/index.html")
	if err != nil {
		log.Println(err)
	}
	_ = view.Execute(w, nil)
}
