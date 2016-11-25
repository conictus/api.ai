package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	MessageTypeText          MessageType = 0
	MessageTypeCard          MessageType = 1
	MessageTypeQuickReply    MessageType = 2
	MessageTypeImage         MessageType = 3
	MessageTypeCustomPayload MessageType = 4
)

type Query struct {
	Query     string `json:"query"`
	SessionID string `json:"sessionId"`
	Lang      string `json:"lang"`
}

type MessageType int

type QueryResponse struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	SessionID string    `json:"sessionId"`
	Result    struct {
		Source           string          `json:"source"`
		ResolvedQuery    string          `json:"resolvedQuery"`
		Action           string          `json:"actin"`
		ActionIncomplete bool            `json:"actionIncomplete"`
		Parameters       json.RawMessage `json:"parameters"`
		Contexts         []Context       `json:"contexts"`
		MetaData         struct {
			IntentID    string `json:"intentId"`
			WebHookUsed string `json:"webhookUsed"`
			IntentName  string `json:"intentName"`
		} `json:"metadata"`
		Fulfillment struct {
			Speech   string `json:"speech"`
			Messages []struct {
				Type   int    `json:"type"`
				Speech string `json:"speech"`
			} `json:"messages"`
		} `json:"fulfillment"`
		Score float64 `json:"score"`
	} `json:"result"`
	Status struct {
		Code         int    `json:"code"`
		ErrorType    string `json:"errorType"`
		ErrorDetails string `json:"errorDetails"`
	} `json:"status"`
}

func (r *QueryResponse) String() string {
	data, _ := json.MarshalIndent(r, "", " ")
	return string(data)
}

func (r *QueryResponse) Context(name string) (*Context, error) {
	for _, ctx := range r.Result.Contexts {
		if ctx.Name == name {
			return &ctx, nil
		}
	}

	return nil, fmt.Errorf("not found")
}

func (c *Client) Query(query Query) (*QueryResponse, error) {
	if query.Lang == "" {
		query.Lang = "en"
	}

	var body bytes.Buffer
	enc := json.NewEncoder(&body)
	if err := enc.Encode(query); err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, c.url(QueryEndpoint), &body)

	response, err := c.do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if err := c.error(response); err != nil {
		return nil, err
	}

	var result QueryResponse
	dec := json.NewDecoder(response.Body)
	if err := dec.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
