package main

import (
    "goblog/bootstrap"
    "goblog/app/http/middlewares"
    "goblog/pkg/logger"
    c "goblog/pkg/config"
    "goblog/config"
    "net/http"
    "github.com/gorilla/mux"
)

var router *mux.Router

func init(){
    config.Initialize()
}

func main() {
    bootstrap.SetupDB()
    router = bootstrap.SetupRoute()

    err := http.ListenAndServe(":"+c.GetString("app.port"), middlewares.RemoveTrailingSlash(router))
    logger.LogError(err)
}
