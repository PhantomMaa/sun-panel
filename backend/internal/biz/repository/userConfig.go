package repository

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

	// 搜索引擎
	SearchEngineJson string                 `json:"-"`
	SearchEngine     map[string]interface{} `gorm:"-" json:"searchEngine"`
}
