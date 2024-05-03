package repositories

import (
	dbModels "stock_broker_application/src/models"

	"gorm.io/gorm"
)

type ForgotPasswordRepository interface {
	VerifyAndUpdatePassword(email string, pancardNumber string, newPassword string) error
}

type userDBRepository struct {
	DB *gorm.DB
}

func NewForgotPasswordRepository(db *gorm.DB) ForgotPasswordRepository {
	return &userDBRepository{DB: db}
}

func (repository *userDBRepository) VerifyAndUpdatePassword(email string, pancardNumber string, newPassword string) error {
	// Update the password for the user if found
	result := repository.DB.Model(&dbModels.Users{}).
		Where("email = ? AND pan_card = ?", email, pancardNumber).
		Update("password", newPassword)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		// No user found with the provided credentials
		return gorm.ErrRecordNotFound
	}

	return nil
}
