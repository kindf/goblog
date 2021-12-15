package models

import (
    "goblog/pkg/types"
    "time"
)

type BaseModel struct {
    ID          uint64      `gorm:"column:id;primaryKey;autoIncrement;not null"`
    CreateAt    time.Time   `gorm:"column:created_at;index"`
    UpdateAt    time.Time   `gorm:"column:updated_at;index"`
}

func (a BaseModel) GetStringID() string {
    return types.Uint64ToString(a.ID)
}
