package controllers

import (
    "goblog/pkg/view"
    "goblog/app/models/user"
    "net/http"
    "fmt"
    "time"
)

type AuthController struct {
}

func (*AuthController) Register(w http.ResponseWriter, r *http.Request){
    view.RenderSimple(w,  view.D{}, "register")
}

func (*AuthController) DoRegister(w http.ResponseWriter, r *http.Request){
    name := r.PostFormValue("name")
    email := r.PostFormValue("email")
    password := r.PostFormValue("password")

    _user := user.User {
        Name: name,
        Email: email,
        Password: password,
    }
    _user.CreateAt = time.Now()
    _user.UpdateAt = time.Now()
    _user.Create()

    if _user.ID > 0 {
        fmt.Fprint(w, "插入成功，ID为"+_user.GetStringID())
    } else {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprint(w, "创建用户失败，请联系管理员")
    }
}
