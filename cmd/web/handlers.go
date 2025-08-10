package main

import (
	"blogplatform/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"blogplatform/conf"
	"time"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
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

func (app *application) Register(w http.ResponseWriter, r *http.Request){
	var input_user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&input_user)
	if err != nil {
		app.errorLog.Printf("Error decoding request: %v", err)
	    http.Error(w, "Invalid request body", http.StatusBadRequest)
    	return
	}
	if input_user.Email == "" || input_user.Password == "" {
    	http.Error(w, "Email and password are required", http.StatusBadRequest)
    	return
	}


	if !strings.Contains(input_user.Email, "@") {
    	http.Error(w, "Invalid email format", http.StatusBadRequest)
    	return
	}


	if len(input_user.Password) < 8 {
    	http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
    	return
	}
	
	if err := models.CreateUser(input_user.Email, input_user.Password); err != nil {
	if errors.Is(err, models.ErrDuplicateEmail) {
        http.Error(w, "Email already exists", http.StatusConflict)
    } else {
        app.errorLog.Printf("User creation failed: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
    }
    	return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
    	"status": "success",
    	"message": "User registered successfully",
	})
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    
    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()
    if err := dec.Decode(&input); err != nil {
        app.errorLog.Printf("JSON decode error: %v", err)
        http.Error(w, "Invalid request format", http.StatusBadRequest)
        return
    }

    if input.Email == "" || input.Password == "" {
        http.Error(w, "Email and password are required", http.StatusBadRequest)
        return
    }

    user, err := models.AuthenticateUser(input.Email, input.Password)
    if err != nil {
        if errors.Is(err, models.ErrInvalidCredentials) {
            http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        } else {
            app.errorLog.Printf("Authentication error: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

    
    expirationDuration, err := time.ParseDuration(conf.Cfg.JwtExpiration)
    if err != nil {
        app.errorLog.Printf("Invalid JWT expiration format: %v", err)
        http.Error(w, "Server configuration error", http.StatusInternalServerError)
        return
    }
	
	if expirationDuration <= 0 {
        app.errorLog.Printf("Invalid JWT expiration duration: %v", expirationDuration)
        http.Error(w, "Server configuration error", http.StatusInternalServerError)
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": user.ID, 
        "exp": time.Now().Add(expirationDuration).Unix(), 
        "iat": time.Now().Unix(),                         
    })

    
    tokenString, err := token.SignedString([]byte(conf.Cfg.JwtSecret))
    if err != nil {
        app.errorLog.Printf("Token signing error: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "token":   tokenString,
        "expires": expirationDuration.String(),
        "user_id": strconv.Itoa(user.ID),
    })
}