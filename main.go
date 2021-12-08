package main

import (
    "net/http"
    "fmt"
    "github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprintf(w, "<h1>Hello, welcome to goblog</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprintf(w, "<h1>Here is about page</h1>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, "<h1>Page not found</h1>")
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    fmt.Fprintf(w, "article ID: "+id)
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "article list")
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "create article")
}

func main() {
    router := mux.NewRouter()

    router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
    router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

    router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
    router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
    router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")

    router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

    http.ListenAndServe(":8080", router)
}
