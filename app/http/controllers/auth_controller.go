package controllers

import (
    "goblog/pkg/view"
    "goblog/app/models/user"
    "goblog/app/requests"
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
    _user := user.User{
        Name:               r.PostFormValue("name"),
        Email:              r.PostFormValue("email"),
        Password:           r.PostFormValue("password"),
        PasswordConfirm:    r.PostFormValue("password_confirm"),
    }

    errs := requests.ValidateRegistrationForm(_user)
    if len(errs) > 0 {
        view.RenderSimple(w, view.D{
            "Errors": errs,
            "User": _user,
        }, "register")
        fmt.Println("error:%T", errs)
    } else {
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
}

func (*AuthController) Login(w http.ResponseWriter, r *http.Request){
    view.RenderSimple(w,  view.D{}, "login")
}

func (*AuthController) DoLogin(w http.ResponseWriter, r *http.Request){
    //
}

