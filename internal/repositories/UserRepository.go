package repositories

import (
	"errors"
	"gamershub/internal/models"
	"gamershub/internal/types"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// create new user
func (repo *UserRepository) CreateUser(user *models.User) error {
	if err := repo.db.Create(user).Error; err != nil {
		return errors.New("failed to create user")
	}
	return nil
}

// find user by id
func (repo *UserRepository) FindUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := repo.db.First(&user, id).Error; err != nil {
		return nil, errors.New("failed to find user by id")
	}
	return &user, nil
}

// find user by email
func (repo *UserRepository) FindUserByEmail(email types.Email) (*models.User, error) {
	var user models.User
	if err := repo.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("failed to find user by email")
	}
	return &user, nil
}

// find user by username
func (repo *UserRepository) FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := repo.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("failed to find user by username")
	}
	return &user, nil
}

// update user
func (repo *UserRepository) UpdateUser(user *models.User) error {
	if err := repo.db.Save(user).Error; err != nil {
		return errors.New("failed to update user")
	}
	return nil
}

// delete user
func (repo *UserRepository) DeleteUser(id uint) error {
	if err := repo.db.Delete(&models.User{}, id).Error; err != nil {
		return errors.New("не удалось удалить пользователя")
	}
	return nil
}

// TODO: team-up methods
