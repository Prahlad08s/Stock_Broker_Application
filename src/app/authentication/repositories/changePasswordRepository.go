package repositories

import (
	genericConstants "stock_broker_application/src/constants"
	"stock_broker_application/src/models"

	"gorm.io/gorm"
)

type ChangePasswordRepository interface {
	CheckUserIDAndPassword(condition map[string]interface{}) bool
	SetNewPassword(PasswordUpdateSQLCondition map[string]interface{}) bool
}

type UserDbRepository struct {
	db *gorm.DB
}

func NewUserDBRepository(dataBase *gorm.DB) *UserDbRepository {
	return &UserDbRepository{db: dataBase}
}

func (userRepository *UserDbRepository) CheckUserIDAndPassword(condition map[string]interface{}) bool {
	var count int64
	if err := userRepository.db.Model(&models.Users{}).Where(condition).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}

func (userRepository *UserDbRepository) SetNewPassword(PasswordUpdateSQLCondition map[string]interface{}) bool {
	var count int64
	if err := userRepository.db.Model(&models.Users{}).Where(genericConstants.Id, PasswordUpdateSQLCondition[genericConstants.Id]).Updates(PasswordUpdateSQLCondition).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}
