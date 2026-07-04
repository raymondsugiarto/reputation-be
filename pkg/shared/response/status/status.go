package status

type ServerErrorCase CaseCode

const (
	ServerErrorCaseInternalServerError ServerErrorCase = 00
	ServerErrorCaseDatabaseError       ServerErrorCase = 01
	ServerErrorCaseCacheError          ServerErrorCase = 02
	ServerErrorCaseEncoderError        ServerErrorCase = 03
	ServerErrorCaseExternalServerError ServerErrorCase = 00
	ServerErrorCaseHTTPClientError     ServerErrorCase = 01
	ServerErrorCaseSMTPError           ServerErrorCase = 02
	ServerErrorCaseGatewayTimeout      ServerErrorCase = 00
)

var statusMap5xx map[int]map[ServerErrorCase][]string = map[int]map[ServerErrorCase][]string{
	500: {
		ServerErrorCaseInternalServerError: {"internal server error", "Unknown internal server error"},
		ServerErrorCaseDatabaseError:       {"database error", "Error in database operation"},
		ServerErrorCaseCacheError:          {"cache error", "Error in cache operation"},
		ServerErrorCaseEncoderError:        {"encoder error", "Error in encode or decode operation"},
	},
	503: {
		ServerErrorCaseExternalServerError: {"service unavailable", "external server error"},
		ServerErrorCaseHTTPClientError:     {"http client error", "Error in http client operation"},
		ServerErrorCaseSMTPError:           {"SMTP error", "Error when using SMTP"},
	},
	504: {
		ServerErrorCaseGatewayTimeout: {"gateway timeout", "Timeout"},
	},
}

func NewServerErrorAppStatus(httpCode int, caseCode ServerErrorCase) *AppStatus {
	test := statusMap5xx[httpCode][caseCode]
	if len(test) == 0 {
		panic("Invalid status code")
	}
	return &AppStatus{
		HTTPCode:    httpCode,
		CaseCode:    int(caseCode),
		Message:     statusMap5xx[httpCode][caseCode][MESSAGE],
		Description: statusMap5xx[httpCode][caseCode][DESCRIPTION],
	}
}

type SuccessCase CaseCode

const (
	SuccessCaseOK        SuccessCase = 00
	SuccessCaseCreated   SuccessCase = 00
	SuccessCaseAccepted  SuccessCase = 00
	SuccessCaseNoContent SuccessCase = 00
)

var statusMap2xx map[int]map[SuccessCase][]string = map[int]map[SuccessCase][]string{
	200: {
		SuccessCaseOK: {"ok", "OK"},
	},
	201: {
		SuccessCaseCreated: {"created", "Resource Created"},
	},
	202: {
		SuccessCaseAccepted: {"accepted", "Accepted"},
	},
	204: {
		SuccessCaseNoContent: {"no content", "No Content"},
	},
}

func NewSuccessAppStatus(httpCode int, caseCode SuccessCase) *AppStatus {
	test := statusMap2xx[httpCode][caseCode]
	if len(test) == 0 {
		panic("Invalid status code")
	}
	return &AppStatus{
		HTTPCode:    httpCode,
		CaseCode:    int(caseCode),
		Message:     statusMap2xx[httpCode][caseCode][MESSAGE],
		Description: statusMap2xx[httpCode][caseCode][DESCRIPTION],
	}
}

type ClientErrorCase CaseCode

const (
	ClientErrorCaseBadRequest                ClientErrorCase = 00
	ClientErrorCaseUnauthorized              ClientErrorCase = 00
	ClientErrorCaseInvalidToken              ClientErrorCase = 01
	ClientErrorCaseInvalidCredential         ClientErrorCase = 02
	ClientErrorCaseMissingOrMalformedToken   ClientErrorCase = 03
	ClientErrorCaseInvalidSession            ClientErrorCase = 04
	ClientErrorCaseMissingAuthenticationData ClientErrorCase = 05
	ClientErrorCaseForbiddenAccess           ClientErrorCase = 00
	ClientErrorCaseRouteNotFound             ClientErrorCase = 00
	ClientErrorCaseEntityNotFound            ClientErrorCase = 01
	ClientErrorMethodNotAllowed              ClientErrorCase = 00
	ClientErrorCaseRequestConflict           ClientErrorCase = 00
	ClientErrorCaseRequestEntityTooLarge     ClientErrorCase = 00
	ClientErrorCaseEntityConflict            ClientErrorCase = 01
	ClientErrorCaseInvalidField              ClientErrorCase = 00
	ClientErrorCaseMissingField              ClientErrorCase = 01
	ClientErrorCaseTooManyRequest            ClientErrorCase = 00
	ClientErrorCaseAttemptTooEarly           ClientErrorCase = 01
	ClientErrorCaseMaxAttemptReached         ClientErrorCase = 02
)

var statusMap4xx map[int]map[ClientErrorCase][]string = map[int]map[ClientErrorCase][]string{
	400: {
		ClientErrorCaseBadRequest: {"bad request", "General request failed error, request message parsing failed."},
	},
	401: {
		ClientErrorCaseUnauthorized:              {"unauthorized", "General unauthorized error"},
		ClientErrorCaseInvalidToken:              {"invalid token", "Token either not exist or invalid"},
		ClientErrorCaseInvalidCredential:         {"invalid credential", "Either username, email or password is invalid"},
		ClientErrorCaseMissingOrMalformedToken:   {"missing or malformed token", "Err missing or malformed JWT"},
		ClientErrorCaseInvalidSession:            {"invalid session", "Session is not registered"},
		ClientErrorCaseMissingAuthenticationData: {"missing authentication data", "Either X-USER-ID, X-ROLE-ID, or X-BRANCH-ID is missing"},
	},
	413: {
		ClientErrorCaseRequestEntityTooLarge: {"request entity too large", "Request entity is too large"},
	},
	403: {
		ClientErrorCaseForbiddenAccess: {"forbidden access", "User is forbidden to access this feature"},
	},
	404: {
		ClientErrorCaseRouteNotFound:  {"route not found", "Route requested is not found"},
		ClientErrorCaseEntityNotFound: {"entity not found", "Entity requested is not found"},
	},
	405: {
		ClientErrorMethodNotAllowed: {"method not allowed", "Method is not allowed for this endpoint"},
	},
	409: {
		ClientErrorCaseRequestConflict: {"request conflict", "Cannot use same X-EXTERNAL-ID in same day"},
		ClientErrorCaseEntityConflict:  {"entity conflict", "Entity in request conflicted with entity in server"},
	},
	422: {
		ClientErrorCaseInvalidField: {"invalid field format", "Invalid format for field %s"},
		ClientErrorCaseMissingField: {"mandatory field missing", "Missing or invalid format on mandatory field %s"},
	},
	429: {
		ClientErrorCaseTooManyRequest:    {"too many requests", "Maximum transaction limit exceeded"},
		ClientErrorCaseAttemptTooEarly:   {"request attempt too early", "Attempt must be delayed for this service. Please wait for %s minutes"},
		ClientErrorCaseMaxAttemptReached: {"maximum attempt reached", "Maximum attempt for this service is reached. Please wait for %s minutes"},
	},
}

func NewClientErrorAppStatus(httpCode int, caseCode ClientErrorCase) *AppStatus {
	test := statusMap4xx[httpCode][caseCode]
	if len(test) == 0 {
		panic("Invalid status code")
	}
	return &AppStatus{
		HTTPCode:    httpCode,
		CaseCode:    int(caseCode),
		Message:     statusMap4xx[httpCode][caseCode][MESSAGE],
		Description: statusMap4xx[httpCode][caseCode][DESCRIPTION],
	}
}
