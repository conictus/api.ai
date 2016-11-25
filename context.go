package ai

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Context struct {
	Name       string          `json:"name"`
	Parameters json.RawMessage `json:"parameters"`
	Lifespan   int             `json:"lifespan,omitempty"`
}

func (ctx *Context) LoadTo(params interface{}) error {
	return json.Unmarshal(ctx.Parameters, &params)
}

func (c *Client) Context(sessionId string, name string) (*Context, error) {
	request, err := http.NewRequest(http.MethodGet,
		c.url(ContextsEndpoint, name).param("sessionId", sessionId).String(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	response, err := c.do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if err := c.error(response); err != nil {
		return nil, err
	}

	var context Context
	dec := json.NewDecoder(response.Body)
	if err := dec.Decode(&context); err != nil {
		return nil, err
	}

	return &context, nil
}

func (c *Client) SetContext(sessionId, name string, parameters interface{}) error {
	context := struct {
		Name       string      `json:"name"`
		Parameters interface{} `json:"parameters"`
	}{
		Name:       name,
		Parameters: parameters,
	}

	var body bytes.Buffer
	enc := json.NewEncoder(&body)
	if err := enc.Encode(&context); err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost,
		c.url(ContextsEndpoint).param("sessionId", sessionId).String(),
		&body,
	)

	if err != nil {
		return err
	}

	response, err := c.do(request)
	defer response.Body.Close()

	if err != nil {
		return err
	}

	if err := c.error(response); err != nil {
		return err
	}

	return nil
}
