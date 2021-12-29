package user

import (
    "goblog/app/models"
)

type User struct {
    models.BaseModel

    Name                string `gorm:"column:name;type:varchar(255);not null;unique" valid:"name"`
    Email               string `gorm:"column:email;type:varchar(255);default:NULL;unique;" valid:"email"`
    Password            string `gorm:"column:password;type:varchar(255)" valid:"password"`
    PasswordConfirm     string `gorm:"-" valid:"password_confirm"`
}

func (user *User) ComparePassword(password string) bool {
    return user.Password == password
}

func (user User) Link() string {
    return ""
}
