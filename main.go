package main

import (
    "net/http"
    "fmt"
)

func handlerfunc(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    if r.URL.Path == "/" {
        fmt.Fprint(w, "<h1>Hello, here is goblog</h1>")
    } else if r.URL.Path == "/about" {
        fmt.Fprint(w, "<h1>about 页面</h1>")
    } else {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprint(w, "<h1>请求页面不存在</h1>")
    }
}

func main() {
    http.HandleFunc("/", handlerfunc)
    http.ListenAndServe(":8080", nil)
}
