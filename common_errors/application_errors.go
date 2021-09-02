package common_errors

type GenericApplicationError interface {
	Error() string
	HttpStatus() int
}

type genericApplicationError struct {
	httpStatus int
	message    string
}

func (error *genericApplicationError) Error() string {
	return error.message
}

func (error *genericApplicationError) HttpStatus() int {
	return error.httpStatus
}

func newGenericError(httpStatus int, message string) GenericApplicationError {
	return &genericApplicationError{
		httpStatus: httpStatus,
		message:    message,
	}
}

func NewBadRequestError(message string) GenericApplicationError {
	return newGenericError(400, message)
}

func NewUnauthorizedError(message string) GenericApplicationError {
	return newGenericError(401, message)
}

func NewForbiddenError(message string) GenericApplicationError {
	return newGenericError(403, message)
}

func NewNotFoundError(message string) GenericApplicationError {
	return newGenericError(404, message)
}

func NewPreconditionFailedError(message string) GenericApplicationError {
	return newGenericError(412, message)
}

func NewInternalServerError(message string) GenericApplicationError {
	return newGenericError(500, message)
}

func NewGenericBadRequestError() GenericApplicationError {
	return newGenericError(400, "bad request")
}

func NewGenericUnauthorizedError() GenericApplicationError {
	return newGenericError(401, "unauthorized")
}

func NewGenericForbiddenError() GenericApplicationError {
	return newGenericError(403, "forbidden")
}

func NewGenericNotFoundError() GenericApplicationError {
	return newGenericError(404, "not found")
}

func NewGenericPreconditionFailedError() GenericApplicationError {
	return newGenericError(412, "precondition failed")
}

func NewGenericInternalServerError() GenericApplicationError {
	return newGenericError(500, "internal server error")
}
