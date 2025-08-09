package main
import("net/http")

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /posts", app.CreateNewPost)
	mux.HandleFunc("GET /posts/id", app.GetPostById)
	mux.HandleFunc("GET /posts", app.GetAllPost)
	mux.HandleFunc("DELETE /posts/id", app.DeletePost)
	mux.HandleFunc("PUT /posts/id", app.UpdatePost)

	return mux
}