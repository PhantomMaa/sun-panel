package repository

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"sun-panel/internal/common"
)

type PanelConfig struct {
	BackgroundImageSrc               string `json:"backgroundImageSrc,omitempty"`
	BackgroundBlur                   *int   `json:"backgroundBlur,omitempty"`
	BackgroundMaskNumber             *int   `json:"backgroundMaskNumber,omitempty"`
	IconStyle                        *int   `json:"iconStyle,omitempty"`
	IconTextColor                    string `json:"iconTextColor,omitempty"`
	IconTextInfoHideDescription      *bool  `json:"iconTextInfoHideDescription,omitempty"`
	IconTextIconHideTitle            *bool  `json:"iconTextIconHideTitle,omitempty"`
	LogoText                         string `json:"logoText,omitempty"`
	LogoImageSrc                     string `json:"logoImageSrc,omitempty"`
	ClockShowSecond                  *bool  `json:"clockShowSecond,omitempty"`
	ClockColor                       string `json:"clockColor,omitempty"`
	SearchBoxShow                    *bool  `json:"searchBoxShow,omitempty"`
	SearchBoxSearchIcon              *bool  `json:"searchBoxSearchIcon,omitempty"`
	MarginTop                        *int   `json:"marginTop,omitempty"`
	MarginBottom                     *int   `json:"marginBottom,omitempty"`
	MaxWidth                         *int   `json:"maxWidth,omitempty"`
	MaxWidthUnit                     string `json:"maxWidthUnit"`
	MarginX                          *int   `json:"marginX,omitempty"`
	FooterHtml                       string `json:"footerHtml,omitempty"`
	SystemMonitorShow                *bool  `json:"systemMonitorShow,omitempty"`
	SystemMonitorShowTitle           *bool  `json:"systemMonitorShowTitle,omitempty"`
	SystemMonitorPublicVisitModeShow *bool  `json:"systemMonitorPublicVisitModeShow,omitempty"`
	NetModeChangeButtonShow          *bool  `json:"netModeChangeButtonShow,omitempty"`
}

type UserConfig struct {
	UserId uint `gorm:"index" json:"userId"`

	// 面板样式数据
	PanelJson string       `json:"-"`
	Panel     *PanelConfig `gorm:"-" json:"panel"`
}

// retrieves user configuration from database by user ID
func GetUserConfig(userId uint) (UserConfig, error) {
	cfg := UserConfig{}
	if err := Db.First(&cfg, "user_id=?", userId).Error; err != nil {
		return cfg, err
	}

	// Process JSON fields
	if err := json.Unmarshal([]byte(cfg.PanelJson), &cfg.Panel); err != nil {
		cfg.Panel = nil
	}

	return cfg, nil
}

// saves user configuration to database
// It will create a new record if not exists, or update existing one
func SaveUserConfig(config *UserConfig) error {
	// Process JSON fields
	config.PanelJson = common.ToJSONString(config.Panel)

	// Check if record exists
	if err := Db.First(&UserConfig{}, "user_id=?", config.UserId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new record
			if err := Db.Create(config).Error; err != nil {
				return err
			}
		} else {
			// Database error
			return err
		}
	} else {
		// Update existing record
		if err := Db.Where("user_id=?", config.UserId).Updates(config).Error; err != nil {
			return err
		}
	}

	return nil
}
