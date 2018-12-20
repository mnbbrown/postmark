package postmark

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// Client is a postmark client
type Client struct {
	Token       string
	FromAddress string
	client      *http.Client
	BaseURL     *url.URL
	UserAgent   string
}

// NewClient creates a new postmark client
func NewClient(token string, fromaddr string) *Client {
	bu, _ := url.Parse("https://api.postmarkapp.com")
	return &Client{
		Token:       token,
		FromAddress: fromaddr,
		client:      http.DefaultClient,
		BaseURL:     bu,
	}
}

func (c *Client) do(req *http.Request) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	pr := &Response{}
	if err = json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}

	if status := resp.StatusCode; status != 200 {
		return pr, errors.New(pr.Message)
	}
	return pr, nil
}

func (c *Client) newRequest(method string, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Postmark-Server-Token", c.Token)
	return req, nil
}

type templateReq struct {
	TemplateID    int
	TemplateModel interface{}
	From          string
	To            string
}

// SendTemplate sends a postmark template
func (c *Client) SendTemplate(templateID int, to string, model interface{}) (*Response, error) {
	r := &templateReq{
		TemplateID:    templateID,
		TemplateModel: model,
		From:          c.FromAddress,
		To:            to,
	}
	req, err := c.newRequest("POST", "/email/withTemplate/", r)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

// SendMessage sends a postmark message
func (c *Client) SendMessage(msg *Message) (*Response, error) {
	pm, err := newPostmarkMessage(msg)
	pm.From = c.FromAddress
	if err != nil {
		return nil, err
	}
	req, err := c.newRequest("POST", "/email", pm)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}
