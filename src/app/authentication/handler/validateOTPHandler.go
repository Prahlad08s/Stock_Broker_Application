package handler

import (
	"authentication/commons/constants"
	"authentication/models"
	"encoding/json"
	genericConstants "stock_broker_application/src/constants"
	genericModel "stock_broker_application/src/models"
	"stock_broker_application/src/utils"

	"github.com/gin-gonic/gin"
)

type OTPService interface {
	OtpVerification(otpData models.ValidateOTPRequest) error
	GenerateAndStoreToken(tokenData genericModel.TokenData, userID uint16) (string, error)
}

type OTPValidationController struct {
	Service OTPService
}

func NewOTPValidationController(service OTPService) *OTPValidationController {
	return &OTPValidationController{
		Service: service,
	}
}

// NewValidateOTPHandler handles the OTP validation request
// @Summary Validate OTP
// @Description Validates the OTP for a user
// @Tags Validate OTP
// @Accept json
// @Produce json
// @Param otpRequest body models.ValidateOTPRequest true "OTP Request"
// @Success 200 {string} string "OTP validated successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "OTP is expired or invalid"
// @Router /v1/auth/validate-otp [post]
func (controller *OTPValidationController) HandleValidateOTP(context *gin.Context) {
	var otpValidationRequest models.ValidateOTPRequest
	if err := context.ShouldBindJSON(&otpValidationRequest); err != nil {
		errorMsgs := genericModel.ErrorMessage{Key: err.(*json.UnmarshalTypeError).Field, ErrorMessage: genericConstants.JsonBindingFieldError}
		utils.SendBadRequest(context, []genericModel.ErrorMessage{errorMsgs})
		return
	}
	if err := controller.Service.OtpVerification(otpValidationRequest); err != nil {
		utils.SendUnauthorizedError(context, err.Error())
		return
	}
	tokenData := genericModel.TokenData{
		UserId: otpValidationRequest.UserID,
	}
	token, err := controller.Service.GenerateAndStoreToken(tokenData, otpValidationRequest.UserID)
	if err != nil {
		utils.SendInternalServerError(context, constants.ErrorGenToken)

		return
	}

	responseModel := genericModel.JWTTokens{
		AccessToken: token,
	}

	// context.Set(genericConstants.GenericTokenMessage, token)
	utils.SendStatusOk(context, responseModel)

}
