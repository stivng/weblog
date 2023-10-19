package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type TemplateData struct {
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Resp struct {
	Success bool `json:"success"`
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

func (app *Application) AboutHandler(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, r, "about_page.gohtml", nil)
}
