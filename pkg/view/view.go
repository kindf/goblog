package view

import (
    "io"
    "html/template"
    "goblog/pkg/logger"
    "goblog/pkg/route"
)

func Render(w io.Writer, data interface{}, tplFiles ...string) {
    viewDir := "./static/"
    files := []string{"./static/app.gohtml", "./static/sidebar.gohtml"}
    for _, f := range tplFiles {
        files = append(files, viewDir+f+".gohtml")
    }
    tmpl, err := template.New("").Funcs(template.FuncMap{"RouteName2URL":route.Name2URL,}).ParseFiles(files...)
    logger.LogError(err)
    err = tmpl.ExecuteTemplate(w, "app", data)
    logger.LogError(err)
}
