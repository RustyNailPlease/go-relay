package server

import (
	"encoding/json"

	"github.com/nbd-wtf/go-nostr"
	"github.com/sirupsen/logrus"
)

// [] => Event
func parseEventMessage(raw []json.RawMessage) nostr.Event {
	var event nostr.Event
	if err := json.Unmarshal(raw[1], &event); err != nil {
		logrus.Error("json.Unmarshal: %v", err)
	}
	return event
}

// [] => Req
func parseReqFilterMessage(raw []json.RawMessage) (subid string, filters []nostr.Filter) {
	var id string
	if err := json.Unmarshal(raw[1], &id); err != nil {
		logrus.Error("json.Unmarshal sub id: %v", err)
	}
	var ff []nostr.Filter
	for i, b := range raw[2:] {
		var f nostr.Filter
		if err := json.Unmarshal(b, &f); err != nil {
			logrus.Error("json.Unmarshal filter %d: %v", i, err)
		}
		ff = append(ff, f)
	}
	return id, ff
}

func SerialMessages(eles ...interface{}) (b []byte) {
	buf, err := json.Marshal(eles)
	if err != nil {
		logrus.Error(err.Error())
		return b
	}
	return buf
}
