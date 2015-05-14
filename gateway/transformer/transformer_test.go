package transformer

import (
	"bytes"
	"net/http"

	. "gopkg.in/check.v1"
)

func (s *S) TestGetTransformer(c *C) {
	c.Check(s.transformers.Get("invalid"), IsNil)
}

func (s *S) TestAddTransformer(c *C) {
	c.Check(s.transformers.Get("AddHeader"), IsNil)
	ah := func(r *http.Request, w *http.Response, buf *bytes.Buffer) {}
	s.transformers.Add("AddHeader", ah)
	c.Check(s.transformers.Get("AddHeader"), NotNil)
}

func (s *S) TestConvertXmlToJson(c *C) {
	req := &http.Request{Header: make(http.Header)}
	resp := &http.Response{Header: make(http.Header)}
	body := bytes.NewBuffer([]byte(`<root><name>Alice</name><list><item>1</item><item>2</item></list></root>`))
	ConvertXmlToJson(req, resp, body)
	c.Assert(body.String(), Equals, `{"root":{"list":{"item":["1","2"]},"name":"Alice"}}`)
	c.Assert(resp.Header.Get("Content-Type"), Equals, "application/json")
}

func (s *S) TestConvertJsonToXml(c *C) {
	req := &http.Request{Header: make(http.Header)}
	resp := &http.Response{Header: make(http.Header)}
	body := bytes.NewBuffer([]byte(`{"root":{"list":{"item":["1","2"]},"name":"Alice"}}`))
	ConvertJsonToXml(req, resp, body)
	c.Assert(body.String(), Equals, `<root><list><item>1</item><item>2</item></list><name>Alice</name></root>`)
	c.Assert(resp.Header.Get("Content-Type"), Equals, "application/xml")
}
