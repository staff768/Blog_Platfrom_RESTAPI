package models

import (
	"blogplatform/conf"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Post struct {
	Id       int `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Tags     string  `json:"tags"`
}

var Db *sql.DB

func Connect() error {
	var err error
	Db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", conf.Cfg.PgHost, conf.Cfg.PgPort, conf.Cfg.PgUser, conf.Cfg.PgPass, conf.Cfg.PgBase))
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе %s", err)
	}

	return nil
}

func Insert(title, content, category, tags string)(int,error){
	var id int
	err := Db.QueryRow("insert into post (title,content,category,tags) values ($1,$2,$3,$4) returning id", title,content,category,tags).Scan(&id)
	if err != nil {
		log.Fatalf("Error while trying to insert your post, %s", err)
		return 0, err
	}
		
	return id, nil
}

func GetById(id int)(*Post, error) {
	row := Db.QueryRow("select * from post where id = $1", id)
	p := &Post{}
	err := row.Scan(&p.Id,&p.Title, &p.Content, &p.Category, &p.Tags)
	if err != nil {
		if errors.Is(err,sql.ErrNoRows) {
			log.Printf("There is no such post in databse %s", err)
			return nil, sql.ErrNoRows
		} else {
			log.Fatalf("Something went wrong while trying to get post by id %s", err)
			return nil, err
		}
	}
	return p, nil
}
func GetAll() ([]*Post, error){
	rows, err := Db.Query("select * from post")
	if err != nil {
		log.Fatalf("Error while getting all post from databse %s", err)
	}
	defer rows.Close()
	posts := []*Post{}

	for rows.Next(){
		p := &Post{}
		err := rows.Scan(&p.Id, &p.Title, &p.Content, &p.Category, &p.Tags)
		if err !=nil {
			log.Fatalf("Error while scaning all posts %s", err)
			return nil, err
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err !=nil {
		log.Fatalf("Something went wrong %s", err) 
		return nil,err
	}
	return posts, nil
}
func Delete(id int) error {
	_, err := Db.Exec("delete from post where id = $1", id)
	if err != nil {
		if errors.Is(err,sql.ErrNoRows) {
			log.Fatalf("There is no such post in databse %s", err)
			return sql.ErrNoRows
		} else {
			log.Fatalf("Something went wrong while trying to delete post %s", err)
			return err
		}
	}
	return nil 	
}

func Update(id int, title,content,category, tags string) (error) {
	_, err := Db.Exec("update post set title = $1, content = $2, category = $3, tags = $4 where id = $5", title, content, category, tags, id)
	if err !=nil {
		log.Fatalf("Error while updating post %s", err)
		return  err
	}
	return nil
}

func GetRandomPost()(*Post,error){
	row := Db.QueryRow("SELECT * FROM post ORDER BY RANDOM() LIMIT 1")
	p := &Post{}
	err := row.Scan(&p.Id,&p.Title, &p.Content, &p.Category, &p.Tags)
	if err != nil {
		if errors.Is(err,sql.ErrNoRows) {
			log.Printf("There is no such post in databse %s", err)
			return nil, sql.ErrNoRows
		} else {
			log.Fatalf("Something went wrong while trying to get post by id %s", err)
			return nil, err
		}
	}
	return p, nil
}