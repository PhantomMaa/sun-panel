package repository

import (
	"time"

	_ "gorm.io/driver/mysql"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createTime"`
	UpdatedAt time.Time `json:"updateTime"`
}

type PagedParam struct {
	Limit int `form:"limit" json:"limit" gorm:"-"`
	Page  int `form:"page" json:"page" gorm:"-"`
}

// 计算分页
func CalcOffset(pagedParam PagedParam) int {
	return pagedParam.Limit * (pagedParam.Page - 1)
}

var Db *gorm.DB
