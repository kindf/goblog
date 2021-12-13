package main

import (
    "goblog/pkg/logger"
    "goblog/pkg/database"
    "goblog/bootstrap"
    "net/http"
    "fmt"
    "strings" // 字符串操作
    "strconv" // 字符串和其他类型转换
    "html/template"
    "unicode/utf8"
    "net/url"
    "database/sql"
    "github.com/gorilla/mux"
)

var router *mux.Router
var db *sql.DB

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

func getArticleByID(id string) (Article, error) {
    article := Article{}
    query := "SELECT * FROM articles WHERE id = ?"
    err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
    return article, err
}

type Article struct {
    Title string
    Body string
    ID int64
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT * from articles")
    logger.LogError(err)
    defer rows.Close()

    var articles []Article
    for rows.Next() {
        var article Article
        err := rows.Scan(&article.ID, &article.Title, &article.Body)
        logger.LogError(err)
        articles = append(articles, article)
    }
    err = rows.Err()
    logger.LogError(err)
    tmpl, err := template.ParseFiles("static/index.gohtml")
    logger.LogError(err)
    err = tmpl.Execute(w, articles)
    logger.LogError(err)
}

// ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
    Title       string
    Body        string
    URL         *url.URL
    Errors      map[string]string
}

func getRouteVariable(param string, r *http.Request) string {
    vars := mux.Vars(r)
    return vars[param]
}


func forceHTMLMiddleware(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        h.ServeHTTP(w, r)
    })
}

func (a Article) Delete() (RowsAffected int64, err error) {
    rs, err := db.Exec("DELETE FROM articles WHERE id = " + strconv.FormatInt(a.ID, 10))
    if err != nil {
        return 0, err
    }

    if n, _ := rs.RowsAffected(); n > 0 {
        return n, nil
    }
    return 0, nil
}

func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {
    id := getRouteVariable("id", r)
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
        RowsAffected, err := article.Delete()
        if err != nil {
            logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "500 internal error")
        } else {
            if RowsAffected > 0 {
                indexURL, _ := router.Get("articles.index").URL()
                http.Redirect(w, r, indexURL.String(), http.StatusFound)
            } else {
                w.WriteHeader(http.StatusNotFound)
                fmt.Fprintf(w, "404 article not found")
            }
        }
    }
}

func removeTrailingSlash(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
        }
        next.ServeHTTP(w, r)
    })
}

func main() {
    database.Initialize()
    db = database.DB
    bootstrap.SetupDB()
    router = bootstrap.SetupRoute()

    router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")

    router.Use(forceHTMLMiddleware)

    err := http.ListenAndServe(":8080", removeTrailingSlash(router))
    if err != nil {
        fmt.Println(err)
    }
}
