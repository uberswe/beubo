package routes

import (
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/uberswe/beubo/pkg/structs"
	"github.com/uberswe/beubo/pkg/utility"
	"net/http"
)

// Login is the default login route
func (br *BeuboRouter) Login(w http.ResponseWriter, r *http.Request) {
	pageData := structs.PageData{
		Template: "login",
		Title:    "Login",
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// LoginPost handles authentication via post request and verifies a username/password via the database
func (br *BeuboRouter) LoginPost(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)

	invalidError := "Email or password incorrect, please try again or contact support"

	email := r.FormValue("email")
	password := r.FormValue("password")

	user := structs.AuthUser(br.DB, email, password)

	if user == nil || user.ID == 0 {
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, "/login", 302)
		return
	}

	ses := structs.CreateSession(br.DB, int(user.ID))

	if ses.ID == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 - Could not create session"))
		return
	}

	session.Set("SES_ID", ses.Token)
	http.Redirect(w, r, "/admin/", 302)
}

// Register renders the default registration page
func (br *BeuboRouter) Register(w http.ResponseWriter, r *http.Request) {
	setting := structs.FetchSettingByKey(br.DB, "enable_user_registration")
	if setting.ID == 0 || setting.Value == "false" {
		// Registration is not allowed
		http.Redirect(w, r, "/login", 302)
		return
	}

	pageData := structs.PageData{
		Template: "register",
		Title:    "Register",
	}

	br.Renderer.RenderHTMLPage(w, r, pageData)
}

// RegisterPost handles a registration request and inserts the user into the database
func (br *BeuboRouter) RegisterPost(w http.ResponseWriter, r *http.Request) {
	setting := structs.FetchSettingByKey(br.DB, "enable_user_registration")
	if setting.ID == 0 || setting.Value == "false" {
		// Registration is not allowed
		http.Redirect(w, r, "/login", 302)
		return
	}

	invalidError := "Please make sure the email is correct or that it does not already belong to a registered account"

	email := r.FormValue("email")
	password := r.FormValue("password")
	roles := []*structs.Role{}
	sites := []*structs.Site{}

	if !utility.IsEmailValid(email) {
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, "/login", 302)
		return
	}

	// TODO default role should be added on register
	if !structs.CreateUser(br.DB, email, password, roles, sites) {
		utility.SetFlash(w, "error", []byte(invalidError))
		http.Redirect(w, r, "/login", 302)
		return
	}

	utility.SetFlash(w, "message", []byte("Registration success, please check your email for further instructions"))

	http.Redirect(w, r, "/login", 302)
}

// Logout handles a GET logout request and removes the user session
func (br *BeuboRouter) Logout(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	session.Delete("SES_ID")
	session.Clear()
	http.Redirect(w, r, "/", 302)
}
