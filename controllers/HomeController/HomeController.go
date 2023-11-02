package HomeController

import (
	"github.com/hudayberdipolatov/go-auth-with-session/helpers/authsession"
	"html/template"
	"net/http"
)

type HomeController struct{}

func (home HomeController) Index(w http.ResponseWriter, r *http.Request) {

	session, _ := authsession.Store.Get(r, authsession.SESSION_ID)
	if len(session.Values) == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		if session.Values["loggedIn"] != true {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			view, _ := template.ParseFiles("views/Home/index.html")
			_ = view.Execute(w, nil)
		}
	}

}
