package api

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"

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
					data, err := parseData(event, hook)
					if err != nil {
						Logger.Warn("Could not parse Event data: %+v.", err)
					}

					go sendWebHook(hook.Config.Address, data)
				}
			}
		}
	}()
}

func parseData(event Event, hook account.Hook) ([]byte, error) {
	var err error
	tmpl := template.New(event.Name())
	data := bytes.NewBufferString("")

	tmpl, err = tmpl.Parse(string(event.Data()))
	if hook.Text != "" {
		tmpl, err = tmpl.Parse(hook.Text)
	}

	err = tmpl.Execute(data, event)
	if err != nil {
		return event.Data(), err
	}

	return data.Bytes(), nil
}

func sendWebHook(Address string, body interface{}) {
	if Address != "" {
		httpClient := requests.NewHTTPClient(Address)
		_, _, _, err := httpClient.MakeRequest(requests.Args{
			AcceptableCode: http.StatusOK,
			Method:         "POST",
			Body:           body,
		})

		if err != nil {
			Logger.Warn(fmt.Sprintf("Failed to call WebHook for %s: %s.", Address, err.Error()))
		}
	}
}
