package main
import("net/http")

func routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /posts", CreateNewPost)
	mux.HandleFunc("GET /posts/id", GetPostById)
	mux.HandleFunc("GET /posts", GetAllPost)
	mux.HandleFunc("DELETE /posts/id", DeletePost)
	mux.HandleFunc("PUT /posts/id", UpdatePost)

	return mux
}