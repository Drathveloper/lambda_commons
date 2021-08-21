package models

import "lambda_commons/custom_errors"

type ValidatedRequest interface {
	Validate() custom_errors.GenericApplicationError
}
