package models

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/lib/pq"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}


func CreateUser(email, password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}
	
	_, err = Db.Exec("INSERT INTO users (email, password_hash) VALUES ($1, $2)", 
		email, hashedPassword)
		if pgErr, ok := err.(*pq.Error); ok {
        if pgErr.Code == "23505" {
            return ErrDuplicateEmail
        }
    }
	return err
}

func AuthenticateUser(email, password string) (*User, error) {
    user := &User{}
    row := Db.QueryRow("SELECT id, email, password_hash FROM users WHERE email = $1", email)
    
    if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrInvalidCredentials
        }
        return nil, err
    }
    
    if !checkPasswordHash(password, user.PasswordHash) {
        return nil, ErrInvalidCredentials
    }
    
    return user, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}


func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}