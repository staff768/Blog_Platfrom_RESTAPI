package main

import (
	"blogplatform/internal/models"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)


func (app *application) CreateNewPost(w http.ResponseWriter, r* http.Request) {
	
	var post struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		Category string `json:"category"`
		Tags     string `json:"tags"`
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&post)

	if err != nil {
		app.infoLog.Printf("Error while decode your post %s", err)
		return
	}

	id, err := models.Insert(post.Title, post.Content, post.Category, post.Tags)
	
	if err != nil {
		app.infoLog.Printf("Error while creating post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	w.WriteHeader(http.StatusCreated) 
	
	newpost, err:= models.GetById(id)
	if err != nil {
		app.infoLog.Printf("Error sending response: %v", err)
	}
	err = json.NewEncoder(w).Encode(newpost)
	if err != nil {
		app.infoLog.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	
}
func(app *application) GetPostById(w http.ResponseWriter, r* http.Request) {
	id,err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1{
		http.NotFound(w,r)
		return
	}
	post, err := models.GetById(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w,r)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			app.infoLog.Printf("Error getting post: %v", err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(post)
	if err != nil {
		app.infoLog.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}
func (app *application) GetAllPost(w http.ResponseWriter, r* http.Request) {
	posts, err := models.GetAll()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		app.infoLog.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
func (app *application) DeletePost(w http.ResponseWriter, r* http.Request) {
	id,err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1{
		http.NotFound(w,r)
		return
	}
	err = models.Delete(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w,r)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			app.infoLog.Printf("Error while deletting post: %v", err)
		}
		return
	}
}
func (app *application) UpdatePost(w http.ResponseWriter, r* http.Request) {
	id,err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1{
		http.NotFound(w,r)
		return
	}

	var post struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		Category string `json:"category"`
		Tags     string `json:"tags"`
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&post)

	if err != nil {
		app.infoLog.Printf("Error while decode your update post %s",err)
	}

	err = models.Update(id,post.Title,post.Content,post.Category,post.Tags)
	if err != nil {
		app.infoLog.Printf("Error while updating your post")
	}
	
	newpost, err:= models.GetById(id)
	if err != nil {
		app.infoLog.Printf("Error sending response: %v", err)
	}
	err = json.NewEncoder(w).Encode(newpost)
	if err != nil {
		app.infoLog.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}