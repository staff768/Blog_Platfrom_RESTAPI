package main
import("net/http")

func (app *application) routes() *http.ServeMux {
    mux := http.NewServeMux()
    
    mux.HandleFunc("POST /register", app.Register)
    mux.HandleFunc("POST /login", app.Login)
    
    protected := http.NewServeMux()
    protected.HandleFunc("POST /posts", app.CreateNewPost)
    protected.HandleFunc("GET /posts/id", app.GetPostById)
    protected.HandleFunc("GET /posts", app.GetAllPost)
    protected.HandleFunc("DELETE /posts/id", app.DeletePost)
    protected.HandleFunc("PUT /posts/id", app.UpdatePost)
    
    protectedWithAuth := app.AuthMiddleware(protected)
    mux.Handle("/", protectedWithAuth)
    
    return mux
}