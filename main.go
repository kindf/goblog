package main

import (
    "goblog/pkg/route"
    "net/http"
    "fmt"
    "strings" // 字符串操作
    "html/template"
    "unicode/utf8"
    "net/url"
    "database/sql"
    "log"
    "time"
    "strconv" // 字符串和其他类型转换
    "github.com/gorilla/mux"
    "github.com/go-sql-driver/mysql"
)

var router *mux.Router
var db *sql.DB

func checkError(err error){
    if err != nil {
        log.Fatal(err)
    }
}

func initDB(){
    var err error
    config := mysql.Config {
        User: "root",
        Passwd: "ljl123456",
        Addr: "127.0.0.1:3306",
        Net: "tcp",
        DBName: "goblog",
        AllowNativePasswords: true,
    }
    db, err = sql.Open("mysql", config.FormatDSN())
    checkError(err)
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)

    err = db.Ping()
    checkError(err)
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

func getRouteVariable(param string, r *http.Request) string {
    vars := mux.Vars(r)
    return vars[param]
}

func getArticleByID(id string) (Article, error) {
    article := Article{}
    query := "SELECT * FROM articles WHERE id = ?"
    err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
    return article, err
}

func createTables() {
    createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
        id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
        title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
        body longtext COLLATE utf8mb4_unicode_ci
    );`

    _, err := db.Exec(createArticlesSQL)
    checkError(err)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<h1>Hello, welcome to goblog</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<h1>Here is about page</h1>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, "<h1>Page not found</h1>")
}

type Article struct {
    Title string
    Body string
    ID int64
}

func Int64ToString(num int64) string {
    return strconv.FormatInt(num, 10)
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
    id := getRouteVariable("id", r)
    article, err := getArticleByID(id)

    if err != nil {
        if err == sql.ErrNoRows {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "404 article not found")
        } else {
            checkError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "500 internal error")
        }
    } else {
        tmpl, err := template.New("show.gohtml").Funcs(template.FuncMap{"RouteName2URL":route.RouteName2URL, "Int64ToString":Int64ToString,}).ParseFiles("static/show.gohtml")
        checkError(err)
        err = tmpl.Execute(w, article)
        checkError(err)
    }
}

func (a Article) Link() string {
    showURL, err := router.Get("articles.show").URL("id", strconv.FormatInt(a.ID, 10))
    if err != nil {
        checkError(err)
        return ""
    }
    return showURL.String()
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT * from articles")
    checkError(err)
    defer rows.Close()

    var articles []Article
    for rows.Next() {
        var article Article
        err := rows.Scan(&article.ID, &article.Title, &article.Body)
        checkError(err)
        articles = append(articles, article)
    }
    err = rows.Err()
    checkError(err)
    tmpl, err := template.ParseFiles("static/index.gohtml")
    checkError(err)
    err = tmpl.Execute(w, articles)
    checkError(err)
}

// ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
    Title       string
    Body        string
    URL         *url.URL
    Errors      map[string]string
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
    title := r.PostFormValue("title")
    body := r.PostFormValue("body")
    errors := validateArticleFormData(title, body)
    if len(errors) == 0 {
        lastInsertId, err := saveArticlesToDB(title, body)
        if lastInsertId > 0 {
            fmt.Fprintf(w, "insert succ. ID:"+strconv.FormatInt(lastInsertId, 10))
        }else {
            checkError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "500 internal error")
        }
    } else {
        storeURL, _ := router.Get("articles.store").URL()
        data := ArticlesFormData {
            Title: title,
            Body: body,
            URL: storeURL,
            Errors: errors,
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
}

func saveArticlesToDB(title string, body string) (int64, error) {
    var (
        id int64
        err error
        rs sql.Result
        stmt *sql.Stmt
    )
    stmt, err = db.Prepare("INSERT INTO articles (title, body) VALUES(?,?)")
    if err != nil {
        return 0, err
    }
    defer stmt.Close()
    rs, err = stmt.Exec(title, body)
    if err != nil {
        return 0, err
    }
    if id, err = rs.LastInsertId(); id > 0 {
        return id, nil
    }

    return 0, err
}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
    storeURL, _ := router.Get("articles.store").URL()
    data := ArticlesFormData{
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

func articlesEditHandler(w http.ResponseWriter, r *http.Request) {
    id := getRouteVariable("id", r)
    article, err := getArticleByID(id)

    if err != nil {
        if err == sql.ErrNoRows {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "404 article not found")
        } else {
            checkError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "500 internal error")
        }
    } else {
        updateURL, _ := router.Get("articles.update").URL("id", id)
        data := ArticlesFormData{
            Title: article.Title,
            Body: article.Body,
            URL: updateURL,
            Errors: nil,
        }
        tmpl, err := template.ParseFiles("static/edit.gohtml")
        checkError(err)
        err = tmpl.Execute(w, data)
        checkError(err)
    }
}

func articlesUpdateHandler(w http.ResponseWriter, r *http.Request) {
    id := getRouteVariable("id", r)
    _, err := getArticleByID(id)
    if err != nil {
        if err == sql.ErrNoRows {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprint(w, "404 文章未找到")
        } else {
            checkError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "500 服务器内部错误")
        }
    } else {
        title := r.PostFormValue("title")
        body := r.PostFormValue("body")

        errors := validateArticleFormData(title, body)
        if len(errors) == 0 {
            query := "UPDATE articles SET title = ?, body = ? WHERE id = ?"
            rs, err := db.Exec(query, title, body, id)

            if err != nil {
                checkError(err)
                w.WriteHeader(http.StatusInternalServerError)
                fmt.Fprint(w, "500 服务器内部错误")
            }
            if n, _ := rs.RowsAffected(); n > 0 {
                showURL, _ := router.Get("articles.show").URL("id", id)
                http.Redirect(w, r, showURL.String(), http.StatusFound)
            } else {
                fmt.Fprint(w, "您没有做任何更改！")
            }
        } else {
            updateURL, _ := router.Get("articles.update").URL("id", id)
            data := ArticlesFormData{
                Title:  title,
                Body:   body,
                URL:    updateURL,
                Errors: errors,
            }
            tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
            checkError(err)

            err = tmpl.Execute(w, data)
            checkError(err)
        }
    }
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
            checkError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "500 internal error")
        }
    } else {
        RowsAffected, err := article.Delete()
        if err != nil {
            checkError(err)
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
    initDB()
    createTables()
    route.Initialize()
    router = route.Router
    router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
    router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

    router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
    router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")
    router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).Methods("POST").Name("articles.update")
    router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")
    router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
    router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
    router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")

    router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
    router.Use(forceHTMLMiddleware)

    err := http.ListenAndServe(":8080", removeTrailingSlash(router))
    if err != nil {
        fmt.Println(err)
    }
}
