package route

import (
    "github.com/gorilla/mux"
    "net/http"
    "goblog/pkg/logger"
)

var route *mux.Router

func Name2URL(routeName string, pairs ...string) string {
    url, err := route.Get(routeName).URL(pairs...)
    if err != nil {
        logger.LogError(err)
        return ""
    }

    return url.String()
}

func GetRouteVariable(param string, r *http.Request) string {
    vars := mux.Vars(r)
    return vars[param]
}

func SetRoute(r *mux.Router) {
    route = r
}
