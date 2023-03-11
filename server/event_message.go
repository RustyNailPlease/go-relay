package server

import (
	"encoding/json"

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
	EVENT_HANDLER[EVENT_KIND_TEXT_NOTE] = publishTextNote
}

/*
*

	{"id":"67ed155dda2b215d8dc66efa40ae04ec05f9db0c54fef3335553ea014666f18c",
	"pubkey":"d7cbfa9c169fe0412e30ad50c7ef5bb578d57ba30abfa22a24eb3455b98c60cb",
	"created_at":1678525968,
	"kind":1,
	"tags":[],
	"content":"我也是",
	"sig":"80b8be15bdf2d6c9b0a1aa18b3e4cc1855c75b69eeb72dbc1d71202c0097c53580ce451a6c4898f2c0b45e5bb4bc26cf7918e8d4a81fe21a021d1ef47884eb9b"}
*/
func publishTextNote(s *melody.Session, event *nostr.Event) {
	buf, _ := json.Marshal(event)
	logrus.Info(string(buf))
	var count int
	dao.DB.Model(&entity.Event{}).Where("id = ?", event.ID).Count(&count)
	if count == 0 {
		// save it
		et, err := entity.FromNostrEvent(event)
		if err != nil {
			logrus.Error("parse event error :", err.Error())
			s.Write(SerialMessages("NOTICE", event.ID, "save note error"))
		} else {
			dao.DB.Model(&entity.Event{}).Create(&et)
			s.Write(SerialMessages("OK", event.ID, true, "note saved."))
		}
	} else {
		et, err := entity.FromNostrEvent(event)
		if err != nil {
			logrus.Error(err.Error())
			s.Write(SerialMessages("NOTICE", event.ID, "save note error"))
		} else {
			dao.DB.Model(&entity.Event{}).Where("id = ?", et.ID).Update(&et)
			s.Write(SerialMessages("OK", event.ID, true, "note saved."))
		}
	}
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
		pu.Relays = make([]entity.Relay, 0)
		dao.DB.Model(&pu).Create(&pu)
		msg := SerialMessages("OK", event.ID, true, "saved.")
		// logrus.Info("msg: ", msg)
		e := s.Write(msg)
		if e != nil {
			logrus.Error(e.Error())
		}
		return
	} else if e.Error != nil {
		logrus.Error(e.Error.Error())
		s.Write(SerialMessages("NOTICE", event.ID, "save meta error"))
		return
	}
	dao.DB.Model(&pu).Where("pubkey = ?", event.PubKey).Update(&pu)
	s.Write(SerialMessages("OK", event.ID, true, "saved."))
}
