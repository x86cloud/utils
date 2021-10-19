package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
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

func (r *request) Do(model interface{}) (err error) {
	// check weather
	rv := reflect.ValueOf(model)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &json.InvalidUnmarshalError{Type: reflect.TypeOf(model)}
	}

	c := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: r.tls,
		},
		Timeout: r.timeout,
	}
	response, err := c.Do(r.request)
	if err != nil {
		return err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode >= 400 {
		return Response{
			Code: response.StatusCode,
			Body: buf.Bytes(),
		}
	}

	//自定义数据预处理方法
	if r.preHandlers != nil {
		data, err := r.preHandlers.PreHandler(buf.Bytes())
		if err != nil {
			return err
		}
		if err = json.Unmarshal(data, model); err != nil {
			return err
		}
	} else {
		if err = json.Unmarshal(buf.Bytes(), model); err != nil {
			return err
		}
	}
	return nil
}
