package errors

type serviceError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (e *serviceError) Error() string {
	return e.Message
}

func NewServiceError(message string, status int) *serviceError {
	return &serviceError{Message: message, Status: status}
}

var (
	RECORD_NOT_FOUND                = NewServiceError("NO_RECORDS_FOUND", 404)
	TRANSACTION_SERIALIZATION_ERROR = NewServiceError("TRANSACTION_SERIALIZATION_ERROR", 500)
	BAD_REQUEST                     = NewServiceError("BAD_REQUEST", 400)
	FAILED_TO_READ_REQUEST_BODY     = NewServiceError("FAILED_TO_READ_REQUEST_BODY", 400)
	UNKNOWN                         = NewServiceError("UNKNOWN", 500)
	RESPONSE_SERIALIZATION_ERROR    = NewServiceError("RESPONSE_SERIALIZATION_ERROR", 400)
	INVALID_UUID                    = NewServiceError("INVALID_UUID", 400)
)
