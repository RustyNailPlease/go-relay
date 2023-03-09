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

		nipsBuf, _ := json.Marshal([]int{1, 2, 3, 4, 11, 20, 65})

		rm = entity.RelayMeta{
			Pubkey:        serverConfig.Relay.AdminPubKey,
			Name:          "RustyWorld",
			Description:   "Welcome Home",
			Contact:       serverConfig.Relay.AdminPubKey,
			SupportedNips: nipsBuf,
			Software:      "http://192.168.2.4:50000",
			Version:       "0.1",
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
