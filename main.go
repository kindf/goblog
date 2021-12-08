package main

import (
    "net/http"
    "fmt"
)

func handlerfunc(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "<h1>Hello, here is goblog</h1>")
}

func main() {
    http.HandleFunc("/", handlerfunc)
    http.ListenAndServe(":8080", nil)
}
