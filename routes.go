package main

import "net/http"

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.HomeHandler)
	mux.HandleFunc("/login", app.LoginHandler)
	mux.HandleFunc("/register", app.RegisterHandler)
	mux.HandleFunc("/blog/view", app.HomeBlogViewHandler)
	mux.Handle("/blog", app.AuthMiddleware(http.HandlerFunc(app.BlogHandler)))
	mux.Handle("/blog/new-blog", app.AuthMiddleware(http.HandlerFunc(app.NewBlogHandler)))
	mux.Handle("/blog/edit", app.AuthMiddleware(http.HandlerFunc(app.EditBlogHandler)))
	mux.Handle("/blog/delete", app.AuthMiddleware(http.HandlerFunc(app.DeleteBlogHandler)))
	mux.HandleFunc("/about", app.AboutHandler)
	mux.HandleFunc("/logout", app.LogoutHandler)

	mux.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	return mux
}
