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
			data := map[string]interface{}{
				"user_fullName": session.Values["FullName"],
				"user_auth":     session.Values["loggedIn"],
			}
			view, _ := template.ParseFiles("views/Home/index.html")
			_ = view.Execute(w, data)
		}
	}

}
