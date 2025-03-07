package repository

import (
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type ModuleConfig struct {
	BaseModel
	UserId    uint                   `gorm:"index" json:"userId"`
	Name      string                 `form:"name" gorm:"type:varchar(255)" json:"name"`
	ValueJson string                 `gorm:"type:text" json:"-"`
	Value     map[string]interface{} `gorm:"-" json:"value"`
}

type ModuleConfigRepo struct {
}

type IModuleConfigRepo interface {
	GetModuleConfigByUserIdAndName(userId uint, name string) (map[string]interface{}, error)
	SaveModuleConfig(config *ModuleConfig) error
}

func NewModuleConfigRepo() IModuleConfigRepo {
	return &ModuleConfigRepo{}
}

// retrieves module configuration by user ID and module name
func (r *ModuleConfigRepo) GetModuleConfigByUserIdAndName(userId uint, name string) (map[string]interface{}, error) {
	cfg := ModuleConfig{}
	if err := Db.First(&cfg, "user_id=? AND name=?", userId, name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	// Process JSON field
	if err := json.Unmarshal([]byte(cfg.ValueJson), &cfg.Value); err != nil {
		cfg.Value = nil
	}
	return cfg.Value, nil
}

// saves module configuration to database
func (r *ModuleConfigRepo) SaveModuleConfig(config *ModuleConfig) error {
	// Process JSON field
	if jb, err := json.Marshal(config.Value); err != nil {
		config.ValueJson = "{}"
	} else {
		config.ValueJson = string(jb)
	}

	// Check if record exists
	if err := Db.First(&ModuleConfig{}, "user_id=? AND name=?", config.UserId, config.Name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new record
			if err := Db.Create(config).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		// Update existing record
		if err := Db.Where("user_id=? AND name=?", config.UserId, config.Name).Updates(config).Error; err != nil {
			return err
		}
	}

	return nil
}
