package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TemplateData struct {
	Url         string
	UserSession string
	Errors      []string
	Blogs       []Blog
	Blog
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
		Expires:  time.Now().Add(12 * time.Minute),
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
	blogs, err := GetBlogs()
	if err != nil {
		log.Println(err)
	}

	app.RenderTemplate(w, r, "index_page.gohtml", &TemplateData{
		Blogs: blogs,
	})
}

func (app *Application) HomeBlogViewHandler(w http.ResponseWriter, r *http.Request) {
	blogId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(blogId)
	if err != nil {
		log.Println(err)
	}

	blog, err := GetBlogById(id)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	app.RenderTemplate(w, r, "blog_view_page.gohtml", &TemplateData{
		Blog: *blog,
	})
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

		// Busca usuario por el email en la BDD.
		dataUserDB, err := GetUserByEmail(user.Email)
		if err != nil {
			log.Printf("Error al extraer el usuario de la base de datos: %s", err)
		}

		if user.Email == dataUserDB.Email && user.Password == dataUserDB.Password {
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

func (app *Application) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.RenderTemplate(w, r, "register_page.gohtml", nil)
	} else if r.Method == http.MethodPost {
		var user User
		var resp Resp

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Println(err)
			return
		}

		resp.Success = true

		_, err := db.Exec("INSERT INTO data_users (email, password) VALUES ($1, $2)", user.Email, user.Password)
		if err != nil {
			log.Println("No fue posible insertar los datos")
			resp.Success = false
			json.NewEncoder(w).Encode(resp)
			return
		}

		resp.Success = true
		setUserSession(w, user.Email)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Println(err)
			return
		}

		// setUserSession(w, user.Email) Esto no es posible hacerlo debido a que en el encoder ya se envio el encabezado
		// y no se puede enviar otro encabezado despues.
	}
}

func (app *Application) BlogHandler(w http.ResponseWriter, r *http.Request) {
	blogs, err := GetBlogs()
	if err != nil {
		log.Println(err)
		return
	}

	errorMessage := r.URL.Query().Get("error") // si no encuentra nada retorna vacio, se podria tratar sie s vacio para mostrarle que el usuario fue elmininado correctamente o algo asi

	app.RenderTemplate(w, r, "blog_page.gohtml", &TemplateData{
		Blogs:  blogs,
		Errors: []string{errorMessage},
	})
}

func slugify(value string) string {
	value = strings.ToLower(value)
	reg := regexp.MustCompile("[^a-z0-9]+")
	value = reg.ReplaceAllString(value, "-")
	value = strings.Trim(value, "-")

	return value
}

func (app *Application) NewBlogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.RenderTemplate(w, r, "blog_new_page.gohtml", nil)
	} else if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		errors := make([]string, 0)

		if len(strings.Trim(title, " ")) == 0 {
			errors = append(errors, "Titulo requerido!")
		}

		if len(strings.Trim(content, " ")) == 0 {
			errors = append(errors, "Contenido requerido!")
		}

		email := getUserSession(r)
		dataUserDB, err := GetUserByEmail(email)
		if err != nil {
			errors = append(errors, "No estas logueado")
		}

		blog := Blog{
			Title:   title,
			Content: template.HTML(content),
			Slug:    slugify(title),
			Author:  dataUserDB.Id,
		}

		if err = CreateBlog(blog); err != nil {
			errors = append(errors, err.Error())
		}

		log.Println(blog)

		if len(errors) > 0 {
			app.RenderTemplate(w, r, "blog_new_page.gohtml", &TemplateData{
				Errors: errors,
			})
			return
		}

		http.Redirect(w, r, "/blog", http.StatusSeeOther)
	}
}

func (app *Application) EditBlogHandler(w http.ResponseWriter, r *http.Request) {
	blogId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(blogId)
	if err != nil {
		return
	}

	if r.Method == http.MethodGet {
		blog, err := GetBlogById(id)
		if err != nil {
			log.Println("error al extraer blog")
			return
		}

		app.RenderTemplate(w, r, "blog_edit.gohtml", &TemplateData{
			Blog: *blog,
		})
	} else if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		// id := r.FormValue("id") tambien funciona para extraer de el id del atributo value de un input hidden. {{.Blog.Id}} que este en el form.

		// Actualizar la base de datos
		blog := Blog{
			Id:      id,
			Title:   title,
			Content: template.HTML(content),
			Slug:    slugify(title),
		}

		if err := UpdateBlog(blog); err != nil {
			log.Println("error ala actualizar blog")
			return
		}

		http.Redirect(w, r, "/blog", http.StatusSeeOther)
		return
	}
}

func (app *Application) DeleteBlogHandler(w http.ResponseWriter, r *http.Request) {
	blogId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(blogId)
	if err != nil {
		return
	}

	rows, err := DeleteBlog(id)
	if err != nil {
		log.Printf("error al borrar el blog. %s", err)
		return
	}

	if rows < 1 {
		http.Redirect(w, r, fmt.Sprintf("/blog?error=%s", "id no valido"), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/blog", http.StatusSeeOther)
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
