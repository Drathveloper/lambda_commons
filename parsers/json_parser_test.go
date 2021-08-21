package parsers_test

import (
	"github.com/stretchr/testify/assert"
	"lambda_commons/custom_errors"
	"lambda_commons/parsers"
	"math"
	"testing"
)

type TestBindModel2 struct {
	SomeX string   `json:"someX"`
	SomeY int      `json:"someY"`
	SomeZ []string `json:"someZ"`
}

type TestBindModel struct {
	Some1 string         `json:"some1"`
	Some2 int            `json:"some2"`
	Some3 TestBindModel2 `json:"some3"`
}

func TestBindRequest_ShouldSucceed(t *testing.T) {
	requestBody := `{
		"some1": "someValue",
		"some2": 1,
		"some3": {
			"someX": "xyz",
            "someY": 2,
			"someZ": ["1", "2", "3"]
		}
	}`
	var model TestBindModel
	actual := parsers.BindRequest(requestBody, &model)
	assert.Nil(t, actual)
	assert.Equal(t, "someValue", model.Some1)
	assert.Equal(t, 1, model.Some2)
	assert.Equal(t, "xyz", model.Some3.SomeX)
	assert.Equal(t, 2, model.Some3.SomeY)
	assert.Equal(t, []string{"1", "2", "3"}, model.Some3.SomeZ)
}

func TestBindRequest_ShouldReturnBadRequestErrorWhenUnmarshalFailed(t *testing.T) {
	requestBody := `{
		"some1": "someValue",
		"some2": 1,
		"some3": {
			"someX": "xyz",
            "someY": 2,
			"someZ": ["1", "2", "3"]
	}`
	var model TestBindModel
	expected := custom_errors.NewBadRequestError("unexpected end of JSON input")
	actual := parsers.BindRequest(requestBody, &model)
	assert.Equal(t, expected, actual)
}

func TestBindResponse_ShouldSucceed(t *testing.T) {
	model := TestBindModel{
		Some1: "x",
		Some2: 5,
		Some3: TestBindModel2{
			SomeX: "xx",
			SomeY: 9,
			SomeZ: []string{"6", "7", "8"},
		},
	}
	expectedModel := `{"some1":"x","some2":5,"some3":{"someX":"xx","someY":9,"someZ":["6","7","8"]}}`
	parsedModel, appErr := parsers.BindResponse(model)
	assert.Nil(t, appErr)
	assert.Equal(t, expectedModel, parsedModel)
}

func TestBindResponse_ShouldReturnInternalServerErrorWhenMarshalFailed(t *testing.T) {
	expected := custom_errors.NewInternalServerError("internal server error")
	_, appErr := parsers.BindResponse(math.Inf(1))
	assert.Equal(t, expected, appErr)
}
