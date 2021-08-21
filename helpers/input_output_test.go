package helpers_test

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"lambda_commons/custom_errors"
	"lambda_commons/helpers"
	"math"
	"testing"
)

func TestMapErrorToAPIGatewayProxyResponse_ShouldSucceed(t *testing.T) {
	customError := custom_errors.NewGenericInternalServerError()
	expected := events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       `{"message":"internal server error"}`,
	}
	actual := helpers.MapErrorToAPIGatewayProxyResponse(customError)
	assert.Equal(t, expected, actual)
}

func TestMapResponseToAPIGatewayProxyResponseWithHeaders_ShouldSucceed(t *testing.T) {
	headers := make(map[string]string, 0)
	headers["someHeader"] = "someValue"
	expectedResponse := events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       `"xx"`,
		Headers:    headers,
	}
	actualResponse := helpers.MapResponseToAPIGatewayProxyResponseWithHeaders(201, "xx", headers)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestMapResponseToAPIGatewayProxyResponseWithHeaders_ShouldReturnInternalServerErrorWhenFailedParsingBody(t *testing.T) {
	headers := make(map[string]string, 0)
	expectedResponse := events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       `{"message":"internal server error"}`,
		Headers:    headers,
	}
	actualResponse := helpers.MapResponseToAPIGatewayProxyResponseWithHeaders(201, math.Inf(1), headers)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestMapResponseToAPIGatewayProxyResponse_ShouldSucceed(t *testing.T) {
	headers := make(map[string]string, 0)
	expectedResponse := events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       `"xx"`,
		Headers:    headers,
	}
	actualResponse := helpers.MapResponseToAPIGatewayProxyResponse(201, "xx")
	assert.Equal(t, expectedResponse, actualResponse)
}