package system

import (
	"fmt"
	"path"
	"strings"
	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/lib/storage"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type FileApi struct{}

func (a *FileApi) UploadImg(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	f, err := c.FormFile("imgfile")
	if err != nil {
		apiReturn.ErrorByCode(c, 1300)
		return
	}

	fileExt := strings.ToLower(path.Ext(f.Filename))
	agreeExts := []string{
		".png",
		".jpg",
		".gif",
		".jpeg",
		".webp",
		".svg",
		".ico",
	}

	if !cmn.InArray(agreeExts, fileExt) {
		apiReturn.ErrorByCode(c, 1301)
		return
	}

	fileName := cmn.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String())) + fileExt

	// 使用存储接口上传文件
	storageInstance := storage.GetStorage()
	filepath, err := storageInstance.Upload(f, fileName)
	if err != nil {
		apiReturn.ErrorByCode(c, 1300)
		return
	}

	// 向数据库添加记录
	mFile := models.File{}
	_, err = mFile.AddFile(userInfo.ID, f.Filename, fileExt, filepath)
	if err != nil {
		global.Logger.Errorf("Failed to add file record to database: %v", err)
		apiReturn.ErrorByCode(c, 1300)
		return
	}

	global.Logger.Infof("Successfully uploaded file %s to %s", f.Filename, filepath)
	apiReturn.SuccessData(c, gin.H{
		"imageUrl": filepath,
	})
}

func (a *FileApi) GetList(c *gin.Context) {
	list := []models.File{}
	userInfo, _ := base.GetCurrentUserInfo(c)
	var count int64
	if err := global.Db.Order("created_at desc").Find(&list, "user_id=?", userInfo.ID).Count(&count).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	data := []map[string]interface{}{}
	for _, v := range list {
		data = append(data, map[string]interface{}{
			"src":        v.Src,
			"fileName":   v.FileName,
			"id":         v.ID,
			"createTime": v.CreatedAt,
			"updateTime": v.UpdatedAt,
		})
	}
	apiReturn.SuccessListData(c, data, count)
}

func (a *FileApi) Deletes(c *gin.Context) {
	req := commonApiStructs.RequestDeleteIds[uint]{}
	userInfo, _ := base.GetCurrentUserInfo(c)
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	global.Db.Transaction(func(tx *gorm.DB) error {
		files := []models.File{}

		if err := tx.Order("created_at desc").Find(&files, "user_id=? AND id in ?", userInfo.ID, req.Ids).Error; err != nil {
			return err
		}

		for _, v := range files {
			if err := storage.GetStorage().Delete(v.Src); err != nil {
				global.Logger.Errorf("Failed to delete file %s: %v", v.Src, err)
				return err
			}
		}

		if err := tx.Order("created_at desc").Delete(&files, "user_id=? AND id in ?", userInfo.ID, req.Ids).Error; err != nil {
			return err
		}

		return nil
	})

	apiReturn.Success(c)

}
