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

// 分页的结构体
type PageLimitStruct struct {
	PageSize  int `gorm:"-"`
	LimitSize int `gorm:"-"`
}

// 计算分页
func calcPage(page_size, limit_size int) (offset, limit int) {
	offset = limit_size * (page_size - 1)
	limit = limit_size
	return
}

var Db *gorm.DB
