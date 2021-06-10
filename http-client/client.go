package client

import (
	"net/http"
	pkgurl "net/url"
)

type client struct {
	urlValues *pkgurl.Values
	headers   *http.Header
	url       string
}

func New(url string) *client {
	return &client{
		url:       url,
		headers:   &http.Header{},
		urlValues: &pkgurl.Values{},
	}
}

func (c *client) AddHeader(key, value string) *client {
	c.headers.Add(key, value)
	return c
}

func (c *client) AddQuery(key, value string) *client {
	c.urlValues.Add(key, value)
	return c
}
