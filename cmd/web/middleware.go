package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"blogplatform/conf"
)

type contextKey string
const userIDKey contextKey = "userID"

func (app *application) AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Логируем входящий запрос
        app.infoLog.Printf("Incoming request: %s %s", r.Method, r.URL.Path)
        
        authHeader := r.Header.Get("Authorization")
        app.infoLog.Printf("Authorization header: '%s'", authHeader) // Логируем заголовок
        
        if authHeader == "" {
            app.errorLog.Println("Authorization header is missing")
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        const prefix = "Bearer "
        if !strings.HasPrefix(authHeader, prefix) {
            app.errorLog.Printf("Invalid Authorization format: '%s'", authHeader)
            http.Error(w, "Bearer token required", http.StatusUnauthorized)
            return
        }
        
        tokenString := strings.TrimPrefix(authHeader, prefix)
        app.infoLog.Printf("Token string: '%s'", tokenString) // Логируем токен

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(conf.Cfg.JwtSecret), nil
        })

        if err != nil {
            app.errorLog.Printf("Token validation error: %v", err)
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        if !token.Valid {
            app.errorLog.Println("Token is invalid")
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            app.errorLog.Println("Invalid token claims")
            http.Error(w, "Invalid token claims", http.StatusUnauthorized)
            return
        }

        userID, ok := claims["sub"].(float64)
        if !ok {
            app.errorLog.Printf("Invalid user ID type: %T", claims["sub"])
            http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), userIDKey, int(userID))
        app.infoLog.Printf("Authentication successful for user ID: %d", int(userID))
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}