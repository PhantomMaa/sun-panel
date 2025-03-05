package repository

type Notice struct {
	BaseModel
	Title       string `gorm:"type:varchar(255)" json:"title"`
	Content     string `gorm:"type:varchar(2000)" json:"content"`
	DisplayType int    `gorm:"type:tinyint(1)" json:"displayType"` // 展示类型 参考常量：NOTICE_DISPLAY_TYPE_XXXXX
	OneRead     int    `gorm:"type:tinyint(1)" json:"oneRead"`     // 1.前端记录读取状态 0.每次都展示
	Url         string `gorm:"type:varchar(255)" json:"url"`       // 跳转地址
	IsLogin     uint   `gorm:"type:tinyint(1)" json:"isLogin"`     // 登录可见
	UserId      uint   `json:"userId"`
	User        User   `json:"user"`
}
