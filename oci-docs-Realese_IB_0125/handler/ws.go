package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/model"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *Handler) waitNotification(l *pq.Listener, conn *websocket.Conn, messageType int, info *model.WaitRegFinish, done <-chan bool) {
	startTime := time.Now()
	for {
		select {
		case n := <-l.Notify:

			h.log.Debugf("Received data from channel [%s] : %s", n.Channel, n.Extra)

			var event model.DBEvent
			if err := json.Unmarshal([]byte(n.Extra), &event); err != nil {
				h.log.Errorf("Error processing JSON DB notification: %v", err)
				return
			}
			if event.Table == "telegram_chats" && strings.ToUpper(event.Action) == "INSERT" && event.Data.Phone == info.Phone {
				conn.WriteMessage(messageType, []byte("success"))
				h.log.Debugf("write message for notification success")
				return
			}
		case d := <-done:
			{
				h.log.Debugf("Done wait notification %t", d)
				return
			}
		case <-time.After(90 * time.Second):
			go func() {
				l.Ping()
			}()
			diff := time.Since(startTime)
			if diff.Hours() > 5 {
				h.log.Debugf("Done wait notification by timeout %f", diff)
				return
			}
		}
	}
}

func (h *Handler) wshandler(w http.ResponseWriter, r *http.Request) {
	wsupgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Errorf("Failed to set websocket upgrade: %+v", err)
		return
	}

	logError := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			h.log.Errorf("DB Listener error: %s", err.Error())
		}
	}
	done := make(chan bool)

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		h.log.Debugf("WS ReadMessage: %s", string(msg))
		var message model.WSMessage
		err = json.Unmarshal(msg, &message)
		if err != nil {
			h.log.Errorf("Unmarshal WS message error: %s", err.Error())
			return
		}
		var info model.WaitRegFinish
		if err := mapstructure.Decode(message.Data, &info); err != nil {
			h.log.Errorf("Decode data error: %s", err.Error())
			return
		}

		if message.Name == "wait_end_registration" {
			go func() {
				listener := pq.NewListener(h.conninfo, 10*time.Second, time.Minute, logError)
				defer func() {
					listener.Close()
				}()
				err = listener.Listen("events")
				if err != nil {
					h.log.Errorf("Listen DB events error: %s", err.Error())
					return
				}
				h.waitNotification(listener, conn, t, &info, done)

			}()
		} else if message.Name == "stop_wait_end_registration" {
			h.log.Debug("stop_wait_end_registration")
			done <- true
		}
	}
}

func (h *Handler) WS(c *gin.Context) {
	h.wshandler(c.Writer, c.Request)
}
