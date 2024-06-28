package handler

import (
	"authentication/business"
	"authentication/commons/constants"
	"authentication/models"
	"encoding/json"
	"net/http"
	genericConstants "stock_broker_application/src/constants"
	genericModel "stock_broker_application/src/models"
	"stock_broker_application/src/utils"
	"stock_broker_application/src/utils/validations"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ForgotPasswordController interface {
	HandleForgotPassword(ctx *gin.Context)
}

type forgotPasswordController struct {
	service business.ForgotPasswordService
}

func NewForgotPasswordController(service business.ForgotPasswordService) ForgotPasswordController {
	return &forgotPasswordController{
		service: service,
	}
}

// HandleForgotPassword updates the user's password.
func (controller *forgotPasswordController) HandleForgotPassword(context *gin.Context) {
	var request models.ForgotPasswordRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		var errorMsgs genericModel.ErrorMessage
		if unmarshalTypeError, ok := err.(*json.UnmarshalTypeError); ok {
			errorMsgs = genericModel.ErrorMessage{Key: unmarshalTypeError.Field, ErrorMessage: genericConstants.JsonBindingFieldError}
		} else {
			errorMsgs = genericModel.ErrorMessage{Key: "error", ErrorMessage: err.Error()}
		}
		utils.SendBadRequest(context, []genericModel.ErrorMessage{errorMsgs})
		return
	}

	if err := validations.GetCustomValidator(context.Request.Context()).Struct(request); err != nil {
		validationErrors := validations.FormatValidationErrors(context.Request.Context(), err.(validator.ValidationErrors))
		context.JSON(http.StatusBadRequest, gin.H{
			genericConstants.GenericJSONErrorMessage: genericConstants.ValidatorError,
			genericConstants.GenericValidationError:  validationErrors,
		})
		return
	}

	err := controller.service.UpdatePassword(request)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			context.JSON(http.StatusNotFound, gin.H{genericConstants.GenericJSONErrorMessage: constants.ErrorInvalidUserData})
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{genericConstants.GenericJSONErrorMessage: constants.ErrorInvalidUserData})
		}
		return
	}

	context.JSON(http.StatusOK, gin.H{genericConstants.GenericJSONMessage: constants.ForgotPasswordSuccessMessage})
}
