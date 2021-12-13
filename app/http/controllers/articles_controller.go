package controllers

import (
    "net/http"
    "html/template"
    "fmt"
    "unicode/utf8"
    "strconv" // 字符串和其他类型转换
    "gorm.io/gorm"
    "goblog/pkg/types"
    "goblog/pkg/route"
    "goblog/app/models/article"
    "goblog/pkg/logger"
)

type ArticlesController struct {
}

type ArticlesFormData struct {
    Title       string
    Body        string
    URL         string
    Errors      map[string]string
}

func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
    id := route.GetRouteVariable("id", r)
    article, err := article.Get(id)

    if err != nil {
        if err == gorm.ErrRecordNotFound {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "404 article not found")
        } else {
            logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "500 internal error")
        }
    } else {
        tmpl, err := template.New("show.gohtml").Funcs(template.FuncMap{"RouteName2URL":route.Name2URL, "Uint64ToString":types.Uint64ToString,}).ParseFiles("static/show.gohtml")
        logger.LogError(err)
        err = tmpl.Execute(w, article)
        logger.LogError(err)
    }
}

func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
    articles, err := article.GetAll()
    
    if err != nil {
        logger.LogError(err)
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "500 服务器内部错误")
    } else {
        tmpl, err := template.ParseFiles("static/index.gohtml")
        logger.LogError(err)
        err = tmpl.Execute(w, articles)
        logger.LogError(err)
    }
}
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
        storeURL := route.Name2URL("articles.store")
        data := ArticlesFormData {
            Title: "",
            Body: "",
            URL: storeURL,
            Errors: nil,
        }
        tmpl, err := template.ParseFiles("static/create.gohtml")
        if err != nil {
            panic(err)
        }
        err = tmpl.Execute(w, data)
        if err != nil {
            panic(err)
        }
}

func validateArticleFormData(title string, body string) map[string]string {
    errors := make(map[string]string)
    // 验证标题
    if title == "" {
        errors["title"] = "标题不能为空"
    } else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
        errors["title"] = "标题长度需介于 3-40"
    }

    // 验证内容
    if body == "" {
        errors["body"] = "内容不能为空"
    } else if utf8.RuneCountInString(body) < 10 {
        errors["body"] = "内容长度需大于或等于 10 个字节"
    }

    return errors
}

func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {
    title := r.PostFormValue("title")
    body := r.PostFormValue("body")
    errors := validateArticleFormData(title, body)
    if len(errors) == 0 {
        _article := article.Article {
            Title: title,
            Body: body,
        }
        _article.Create()
        if _article.ID > 0 {
            fmt.Fprintf(w, "插入成功, ID为"+strconv.FormatUint(_article.ID, 10))
        } else {
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "创建文章失败，请联系管理员")
        }
    } else {
        storeURL := route.Name2URL("articles.store")
        data := ArticlesFormData {
            Title: title,
            Body: body,
            URL: storeURL,
            Errors: errors,
        }
        tmpl, err := template.ParseFiles("static/create.gohtml")
        logger.LogError(err)

        err = tmpl.Execute(w, data)
        logger.LogError(err)
    }
}

func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
    id := route.GetRouteVariable("id", r)

    _article, err := article.Get(id)

    if err != nil {
        if err == gorm.ErrRecordNotFound {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "404 article not found")
        } else {
            logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "500 internal error")
        }
    } else {
        updateURL := route.Name2URL("articles.update", "id", id)
        data := ArticlesFormData{
            Title: _article.Title,
            Body: _article.Body,
            URL: updateURL,
            Errors: nil,
        }
        tmpl, err := template.ParseFiles("static/edit.gohtml")
        logger.LogError(err)
        err = tmpl.Execute(w, data)
        logger.LogError(err)
    }
}

func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
    id := route.GetRouteVariable("id", r)
    _article, err := article.Get(id)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprint(w, "404 文章未找到")
        } else {
            logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "500 服务器内部错误")
        }
    } else {
        title := r.PostFormValue("title")
        body := r.PostFormValue("body")

        errors := validateArticleFormData(title, body)
        if len(errors) == 0 {
            _article.Title = title
            _article.Body = body

            rowsAffected, err := _article.Update()

            if err != nil {
                logger.LogError(err)
                w.WriteHeader(http.StatusInternalServerError)
                fmt.Fprint(w, "500 服务器内部错误")
            }
            if rowsAffected > 0 {
                showURL := route.Name2URL("articles.show", "id", id)
                http.Redirect(w, r, showURL, http.StatusFound)
            } else {
                fmt.Fprint(w, "您没有做任何更改！")
            }
        } else {
            updateURL := route.Name2URL("articles.update", "id", id)
            data := ArticlesFormData{
                Title:  title,
                Body:   body,
                URL:    updateURL,
                Errors: errors,
            }
            tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
            logger.LogError(err)

            err = tmpl.Execute(w, data)
            logger.LogError(err)
        }
    }
}

func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
    id := route.GetRouteVariable("id", r)
    _article, err := article.Get(id)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "404 文章未找到")
        } else {
            logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "500 服务器内部错误")
        }
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
