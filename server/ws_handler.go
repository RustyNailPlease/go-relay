package server

import (
	"encoding/json"

	cacheutils "github.com/RustyNailPlease/CacheUtil"
	"github.com/RustyNailPlease/go-relay/cache"
	"github.com/RustyNailPlease/go-relay/dao"
	"github.com/RustyNailPlease/go-relay/entity"
	"github.com/olahol/melody"
	"github.com/sirupsen/logrus"
)

const (
	ONLINE_USERS string = "online_users"
)

var spammerKeyLRU *cacheutils.LRUCache[string]

func init() {
	spammerKeyLRU = cacheutils.NewLRU[string](5000)
}

func initWSHandlers() {
	wsServer.HandleConnect(func(s *melody.Session) {
		logrus.Info("==================connected=======================")
		logrus.Info(s.Request.Header.Get("Cf-Connecting-Ip") + " connected.")
		// logrus.Info(s.Request.RemoteAddr)
		// for k, v := range s.Request.Header {
		// 	logrus.Info(k, " --> ", strings.Join(v, ""))
		// }
		cache.RClient.SAdd(ONLINE_USERS, s.Request.Header.Get("Cf-Connecting-Ip"))
		// cache.RClient.ZAdd(ONLINE_USERS, &redis.Z{Score: float64(time.Now().UnixMilli()), Member: s.Request.Header.Get("Cf-Connecting-Ip")})
		// now users
		cache.PrintOnline()
		logrus.Info("==================connected=======================")
	})

	wsServer.HandleDisconnect(func(s *melody.Session) {
		logrus.Info("================== disconnected =======================")
		logrus.Info(s.Request.Header.Get("Cf-Connecting-Ip") + " disconnected.")
		cache.RClient.SRem(ONLINE_USERS, s.Request.Header.Get("Cf-Connecting-Ip"))
		logrus.Info("================== disconnected =======================")
	})

	wsServer.HandleMessage(func(s *melody.Session, b []byte) {
		midJson := make([]json.RawMessage, 0)

		if err := json.Unmarshal(b, &midJson); err != nil {
			logrus.Error(err.Error())
			s.Write(SerialMessages("NOTICE", "", "error handlering message"))
			return
		}

		var typ string
		if err := json.Unmarshal(midJson[0], &typ); err != nil {
			logrus.Error(err.Error())
			s.Write(SerialMessages("NOTICE", "", "error handlering message."))
			return
		}

		switch typ {
		case "EVENT":
			event := parseEventMessage(midJson)

			// check spammer
			if isSpamUser(event.PubKey) {
				s.Write(SerialMessages("NOTICE", event.ID, "banned"))
				return
			}

			logrus.Info("event: ", event.ID, "[", event.Kind, "]", " pub: ", event.PubKey)
			if check, err := event.CheckSignature(); (!check) || err != nil {
				if err != nil {
					s.Write(SerialMessages("NOTICE", event.ID, "check event sig error."))
					return
				}
				if !check {
					s.Write(SerialMessages("NOTICE", event.ID, "event sig invalid."))
					return
				}
			}

			if h, ok := EVENT_HANDLER[event.Kind]; ok {
				h(s, &event)
			} else {
				s.Write(SerialMessages("NOTICE", event.ID, "this event kind not implemented yet."))
			}

			s.Write(SerialMessages("OK", event.PubKey, ""))
		case "REQ":
			// logrus.Info("==============")
			// for _, s := range midJson {
			// 	logrus.Info(string(s))
			// }
			// logrus.Info("==============")
			subId, filters := parseReqFilterMessage(midJson)
			// if filters != nil {
			// 	logrus.Info()
			// }

			reqlog := entity.ReqLog{
				AcceptLanguage: s.Request.Header.Get("Accept-Language"),
				UserAgent:      s.Request.Header.Get("User-Agent"),
				Origin:         s.Request.Header.Get("Origin"),
				CFIPCountry:    s.Request.Header.Get("Cf-Ipcountry"),
				CFConnectingIP: s.Request.Header.Get("Cf-Connecting-Ip"),
				ReqBody:        json.RawMessage(b),
			}
			dao.DB.Create(&reqlog)

			handleReqRequest(s, subId, filters)
			// logrus.Info("filter req: ", subId,  " filters: "m f)
			// s.Write(SerialMessages("OK", subId, ""))
		default:
			break
		}

	})
}

func isSpamUser(pubKey string) (is bool) {
	is = false
	if spammerKeyLRU.Contains(pubKey) {
		return true
	}
	var count int
	dao.DB.Model(&entity.SpamUser{}).Where("user = ?", pubKey).Count(&count)
	if count > 0 {
		spammerKeyLRU.Set(pubKey, "")
		return true
	}
	return is
}
