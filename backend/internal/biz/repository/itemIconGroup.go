package repository

import (
	"gorm.io/gorm"
	"sun-panel/internal/web/model/param/commonApiStructs"
)

type ItemIconGroupRepo struct {
}

type IItemIconGroupRepo interface {
	Save(itemIconGroup *ItemIconGroup) error
	GetList(userId uint) ([]ItemIconGroup, error)
	Count(userId uint) (int, error)
	Deletes(userId uint, ids []uint) error
	BatchSaveSort(userId uint, sortItems []commonApiStructs.SortRequestItem) error
}

func NewItemIconGroupRepo() IItemIconGroupRepo {
	return &ItemIconGroupRepo{}
}

type ItemIconGroup struct {
	BaseModel
	Icon        string `json:"icon"`
	Title       string `gorm:"type:varchar(50)" json:"title"`
	Description string `gorm:"type:varchar(1000)" json:"description"`
	Sort        int    `gorm:"type:int(11)" json:"sort"`
	UserId      uint   `json:"userId"`
	User        User   `json:"user"`
}

func (r *ItemIconGroupRepo) Save(itemIconGroup *ItemIconGroup) error {
	if itemIconGroup.ID != 0 {
		updateField := []string{"IconJson", "Icon", "Title", "Url", "LanUrl", "Description", "OpenMethod", "GroupId", "UserId"}
		if itemIconGroup.Sort != 0 {
			updateField = append(updateField, "Sort")
		}
		if err := Db.Select(updateField).Where("id=?", itemIconGroup.ID).Updates(&itemIconGroup).Error; err != nil {
			return err
		}
	} else {
		if err := Db.Create(&itemIconGroup).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *ItemIconGroupRepo) GetList(userId uint) ([]ItemIconGroup, error) {
	var groups []ItemIconGroup

	if err := Db.Where("user_id=?", userId).Order("sort ,created_at").Find(&groups).Error; err != nil {
		return nil, err
	}

	return groups, nil
}

func (r *ItemIconGroupRepo) Count(userId uint) (int, error) {
	var count int64
	if err := Db.Model(&ItemIconGroup{}).Where(" user_id=?", userId).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *ItemIconGroupRepo) Deletes(userId uint, ids []uint) error {
	txErr := Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&ItemIconGroup{}, "id in ? AND user_id=?", ids, userId).Error; err != nil {
			return err
		}

		if err := tx.Delete(&ItemIcon{}, "item_icon_group_id in ? AND user_id=?", ids, userId).Error; err != nil {
			return err
		}

		return nil
	})

	return txErr
}

func (r *ItemIconGroupRepo) BatchSaveSort(userId uint, sortItems []commonApiStructs.SortRequestItem) error {
	return Db.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		for _, v := range sortItems {
			if err := tx.Model(&ItemIconGroup{}).Where("user_id=? AND id=?", userId, v.Id).Update("sort", v.Sort).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
		}

		// 返回 nil 提交事务
		return nil
	})
}
