package transformer

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/clbanning/mxj/j2x"
	"github.com/clbanning/mxj/x2j"
)

// Function which modifies the response.
type Transformer func(*http.Request, *http.Response, *bytes.Buffer)

// An array of Transformer with key to be used by the gateway and service.
type Transformers map[string]Transformer

func (f Transformers) Add(key string, value Transformer) {
	f[key] = value
}

func (f Transformers) Get(key string) Transformer {
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
