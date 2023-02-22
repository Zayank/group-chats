package helpers

import (
	"net/http"
)

type APIErrors struct {
	Errors []*APIError `json:"errors"`
}

func (errors *APIErrors) Status() int {
	return errors.Errors[0].Status
}

type APIError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Title   string `json:"title"`
	Details string `json:"details"`
}

func newAPIError(status int, code string, title string, details string) *APIError {
	return &APIError{
		Status:  status,
		Code:    code,
		Title:   title,
		Details: details,
	}
}

var (
	ErrUnauthorzied = newAPIError(http.StatusUnauthorized, "unauthorized", "Unauthorized Error", "Please login to proceed")
	ErrForbidden    = newAPIError(http.StatusForbidden, "forbidden", "Forbidden Error", "access denied")

	ErrBadRequest = newAPIError(http.StatusBadRequest, "bad_request", "Forbidden Error", "")
	ErrDB         = newAPIError(http.StatusInternalServerError, "internal_error", "Internal Error", "try again later")
	ErrInternal   = newAPIError(http.StatusInternalServerError, "internal_error", "Internal Error", "try again later")

	SuccessStatusOk      = newAPIError(http.StatusOK, "success", "Success", "")
	SuccessStatusCreated = newAPIError(http.StatusCreated, "created", "Created", "")
)

func DynamicError(status int64, message string) *APIError {
	switch status {
	case http.StatusBadRequest:
		errorMessage := ErrBadRequest
		errorMessage.Details = message
		return errorMessage
	default:
		return ErrInternal
	}
}

func DynamicSuccessMessage(status int64, message string) *APIError {
	var response *APIError

	switch status {
	case http.StatusOK:
		response = SuccessStatusOk

	case http.StatusCreated:
		response = SuccessStatusCreated

	default:
		response = ErrInternal
	}

	response.Details = message

	return response
}
