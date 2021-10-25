package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PreHandler interface {
	PreHandler(body []byte) ([]byte, error)
}

type Response struct {
	Code int
	Body []byte
}

func (s Response) String() string {
	resp, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(resp)
}

func (s Response) Error() string {
	return s.String()
}

type request struct {
	tls         *tls.Config
	request     *http.Request
	preHandlers PreHandler
	timeout     time.Duration
}

func (c *client) Post(body []byte) *request {
	buffer := bytes.NewBuffer(body)
	url := fmt.Sprintf("%s?%s", c.url, c.urlValues.Encode())
	req, _ := http.NewRequest("POST", url, buffer)

	req.Header = c.headers.Clone()
	return &request{request: req}
}

func (c *client) Get() *request {
	url := fmt.Sprintf("%s?%s", c.url, c.urlValues.Encode())
	req, _ := http.NewRequest("GET", url, nil)
	req.Header = c.headers.Clone()
	return &request{request: req}
}

func (c *client) Delete() *request {
	url := fmt.Sprintf("%s?%s", c.url, c.urlValues.Encode())
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header = c.headers.Clone()
	return &request{request: req}
}

func (c *client) Put(body []byte) *request {
	buffer := bytes.NewBuffer(body)
	url := fmt.Sprintf("%s?%s", c.url, c.urlValues.Encode())
	req, _ := http.NewRequest("PUT", url, buffer)
	req.Header = c.headers.Clone()
	return &request{request: req}
}

func (r *request) TLS(tls *tls.Config) *request {
	r.tls = tls
	return r
}

func (r *request) AddPreHandler(handler PreHandler) *request {
	r.preHandlers = handler
	return r
}
func (r *request) SetTimeout(sec time.Duration) *request {
	r.timeout = sec * time.Second
	return r
}

func (r *request) Do() ([]byte, error) {
	c := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: r.tls,
		},
		Timeout: r.timeout,
	}
	response, err := c.Do(r.request)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 400 {
		return nil, Response{
			Code: response.StatusCode,
			Body: buf.Bytes(),
		}
	}

	return buf.Bytes(), nil
}
