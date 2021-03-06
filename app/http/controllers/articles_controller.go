package controllers

import (
    "net/http"
    "fmt"
    "time"
    // "path/filepath"
    // "gorm.io/gorm"
    "goblog/pkg/route"
    "goblog/pkg/auth"
    "goblog/app/models/article"
    "goblog/pkg/logger"
    "goblog/pkg/view"
    "goblog/app/requests"
    "goblog/app/policies"
    // "goblog/pkg/flash"
)

type ArticlesController struct {
    BaseController
}

func (ac *ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
    id := route.GetRouteVariable("id", r)
    article, err := article.Get(id)

    if err != nil {
        ac.ResponseForSQLError(w, err)
    } else {
        view.Render(w, view.D{
            "Article": article,
            "CanModifyArticle": policies.CanModifyArticle(article),
        }, "show", "_article_meta")
    }
}

func (ac *ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
    articles, err := article.GetAll()
    
    if err != nil {
        ac.ResponseForSQLError(w, err)
    } else {
        view.Render(w, view.D{
            "Articles": articles,
        }, "index", "_article_meta")
    }
}
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
    view.Render(w, view.D{}, "create", "_form_field")
}

func (ac *ArticlesController) Store(w http.ResponseWriter, r *http.Request) {
    currentUser := auth.User()
    _article := article.Article{
        Title: r.PostFormValue("title"),
        Body: r.PostFormValue("body"),
        UserID: currentUser.ID,
    }
    errors := requests.ValidateArticleForm(_article)
    if len(errors) == 0 {
        _article.CreateAt = time.Now()
        _article.UpdateAt = time.Now()
        _article.Create()
        if _article.ID > 0 {
            indexURL := route.Name2URL("articles.show", "id", _article.GetStringID())
            http.Redirect(w, r, indexURL, http.StatusFound)
        } else {
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "创建文章失败，请联系管理员")
        }
    } else {
        view.Render(w, view.D{
            "Errors": errors,
        }, "create", "_form_field")
    }
}

func (ac *ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {

    id := route.GetRouteVariable("id", r)
    _article, err := article.Get(id)

    if err != nil {
        ac.ResponseForSQLError(w, err)
    } else {
        if !policies.CanModifyArticle(_article) {
            ac.ResponseForUnauthorized(w, r)
        } else {
            view.Render(w, view.D{
                "Article": _article,
                "Errors": view.D{},
            }, "edit", "_form_field")
        }
    }
}

func (ac *ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
    id := route.GetRouteVariable("id", r)
    _article, err := article.Get(id)
    if err != nil {
        ac.ResponseForSQLError(w, err)
    } else {
        if !policies.CanModifyArticle(_article) {
            ac.ResponseForUnauthorized(w, r)
        } else {
            _article.Title = r.PostFormValue("title")
            _article.Body = r.PostFormValue("body")

            errors := requests.ValidateArticleForm(_article)
            if len(errors) == 0 {

                rowsAffected, err := _article.Update()

                if err != nil {
                    w.WriteHeader(http.StatusInternalServerError)
                    fmt.Fprint(w, "500 服务器内部错误")
                    return
                }
                if rowsAffected > 0 {
                    showURL := route.Name2URL("articles.show", "id", id)
                    http.Redirect(w, r, showURL, http.StatusFound)
                } else {
                    fmt.Fprint(w, "您没有做任何更改！")
                }
            } else {
                view.Render(w, view.D{
                    "Article": _article,
                    "Errors": errors,
                }, "edit")
            }
        }
    }
}

func (ac *ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
    id := route.GetRouteVariable("id", r)
    _article, err := article.Get(id)
    if err != nil {
        ac.ResponseForSQLError(w, err)
    } else {
        if !policies.CanModifyArticle(_article) {
            ac.ResponseForUnauthorized(w, r)
        } else {
            rowsAffected, err := _article.Delete()
            if err != nil {
                logger.LogError(err)
                w.WriteHeader(http.StatusInternalServerError)
                fmt.Fprintf(w, "500 服务器内部错误")
            } else {
                if rowsAffected > 0 {
                    indexURL := route.Name2URL("articles.index")
                    http.Redirect(w, r, indexURL, http.StatusFound)
                } else {
                    w.WriteHeader(http.StatusNotFound)
                    fmt.Fprintf(w, "404 文章未找到")
                }
            }
        }
    }
}
