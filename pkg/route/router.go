package route

import (
    "github.com/gorilla/mux"
    "net/http"
)

func RouteName2URL(routeName string, pairs ...string) string {
    var route *mux.Router
    url, err := route.Get(routeName).URL(pairs...)
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


