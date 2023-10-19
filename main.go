package main

import (
	"flag"
	"html/template"
	"log"
)

func main() {
	templatesCache := make(map[string]*template.Template)

	config := Config{
		Version: "1.0.0",
	}

	flag.StringVar(&config.Port, "port", "8080", "Server port")
	flag.StringVar(&config.Env, "env", "dev", "Application environment")
	flag.Parse()

	application := Application{
		Config:         config,
		TemplatesCache: templatesCache,
	}

	log.Printf("server start in port: %s - mode: %s, version: %s", config.Port, config.Env, config.Version)
	if err := application.StartServer(); err != nil {
		log.Fatalf("Error al iniciar servidor: %s", err)
	}
}
