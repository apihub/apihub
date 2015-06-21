package api

import (
	"fmt"
	"net/http"

	"github.com/backstage/maestro/account"
	. "github.com/backstage/maestro/log"
	"github.com/backstage/maestro/requests"
)

type Event interface {
	Data() []byte
	Name() string
}

func (api *Api) EventNotifier(ev Event) {
	api.Events <- ev
}

func (api *Api) ListenEvents() {
	Logger.Info("Started listening to events in background.")

	go func() {
		for event := range api.Events {
			// TODO: Need to load team specific!
			allHookw, err := account.FindHooksByEvent(event.Name())

			if len(allHookw) > 0 && err == nil {
				Logger.Debug(fmt.Sprintf("Start sending hooks to the following list: %+v.", allHookw))
				for _, hook := range allHookw {
					go sendWebHook(hook.Config.URL, string(event.Data()))
				}
			}
		}
	}()
}

func sendWebHook(URL string, body interface{}) {
	if URL != "" {
		httpClient := requests.NewHTTPClient(URL)
		_, _, _, err := httpClient.MakeRequest(requests.Args{
			AcceptableCode: http.StatusOK,
			Method:         "POST",
			Path:           "",
			Body:           body,
		})

		if err != nil {
			Logger.Warn(fmt.Sprintf("Failed to call WebHook for %s: %s.", URL, err.Error()))
		}
	}
}
