package v2

import "fmt"

// APIError
// https://docs.signeasy.com/docs/error-codes
type APIError struct {
	Status  int32  `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("status: %v, message: %v", e.Status, e.Message)
}

// relevantError returns any http-related error if it exists
// if http-related errors don't exist, it returns apiError if it exists
// else it returns nil
func relevantError(httpError error, apiError APIError) error {
	if httpError != nil {
		return httpError
	}

	if (APIError{}) == apiError {
		return nil
	}

	return apiError
}
