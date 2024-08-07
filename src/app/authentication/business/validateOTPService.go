package business

import (
	"authentication/commons/constants"
	"authentication/models"
	"authentication/repositories"
	"errors"
	genericConstants "stock_broker_application/src/constants"
	genericModel "stock_broker_application/src/models"
	"stock_broker_application/src/utils/authorization"
)

type OTPService struct {
	UserRepository repositories.UserRepository
}

func NewOTPService(userRepository repositories.UserRepository) *OTPService {
	return &OTPService{
		UserRepository: userRepository,
	}
}

func (service *OTPService) OtpVerification(otpData models.ValidateOTPRequest) error {
	if otpData.UserID == 0 {
		return errors.New(constants.ErrorRequiredUserID)
	}
	if otpData.OTP < 1000 || otpData.OTP > 9999 {
		return errors.New(constants.ErrorInvalidOtpFormat)
	}

	userExists, err := service.UserRepository.CheckUserExists(otpData.UserID)
	if err != nil {
		return err
	}
	if !userExists {
		return errors.New(constants.ErrorInvalidUserID)
	}

	isValid, err := service.UserRepository.CheckOtp(otpData.UserID, otpData.OTP)
	if err != nil {
		return err
	}
	if !isValid {
		return errors.New(constants.ErrorInvalidOtpFormat) // This line was updated
	}
	return nil
}

func (service *OTPService) GenerateAndStoreToken(tokenData genericModel.TokenData, userID uint16) (string, error) {
	token, err := authorization.GenerateJWTToken(tokenData)
	if err != nil {
		return genericConstants.EmptySpace, err
	}

	if err := service.UserRepository.UpdateUserToken(userID, token); err != nil {
		return genericConstants.EmptySpace, err
	}

	return token, nil
}
