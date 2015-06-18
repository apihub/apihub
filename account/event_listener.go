package account

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/backstage/maestro/log"
)

const (
	DEFAULT_EVENTS_CHANNEL_LEN = 100
)

var listener *WebHookListener = newWebHookAndListen()

type Event interface {
	Data() *EventData
}

type EventData struct {
	name string
	team string
}

type WebHookListener struct {
	C chan Event
}

func newWebHookAndListen() *WebHookListener {
	l := &WebHookListener{C: make(chan Event, DEFAULT_EVENTS_CHANNEL_LEN)}
	l.Listen()
	return l
}

func (wh *WebHookListener) Listen() {
	go func() {
		for {

			select {
			case event := <-wh.C:
				e, err := json.Marshal(event)
				if err != nil {
					Logger.Error(fmt.Sprintf("Failed to marshal the event content: %s.", err.Error()))
					continue
				}

				whs, err := store.FindWebhooksByEventAndTeam(event.Data().name, event.Data().team)
				Logger.Debug(fmt.Sprintf("Hooks for: %+v.", whs))
				if err == nil {
					Logger.Debug("Start sending...")
					for _, hook := range whs {
						if hook.Config.Url != "" {

							httpClient := NewHTTPClient(hook.Config.Url)
							_, _, _, err = httpClient.MakeRequest(RequestArgs{
								AcceptableCode: http.StatusOK,
								Method:         "POST",
								Path:           "",
								Body:           string(e),
							})
							if err != nil {
								Logger.Error(fmt.Sprintf("Failed to call WebHook: %s.", err.Error()))
							}
							Logger.Debug(fmt.Sprintf("Sent to: %s", string(e)))

						}
					}
					Logger.Debug("Hooks sent.")
				}
			}

		}
	}()
}

func sendHook(event Event) {
	listener.C <- event
}
