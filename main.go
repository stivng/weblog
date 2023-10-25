package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

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

	// Conexion a la base de datos postgreSQL
	var err error

	connStr := "user=alumno dbname=course-db password=123456 host=localhost port=5433 sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error conexion BDD: %s", err)
	}
	defer db.Close() // Con este defer toca tener cuidado cuando la conexion se hace en otro paquete que no sea main

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Conexion exitosa")

	log.Printf("server start in port: %s - mode: %s, version: %s", config.Port, config.Env, config.Version)
	if err = application.StartServer(); err != nil {
		log.Fatalf("Error al iniciar servidor: %s", err)
	}
}
