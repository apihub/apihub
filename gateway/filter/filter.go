package filter

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/clbanning/mxj/j2x"
	"github.com/clbanning/mxj/x2j"
)

// Function which modify the response.
type Filter func(*http.Request, *http.Response, *bytes.Buffer)

// An array of Filter with key to be used by the gateway and service.
type Filters map[string]Filter

func (f Filters) Add(key string, value Filter) {
	f[key] = value
}

func (f Filters) Get(key string) Filter {
	return f[key]
}

// Given a xml as response body, convert it to json.
func ConvertXmlToJson(r *http.Request, w *http.Response, body *bytes.Buffer) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		log.Print("err:", err.Error())
		body.Write([]byte(err.Error()))
		return
	}
	m, err := x2j.XmlToJson(b)
	if err != nil {
		log.Print("err:", err.Error())
		body.Write([]byte(err.Error()))
		return
	}
	w.Header.Set("Content-Type", "application/json")
	body.Reset()
	body.Write(m)
}

// Given a json as response body, convert it to xml.
func ConvertJsonToXml(r *http.Request, w *http.Response, body *bytes.Buffer) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		log.Print("err: ", err.Error())
		body.Write([]byte(err.Error()))
		return
	}
	m, err := j2x.JsonToXml(b)
	if err != nil {
		log.Print("err:", err.Error())
		body.Write([]byte(err.Error()))
		return
	}
	w.Header.Set("Content-Type", "application/xml")
	body.Reset()
	body.Write(m)
}
