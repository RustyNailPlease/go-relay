package dao

import (
	"fmt"

	"github.com/RustyNailPlease/go-relay/entity"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func InitDB(host string, port string, dbName string, username string, password string) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, username, password, dbName, port,
	)
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Errorf("连接数据库失败: %s", err))

	}
	DB = db
	DB.LogMode(true)
	db.AutoMigrate(&entity.Event{}, &entity.User{})
}
