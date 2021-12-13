package bootstrap

import (
    "goblog/pkg/routes"
    "goblog/pkg/route"
    "github.com/gorilla/mux"
)

func SetupRoute() *mux.Router{
    router := mux.NewRouter()
    routes.RegisterWebRoutes(router)
    route.SetRoute(router)
    return router
}

