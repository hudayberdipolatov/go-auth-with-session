package authController

import (
	"errors"
	"github.com/hudayberdipolatov/go-auth-with-session/helpers/authsession"
	"github.com/hudayberdipolatov/go-auth-with-session/models"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
)

type AuthController struct{}

type LoginInput struct {
	Username string
	Password string
}

var userModel models.Users

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
	r.ParseForm()
	loginInput := LoginInput{
		Username: r.PostForm.Get("username"),
		Password: r.PostForm.Get("password"),
	}
	var message error
	if loginInput.Username == "" {
		message = errors.New("Username ya-da password yalnys!!!")

	} else if loginInput.Password == "" {
		message = errors.New("Username ya-da password yalnys!!!")
	} else {
		user := userModel.GetUser(loginInput.Username)
		if user.ID == 0 {
			message = errors.New("Username ya-da password yalnys!!!")
		} else {
			errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInput.Password))
			if errPassword != nil {
				message = errors.New("Username ya-da password yalnys!!!")
			}
		}
	}
	if message != nil {
		data := map[string]interface{}{
			"error": message,
		}
		view, _ := template.ParseFiles("views/auth/login.html")
		_ = view.Execute(w, data)
	} else {
		user := userModel.GetUser(loginInput.Username)
		session, _ := authsession.Store.Get(r, authsession.SESSION_ID)
		session.Values["loggedIn"] = true
		session.Values["Username"] = user.Username
		session.Values["FullName"] = user.FullName
		session.Values["Email"] = user.Email
		_ = session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (auth AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := authsession.Store.Get(r, authsession.SESSION_ID)
	session.Options.MaxAge = -1
	_ = session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
