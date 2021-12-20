package view

import (
    "io"
    "html/template"
    "goblog/pkg/logger"
    "goblog/pkg/route"
    "goblog/pkg/auth"
    "goblog/pkg/flash"
)

type D map[string]interface{}

func Render(w io.Writer, data D, tplFiles ...string) {
    RenderTemplate(w, "app", data, tplFiles...)
}

func RenderSimple(w io.Writer, data D, tplFiles ...string) {
    RenderTemplate(w, "simple", data, tplFiles...)
}

func RenderTemplate(w io.Writer, name string, data D, tplFiles ...string) {
    data["isLogined"] = auth.Check()
    data["loginUser"] = auth.User
    data["flash"] = flash.All()

    allFiles := getTemplateFiles(tplFiles...)

    tmpl, err := template.New("").Funcs(template.FuncMap{"RouteName2URL":route.Name2URL,}).ParseFiles(allFiles...)
    logger.LogError(err)
    err = tmpl.ExecuteTemplate(w, name, data)
    logger.LogError(err)
}

func getTemplateFiles(tplFiles ...string) []string {
    viewDir := "./static/"
    files := []string{"./static/app.gohtml", "./static/sidebar.gohtml", "./static/simple.gohtml", "./static/_form_error_feedback.gohtml", "./static/_form_field.gohtml", "./static/_messages.gohtml"}
    for _, f := range tplFiles {
        files = append(files, viewDir+f+".gohtml")
    }
    return files
}


