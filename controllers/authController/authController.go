package authController

import (
	"errors"
	"github.com/hudayberdipolatov/go-auth-with-session/helpers/authsession"
	"github.com/hudayberdipolatov/go-auth-with-session/libraries"
	"github.com/hudayberdipolatov/go-auth-with-session/models"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
)

type AuthController struct{}

type LoginInput struct {
	Username string `validate:"required,gte=3" label:"Username"`
	Password string `validate:"required,gte=4" label:"Password"`
}

type RegisterInput struct {
	Username        string `validate:"required,gte=3" label:"Username"`
	FullName        string `validate:"required,gte=3" label:"Full name"`
	Email           string `validate:"required,email" label:"Email"`
	Password        string `validate:"required,gte=4" label:"Password"`
	ConfirmPassword string `validate:"required,eqfield=Password" label:"Confirm Password"`
}

var userModel models.Users
var validation libraries.Validation

func (auth AuthController) RegisterPage(w http.ResponseWriter, r *http.Request) {
	session, _ := authsession.Store.Get(r, authsession.SESSION_ID)
	if len(session.Values) != 0 {

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	view, err := template.ParseFiles("views/auth/register.html")
	if err != nil {
		log.Println(err)
		return
	}
	_ = view.Execute(w, nil)
}
func (auth AuthController) Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	registerInput := RegisterInput{
		Username:        r.PostForm.Get("username"),
		FullName:        r.PostForm.Get("full_name"),
		Email:           r.PostForm.Get("email"),
		Password:        r.PostForm.Get("password"),
		ConfirmPassword: r.PostForm.Get("confirm_password"),
	}
	errorValidate := validation.Struct(registerInput)
	if errorValidate != nil {
		errorMessages := validation.Struct(registerInput)
		data := map[string]interface{}{
			"validation": errorMessages,
			"user":       registerInput,
		}
		view, _ := template.ParseFiles("views/auth/register.html")
		_ = view.Execute(w, data)
	} else {
		getUser := userModel.GetUser(registerInput.Username)
		if getUser.ID == 0 {
			hashPassword, _ := bcrypt.GenerateFromPassword([]byte(registerInput.Password), bcrypt.DefaultCost)
			password := string(hashPassword)
			models.Users{
				Username: registerInput.Username,
				FullName: registerInput.FullName,
				Email:    registerInput.Email,
				Password: password,
			}.CreateUser()
			data := make(map[string]interface{})
			data["success"] = "Register Successfully!!! Ulgama girip bilersiňiz!!!"
			view, _ := template.ParseFiles("views/auth/login.html")
			_ = view.Execute(w, data)
		} else {
			data := make(map[string]interface{})
			data["user_exists"] = registerInput.Username + " ulanyjy ady öň hem ulanylýar!!!"
			data["user"] = registerInput
			view, _ := template.ParseFiles("views/auth/register.html")
			_ = view.Execute(w, data)
		}
	}
}

func (auth AuthController) LoginPage(w http.ResponseWriter, r *http.Request) {
	session, _ := authsession.Store.Get(r, authsession.SESSION_ID)
	if len(session.Values) != 0 {

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
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
	errorValidate := validation.Struct(loginInput)
	if errorValidate != nil {
		//log.Println(errorValidate)
		data := map[string]interface{}{
			"validation": errorValidate,
			"user":       loginInput,
		}
		view, _ := template.ParseFiles("views/auth/login.html")
		_ = view.Execute(w, data)
	} else {
		var message error
		user := userModel.GetUser(loginInput.Username)
		if user.ID == 0 {
			message = errors.New("Username ýa-da password yalnyş!!!")
		} else {
			errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInput.Password))
			if errPassword != nil {
				message = errors.New("Username ýa-da password yalnyş!!!")
			}
		}
		if message != nil {
			data := map[string]interface{}{
				"error": message,
			}
			view, _ := template.ParseFiles("views/auth/login.html")
			_ = view.Execute(w, data)

		} else {
			user = userModel.GetUser(loginInput.Username)
			session, _ := authsession.Store.Get(r, authsession.SESSION_ID)
			session.Values["loggedIn"] = true
			session.Values["Username"] = user.Username
			session.Values["FullName"] = user.FullName
			session.Values["Email"] = user.Email
			_ = session.Save(r, w)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}
func (auth AuthController) Logout(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		session, _ := authsession.Store.Get(r, authsession.SESSION_ID)
		session.Options.MaxAge = -1
		_ = session.Save(r, w)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		view, _ := template.ParseFiles("views/errors/error.html")
		_ = view.Execute(w, nil)
	}

}
