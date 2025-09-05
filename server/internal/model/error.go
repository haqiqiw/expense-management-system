package model

import "net/http"

var (
	ErrInternalServerError        = NewCustomError(http.StatusInternalServerError, 100, "Internal server error")
	ErrUnauthorized               = NewCustomError(http.StatusUnauthorized, 101, "Unauthorized")
	ErrBadRequest                 = NewCustomError(http.StatusBadRequest, 102, "Bad request")
	ErrForbidden                  = NewCustomError(http.StatusForbidden, 103, "Forbidden")
	ErrMissingOrInvalidAuthHeader = NewCustomError(http.StatusUnauthorized, 104, "Missing or invalid auth header")
	ErrInvalidAuthToken           = NewCustomError(http.StatusUnauthorized, 105, "Invalid auth token")
	ErrTokenRevoked               = NewCustomError(http.StatusUnauthorized, 106, "Token revoked")

	ErrEmailAlreadyExist       = NewCustomError(http.StatusBadRequest, 1000, "Email already exist")
	ErrUserNotFound            = NewCustomError(http.StatusNotFound, 1001, "User not found")
	ErrInvalidPassword         = NewCustomError(http.StatusUnauthorized, 1002, "Invalid password")
	ErrExpenseNotFound         = NewCustomError(http.StatusNotFound, 1003, "Expense not found")
	ErrExpenseMinAmount        = NewCustomError(http.StatusBadRequest, 1004, "Amount can't be less than Rp 10.000")
	ErrExpenseMaxAmount        = NewCustomError(http.StatusBadRequest, 1005, "Amount can't be greater than Rp 50.000.000")
	ErrExpenseAlreadyProcessed = NewCustomError(http.StatusUnprocessableEntity, 1006, "Expense already processed")
)

type ErrorItem struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CustomError struct {
	HTTPStatus int         `json:"http_status"`
	Errors     []ErrorItem `json:"errors"`
}

func (c *CustomError) Append(ei ErrorItem) {
	c.Errors = append(c.Errors, ei)
}

func (c *CustomError) Error() string {
	if len(c.Errors) == 0 {
		return "empty error"
	}

	return c.Errors[0].Message
}

func NewCustomError(httpStatus int, code int, msg string) *CustomError {
	return &CustomError{
		HTTPStatus: httpStatus,
		Errors: []ErrorItem{
			{
				Code:    code,
				Message: msg,
			},
		},
	}
}
