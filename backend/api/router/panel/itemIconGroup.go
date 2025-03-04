package panel

import (
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"math"
	"sun-panel/api/common/apiData/commonApiStructs"
	"sun-panel/api/common/apiReturn"
	"sun-panel/api/common/base"
	"sun-panel/api/middleware"
	"sun-panel/internal/global"
	"sun-panel/internal/repository"

	"github.com/gin-gonic/gin"
)

type ItemIconGroup struct {
}

func NewItemIconGroupRouter() *ItemIconGroup {
	return &ItemIconGroup{}
}

func (a *ItemIconGroup) InitRouter(router *gin.RouterGroup) {
	r := router.Group("")
	r.Use(middleware.JWTAuth())
	{
		r.POST("/panel/itemIconGroup/edit", a.Edit)
		r.POST("/panel/itemIconGroup/deletes", a.Deletes)
		r.POST("/panel/itemIconGroup/saveSort", a.SaveSort)
		r.GET("/panel/itemIconGroup/getList", a.GetList)
	}
}

func (a *ItemIconGroup) Edit(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	req := repository.ItemIconGroup{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	req.UserId = userInfo.ID

	if req.ID != 0 {
		// 修改
		updateField := []string{"IconJson", "Icon", "Title", "Url", "LanUrl", "Description", "OpenMethod", "GroupId", "UserId"}
		if req.Sort != 0 {
			updateField = append(updateField, "Sort")
		}
		global.Db.Model(&repository.ItemIconGroup{}).
			Select(updateField).
			Where("id=?", req.ID).Updates(&req)
	} else {
		// 创建
		global.Db.Create(&req)
	}

	apiReturn.SuccessData(c, req)
}

func (a *ItemIconGroup) GetList(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	var groups []repository.ItemIconGroup

	err := global.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Order("sort ,created_at").Where("user_id=?", userInfo.ID).Find(&groups).Error; err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return err
		}

		// 判断分组是否为空，为空将自动创建默认分组
		if len(groups) == 0 {
			defaultGroup := repository.ItemIconGroup{
				Title:  "APP",
				UserId: userInfo.ID,
				Icon:   "material-symbols:ad-group-outline",
			}
			if err := tx.Create(&defaultGroup).Error; err != nil {
				apiReturn.ErrorDatabase(c, err.Error())
				return err
			}

			// 并将当前账号下所有无分组的图标更新到当前组
			if err := tx.Model(&repository.ItemIcon{}).Where("user_id=?", userInfo.ID).Update("item_icon_group_id", defaultGroup.ID).Error; err != nil {
				apiReturn.ErrorDatabase(c, err.Error())
				return err
			}

			groups = append(groups, defaultGroup)
		}

		// 返回 nil 提交事务
		return nil
	})

	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	} else {
		apiReturn.SuccessListData(c, groups, 0)
	}
}

func (a *ItemIconGroup) Deletes(c *gin.Context) {
	req := commonApiStructs.RequestDeleteIds[uint]{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	var count int64
	if err := global.Db.Model(&repository.ItemIconGroup{}).Where(" user_id=?", userInfo.ID).Count(&count).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	} else {
		if math.Abs(float64(len(req.Ids))-float64(count)) < 1 {
			apiReturn.ErrorCode(c, 1201, "At least one must be retained", nil)
			return
		}

	}

	txErr := global.Db.Transaction(func(tx *gorm.DB) error {
		mitemIcon := repository.ItemIcon{}
		if err := tx.Delete(&repository.ItemIconGroup{}, "id in ? AND user_id=?", req.Ids, userInfo.ID).Error; err != nil {
			return err
		}

		if err := mitemIcon.DeleteByItemIconGroupIds(tx, userInfo.ID, req.Ids); err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		apiReturn.ErrorDatabase(c, txErr.Error())
		return
	}

	apiReturn.Success(c)
}

func (a *ItemIconGroup) SaveSort(c *gin.Context) {
	req := commonApiStructs.SortRequest{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	transactionErr := global.Db.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		for _, v := range req.SortItems {
			if err := tx.Model(&repository.ItemIconGroup{}).Where("user_id=? AND id=?", userInfo.ID, v.Id).Update("sort", v.Sort).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
		}

		// 返回 nil 提交事务
		return nil
	})

	if transactionErr != nil {
		apiReturn.ErrorDatabase(c, transactionErr.Error())
		return
	}

	apiReturn.Success(c)
}
