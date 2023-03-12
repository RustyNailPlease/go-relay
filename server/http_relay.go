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

func onNip05(ctx *gin.Context) {
	ctx.Request.ParseForm()
	name, ok := ctx.Params.Get("name")
	if !ok {
		ctx.JSON(http.StatusOK, make(map[string]interface{}))
		return
	}

	var users []entity.User
	err := dao.DB.Model(&entity.User{}).Where("name = ? and signed_nip5 = ?", name, true).Find(&users)
	if err.Error != nil && gorm.IsRecordNotFoundError(err.Error) {
		ctx.JSON(http.StatusOK, make(map[string]interface{}))
		return
	}

	var r entity.Nip5Response

	for _, u := range users {
		r.Names[u.Name] = u.Pubkey

		if u.Relays != nil && len(u.Relays) > 0 {
			var relays []entity.Relay
			if err := json.Unmarshal(u.Relays, &relays); err == nil && len(relays) > 0 {
				var ru []string
				for _, relay := range relays {
					ru = append(ru, relay.Url)
				}
				r.Relays[u.Pubkey] = ru
			}

		}
	}

	ctx.JSON(http.StatusOK, r)
}
