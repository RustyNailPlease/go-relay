package dao

import (
	"fmt"

	"github.com/RustyNailPlease/go-relay/entity"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
)

var DB *gorm.DB

func InitDB(host string, port string, dbName string, username string, password string, debug bool) {
	logrus.Info("数据库初始化。")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, username, password, dbName, port,
	)
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Errorf("连接数据库失败: %s", err))

	}
	DB = db
	DB.LogMode(debug)
	db.AutoMigrate(&entity.Event{}, &entity.User{}, &entity.RelayMeta{}, &entity.SpamUser{})
	logrus.Info("数据库初始化完成。")
}
