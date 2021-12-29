package article

import (
    "goblog/pkg/route"
    "goblog/pkg/types"
    "goblog/app/models"
    "goblog/app/models/user"
)

type Article struct {
    models.BaseModel

    Title string `gorm:"type:varchar(255);not null;" valid:"title"`
    Body string `gorm:"type:longtext;not null;" valid:"body"`
    UserID uint64 `gorm:"not null;index"`
    User user.User
}

func (article Article) Link() string {
    return route.Name2URL("articles.show", "id", types.Uint64ToString(article.ID))
}

func (article Article) CreateAtDate() string {
    return article.CreateAt.Format("2006-01-02")
}
