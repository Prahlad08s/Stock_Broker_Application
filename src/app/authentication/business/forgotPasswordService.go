package business

import (
	"authentication/models"
	"authentication/repositories"
)

type ForgotPasswordService interface {
	UpdatePassword(credentials models.ForgotPasswordRequest) error
}

type forgotPasswordService struct {
	userCredentialsInterface repositories.ForgotPasswordRepository
}

func NewForgotPasswordService(userData repositories.ForgotPasswordRepository) ForgotPasswordService {
	return &forgotPasswordService{
		userCredentialsInterface: userData,
	}
}

func (service *forgotPasswordService) UpdatePassword(credentials models.ForgotPasswordRequest) error {
	return service.userCredentialsInterface.VerifyAndUpdatePassword(credentials.Email, credentials.PanCardNumber, credentials.NewPassword)
}
