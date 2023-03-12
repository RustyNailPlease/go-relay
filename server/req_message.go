package server

import (
	"encoding/json"
	"sort"

	"github.com/RustyNailPlease/go-relay/dao"
	"github.com/RustyNailPlease/go-relay/entity"
	"github.com/nbd-wtf/go-nostr"
	"github.com/olahol/melody"
	"github.com/sirupsen/logrus"
)

func handleReqRequest(s *melody.Session, subid string, filters nostr.Filters) {

	if len(filters) == 0 {
		s.Write(SerialMessages("EOSE", subid))
	}

	for _, filter := range filters {
		query := dao.DB.Model(&entity.Event{})
		conditions := make(map[string]interface{})

		if filter.IDs != nil && len(filter.IDs) > 0 {
			conditions["id in (?)"] = filter.IDs
		}

		if filter.Authors != nil && len(filter.Authors) > 0 {
			conditions["pub_key in (?)"] = filter.Authors
		}

		if filter.Kinds != nil && len(filter.Kinds) > 0 {
			conditions["kind in (?)"] = filter.Kinds
		}

		if filter.Since != nil {
			conditions["created_at >= ?"] = filter.Since.Unix()
		}

		if filter.Until != nil {
			conditions["created_at < ?"] = filter.Until.Unix()
		}

		if len(filter.Tags) > 0 {
			tb, _ := json.Marshal(filter.Tags)
			logrus.Info("tag query: ", string(tb))
		}

		if len(conditions) > 0 {
			sql := "("
			args := make([]interface{}, 0)
			for k, v := range conditions {
				sql += k + " and "
				args = append(args, v)
			}
			sql += "1=1)"

			query = query.Where(sql, args...)
		}

		if filter.Limit > 0 && filter.Limit < 200 {
			query = query.Limit(filter.Limit)
		} else {
			query = query.Limit(20)
		}

		var es []entity.Event
		query.Find(&es)

		sort.Sort(entity.Events(es))

		for _, e := range es {
			// logrus.Info("下发： ", e.ID, "[", e.Content, "] --> ", subid)
			s.Write(SerialMessages("EVENT", subid, e))
		}
	}
	s.Write(SerialMessages("EOSE", subid))
}
