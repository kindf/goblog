package view

import (
    "io"
    "html/template"
    "goblog/pkg/logger"
    "goblog/pkg/route"
)

type D map[string]interface{}

func Render(w io.Writer, data interface{}, tplFiles ...string) {
    RenderTemplate(w, "app", data, tplFiles...)
}

func RenderSimple(w io.Writer, data interface{}, tplFiles ...string) {
    RenderTemplate(w, "simple", data, tplFiles...)
}

func RenderTemplate(w io.Writer, name string, data interface{}, tplFiles ...string) {
    viewDir := "./static/"
    files := []string{"./static/app.gohtml", "./static/sidebar.gohtml", "./static/simple.gohtml", "./static/_form_error_feedback.gohtml", "./static/_form_field.gohtml"}
    for _, f := range tplFiles {
        files = append(files, viewDir+f+".gohtml")
    }
    tmpl, err := template.New("").Funcs(template.FuncMap{"RouteName2URL":route.Name2URL,}).ParseFiles(files...)
    logger.LogError(err)
    err = tmpl.ExecuteTemplate(w, name, data)
    logger.LogError(err)
}
