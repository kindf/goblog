package article

import (
    "goblog/pkg/route"
    "goblog/pkg/types"
    "goblog/app/models"
)

type Article struct {
    models.BaseModel

    ID uint64
    Title string
    Body string
}

func (article Article) Link() string {
    return route.Name2URL("articles.show", "id", types.Uint64ToString(article.ID))
}
