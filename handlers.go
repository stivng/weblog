package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type TemplateData struct {
	Url         string
	UserSession string
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Resp struct {
	Success bool `json:"success"`
}

func setUserSession(w http.ResponseWriter, email string) {
	session := http.Cookie{
		Name:     "session",
		Value:    email,
		Path:     "/",
		Expires:  time.Now().Add(2 * time.Minute),
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &session)
}

func getUserSession(r *http.Request) string {
	session, err := r.Cookie("session")
	if err != nil {
		return ""
	}

	email := session.Value
	return email
}

func (app *Application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, r, "index_page.gohtml", nil)
}

func (app *Application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.RenderTemplate(w, r, "login_page.gohtml", nil)
	} else if r.Method == http.MethodPost {
		var user User
		var resp Resp

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Println(err)
			return
		}

		if user.Email == "majestiic-10@hotmail.com" && user.Password == "12345" {
			resp.Success = true
			setUserSession(w, user.Email)

			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Println(err)
				return
			}
		} else {
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (app *Application) BlogHandler(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, r, "blog_page.gohtml", nil)
}

func (app *Application) AboutHandler(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, r, "about_page.gohtml", nil)
}

func (app *Application) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session := http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, &session)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return
}
