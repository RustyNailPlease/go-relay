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
	EVENT_KIND_CONTACTS         int = 3
	EVENT_KIND_DELETION         int = 5
	EVENT_KIND_REACTION         int = 7
	EVENT_KIND_RELAY_LIST       int = 10002
)

func init() {
	EVENT_HANDLER = make(map[int]func(*melody.Session, *nostr.Event))

	EVENT_HANDLER[EVENT_KIND_SET_METADATA] = setMetaData
	EVENT_HANDLER[EVENT_KIND_TEXT_NOTE] = publishTextNote
	EVENT_HANDLER[EVENT_KIND_RELAY_LIST] = setRelays
	EVENT_HANDLER[EVENT_KIND_CONTACTS] = setContacts
	EVENT_HANDLER[EVENT_KIND_DELETION] = deleteEvent
	EVENT_HANDLER[EVENT_KIND_REACTION] = saveReaction
}

func deleteEvent(s *melody.Session, event *nostr.Event) {
	// todo
	e, err := entity.FromNostrEvent(event)
	if err != nil {
		logrus.Error(err.Error())
		s.Write(SerialMessages("NOTICE", event.ID, "parse event error"))
		return
	}
	dao.DB.Model(&e).Create(&e)

	for _, tag := range event.Tags {
		dao.DB.Model(&entity.Event{}).Where("pub_key = ? and id = ?", event.PubKey, tag[1]).UpdateColumn("content", e.Content)
	}
	s.Write(SerialMessages("OK", event.ID, true, "event deleted saved."))
}

func saveReaction(s *melody.Session, event *nostr.Event) {
	// b, _ := json.Marshal(event)
	// logrus.Info("contract: ", string(b))
	var count int
	dao.DB.Model(&entity.Event{}).Where("id = ? and kind = ?", event.ID, event.Kind).Count(&count)
	if count == 0 {
		e, err := entity.FromNostrEvent(event)
		if err != nil {
			logrus.Error(err.Error())
			s.Write(SerialMessages("NOTICE", event.ID, "save reaction error"))
			return
		}
		dao.DB.Model(&entity.Event{}).Create(&e)
		s.Write(SerialMessages("OK", event.ID, true, "reaction saved."))
	} else {
		e, err := entity.FromNostrEvent(event)
		if err != nil {
			logrus.Error(err.Error())
			s.Write(SerialMessages("NOTICE", event.ID, "save reaction error"))
			return
		}
		dao.DB.Model(&entity.Event{}).Where("id = ? and kind = ?", event.ID, event.Kind).Update(&e)
		s.Write(SerialMessages("OK", event.ID, true, "reaction saved."))
	}
}

func setContacts(s *melody.Session, event *nostr.Event) {
	// b, _ := json.Marshal(event)
	// logrus.Info("contract: ", string(b))

	var count int
	dao.DB.Model(&entity.Event{}).Where("pub_key = ? and kind = ?", event.PubKey, event.Kind).Count(&count)
	if count == 0 {
		e, err := entity.FromNostrEvent(event)
		if err != nil {
			logrus.Error(err.Error())
			s.Write(SerialMessages("NOTICE", event.ID, "save contract error"))
			return
		}
		dao.DB.Model(&entity.Event{}).Create(&e)
	} else {
		e, err := entity.FromNostrEvent(event)
		if err != nil {
			logrus.Error(err.Error())
			s.Write(SerialMessages("NOTICE", event.ID, "save contract error"))
			return
		}
		dao.DB.Model(&entity.Event{}).Where("pub_key = ? and kind = ?", event.PubKey, event.Kind).Update(&e)
	}
	s.Write(SerialMessages("OK", event.ID, true, "contract saved."))
}

func setRelays(s *melody.Session, event *nostr.Event) {
	logrus.Info("event: ", event)
	var user entity.User

	rs := make([]entity.Relay, 0)
	for _, t := range event.Tags {
		tmp := entity.Relay{
			Url: t[1],
		}
		if len(t) > 2 {
			tmp.Read = t[2]
		} else {
			tmp.Read = "write"
		}
		rs = append(rs, tmp)
	}

	buf, _ := json.Marshal(rs)

	o := dao.DB.Model(&entity.User{}).Where("pubkey = ?", event.PubKey).Find(&user)
	if o.Error != nil && !gorm.IsRecordNotFoundError(o.Error) {
		s.Write(SerialMessages("NOTICE", event.ID, "save relays error"))
		return
	} else if o.Error != nil && gorm.IsRecordNotFoundError(o.Error) {
		user = entity.User{}
		user.Pubkey = event.PubKey
		user.Relays = buf

		dao.DB.Model(&entity.User{}).Create(&user)
	}

	user.Relays = buf
	dao.DB.Model(&entity.User{}).Where("pubkey = ?", event.PubKey).Update(&user)
	s.Write(SerialMessages("OK", event.ID, true, "relays saved."))
}

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
	if e.Error != nil && gorm.IsRecordNotFoundError(e.Error) {
		pu.Relays, _ = json.Marshal(make([]entity.Relay, 0))

		dao.DB.Model(&pu).Create(&pu)
		msg := SerialMessages("OK", event.ID, true, "saved.")
		// logrus.Info("msg: ", msg)
		e := s.Write(msg)
		if e != nil {
			logrus.Error(e.Error())
		}
		return
	} else if e.Error != nil && !gorm.IsRecordNotFoundError(e.Error) {
		logrus.Error(e.Error.Error())
		s.Write(SerialMessages("NOTICE", event.ID, "save meta error"))
		return
	}
	dao.DB.Model(&pu).Where("pubkey = ?", event.PubKey).Update(&pu)
	s.Write(SerialMessages("OK", event.ID, true, "saved."))
}
