package ai

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
)

const (
	ApiVersion = "20150910"
	ApiURL     = "https://api.api.ai/v1"

	QueryEndpoint    = "query"
	ContextsEndpoint = "contexts"
)

var (
	ErrContextNotFound = fmt.Errorf("not found")
)

type url string

func (u url) param(name, value string) url {
	return url(fmt.Sprintf("%s&%s=%s", u, name, value))
}

func (u url) String() string {
	return string(u)
}

type Client struct {
	version string
	key     string
	cl      *http.Client
}

func New(key string) *Client {
	return &Client{
		key:     key,
		version: ApiVersion,
		cl:      &http.Client{},
	}
}

func (c *Client) url(endpoint ...string) url {
	p := path.Join(endpoint...)
	return url(fmt.Sprintf("%s/%s?v=%s", ApiURL, p, ApiVersion))
}

func (c *Client) do(request *http.Request) (*http.Response, error) {
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.key))
	request.Header.Set("Content-Type", "application/json")

	return c.cl.Do(request)
}

func (c *Client) error(response *http.Response) error {
	if response.StatusCode != http.StatusOK {
		resp, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf("API.AI [%d] %s: %s", response.StatusCode, response.Status, string(resp))
	}
	return nil
}
