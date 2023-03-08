package server

import (
	"github.com/RustyNailPlease/go-relay/dao"
	"github.com/RustyNailPlease/go-relay/entity"
	"github.com/jinzhu/gorm"
	"github.com/nbd-wtf/go-nostr"
	"github.com/olahol/melody"
	"github.com/sirupsen/logrus"
)

// handlers by event kinds
var EVENT_HANDLER map[int]func(*melody.Session, *nostr.Event)

const (
	EVENT_KIND_SET_METADATA     int = 0
	EVENT_KIND_TEXT_NOTE        int = 1
	EVENT_KIND_RECOMMEND_SERVER int = 2
	EVENT_KIND_RELAY_LIST       int = 10002
)

func init() {
	EVENT_HANDLER = make(map[int]func(*melody.Session, *nostr.Event))

	EVENT_HANDLER[EVENT_KIND_SET_METADATA] = setMetaData
}

func setMetaData(s *melody.Session, event *nostr.Event) {
	var user entity.User

	pu, err := entity.GetUserFromProtocol(event)
	if err != nil {
		s.Write(SerialMessages("NOTICE", event.ID, "parse metadata error"))
		return
	}

	e := dao.DB.Model(&pu).Where("pubkey = ?", event.PubKey).First(&user)
	if gorm.IsRecordNotFoundError(e.Error) {
		dao.DB.Model(&pu).Create(&pu)
		s.Write(SerialMessages("NOTICE", event.ID, "saved."))
		return
	} else if e.Error != nil {
		logrus.Error(e.Error.Error())
		s.Write(SerialMessages("NOTICE", event.ID, "save meta error"))
		return
	}
	dao.DB.Model(&pu).Where("pubkey = ?", event.PubKey).Update(&pu)
	s.Write(SerialMessages("NOTICE", event.ID, "saved."))
}
