package article

import (
    "goblog/pkg/route"
)

type Article struct {
    ID uint64
    Title string
    Body string
}

func (article Article) Link() string {
    return route.Name2URL("articles.show", "id", article.GetStringID())
}

func (article Article) GetStringID() string {
    return "1"
}
