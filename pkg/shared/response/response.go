package response

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/response/status"
)

const RESPONSE_CODE_FORMAT = "%03d%02s%03d"

var AppCode = "11"

func SetAppCode(appCode string) {
	AppCode = appCode
}

type Response[T any] struct {
	HTTPCode  int       `json:"-"`
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	Errors    []error   `json:"errors"`
	Data      T         `json:"data,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func (r *Response[any]) Error() string {
	if len(r.Errors) > 0 {
		return r.Errors[0].Error()
	}
	return r.Message
}

type ErrorResponse struct {
	Field   string      `json:"field"`
	Tag     string      `json:"tag"`
	Value   interface{} `json:"value"`
	Param   string      `json:"param"`
	Message string      `json:"-"`
}

func (v ErrorResponse) Error() string {
	return v.Message
}

// NewSuccess initializes new response.
func NewSuccess(httpCode int, data any) *Response[any] {
	var appStatus *status.AppStatus = status.OK
	// TODO : fix this shit
	if httpCode == 201 {
		appStatus = status.Created
	} else if httpCode == 202 {
		appStatus = status.Accepted
	} else if httpCode == 204 {
		appStatus = status.NoContent
	}

	return &Response[any]{
		HTTPCode:  appStatus.HTTPCode,
		Code:      fmt.Sprintf(RESPONSE_CODE_FORMAT, appStatus.HTTPCode, AppCode, appStatus.CaseCode),
		Message:   appStatus.Error(),
		Data:      data,
		Timestamp: time.Now(),
	}
}

// NewError initializes new error.
func NewError(err error) *Response[any] {
	var appError *status.AppStatus

	// Fallback to internal server error if not specified
	if !errors.As(err, &appError) {
		var fiberError *fiber.Error
		if errors.As(err, &fiberError) {
			appError = status.FromFiberError(fiberError)
		} else {
			appError = status.New(status.InternalServerError, err)
		}
	}

	return &Response[any]{
		HTTPCode:  appError.HTTPCode,
		Code:      fmt.Sprintf(RESPONSE_CODE_FORMAT, appError.HTTPCode, AppCode, appError.CaseCode),
		Message:   appError.Error(),
		Data:      nil,
		Errors:    checkError(appError.Errors),
		Timestamp: time.Now(),
	}
}

// FromAppStatus initializes new response from app status.
func FromAppStatus(errType *status.AppStatus, errs ...error) *Response[any] {
	errType.Errors = errs
	return &Response[any]{
		HTTPCode:  errType.HTTPCode,
		Code:      fmt.Sprintf(RESPONSE_CODE_FORMAT, errType.HTTPCode, AppCode, errType.CaseCode),
		Message:   errType.Error(),
		Data:      nil,
		Errors:    checkError(errType.Errors),
		Timestamp: time.Now(),
	}
}

func checkError(errs []error) []error {
	emptyError := []error{}

	if len(errs) > 0 {
		// Add error type check if new error type want to be returned in response
		var errResp ErrorResponse
		if !errors.As(errs[0], &errResp) {
			return emptyError
		}
		return errs
	}

	return emptyError
}
