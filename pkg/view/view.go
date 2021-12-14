package view

import (
    "io"
    "html/template"
    "goblog/pkg/logger"
    "goblog/pkg/route"
)

func Render(w io.Writer, name string, data interface{}) {
    viewDir := "./static/"
    files := []string{"./static/app.gohtml", "./static/sidebar.gohtml"}
    files = append(files, viewDir+name+".gohtml")


    tmpl, err := template.New(name+".gohtml").Funcs(template.FuncMap{"RouteName2URL":route.Name2URL,}).ParseFiles(files...)
    logger.LogError(err)
    err = tmpl.ExecuteTemplate(w, "app", data)
    logger.LogError(err)
}
