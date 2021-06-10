package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

type PreHandler interface {
	PreHandler(body []byte) ([]byte, error)
}

type request struct {
	request     *http.Request
	preHandlers PreHandler
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

func (r *request) AddPreHandler(handler PreHandler) *request {
	r.preHandlers = handler
	return r
}

func (r *request) Do(model interface{}) (err error) {
	// check wether
	rv := reflect.ValueOf(model)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &json.InvalidUnmarshalError{Type: reflect.TypeOf(model)}
	}

	c := http.Client{}
	response, err := c.Do(r.request)
	if err != nil {
		return err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(response.Body)
	if response.StatusCode >= 400 {
		return fmt.Errorf(buf.String())
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
