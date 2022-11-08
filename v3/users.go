package v3

import (
	"net/http"
	"signeasygo/hsend"
)

type UserService struct {
	hsend *hsend.HSend
}

func newUserService(hsend *hsend.HSend) *UserService {
	return &UserService{
		hsend: hsend.Path(""),
	}
}

type User struct {
	IsAutoRenew               int         `json:"is_auto_renew"`
	LastName                  string      `json:"last_name"`
	RsAllowed                 bool        `json:"rs_allowed"`
	CreatedTime               int         `json:"created_time"`
	BillingCycle              string      `json:"billing_cycle"`
	SignedFileCount           int         `json:"signed_file_count"`
	DocumentCredits           int         `json:"document_credits"`
	Id                        int         `json:"id"`
	IsSocial                  int         `json:"is_social"`
	FirstName                 string      `json:"first_name"`
	ReferrerCode              string      `json:"referrer_code"`
	RemainingRsCount          interface{} `json:"remaining_rs_count"`
	ReferrerPromoLeft         int         `json:"referrer_promo_left"`
	ImportedFileCount         int         `json:"imported_file_count"`
	RemainingTemplateCount    interface{} `json:"remaining_template_count"`
	CardBrand                 interface{} `json:"card_brand"`
	SubscriptionStatus        interface{} `json:"subscription_status"`
	IsTemplateCreationAllowed bool        `json:"is_template_creation_allowed"`
	ActivationTime            int         `json:"activation_time"`
	Email                     string      `json:"email"`
	Status                    int         `json:"status"`
	AccountType               int         `json:"account_type"`
	ImportStatus              bool        `json:"import_status"`
	Company                   interface{} `json:"company"`
	PlanSubtype               int         `json:"plan_subtype"`
	CardLast4                 interface{} `json:"card_last4"`
	CompanySize               interface{} `json:"company_size"`
	ReferrerPromoUsed         int         `json:"referrer_promo_used"`
	EmailVerified             int         `json:"email_verified"`
	IsPaid                    int         `json:"is_paid"`
	LastModifiedTime          int         `json:"last_modified_time"`
	SubscriptionExpiryTime    int         `json:"subscription_expiry_time"`
	EligibleForYearlyPlan     bool        `json:"eligible_for_yearly_plan"`
	BilledBy                  interface{} `json:"billed_by"`
}

type UserResponse struct {
	User
}

func (u *UserService) FetchSelf() (*UserResponse, *http.Response, error) {
	user := new(UserResponse)
	apiError := new(APIError)
	resp, err := u.hsend.New().Get("me/").Receive(user, apiError)
	return user, resp, relevantError(err, *apiError)
}
