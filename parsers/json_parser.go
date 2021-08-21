package parsers

import (
	"encoding/json"
	"github.com/Drathveloper/lambda_commons/custom_errors"
)

func BindRequest(body string, bind interface{}) custom_errors.GenericApplicationError {
	bodyBytes := []byte(body)
	unmarshalErr := json.Unmarshal(bodyBytes, &bind)
	if unmarshalErr != nil {
		return custom_errors.NewBadRequestError(unmarshalErr.Error())
	}
	return nil
}

func BindResponse(bind interface{}) (string, custom_errors.GenericApplicationError) {
	bodyBytes, marshalErr := json.Marshal(bind)
	if marshalErr != nil {
		return "", custom_errors.NewGenericInternalServerError()
	}
	return string(bodyBytes), nil
}
