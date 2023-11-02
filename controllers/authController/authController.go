package authController

import (
	"errors"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/hudayberdipolatov/go-auth-with-session/helpers/authsession"
	"github.com/hudayberdipolatov/go-auth-with-session/models"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"reflect"
)

type AuthController struct{}

type LoginInput struct {
	Username string
	Password string
}

type RegisterInput struct {
	Username        string `validate:"required,gte=3" label:"Username"`
	FullName        string `validate:"required,gte=3" label:"Full name"`
	Email           string `validate:"required,email" label:"Email"`
	Password        string `validate:"required,gte=4" label:"Password"`
	ConfirmPassword string `validate:"required,eqfield=Password" label:"Confirm Password"`
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
	r.ParseForm()
	registerInput := RegisterInput{
		Username:        r.PostForm.Get("username"),
		FullName:        r.PostForm.Get("full_name"),
		Email:           r.PostForm.Get("email"),
		Password:        r.PostForm.Get("password"),
		ConfirmPassword: r.PostForm.Get("confirm_password"),
	}
	//errorMessages := make(map[string]interface{})
	//if registerInput.Username == "" {
	//	errorMessages["username_validate"] = "Username hökmany"
	//}
	//if registerInput.FullName == "" {
	//	errorMessages["FullName_validate"] = "FullName hökmany"
	//}
	//if registerInput.Email == "" {
	//	errorMessages["email_validate"] = "Email hökmany"
	//}
	//if registerInput.Password == "" {
	//	errorMessages["password_validate"] = "Password hökmany"
	//}
	//if registerInput.ConfirmPassword == "" {
	//	errorMessages["ConfirmPassword_validate"] = "Confirm Password hökmany"
	//} else {
	//	if registerInput.Password != registerInput.ConfirmPassword {
	//		errorMessages["ConfirmPassword_validate"] = "Password-lar biri-birine gabat gelenok!!!"
	//	}
	//}
	//
	//if len(errorMessages) > 0 {
	//	data := map[string]interface{}{
	//		"validation": errorMessages,
	//	}
	//	//log.Fatal(errorMessages)
	//	view, _ := template.ParseFiles("views/auth/register.html")
	//	_ = view.Execute(w, data)
	//} else {
	//
	//	getUser := userModel.GetUser(registerInput.Username)
	//	if getUser.ID == 0 {
	//		hashPassword, _ := bcrypt.GenerateFromPassword([]byte(registerInput.Password), bcrypt.DefaultCost)
	//		password := string(hashPassword)
	//		models.Users{
	//			Username: registerInput.Username,
	//			FullName: registerInput.FullName,
	//			Email:    registerInput.Email,
	//			Password: password,
	//		}.CreateUser()
	//		data := make(map[string]interface{})
	//		data["success"] = "Registered Successfully!!! Ulgama girip bilersiňiz!!!"
	//		view, _ := template.ParseFiles("views/auth/login.html")
	//		_ = view.Execute(w, data)
	//	} else {
	//		data := make(map[string]interface{})
	//		data["user_exists"] = registerInput.Username + " ulanyjy ady öň hem ulanylýar!!!"
	//		view, _ := template.ParseFiles("views/auth/register.html")
	//		_ = view.Execute(w, data)
	//	}
	//
	//}

	translator := en.New()
	uni := ut.New(translator, translator)
	trans, _ := uni.GetTranslator("en")

	validate := validator.New()
	// register default translation (en)
	en_translations.RegisterDefaultTranslations(validate, trans)

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		labelName := field.Tag.Get("label")
		return labelName
	})
	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} meýdany hökmany", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	validateErrors := validate.Struct(registerInput)
	errorMessages := make(map[string]interface{})

	if validateErrors != nil {
		for _, e := range validateErrors.(validator.ValidationErrors) {
			errorMessages[e.StructField()] = e.Translate(trans)

		}
		//fmt.Println(errorMessages)
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
