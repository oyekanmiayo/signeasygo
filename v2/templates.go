package v2

import (
	"fmt"
	"net/http"
	"signeasygo/hsend"
)

type TemplateService struct {
	hsend *hsend.HSend
}

func newTemplateService(hsend *hsend.HSend) *TemplateService {
	return &TemplateService{
		hsend: hsend.Path("template/"),
	}
}

type ListTemplateResponse []Template

type Template struct {
	CreatedTime  float64          `json:"created_time"`
	Dirty        bool             `json:"dirty"`
	Hash         string           `json:"hash"`
	ID           int64            `json:"id"`
	IsOrdered    bool             `json:"is_ordered"`
	IsOwned      bool             `json:"is_owned"`
	IsPublic     bool             `json:"is_public"`
	IsShared     bool             `json:"is_shared"`
	Link         string           `json:"link"`
	Message      string           `json:"message"`
	Metadata     *MetadataDetails `json:"metadata"`
	ModifiedTime float64          `json:"modified_time"`
	Name         string           `json:"name"`
}

type MetadataDetails struct {
	Fields []MetadataField `json:"fields"`
	Roles  []MetadataRole  `json:"roles"`
}

type MetadataField struct {
	AdditionalInfo string `json:"additional_info"`
	ID             int64  `json:"id"`
	Required       bool   `json:"required"`
	RoleID         int64  `json:"role_id"`
	SubType        string `json:"sub_type"`
	Type           string `json:"type"`
}

type MetadataRole struct {
	Color  string  `json:"color"`
	Fields []int64 `json:"fields"`
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
}

func (t *TemplateService) ListTemplates() (*ListTemplateResponse, *http.Response, error) {
	templates := new(ListTemplateResponse)
	apiError := new(APIError)
	resp, err := t.hsend.New().Get("").Receive(templates, apiError)
	return templates, resp, relevantError(err, *apiError)
}

func (t *TemplateService) GetTemplate(id int32) (*Template, *http.Response, error) {
	template := new(Template)
	apiError := new(APIError)
	resp, err := t.hsend.New().Get(fmt.Sprintf("%v", id)).Receive(template, apiError)
	return template, resp, relevantError(err, *apiError)
}
