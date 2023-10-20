package main

import "net/http"

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.HomeHandler)
	mux.HandleFunc("/login", app.LoginHandler)
	mux.Handle("/blog", app.AuthMiddleware(http.HandlerFunc(app.BlogHandler)))
	mux.HandleFunc("/about", app.AboutHandler)
	mux.HandleFunc("/logout", app.LogoutHandler)

	mux.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	return mux
}
