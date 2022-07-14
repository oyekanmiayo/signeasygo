package v2

import (
	"fmt"
	"net/http"
	"signeasygo/hsend"
)

type EmbeddedService struct {
	hsend *hsend.HSend
}

func newEmbeddedService(hsend *hsend.HSend) *EmbeddedService {
	return &EmbeddedService{
		hsend: hsend.Path("me/"),
	}
}

type FetchSelfSignURLBodyParam struct {
	FileID  int32  `json:"file_id"`
	Message string `json:"message"`

	// URL that the user would be redirected to, once they sign the document.
	// The pending_file_id of the signature requests will be added as a query parameter.
	RedirectURL string `json:"redirect_url"`
}

type FetchSelfSignURLResponse struct {
	FileID int32  `json:"file_id"`
	URL    string `json:"url"`
}

// FetchSelfSignURL is for https://docs.signeasy.com/reference/fetch-embedded-self-signing-url
func (e *EmbeddedService) FetchSelfSignURL(bodyParams *FetchSelfSignURLBodyParam) (*FetchSelfSignURLResponse, *http.Response, error) {
	response := new(FetchSelfSignURLResponse)
	apiError := new(APIError)
	httpResp, httpErr := e.hsend.New().Post("embedded/url/").BodyJSON(bodyParams).Receive(response, apiError)
	return response, httpResp, relevantError(httpErr, *apiError)
}

type FetchSelfSignedFilesResponse struct {
	Count int64            `json:"count"`
	Files []SelfSignedFile `json:"files"`
}

type SelfSignedFile struct {
	// Should I use string or *string to allow nil?
	AadhaarEnabled   string `json:"aadhaar_enabled"`
	CheckSum         string `json:"check_sum"`
	CreatedTime      int64  `json:"created_time"`
	FaAadhaarEnabled string `json:"fa_aadhaar_enabled"`
	ID               int64  `json:"id"`
	LastModifiedTime int64  `json:"last_modified_time"`
	Name             string `json:"name"`
	// PendingFile <unknown_type>
	PublicIdentifier string `json:"public_identifier"`
}

// FetchSelfSignedFiles is for https://docs.signeasy.com/reference/fetch-all-self-signed-files
func (e *EmbeddedService) FetchSelfSignedFiles() (*FetchSelfSignedFilesResponse, *http.Response, error) {
	response := new(FetchSelfSignedFilesResponse)
	apiError := new(APIError)
	httpResp, httpErr := e.hsend.New().Get("signed/").Receive(response, apiError)
	return response, httpResp, relevantError(httpErr, *apiError)
}

// FetchSelfSignedFile is for https://docs.signeasy.com/reference/get-self-signed-document-details
func (e *EmbeddedService) FetchSelfSignedFile(signedID int32) (*SelfSignedFile, *http.Response, error) {
	response := new(SelfSignedFile)
	apiError := new(APIError)
	httpResp, httpErr := e.hsend.New().Get(fmt.Sprintf("signed/%v", signedID)).
		Receive(response, apiError)
	return response, httpResp, relevantError(httpErr, *apiError)
}

// DownloadSelfSignedFile is for https://docs.signeasy.com/reference/download-self-signed-document
// How do I download a pdf document? I need to test this and provide an example
func (e *EmbeddedService) DownloadSelfSignedFile(signedID int32) (interface{}, *http.Response, error) {
	var response interface{}
	apiError := new(APIError)
	httpResp, httpErr := e.hsend.New().Get(fmt.Sprintf("signed/%v/download", signedID)).
		Receive(response, apiError)
	return response, httpResp, relevantError(httpErr, *apiError)
}

// DownloadSelfSignedFileCertificate is for https://docs.signeasy.com/reference/download-certificate-of-self-signed-document
func (e *EmbeddedService) DownloadSelfSignedFileCertificate(signedID int32) (interface{}, *http.Response, error) {
	var response interface{}
	apiError := new(APIError)
	httpResp, httpErr := e.hsend.New().Get(fmt.Sprintf("signed/%v/certificate", signedID)).
		Receive(response, apiError)
	return response, httpResp, relevantError(httpErr, *apiError)
}

// DeleteSelfSignedFile is for https://docs.signeasy.com/reference/delete-self-signed-document
// Random thought: I need consistent name throughout this lib. Should we use file or document everywhere?
// The documentation mixes and matches :)
func (e *EmbeddedService) DeleteSelfSignedFile(signedID int32) (*http.Response, error) {
	apiError := new(APIError)
	httpResp, httpErr := e.hsend.New().Delete(fmt.Sprintf("signed/%v", signedID)).
		Receive(nil, apiError)
	return httpResp, relevantError(httpErr, *apiError)
}
