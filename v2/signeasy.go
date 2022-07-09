package v2

import (
	"net/http"
	"signeasygo/hsend"
)

const signeasyV2API = "https://api.signeasy.com/v2.1/"

type Client struct {
	Originals *OriginalService
	Templates *TemplateService
}

func NewClient(client *http.Client, accessToken string) *Client {
	baseHSend := hsend.New().Client(client).Base(signeasyV2API)
	baseHSend.Add("Authorization", "Bearer "+accessToken)

	return &Client{
		Originals: newOriginalService(baseHSend.New()),
		Templates: newTemplateService(baseHSend.New()),
	}
}
