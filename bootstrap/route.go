package bootstrap

import (
    "goblog/pkg/routes"
    "github.com/gorilla/mux"
)

func SetupRoute() *mux.Router{
    router := mux.NewRouter()
    routes.RegisterWebRoutes(router)
    return router
}

