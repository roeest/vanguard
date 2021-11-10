package vanguard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	apiHost         = "api.vanguard.com"
	refererTemplate = "https://investor.vanguard.com/etf/profile/overview/%s"
)

type client struct {
	*http.Client
	format string
	header http.Header
	debug  bool
}

type Option func(c *client)

func DebugOption(c *client) {
	c.debug = true
}

func NewClient(opts ...Option) *client {
	c := &client{Client: http.DefaultClient, format: "json", header: make(http.Header, 1)}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (c *client) GetEtf(symbol string) (*Etf, error) {
	c.header.Set("Referer", "https://api.vanguard.com")
	etf, err := newEtf(c, symbol)
	if err != nil {
		return nil, err
	}
	return etf, nil
}

func (c *client) generateUrl(symbol, resource, q string) string {
	u := url.URL{
		Scheme:   "https",
		Host:     apiHost,
		Path:     fmt.Sprintf("rs/ire/01/ind/fund/%s/%s.%s", symbol, resource, c.format),
		RawQuery: q,
	}
	return u.String()
}

func (c *client) getResource(symbol, resource string, result interface{}) error {
	return c.getResourceWithQueryParams(symbol, resource, result, "")
}

func (c *client) getResourceWithQueryParams(symbol, resource string, result interface{}, queryParams string) error {
	u := c.generateUrl(symbol, resource, queryParams)
	r, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	r.Header = c.header
	resp, err := c.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if c.debug {
		fmt.Println(string(respData))
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got non ok status %d - body: %s", resp.StatusCode, string(respData))
	}
	return json.Unmarshal(respData, result)
}
