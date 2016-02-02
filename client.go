package postmanpat

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type Client struct {
	Token       string
	FromAddress string
	client      *http.Client
	BaseURL     *url.URL
	UserAgent   string
}

func NewClient(token string, fromaddr string) *Client {
	bu, _ := url.Parse(DEFAULT_BASE_URL)
	return &Client{
		Token:       token,
		FromAddress: fromaddr,
		client:      http.DefaultClient,
		BaseURL:     bu,
	}
}

func (c *Client) Do(req *http.Request) (*PostmarkResponse, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	pr := &PostmarkResponse{}
	if err = json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}

	if status := resp.StatusCode; status != 200 {
		return pr, errors.New(pr.Message)
	}
	return pr, nil
}

func (c *Client) NewRequest(method string, urlStr string, body interface{}) (*http.Request, error) {
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

func (c *Client) SendTemplate(templateId int, to string, model interface{}) (*PostmarkResponse, error) {
	r := &templateReq{
		TemplateID:    templateId,
		TemplateModel: model,
		From:          c.FromAddress,
		To:            to,
	}
	req, err := c.NewRequest("POST", "/email/withTemplate/", r)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) SendMessage(msg *Message) (*PostmarkResponse, error) {
	pm, err := msg.Prepare()
	pm.From = c.FromAddress
	if err != nil {
		return nil, err
	}
	req, err := c.NewRequest("POST", "/email", pm)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}
