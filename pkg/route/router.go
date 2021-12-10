package route

import (
    "github.com/gorilla/mux"
    "net/http"
    "goblog/pkg/routes"
)

var Router *mux.Router

func Initialize() {
    Router = mux.NewRouter()
    routes.RegisterWebRoutes(Router)
}

func RouteName2URL(routeName string, pairs ...string) string {
    url, err := Router.Get(routeName).URL(pairs...)
    if err != nil {
        // checkError(err)
        return ""
    }

    return url.String()
}

func GetRouteVariable(param string, r *http.Request) string {
    vars := mux.Vars(r)
    return vars[param]
}


