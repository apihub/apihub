package filter

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/clbanning/mxj/j2x"
	"github.com/clbanning/mxj/x2j"
)

type Filter func(*http.Request, *http.Response, *bytes.Buffer)
type Filters map[string]Filter

func (f Filters) Add(key string, value Filter) {
	f[key] = value
}

func (f Filters) Get(key string) Filter {
	return f[key]
}

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
