package common_models

import "github.com/Drathveloper/lambda_commons/common_errors"

type ValidatedRequest interface {
	Validate() common_errors.GenericApplicationError
}
