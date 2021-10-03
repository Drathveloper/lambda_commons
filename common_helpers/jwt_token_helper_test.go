package common_helpers_test

import (
	"github.com/Drathveloper/lambda_commons/common_errors"
	"github.com/Drathveloper/lambda_commons/common_helpers"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type JwtHelperTestSuite struct {
	suite.Suite
	jwtHelper common_helpers.JwtHelper
}

func TestJwtHelperTestSuite(t *testing.T) {
	suite.Run(t, new(JwtHelperTestSuite))
}

func (suite *JwtHelperTestSuite) SetupTest() {
	suite.jwtHelper = common_helpers.NewJwtHelper("-----BEGIN RSA PRIVATE KEY-----\n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\n-----END RSA PRIVATE KEY-----")
}

func (suite *JwtHelperTestSuite) TestNewJwtHelperShouldPanicWhenInvalidRsaPrivateKeyLoaded() {
	suite.Assert().Panics(func() {
		common_helpers.NewJwtHelper("someInvalidPrivateKey")
	})
}

func (suite *JwtHelperTestSuite) TestGenerateJwtTokenShouldSucceed() {
	claims := jwt.MapClaims{
		"someClaim": "someValue",
	}

	token, err := suite.jwtHelper.GenerateJwtToken(claims)

	suite.NoError(err)
	suite.Assert().NotEmpty(token)
}

func (suite *JwtHelperTestSuite) TestValidateJwtTokenShouldSucceed() {
	expectedClaims := jwt.MapClaims{
		"someClaim": "someValue",
	}
	jwtToken, _ := suite.jwtHelper.GenerateJwtToken(expectedClaims)

	claims, err := suite.jwtHelper.ValidateJwtToken(jwtToken)

	suite.NoError(err)
	suite.Assert().Equal(expectedClaims, claims)
}

func (suite *JwtHelperTestSuite) TestValidateJwtTokenShouldReturnErrorWhenTokenIsNotValid() {
	jwtToken := "invalidJwtToken"
	expectedErr := common_errors.NewUnauthorizedError("given jwt token is not valid or is expired")

	_, err := suite.jwtHelper.ValidateJwtToken(jwtToken)

	suite.Assert().Equal(expectedErr, err)
}
