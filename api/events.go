package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"text/template"

	"github.com/apihub/apihub/account"
	. "github.com/apihub/apihub/log"
	"github.com/apihub/apihub/requests"
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
			// Send Hook events
			go func(event Event) {
				sendHooks(event)
			}(event)

		}
	}()
}

func sendHooks(event Event) {
	// TODO: Need to load team specific!
	allHookw, err := account.FindHooksByEvent(event.Name())

	if len(allHookw) > 0 && err == nil {
		Logger.Debug(fmt.Sprintf("Start sending hooks to the following list: %+v.", allHookw))
		for _, hook := range allHookw {
			data, err := parseData(event, hook)
			if err != nil {
				Logger.Warn("Could not parse Event data: %+v. Default format will be delivered.", err)
				data = event.Data()
			}

			go func(config account.HookConfig, data []byte) {
				sendWebHook(config, data)
			}(hook.Config, data)
		}
	}
}

func parseData(event Event, hook account.Hook) ([]byte, error) {
	var err error
	tmpl := template.New(event.Name())
	data := bytes.NewBufferString("")

	if hook.Text != "" {
		tmpl, err = tmpl.Parse(hook.Text)
	}

	err = tmpl.Execute(data, event)
	if err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

func sendWebHook(config account.HookConfig, body interface{}) {
	if config.Address != "" {
		url, err := url.Parse(config.Address)
		if config.Method == "" {
			config.Method = "POST"
		}
		b := body.([]byte)
		httpClient := requests.NewHTTPClient(fmt.Sprintf("%s://%s", url.Scheme, url.Host))
		_, _, _, err = httpClient.MakeRequest(requests.Args{
			AcceptableCode: http.StatusOK,
			Method:         config.Method,
			Path:           url.RequestURI(),
			Body:           string(b),
		})

		if err != nil {
			Logger.Warn(fmt.Sprintf("Failed to call WebHook for %s: %s.", config.Address, err.Error()))
		}
	}
}
