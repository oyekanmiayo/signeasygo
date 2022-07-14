package v2

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
	resp, err := rs.hsend.New().BodyJSON(bodyParams).Receive(rsResp, apiError)
	return rsResp, resp, relevantError(err, *apiError)
}
