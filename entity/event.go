package entity

import (
	"encoding/json"
	"errors"

	"github.com/nbd-wtf/go-nostr"
)

type Event struct {
	ID        string
	PubKey    string `json:"pubkey"`
	CreatedAt int64  `json:"created_at"` // secs
	Kind      int
	Tags      json.RawMessage `gorm:"type:jsonb"`
	Content   string
	Sig       string
}

func FromNostrEvent(e *nostr.Event) (Event, error) {
	if e == nil {
		return Event{}, errors.New("empty msg")
	}
	tags, _ := json.Marshal(e.Tags)
	return Event{
		ID:        e.ID,
		PubKey:    e.PubKey,
		CreatedAt: e.CreatedAt.Unix(),
		Kind:      e.Kind,
		Tags:      tags,
		Content:   e.Content,
		Sig:       e.Sig,
	}, nil
}
