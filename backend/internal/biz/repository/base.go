package repository

import (
	"time"

	_ "gorm.io/driver/mysql"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createTime"`
	UpdatedAt time.Time      `json:"updateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type BaseModelNoId struct {
	CreatedAt time.Time      `json:"createTime"`
	UpdatedAt time.Time      `json:"updateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// PageLimitStruct 分页的结构体
type PageLimitStruct struct {
	PageSize  int `gorm:"-"`
	LimitSize int `gorm:"-"`
}

// 计算分页
func calcPage(pageSize, limitSize int) (offset, limit int) {
	offset = limitSize * (pageSize - 1)
	limit = limitSize
	return
}

var Db *gorm.DB
