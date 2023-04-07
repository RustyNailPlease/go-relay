package main

import (
	"flag"

	"github.com/RustyNailPlease/go-relay/cache"
	"github.com/RustyNailPlease/go-relay/config"
	"github.com/RustyNailPlease/go-relay/dao"
	"github.com/RustyNailPlease/go-relay/server"
	"github.com/spf13/viper"
)

var configFile *string

func init() {
	viper.AddConfigPath(".")
	configFile = flag.String("c", "./.config.toml", "配置文件路径")
	flag.Parse()
}

func main() {
	viper.SetConfigFile(*configFile)
	err := viper.ReadInConfig()

	if err != nil {
		panic("读取配置文件失败:" + err.Error())
	}

	var config config.ServerConfig
	err = viper.Unmarshal(&config)
	if err != nil {
		panic("解析配置文件失败:" + err.Error())
	}

	go dao.InitDB(
		config.PGSQL.Host,
		config.PGSQL.Port,
		config.PGSQL.DBName,
		config.PGSQL.Username,
		config.PGSQL.Password,
		config.ServerMode == "debug",
	)

	go cache.InitRedis(
		config.Redis.Host,
		config.Redis.Port,
		config.Redis.Username,
		config.Redis.Password,
		config.Redis.DB,
	)

	go cache.CleanOnlineZSet()

	server.InitServer(&config)
}
