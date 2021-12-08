package main

import (
    "net/http"
    "fmt"
    "strings"
    "html/template"
    "unicode/utf8"
    "net/url"
    "github.com/gorilla/mux"
)

 var router = mux.NewRouter()

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
    fmt.Fprintf(w, "create article")

    title := r.PostFormValue("title")
    body := r.PostFormValue("body")

    errors := make(map[string]string)

    if title == "" {
        errors["title"] = "标题不存在"
    } else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
        errors["title"] = "标题长度需介于 3-40"
    }

    if body == "" {
        errors["body"] = "标题不存在"
    } else if utf8.RuneCountInString(body) < 10 {
        errors["body"] = "内容长度需大于等于10个字节"
    }

    if len(errors) == 0 {
    } else {
        html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <title>创建文章 —— 我的技术博客</title>
    <style type="text/css">.error {color: red;}</style>
</head>
<body>
    <form action="{{ .URL }}" method="post">
        <p><input type="text" name="title" value="{{ .Title }}"></p>
        {{ with .Errors.title }}
        <p class="error">{{ . }}</p>
        {{ end }}
        <p><textarea name="body" cols="30" rows="10">{{ .Body }}</textarea></p>
        {{ with .Errors.body }}
        <p class="error">{{ . }}</p>
        {{ end }}
        <p><button type="submit">提交</button></p>
    </form>
</body>
</html>
`
        storeURL, _ := router.Get("articles.store").URL()
        data := ArticlesFormData {
            Title: title,
            Body: body,
            URL: storeURL,
            Errors: errors,
        }
        tmpl, err := template.New("create-form").Parse(html)
        if err != nil {
            panic(err)
        }

        err = tmpl.Execute(w, data)
        if err != nil {
            panic(err)
        }
    }
}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "create blog")
    html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <title>创建文章 —— 我的技术博客</title>
</head>
<body>
    <form action="%s" method="post">
        <p><input type="text" name="title"></p>
        <p><textarea name="body" cols="30" rows="10"></textarea></p>
        <p><button type="submit">提交</button></p>
    </form>
</body>
</html>
`
    storeURL, _ := router.Get("articles.store").URL()
    fmt.Fprintf(w, html, storeURL)
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
