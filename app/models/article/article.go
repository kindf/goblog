package article

import (
    "goblog/pkg/route"
    "goblog/pkg/types"
    "goblog/app/models"
)

type Article struct {
    models.BaseModel

    Title string `gorm:"type:varchar(255);not null;" valid:"title"`
    Body string `gorm:"type:longtext;not null;" valid:"body"`
}

func (article Article) Link() string {
    return route.Name2URL("articles.show", "id", types.Uint64ToString(article.ID))
}
