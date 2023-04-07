package cache

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
)

const (
	ONLINE_USERS string = "online_users"
)

var RClient *redis.Client

func InitRedis(ip string, port string, username string, password string, db int) {
	rRedisClient := redis.NewClient(&redis.Options{
		Addr:     ip + port,
		Username: username,
		Password: password,
		DB:       db,
	})
	RClient = rRedisClient

	p := rRedisClient.Ping()
	if p.Err() != nil {
		logrus.Info("error init redis --> ", p.Err().Error())
		return
	}
	logrus.Info("redis pong: ", p.Val())
}

func CleanOnlineZSet() {
	rmTicker := time.NewTicker(10 * time.Minute)
	prtTicker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-rmTicker.C:
			RClient.ZRemRangeByScore(ONLINE_USERS, "0", fmt.Sprintf("%d", time.Now().Add(-1*10*time.Minute).UnixMilli()))
		case <-prtTicker.C:
			PrintOnline()
		}
	}
}

func PrintOnline() {
	rc := RClient.ZCount(ONLINE_USERS, fmt.Sprintf("%d", time.Now().Add(-1*10*time.Minute).UnixMilli()), fmt.Sprintf("%d", time.Now().UnixMilli()))
	if rc.Err() != nil {
		logrus.Error("获取在线人数失败：", rc.Err().Error())
	} else {
		logrus.Info("当前【", rc.Val(), "】连接")
	}
}
