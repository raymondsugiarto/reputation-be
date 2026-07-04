package status

import (
	"github.com/gofiber/fiber/v3"
)

type CaseCode int

var mergedStatusMap = mergeMaps()

type AppStatus struct {
	Message     string
	Description string
	CaseCode    int
	HTTPCode    int
	Errors      []error
}

func (a *AppStatus) Error() string {
	if len(a.Errors) > 0 {
		return a.Errors[0].Error()
	}
	return a.Message
}

const (
	MESSAGE     = 0
	DESCRIPTION = 1
)

const (
	DEFAULT_CASE_CODE = 0
)

func FromFiberError(fiberErr *fiber.Error) *AppStatus {
	return &AppStatus{
		HTTPCode:    fiberErr.Code,
		CaseCode:    DEFAULT_CASE_CODE,
		Message:     mergedStatusMap[fiberErr.Code][DEFAULT_CASE_CODE][MESSAGE],
		Description: mergedStatusMap[fiberErr.Code][DEFAULT_CASE_CODE][DESCRIPTION],
		Errors:      []error{fiberErr},
	}
}

var (
	DefaultAppStatus             *AppStatus = &AppStatus{}
	OK                           *AppStatus = NewSuccessAppStatus(fiber.StatusOK, SuccessCaseOK)
	Created                      *AppStatus = NewSuccessAppStatus(fiber.StatusCreated, SuccessCaseCreated)
	Accepted                     *AppStatus = NewSuccessAppStatus(fiber.StatusAccepted, SuccessCaseAccepted)
	NoContent                    *AppStatus = NewSuccessAppStatus(fiber.StatusNoContent, SuccessCaseNoContent)
	BadRequest                   *AppStatus = NewClientErrorAppStatus(fiber.StatusBadRequest, ClientErrorCaseBadRequest)
	Unauthorized                 *AppStatus = NewClientErrorAppStatus(fiber.StatusUnauthorized, ClientErrorCaseUnauthorized)
	InvalidToken                 *AppStatus = NewClientErrorAppStatus(fiber.StatusUnauthorized, ClientErrorCaseInvalidToken)
	InvalidCredential            *AppStatus = NewClientErrorAppStatus(fiber.StatusUnauthorized, ClientErrorCaseInvalidCredential)
	ErrMissingOrMalformedJWT     *AppStatus = NewClientErrorAppStatus(fiber.StatusUnauthorized, ClientErrorCaseMissingOrMalformedToken)
	InvalidSession               *AppStatus = NewClientErrorAppStatus(fiber.StatusUnauthorized, ClientErrorCaseInvalidSession)
	ErrMissingAuthenticationData *AppStatus = NewClientErrorAppStatus(fiber.StatusUnauthorized, ClientErrorCaseMissingAuthenticationData)
	Forbidden                    *AppStatus = NewClientErrorAppStatus(fiber.StatusForbidden, ClientErrorCaseForbiddenAccess)
	RequestNotFound              *AppStatus = NewClientErrorAppStatus(fiber.StatusNotFound, ClientErrorCaseRouteNotFound)
	EntityNotFound               *AppStatus = NewClientErrorAppStatus(fiber.StatusNotFound, ClientErrorCaseEntityNotFound)
	RequestConflict              *AppStatus = NewClientErrorAppStatus(fiber.StatusConflict, ClientErrorCaseRequestConflict)
	RequestEntityTooLarge        *AppStatus = NewClientErrorAppStatus(fiber.StatusRequestEntityTooLarge, ClientErrorCaseRequestEntityTooLarge)
	EntityConflict               *AppStatus = NewClientErrorAppStatus(fiber.StatusConflict, ClientErrorCaseEntityConflict)
	InvalidFieldFormat           *AppStatus = NewClientErrorAppStatus(fiber.StatusUnprocessableEntity, ClientErrorCaseInvalidField)
	MandatoryFieldMissing        *AppStatus = NewClientErrorAppStatus(fiber.StatusUnprocessableEntity, ClientErrorCaseMissingField)
	TooManyRequest               *AppStatus = NewClientErrorAppStatus(fiber.StatusTooManyRequests, ClientErrorCaseTooManyRequest)
	AttemptTooEarly              *AppStatus = NewClientErrorAppStatus(fiber.StatusTooManyRequests, ClientErrorCaseAttemptTooEarly)
	MaxAttemptReached            *AppStatus = NewClientErrorAppStatus(fiber.StatusTooManyRequests, ClientErrorCaseMaxAttemptReached)
	InternalServerError          *AppStatus = NewServerErrorAppStatus(fiber.StatusInternalServerError, ServerErrorCaseInternalServerError)
	DatabaseError                *AppStatus = NewServerErrorAppStatus(fiber.StatusInternalServerError, ServerErrorCaseDatabaseError)
	CacheError                   *AppStatus = NewServerErrorAppStatus(fiber.StatusInternalServerError, ServerErrorCaseCacheError)
	EncoderError                 *AppStatus = NewServerErrorAppStatus(fiber.StatusInternalServerError, ServerErrorCaseEncoderError)
	ExternalServerError          *AppStatus = NewServerErrorAppStatus(fiber.StatusServiceUnavailable, ServerErrorCaseExternalServerError)
	HTTPClientError              *AppStatus = NewServerErrorAppStatus(fiber.StatusServiceUnavailable, ServerErrorCaseHTTPClientError)
	SMTPError                    *AppStatus = NewServerErrorAppStatus(fiber.StatusServiceUnavailable, ServerErrorCaseSMTPError)
	GatewayTimeout               *AppStatus = NewServerErrorAppStatus(fiber.StatusGatewayTimeout, ServerErrorCaseGatewayTimeout)
)

func New(appType *AppStatus, errs ...error) *AppStatus {
	if len(errs) > 0 && errs[0] == appType {
		return appType
	}
	appType.Errors = errs
	return appType
}

func NewErrorAppStatus(httpCode int, caseCode CaseCode) *AppStatus {
	test := mergedStatusMap[httpCode][caseCode]
	if len(test) == 0 {
		panic("Invalid status code")
	}
	return &AppStatus{
		HTTPCode:    httpCode,
		CaseCode:    int(caseCode),
		Message:     mergedStatusMap[httpCode][caseCode][MESSAGE],
		Description: mergedStatusMap[httpCode][caseCode][DESCRIPTION],
	}
}

func mergeMaps() map[int]map[CaseCode][]string {
	merged := make(map[int]map[CaseCode][]string)
	for k := range statusMap2xx {
		if _, ok := merged[k]; !ok {
			merged[k] = make(map[CaseCode][]string)
		}
		for k2, v2 := range statusMap2xx[k] {
			merged[k][CaseCode(k2)] = v2
		}
	}
	for k := range statusMap4xx {
		if _, ok := merged[k]; !ok {
			merged[k] = make(map[CaseCode][]string)
		}
		for k2, v2 := range statusMap4xx[k] {
			merged[k][CaseCode(k2)] = v2
		}
	}
	for k := range statusMap5xx {
		if _, ok := merged[k]; !ok {
			merged[k] = make(map[CaseCode][]string)
		}
		for k2, v2 := range statusMap5xx[k] {
			merged[k][CaseCode(k2)] = v2
		}
	}
	return merged
}

func MergeStatusMap(errorMap map[int]map[CaseCode][]string) {
	for k := range errorMap {
		if _, ok := mergedStatusMap[k]; !ok {
			mergedStatusMap[k] = make(map[CaseCode][]string)
		}
		for k2, v2 := range errorMap[k] {
			mergedStatusMap[k][CaseCode(k2)] = v2
		}
	}
}
