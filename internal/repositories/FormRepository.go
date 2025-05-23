package repositories

import (
	"encoding/json"
	"fmt"
	"gamershub/internal/models"
	"gorm.io/gorm"
)

type FormRepository struct {
	db *gorm.DB
}

func NewFormRepository(db *gorm.DB) *FormRepository {
	return &FormRepository{db: db}
}

func (r *FormRepository) Create(form *models.PlayerForm) error {
	// Сериализация данных в JSON
	if err := serializeFormData(form); err != nil {
		return fmt.Errorf("serialization error: %w", err)
	}

	// Создаем запись в транзакции
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Create(form).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("database create error: %w", err)
	}

	// После создания - десериализуем данные обратно
	if err := deserializeFormData(form); err != nil {
		tx.Rollback()
		return fmt.Errorf("deserialization error: %w", err)
	}

	return tx.Commit().Error
}

func (r *FormRepository) GetByUserID(userID uint) (*models.PlayerForm, error) {
	var form models.PlayerForm
	if err := r.db.Where("user_id = ?", userID).First(&form).Error; err != nil {
		return nil, err
	}

	// Десериализуем JSON данные
	if err := deserializeFormData(&form); err != nil {
		return nil, err
	}

	return &form, nil
}

// Вспомогательные функции для работы с JSON
func serializeFormData(form *models.PlayerForm) error {
	statsJSON, err := json.Marshal(form.Statistics)
	if err != nil {
		return err
	}
	form.StatisticsJSON = string(statsJSON)

	commJSON, err := json.Marshal(form.CommunicationWays)
	if err != nil {
		return err
	}
	form.CommunicationJSON = string(commJSON)

	regionJSON, err := json.Marshal(form.Region)
	if err != nil {
		return err
	}
	form.RegionJSON = string(regionJSON)

	return nil
}

func deserializeFormData(form *models.PlayerForm) error {
	if err := json.Unmarshal([]byte(form.StatisticsJSON), &form.Statistics); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(form.CommunicationJSON), &form.CommunicationWays); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(form.RegionJSON), &form.Region); err != nil {
		return err
	}

	return nil
}

func (r *FormRepository) Delete(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.PlayerForm{}).Error
}
