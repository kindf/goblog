package controllers

import (
    "net/http"
    "html/template"
    "goblog/pkg/logger"
    "database/sql"
    "fmt"
    "goblog/pkg/types"
    "goblog/pkg/route"
    "goblog/app/models/article"
)

type ArticlesController struct {
}

func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
    id := route.GetRouteVariable("id", r)
    article, err := article.Get(id)

    if err != nil {
        if err == sql.ErrNoRows {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "404 article not found")
        } else {
            logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "500 internal error")
        }
    } else {
        tmpl, err := template.New("show.gohtml").Funcs(template.FuncMap{"RouteName2URL":route.RouteName2URL, "Uint64ToString":types.Uint64ToString,}).ParseFiles("static/show.gohtml")
        logger.LogError(err)
        err = tmpl.Execute(w, article)
        logger.LogError(err)
    }
}


