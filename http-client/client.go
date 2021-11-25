package client

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

type request struct {
	urlValues *url.Values
	headers   *http.Header
	url       string

	tls     *tls.Config
	request *http.Request
	timeout time.Duration
}

func New(uri string) *request {
	query, _ := url.ParseQuery(uri)
	return &request{
		url:       uri,
		headers:   &http.Header{},
		urlValues: &query,
	}
}

func (c *request) AddHeader(key, value string) *request {
	c.headers.Add(key, value)
	return c
}

func (c *request) AddHeaders(headers map[string]string) *request {
	for key, value := range headers {
		c.headers.Add(key, value)
	}
	return c
}

func (c *request) AddQuery(key, value string) *request {
	c.urlValues.Add(key, value)
	return c
}

func (c *request) AddQueries(queries map[string]string) *request {
	for key, value := range queries {
		c.urlValues.Add(key, value)
	}
	return c
}
