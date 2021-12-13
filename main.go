package main

import (
    "goblog/pkg/database"
    "goblog/bootstrap"
    "net/http"
    "fmt"
    "strings" // 字符串操作
    "database/sql"
    "github.com/gorilla/mux"
)

var router *mux.Router
var db *sql.DB

func forceHTMLMiddleware(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        h.ServeHTTP(w, r)
    })
}

func removeTrailingSlash(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
        }
        next.ServeHTTP(w, r)
    })
}

func main() {
    database.Initialize()
    db = database.DB
    bootstrap.SetupDB()
    router = bootstrap.SetupRoute()
    router.Use(forceHTMLMiddleware)
    err := http.ListenAndServe(":8080", removeTrailingSlash(router))
    if err != nil {
        fmt.Println(err)
    }
}
