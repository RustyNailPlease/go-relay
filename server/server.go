package server

import (
	"net/http"

	"github.com/RustyNailPlease/go-relay/config"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"github.com/sirupsen/logrus"
)

var httpServer *gin.Engine
var wsServer *melody.Melody
var serverConfig *config.ServerConfig

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{})
}

func InitServer(config *config.ServerConfig) {
	serverConfig = config
	// gin mode
	if config.ServerMode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	// init http & ws server
	httpServer = gin.New()
	httpServer.Use(gin.Logger(), gin.Recovery())
	httpServer.SetTrustedProxies(make([]string, 0))

	wsServer = melody.New()
	wsServer.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// 100 Mb
	wsServer.Config.MaxMessageSize = serverConfig.MaxMessageSize
	// init
	initHandlers()

	httpServer.Run(serverConfig.ServerPort)
}

func initHandlers() {
	httpServer.GET(serverConfig.ServerPath, func(ctx *gin.Context) {

		// return relay meta data
		acceptHeader := ctx.GetHeader("Accept")
		if acceptHeader != "" && acceptHeader == "application/nostr+json" {
			onNip11(ctx)
			return
		}

		err := wsServer.HandleRequest(ctx.Writer, ctx.Request)
		if err != nil {
			logrus.Panic(err.Error())
		}
	})
	// ws handlers' entry
	initWSHandlers()
}
