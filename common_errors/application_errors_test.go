package common_errors_test

import (
	"github.com/Drathveloper/lambda_commons/v2/common_errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBadRequestError(t *testing.T) {
	actual := common_errors.NewBadRequestError("someErr")
	assert.Equal(t, "someErr", actual.Error())
	assert.Equal(t, 400, actual.HttpStatus())
}

func TestNewUnauthorizedError(t *testing.T) {
	actual := common_errors.NewUnauthorizedError("someErr")
	assert.Equal(t, "someErr", actual.Error())
	assert.Equal(t, 401, actual.HttpStatus())
}

func TestNewForbiddenError(t *testing.T) {
	actual := common_errors.NewForbiddenError("someErr")
	assert.Equal(t, "someErr", actual.Error())
	assert.Equal(t, 403, actual.HttpStatus())
}

func TestNewNotFoundError(t *testing.T) {
	actual := common_errors.NewNotFoundError("someErr")
	assert.Equal(t, "someErr", actual.Error())
	assert.Equal(t, 404, actual.HttpStatus())
}

func TestNewPreconditionFailedError(t *testing.T) {
	actual := common_errors.NewPreconditionFailedError("someErr")
	assert.Equal(t, "someErr", actual.Error())
	assert.Equal(t, 412, actual.HttpStatus())
}

func TestNewInternalServerError(t *testing.T) {
	actual := common_errors.NewInternalServerError("someErr")
	assert.Equal(t, "someErr", actual.Error())
	assert.Equal(t, 500, actual.HttpStatus())
}

func TestNewGenericBadRequestError(t *testing.T) {
	actual := common_errors.NewGenericBadRequestError()
	assert.Equal(t, "bad request", actual.Error())
	assert.Equal(t, 400, actual.HttpStatus())
}

func TestNewGenericUnauthorizedError(t *testing.T) {
	actual := common_errors.NewGenericUnauthorizedError()
	assert.Equal(t, "unauthorized", actual.Error())
	assert.Equal(t, 401, actual.HttpStatus())
}

func TestNewGenericForbiddenError(t *testing.T) {
	actual := common_errors.NewGenericForbiddenError()
	assert.Equal(t, "forbidden", actual.Error())
	assert.Equal(t, 403, actual.HttpStatus())
}

func TestNewGenericNotFoundError(t *testing.T) {
	actual := common_errors.NewGenericNotFoundError()
	assert.Equal(t, "not found", actual.Error())
	assert.Equal(t, 404, actual.HttpStatus())
}

func TestNewGenericPreconditionFailedError(t *testing.T) {
	actual := common_errors.NewGenericPreconditionFailedError()
	assert.Equal(t, "precondition failed", actual.Error())
	assert.Equal(t, 412, actual.HttpStatus())
}

func TestNewGenericInternalServerError(t *testing.T) {
	actual := common_errors.NewGenericInternalServerError()
	assert.Equal(t, "internal server error", actual.Error())
	assert.Equal(t, 500, actual.HttpStatus())
}
