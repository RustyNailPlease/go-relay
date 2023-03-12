package server

import (
	"encoding/json"
	"net/http"

	"github.com/RustyNailPlease/go-relay/dao"
	"github.com/RustyNailPlease/go-relay/entity"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func onNip11(ctx *gin.Context) {
	var rm entity.RelayMeta
	errO := dao.DB.Model(&entity.RelayMeta{}).Where("pubkey = ?", serverConfig.Relay.AdminPubKey).First(&rm)
	if gorm.IsRecordNotFoundError(errO.Error) {

		nipsBuf, _ := json.Marshal([]int{1})
		if serverConfig.Relay.Nips != nil || len(serverConfig.Relay.Nips) > 0 {
			nipsBuf, _ = json.Marshal(serverConfig.Relay.Nips)
		}

		rm = entity.RelayMeta{
			Pubkey:        serverConfig.Relay.AdminPubKey,
			Name:          serverConfig.Relay.Name,
			Description:   serverConfig.Relay.Description,
			Contact:       serverConfig.Relay.Contract,
			SupportedNips: nipsBuf,
			Software:      serverConfig.Relay.Software,
			Version:       serverConfig.Relay.Version,
		}
		dao.DB.Model(&entity.RelayMeta{}).Create(&rm)
		ctx.JSON(http.StatusOK, rm)
	} else {
		if errO.Error == nil {
			ctx.JSON(http.StatusOK, rm)
			return
		}
		ctx.JSON(http.StatusNotFound, nil)
	}
}
