package repositories

import (
	genericConstants "stock_broker_application/src/constants"
	dbModels "stock_broker_application/src/models"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type ForgotPasswordRepository interface {
	VerifyAndUpdatePassword(email string, panCardNumber string, newPassword string) error
}

type userDBRepository struct {
	DB *gorm.DB
}

// Mock repository for testing
type MockForgotPasswordRepository struct {
	mock.Mock
}

func (m *MockForgotPasswordRepository) VerifyAndUpdatePassword(email string, panCardNumber string, newPassword string) error {
	m.Called(email, panCardNumber, newPassword)
	return nil
}

func NewForgotPasswordRepository(db *gorm.DB) ForgotPasswordRepository {
	return &userDBRepository{DB: db}
}

func (repository *userDBRepository) VerifyAndUpdatePassword(email string, pancardNumber string, newPassword string) error {
	result := repository.DB.Model(&dbModels.Users{}).
		Where("email = ? AND pan_card = ?", email, pancardNumber).
		Update(genericConstants.Password, newPassword)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
