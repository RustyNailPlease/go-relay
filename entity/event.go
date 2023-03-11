package entity

import (
	"encoding/json"
	"errors"

	"github.com/nbd-wtf/go-nostr"
)

type Event struct {
	ID        string          `json:"id"`
	PubKey    string          `json:"pubkey"`
	CreatedAt int64           `json:"created_at"` // secs
	Kind      int             `json:"kind"`
	Tags      json.RawMessage `gorm:"type:jsonb" json:"tags"`
	Content   string          `json:"content"`
	Sig       string          `json:"sig"`
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

type Events []Event

func (a Events) Len() int           { return len(a) }
func (a Events) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Events) Less(i, j int) bool { return a[i].CreatedAt > a[j].CreatedAt }
