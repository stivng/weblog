package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// Config contiene la configuracion basica del servidor.
type Config struct {
	Port    string
	Env     string
	Version string
}

// Application contiene las configuraciones generales que son necesarias en la aplicacion.
type Application struct {
	Config
	TemplatesCache map[string]*template.Template
}

// StartServer inicia el servidor.
func (app *Application) StartServer() error {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", app.Config.Port),
	}

	return srv.ListenAndServe()
}
