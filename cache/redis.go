package cache

import (
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
	// rmTicker := time.NewTicker(100 * time.Minute)
	prtTicker := time.NewTicker(1 * time.Minute)
	for {
		select {
		// case <-rmTicker.C:
		// 	continue
		// RClient.ZRemRangeByScore(ONLINE_USERS, "0", fmt.Sprintf("%d", time.Now().Add(-1*10*time.Minute).UnixMilli()))
		case <-prtTicker.C:
			PrintOnline()
		}
	}
}

// count set
func PrintOnline() {
	count := RClient.SCard(ONLINE_USERS)
	if count.Err() != nil {
		logrus.Error(count.Err().Error())
		return
	}
	logrus.Info("online user connection: ", count.Val())
}
