package v3

import (
	"net/http"
	"signeasygo/hsend"
)

type RequestSignatureService struct {
	hsend *hsend.HSend
}

func newRequestSignatureService(hsend *hsend.HSend) *RequestSignatureService {
	return &RequestSignatureService{
		hsend: hsend.Path("template/"),
	}
}

type RequestSignatureWithoutMarkersBodyParams struct {
	OriginalFileId  int32        `json:"original_file_id"`
	Recipients      []Recipient  `json:"recipients"`
	CC              []CarbonCopy `json:"cc"`
	Message         string       `json:"message"`
	Name            string       `json:"name"`
	EmbeddedSigning bool         `json:"embedded_signing"`
	IsOrdered       bool         `json:"is_ordered"`
}

type Recipient struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CarbonCopy struct {
	Email string `json:"email"`
}

type RequestSignatureWithoutMarkersResponse struct {
	PendingFileId string `json:"pending_file_id"`
}

func (rs *RequestSignatureService) RequestSignatureWithoutMarkers(bodyParams *RequestSignatureWithoutMarkersBodyParams) (*RequestSignatureWithoutMarkersResponse, *http.Response, error) {
	rsResp := new(RequestSignatureWithoutMarkersResponse)
	apiError := new(APIError)
	resp, err := rs.hsend.New().Post("").BodyJSON(bodyParams).
		Receive(rsResp, apiError)
	return rsResp, resp, relevantError(err, *apiError)
}

type ListSignatureRequestsWithoutMarkersResponse struct {
	Count int64                  `json:"count"`
	Files []SignatureRequestFile `json:"files"`
}

type SignatureRequestFile struct {
	AadhaarEnabled   *string                         `json:"aadhaar_enabled"`
	CreatedTime      int64                           `json:"created_time"`
	HasMarkers       int32                           `json:"has_markers"`
	ID               int64                           `json:"id"`
	IsInPerson       int32                           `json:"is_in_person"`
	IsOrdered        int32                           `json:"is_ordered"`
	LastModifiedTime int64                           `json:"last_modified_time"`
	Logo             *string                         `json:"logo"`
	Name             string                          `json:"name"`
	Next             int64                           `json:"next"`
	OwnerCompany     string                          `json:"owner_company"`
	OwnerEmail       string                          `json:"owner_email"`
	OwnerFirstName   string                          `json:"owner_first_name"`
	OwnerLastName    string                          `json:"owner_last_name"`
	OwnerUserID      string                          `json:"owner_user_id"`
	Recipients       []SignatureRequestFileRecipient `json:"recipients"`
}

type SignatureRequestFileRecipient struct {
	CreatedTime      int64  `json:"created_time"`
	Email            string `json:"email"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	LastModifiedTime int64  `json:"last_modified_time"`
	OrderID          int64  `json:"order_id"`
	RecipientID      int64  `json:"recipient_id"`
	RecipientUserID  int64  `json:"recipient_user_id"`
	Status           string `json:"status"`
}

func (rs *RequestSignatureService) ListSignatureRequestsWithoutMarkers() (*ListSignatureRequestsWithoutMarkersResponse, *http.Response, error) {
	rsResp := new(ListSignatureRequestsWithoutMarkersResponse)
	apiError := new(APIError)
	resp, err := rs.hsend.New().Get("").Receive(rsResp, apiError)
	return rsResp, resp, relevantError(err, *apiError)
}
