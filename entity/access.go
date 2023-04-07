package entity

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
)

type ReqLog struct {
	gorm.Model

	AcceptLanguage string
	UserAgent      string
	Origin         string
	CFIPCountry    string
	CFConnectingIP string
	ReqBody        json.RawMessage `gorm:"type:jsonb"`
}
