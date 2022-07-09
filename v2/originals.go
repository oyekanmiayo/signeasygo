package v2

import (
	"bytes"
	"net/http"
	"signeasygo/hsend"
)

type OriginalService struct {
	hsend *hsend.HSend
}

func newOriginalService(hsend *hsend.HSend) *OriginalService {
	return &OriginalService{
		hsend: hsend.Path("original/"),
	}
}

/*
ImportDocumentBodyParams
In order to abide by a security measure to the file import/upload API,
we're removing the following characters in the file name at our end /, :, *, <, >, |.
Please ensure that these characters are not present in the "name" property in the file
import/upload API.

The payload should contain values for "file", "name" and "rename_if_exists"
*/
type ImportDocumentBodyParams struct {
	Payload              *bytes.Buffer
	MultipartContentType string
}

type ImportDocumentResponse struct {
	CreatedTime      int64  `json:"created_time"`
	ID               int64  `json:"id"`
	LastModifiedTime int64  `json:"last_modified_time"`
	Name             string `json:"name"`
}

func (o *OriginalService) ImportDocument(bodyParams *ImportDocumentBodyParams) (*ImportDocumentResponse, *http.Response, error) {
	response := new(ImportDocumentResponse)
	apiError := new(APIError)
	resp, err := o.hsend.New().Post("").
		BodyMultiPartForm(bodyParams.Payload, bodyParams.MultipartContentType).
		Receive(response, apiError)
	return response, resp, relevantError(err, *apiError)
}

type ListOriginalsResponse struct {
	Count int64          `json:"count"`
	Files []FileMetadata `json:"files"`
}

type FileMetadata struct {
	CreatedTime      int64  `json:"created_time"`
	ID               int64  `json:"id"`
	LastModifiedTime int64  `json:"last_modified_time"`
	Name             string `json:"name"`
}

func (o *OriginalService) ListOriginals() (*ListOriginalsResponse, *http.Response, error) {
	lor := new(ListOriginalsResponse)
	apiError := new(APIError)
	httpResp, httpErr := o.hsend.New().Get("").Receive(lor, apiError)
	return lor, httpResp, relevantError(httpErr, *apiError)
}
