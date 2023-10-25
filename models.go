package main

import (
	"fmt"
	"time"
)

type DataUserDB struct {
	Id       int
	Email    string
	Password string
}

type Blog struct {
	Id        int
	Title     string
	Content   string
	Slug      string
	Author    int
	CreatedAt time.Time
	UpdateAt  time.Time
}

func GetUserByEmail(email string) (*DataUserDB, error) {
	var dataUserDB DataUserDB
	row := db.QueryRow("SELECT id, email, password FROM data_users WHERE email = $1", email)

	if err := row.Scan(&dataUserDB.Id, &dataUserDB.Email, &dataUserDB.Password); err != nil {
		fmt.Println("Error al extraer usuario")
		return nil, err
	}

	return &dataUserDB, nil
}

func CreateBlog(blog Blog) error {
	_, err := db.Exec("INSERT INTO blogs (title, slug, content, user_id) VALUES ($1, $2, $3, $4)", blog.Title, blog.Slug, blog.Content, blog.Author)
	if err != nil {
		return err
	}

	return nil
}

func GetBlogs() ([]Blog, error) {
	var blogs []Blog

	rows, err := db.Query("SELECT id, title, content, slug, user_id, created_at, update_at FROM blogs")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var blog Blog

		if err = rows.Scan(
			&blog.Id,
			&blog.Title,
			&blog.Content,
			&blog.Slug,
			&blog.Author,
			&blog.CreatedAt,
			&blog.UpdateAt,
		); err != nil {
			return nil, err
		}

		blogs = append(blogs, blog)
	}

	return blogs, nil
}
