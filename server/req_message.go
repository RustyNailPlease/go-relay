package server

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

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

		if filter.Since == nil && filter.Until == nil {
			conditions["created_at >= ?"] = time.Now().Add(-12 * time.Hour).Unix()
		}

		if len(filter.Tags) > 0 {
			tb, _ := json.Marshal(filter.Tags)
			logrus.Info("tag query: ", string(tb))

			if eid, ok := filter.Tags["e"]; ok {
				conditions["tags @> ?"] = fmt.Sprintf("[[\"e\", \"%s\"]]", eid[0])
			}

			if pid, ok := filter.Tags["p"]; ok {
				conditions["tags @> ?"] = fmt.Sprintf("[[\"p\", \"%s\"]]", pid[0])
			}

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

		query = query.Order("created_at desc")

		if filter.Limit > 0 && filter.Limit <= serverConfig.MaxRows {
			query = query.Limit(filter.Limit)
		} else {
			if serverConfig.MaxRows == -1 {
				//query = query.Limit()
			} else {
				query = query.Limit(serverConfig.MaxRows)
			}
		}

		var es []entity.Event
		query.Find(&es)

		sort.Sort(entity.Events(es))

		for _, e := range es {
			if eventDeleted(e.ID) {
				continue
			}
			s.Write(SerialMessages("EVENT", subid, e))
		}
	}
	s.Write(SerialMessages("EOSE", subid))
}

func eventDeleted(id string) bool {

	if deletedCache.Contains(id) {
		return true
	}

	var count int
	dao.DB.Model(&entity.Event{}).Where("id = ? and kind = 5", id).Count(&count)
	if count > 0 {
		logrus.Info(id, " skip for deleted")
		deletedCache.Set(id, id)
		return true
	}
	return false
}
