package common_parsers

import (
	"encoding/json"
	"github.com/Drathveloper/lambda_commons/v2/common_errors"
)

func BindRequest(body string, bind interface{}) common_errors.GenericApplicationError {
	bodyBytes := []byte(body)
	unmarshalErr := json.Unmarshal(bodyBytes, &bind)
	if unmarshalErr != nil {
		return common_errors.NewBadRequestError(unmarshalErr.Error())
	}
	return nil
}

func BindResponse(bind interface{}) (string, common_errors.GenericApplicationError) {
	bodyBytes, marshalErr := json.Marshal(bind)
	if marshalErr != nil {
		return "", common_errors.NewGenericInternalServerError()
	}
	return string(bodyBytes), nil
}
