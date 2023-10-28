package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

//go:embed templates
var FilesFS embed.FS

var functions = template.FuncMap{
	"GetYear":       getYear,
	"FormatedDate":  FormatedDate,
	"FormatedDate2": FormatedDate2,
}

func getYear() int {
	year := time.Now().Year()
	return year
}

func FormatedDate(date time.Time) string {
	dateStr := date.Format("020106 15:04:05")

	return dateStr
}

func FormatedDate2(date time.Time) string {
	dateStr := date.Format("Jan 02")

	return dateStr
}

func addDefaultData(r *http.Request, templateData *TemplateData) *TemplateData {
	templateData.Url = r.URL.Path
	templateData.UserSession = getUserSession(r)

	return templateData
}

func (app *Application) RenderTemplate(w http.ResponseWriter, r *http.Request, templName string, templateData *TemplateData) {
	var templ *template.Template
	var err error

	_, ok := app.TemplatesCache[templName]
	if !ok || app.Config.Env == "dev" {
		templ, err = parseTemplate(templName, app.Config.Env)
		if err != nil {
			log.Fatalf("Error al parsear la plantilla: %s", err)
		}

		app.TemplatesCache[templName] = templ
	} else {
		templ = app.TemplatesCache[templName]
	}

	if templateData == nil {
		templateData = &TemplateData{}
	}

	templateData = addDefaultData(r, templateData)

	if err = templ.ExecuteTemplate(w, "base", templateData); err != nil {
		log.Fatalf("Error al executar el template: %s", err)
	}
}

func parseTemplate(templName, env string) (*template.Template, error) {
	templ := template.New("").Funcs(functions)

	if env == "dev" {
		return templ.ParseFiles(
			"templates/base_layout.gohtml",
			fmt.Sprintf("templates/%s", templName),
			"templates/navbar_layout.gohtml",
		)
	}

	return templ.ParseFS(
		FilesFS,
		"templates/base_layout.gohtml",
		fmt.Sprintf("templates/%s", templName),
		"templates/navbar_layout.gohtml",
	)
}
