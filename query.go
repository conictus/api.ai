package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

type Event struct {
	Name string            `json:"name"`
	Data map[string]string `json:"data"`
}

type Query struct {
	Query     string `json:"query,omitempty"`
	Event     *Event `json:"event,omitempty"`
	SessionID string `json:"sessionId"`
	Lang      string `json:"lang"`
}

type Button struct {
	Text     string `json:"text"`
	PostBack string `json:"postback"`
}

type MessageType int

type Message struct {
	Type     MessageType `json:"type"`
	Speech   string      `json:"speech,omitempty"`
	Platform string      `json:"platform"`
	ImageURL string      `json:"imageUrl,omitempty"`

	Title    string   `json:"title,omitempty"`
	Subtitle string   `json:"subtitle,omitempty"`
	Buttons  []Button `json:"buttons,omitempty"`

	Replies []string `json:"replies,omitempty"`

	Payload json.RawMessage `json:"payload,omitempty"`
}

type Response struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	SessionID string    `json:"sessionId"`
	Result    struct {
		Source           string          `json:"source"`
		ResolvedQuery    string          `json:"resolvedQuery"`
		Action           string          `json:"action"`
		ActionIncomplete bool            `json:"actionIncomplete"`
		Parameters       json.RawMessage `json:"parameters"`
		Contexts         []Context       `json:"contexts"`
		MetaData         struct {
			IntentID    string `json:"intentId"`
			WebHookUsed string `json:"webhookUsed"`
			IntentName  string `json:"intentName"`
		} `json:"metadata"`
		Fulfillment struct {
			Speech   string    `json:"speech"`
			Messages []Message `json:"messages"`
		} `json:"fulfillment"`
		Score float64 `json:"score"`
	} `json:"result"`
	Status struct {
		Code         int    `json:"code"`
		ErrorType    string `json:"errorType"`
		ErrorDetails string `json:"errorDetails"`
	} `json:"status"`
}

func (r *Response) ToJson(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func (r *Response) String() string {
	data, _ := json.MarshalIndent(r, "", " ")
	return string(data)
}

func (r *Response) DialogContext(name string) *Context {
	return r.Context(fmt.Sprintf("%s_dialog_context", name))
}

func (r *Response) Context(name string) *Context {
	for _, ctx := range r.Result.Contexts {
		if ctx.Name == name {
			return &ctx
		}
	}

	return nil
}

func (c *Client) Query(query Query) (*Response, error) {
	if query.Lang == "" {
		query.Lang = "en"
	}

	var body bytes.Buffer
	enc := json.NewEncoder(&body)
	if err := enc.Encode(query); err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, c.url(QueryEndpoint).String(), &body)

	response, err := c.do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if err := c.error(response); err != nil {
		return nil, err
	}

	var result Response
	dec := json.NewDecoder(response.Body)
	if err := dec.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
