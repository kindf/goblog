package main

import (
    "goblog/bootstrap"
    "goblog/app/http/middlewares"
    "goblog/pkg/logger"
    "net/http"
    "github.com/gorilla/mux"
)

var router *mux.Router

func main() {
    bootstrap.SetupDB()
    router = bootstrap.SetupRoute()

    err := http.ListenAndServe(":8080", middlewares.RemoveTrailingSlash(router))
    logger.LogError(err)
}
