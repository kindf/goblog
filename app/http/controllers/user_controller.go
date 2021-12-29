package controllers

import (
    "net/http"
    "fmt"
    "gorm.io/gorm"
    "goblog/pkg/route"
    "goblog/app/models/article"
    "goblog/pkg/logger"
    "goblog/pkg/view"
    "goblog/app/models/user"
)

type UserController struct {
}


func (*UserController) Show(w http.ResponseWriter, r *http.Request) {
    id := route.GetRouteVariable("id", r)
    _user, err := user.Get(id)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprint(w, "404 文章未找到")
        } else {
            logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "500 服务器内部错误")
        }
    } else {
        articles, err := article.GetByUserID(_user.GetStringID())
        if err != nil {
            logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "500 服务器内部错误")
        } else {
            view.Render(w, view.D{
                "Articles": articles,
            }, "index", "_article_meta")
        }
    }
}
