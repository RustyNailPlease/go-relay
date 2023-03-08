package entity

import "encoding/json"

type Event struct {
	ID        string
	PubKey    string `json:"pubkey"`
	CreatedAt int64  `json:"created_at"` // secs
	Kind      int
	Tags      json.RawMessage `gorm:"type:json"`
	Content   string
	Sig       string
}
