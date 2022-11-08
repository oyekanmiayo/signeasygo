package v3

import (
	"net/http"
	"signeasygo/hsend"
)

const signeasyV2API = "https://api.signeasy.com/v2.1/" /* Unstable */
const signeasyV3API = "https://api.signeasy.com/v3/"   /* Stable */

type Client struct {
	Originals        *OriginalService
	Templates        *TemplateService
	RequestSignature *RequestSignatureService
	Embedded         *EmbeddedSelfSignService
	Users            *UserService
}

func NewClient(client *http.Client, accessToken string) *Client {
	baseHSend := hsend.New().Client(client).Base(signeasyV3API)
	baseHSend.Add("Authorization", "Bearer "+accessToken)

	return &Client{
		Originals:        newOriginalService(baseHSend.New()),
		Templates:        newTemplateService(baseHSend.New()),
		RequestSignature: newRequestSignatureService(baseHSend.New()),
		Embedded:         newEmbeddedSelfSignService(baseHSend.New()),
		Users:            newUserService(baseHSend.New()),
	}
}
