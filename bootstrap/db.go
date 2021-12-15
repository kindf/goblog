package bootstrap

import (
    "goblog/pkg/model"
    "goblog/app/models/user"
    "goblog/app/models/article"
    "time"
    "gorm.io/gorm"
)

func SetupDB() {
    db := model.ConnectDB()
    
    sqlDB, _ := db.DB()

    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetMaxIdleConns(25)
    sqlDB.SetConnMaxLifetime(5 * time.Minute)
    migration(db)
}

func migration(db *gorm.DB){
    db.AutoMigrate(
        &user.User{},
        &article.Article{},
    )
}
