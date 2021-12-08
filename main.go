package main

import (
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

var router = mux.NewRouter()
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

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    fmt.Fprintf(w, "article ID: "+id)
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "article list")
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

    errors := make(map[string]string)

    if title == "" {
        errors["title"] = "标题不存在"
    } else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
        errors["title"] = "标题长度需介于 3-40"
    }

    if body == "" {
        errors["body"] = "内容不存在"
    } else if utf8.RuneCountInString(body) < 10 {
        errors["body"] = "内容长度需大于等于10个字节"
    }

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

func forceHTMLMiddleware(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        h.ServeHTTP(w, r)
    })
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
    router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
    router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

    router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
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
