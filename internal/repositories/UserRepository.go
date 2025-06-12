package repositories

import (
	"errors"
	"fmt"
	"gamershub/internal/models"
	"gamershub/internal/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser godoc
// @Summary      Создать нового пользователя
// @Description  Добавить пользователя в систему
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user  body      UserRequest  true  "Данные пользователя"
// @Success      201   {object}  UserResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/v1/auth/register [post]
func (repo *UserRepository) CreateUser(user *models.User) error {
	var existing models.User
	result := repo.db.Where("email = ? OR username = ?",
		user.Email,
		user.Username).First(&existing)

	if result.Error == nil {
		if existing.Email == user.Email {
			return fmt.Errorf("email already exists")
		}
		if existing.Username == user.Username && user.Username != "" {
			return fmt.Errorf("username already exists")
		}
	}

	if err := repo.db.Create(user).Error; err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	return nil
}

func (repo *UserRepository) FindUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := repo.db.First(&user, id).Error; err != nil {
		return nil, errors.New("failed to find user by id")
	}
	return &user, nil
}

func (repo *UserRepository) FindUserByEmail(email types.Email) (*models.User, error) {
	var user models.User
	if err := repo.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("failed to find user by email")
	}
	return &user, nil
}

func (repo *UserRepository) FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := repo.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("failed to find user by username")
	}
	return &user, nil
}

func (repo *UserRepository) UpdateUser(user *models.User) error {
	return repo.db.Save(user).Error
}

func (repo *UserRepository) DeleteUser(id uint) error {
	if err := repo.db.Delete(&models.User{}, id).Error; err != nil {
		return errors.New("не удалось удалить пользователя")
	}
	return nil
}

func (r *UserRepository) FindByRefreshToken(token string) (*models.User, error) {
	var user models.User
	err := r.db.Where("refresh_token IS NOT NULL").First(&user).Error
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.RefreshToken), []byte(token)); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) RevokeRefreshToken(userId uint) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userId).
		Updates(map[string]interface{}{
			"refresh_token": nil,
			"token_expiry":  nil,
		}).Error
}
