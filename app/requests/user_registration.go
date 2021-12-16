package requests

import (
    "github.com/thedevsaddam/govalidator"
    "goblog/app/models/user"
)

func ValidateRegistrationForm(data user.User) map[string][]string {
    rules := govalidator.MapData{
        "name":             []string{"required", "alpha_num", "between:3,20"},
        "email":            []string{"required", "min:4", "max:30", "email"},
        "password":         []string{"required", "min:6"},
        "password_confirm": []string{"required"},
    }

    messages := govalidator.MapData{
        "name":             []string{
            "required:用户名必填项",
            "alpha_num:格式错误，只允许数字和英文",
            "between:用户名长度需在3~20之间",
        },

        "email":            []string{
            "required:邮箱必填项",
            "min:Email长度需大于4",
            "max:Email长度需小于30",
            "email:Email格式不正确",
        },

        "password":         []string{
            "required:密码必填项",
            "min:长度需大于6",
        },

        "password_confirm": []string{
            "required:确认密码必填项",
        },
    }

    opts := govalidator.Options{
        Data:           &data,
        Rules:          rules,
        TagIdentifier:  "valid",
        Messages:       messages,
    }

    errs := govalidator.New(opts).ValidateStruct()

    if data.Password != data.PasswordConfirm{
        errs["password_confirm"] = append(errs["password_confirm"], "两次输入密码不匹配")
    }
    return errs
}
