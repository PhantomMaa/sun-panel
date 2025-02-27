package panel

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"net/url"
	"path"
	"strings"
	"sun-panel/api/common/apiData/commonApiStructs"
	"sun-panel/api/common/apiData/panelApiStructs"
	"sun-panel/api/common/apiReturn"
	"sun-panel/api/common/base"
	"sun-panel/internal/global"
	repository2 "sun-panel/internal/repository"
	"sun-panel/internal/siteFavicon"
	"sun-panel/internal/storage"
)

type ItemIcon struct {
	storage storage.RcloneStorage
}

var filePrefix string

func NewItemIcon(s storage.RcloneStorage) *ItemIcon {
	source_path := global.Config.GetValueString("base", "source_path")
	filePrefix = fmt.Sprintf("/%s/", source_path)
	return &ItemIcon{
		storage: s,
	}
}

func (a *ItemIcon) Edit(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	req := repository2.ItemIcon{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if req.ItemIconGroupId == 0 {
		// apiReturn.Error(c, "Group is mandatory")
		apiReturn.ErrorParamFomat(c, "Group is mandatory")
		return
	}

	req.UserId = userInfo.ID

	// json转字符串
	if j, err := json.Marshal(req.Icon); err == nil {
		req.IconJson = string(j)
	}

	if req.ID != 0 {
		// 修改
		updateField := []string{"IconJson", "Icon", "Title", "Url", "LanUrl", "Description", "OpenMethod", "GroupId", "UserId", "ItemIconGroupId"}
		if req.Sort != 0 {
			updateField = append(updateField, "Sort")
		}
		global.Db.Model(&repository2.ItemIcon{}).
			Select(updateField).
			Where("id=?", req.ID).Updates(&req)
	} else {
		req.Sort = 9999
		// 创建
		global.Db.Create(&req)
	}

	apiReturn.SuccessData(c, req)
}

// 添加多个图标
func (a *ItemIcon) AddMultiple(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	// type Request
	var req []repository2.ItemIcon

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	for i := 0; i < len(req); i++ {
		if req[i].ItemIconGroupId == 0 {
			apiReturn.ErrorParamFomat(c, "Group is mandatory")
			return
		}
		req[i].UserId = userInfo.ID
		// json转字符串
		if j, err := json.Marshal(req[i].Icon); err == nil {
			req[i].IconJson = string(j)
		}
	}

	global.Db.Create(&req)

	apiReturn.SuccessData(c, req)
}

func (a *ItemIcon) GetListByGroupId(c *gin.Context) {
	req := repository2.ItemIcon{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	var itemIcons []repository2.ItemIcon

	if err := global.Db.Order("sort ,created_at").Find(&itemIcons, "item_icon_group_id = ? AND user_id=?", req.ItemIconGroupId, userInfo.ID).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	for k, v := range itemIcons {
		json.Unmarshal([]byte(v.IconJson), &itemIcons[k].Icon)
	}

	apiReturn.SuccessListData(c, itemIcons, 0)
}

func (a *ItemIcon) Deletes(c *gin.Context) {
	req := commonApiStructs.RequestDeleteIds[uint]{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	// Start a transaction to ensure data consistency
	err := global.Db.Transaction(func(tx *gorm.DB) error {
		// First find all items to get their icon paths
		var items []repository2.ItemIcon
		if err := tx.Find(&items, "id in ? AND user_id=?", req.Ids, userInfo.ID).Error; err != nil {
			return err
		}

		// Delete associated files
		for _, item := range items {
			var icon map[string]interface{}
			if err := json.Unmarshal([]byte(item.IconJson), &icon); err != nil {
				global.Logger.Errorf("Failed to unmarshal icon JSON: %v", err)
				continue
			}

			// Check if the icon has a src field indicating a file path
			if src, ok := icon["src"].(string); ok && strings.HasPrefix(src, filePrefix) {
				// Extract the file path from the URL
				filePath := strings.TrimPrefix(src, filePrefix)

				// Find and delete the file record
				var file repository2.File
				if err := tx.Where("src = ? AND user_id = ?", filePath, userInfo.ID).First(&file).Error; err == nil {
					if err := tx.Delete(&file).Error; err != nil {
						return err
					}

					// Delete the actual file using storage interface
					if err := a.storage.Delete(c.Request.Context(), filePath); err != nil {
						global.Logger.Errorf("Failed to delete file %s: %v", filePath, err)
						// Continue with deletion even if file removal fails
					}
				}
			}
		}

		// Finally delete the item icons
		if err := tx.Delete(&repository2.ItemIcon{}, "id in ? AND user_id=?", req.Ids, userInfo.ID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

// GetSiteFavicon 支持获取并直接下载对方网站图标到服务器
func (a *ItemIcon) GetSiteFavicon(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	req := panelApiStructs.ItemIconGetSiteFaviconReq{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	resp := panelApiStructs.ItemIconGetSiteFaviconResp{}
	fullUrl := ""
	if iconUrl, err := siteFavicon.GetOneFaviconURL(req.Url); err != nil {
		apiReturn.Error(c, "acquisition failed: get ico error:"+err.Error())
		return
	} else {
		fullUrl = iconUrl
	}

	parsedURL, err := url.Parse(req.Url)
	if err != nil {
		apiReturn.Error(c, "acquisition failed:"+err.Error())
		return
	}

	protocol := parsedURL.Scheme
	global.Logger.Debug("protocol:", protocol)
	global.Logger.Debug("fullUrl:", fullUrl)

	// 如果URL以双斜杠（//）开头，则使用当前页面协议
	if strings.HasPrefix(fullUrl, "//") {
		fullUrl = protocol + "://" + fullUrl[2:]
	} else if !strings.HasPrefix(fullUrl, "http://") && !strings.HasPrefix(fullUrl, "https://") {
		// 如果URL既不以http://开头也不以https://开头，则默认为http协议
		fullUrl = "http://" + fullUrl
	}
	global.Logger.Debug("fullUrl:", fullUrl)
	// 去除图标的get参数
	{
		parsedIcoURL, err := url.Parse(fullUrl)
		if err != nil {
			apiReturn.Error(c, "acquisition failed: parsed ico URL :"+err.Error())
			return
		}
		fullUrl = parsedIcoURL.Scheme + "://" + parsedIcoURL.Host + parsedIcoURL.Path
	}
	global.Logger.Debug("fullUrl:", fullUrl)

	// 下载图标
	filepath, err := siteFavicon.DownloadImage(c.Request.Context(), fullUrl, a.storage)
	if err != nil {
		apiReturn.Error(c, "acquisition failed: download "+err.Error())
		return
	}

	// 保存到数据库
	ext := path.Ext(fullUrl)
	if ext == "" {
		ext = ".ico"
	}
	mFile := repository2.File{}
	if _, err := mFile.AddFile(userInfo.ID, parsedURL.Host, ext, filepath); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	resp.IconUrl = filePrefix + filepath
	apiReturn.SuccessData(c, resp)
}

// 保存排序
func (a *ItemIcon) SaveSort(c *gin.Context) {
	req := panelApiStructs.ItemIconSaveSortRequest{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	transactionErr := global.Db.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		for _, v := range req.SortItems {
			if err := tx.Model(&repository2.ItemIcon{}).Where("user_id=? AND id=? AND item_icon_group_id=?", userInfo.ID, v.Id, req.ItemIconGroupId).Update("sort", v.Sort).Error; err != nil {
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
