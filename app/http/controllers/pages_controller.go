package controllers

import (
    "fmt"
    "net/http"
)

type PagesController struct {
}

func (*PagesController) Home(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<h1>Hello, welcome to goblog</h1>")
}

func (*PagesController) About(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<h1>Here is about page</h1>")
}

func (*PagesController) NotFound(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, "<h1>Page not found</h1>")
}

