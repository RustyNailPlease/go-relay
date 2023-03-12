package entity

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
	"github.com/nbd-wtf/go-nostr"
	"github.com/sirupsen/logrus"
)

type Relay struct {
	Url  string
	Read string
}

type User struct {
	gorm.Model
	Pubkey     string `gorm:"primaryKey"`
	Name       string
	About      string
	Picture    string
	SignedNip5 bool
	Relays     json.RawMessage `gorm:"type:jsonb"`
}

func GetUserFromProtocol(event *nostr.Event) (user User, e error) {
	if pm, err := nostr.ParseMetadata(*event); err == nil {
		user.Pubkey = event.PubKey
		user.Name = pm.Name
		user.About = pm.About
		user.Picture = pm.Picture
		return user, nil
	} else {
		logrus.Error(err.Error())
		return user, err
	}
}
