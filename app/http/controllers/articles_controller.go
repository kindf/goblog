package controllers

import (
    "net/http"
    "html/template"
    "goblog/pkg/logger"
    "database/sql"
    "fmt"
    "goblog/pkg/types"
    "goblog/pkg/route"
)

type ArticlesController struct {
}

func getArticleByID(id string) (int, error) {
    if id == "todo" {
        return 1, sql.ErrNoRows
    }
    return 2, sql.ErrNoRows
}

func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
    id := route.GetRouteVariable("id", r)
    // TODO:
    article, err := getArticleByID(id)

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
        tmpl, err := template.New("show.gohtml").Funcs(template.FuncMap{"RouteName2URL":route.RouteName2URL, "Int64ToString":types.Int64ToString,}).ParseFiles("static/show.gohtml")
        logger.LogError(err)
        err = tmpl.Execute(w, article)
        logger.LogError(err)
    }
}


