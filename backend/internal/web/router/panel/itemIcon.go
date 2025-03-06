package panel

import (
	"encoding/json"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"net/url"
	"path"
	"strings"
	"sun-panel/internal/biz/repository"
	"sun-panel/internal/common"
	"sun-panel/internal/common/favicon"
	"sun-panel/internal/global"
	"sun-panel/internal/infra/storage"
	"sun-panel/internal/web/interceptor"
	"sun-panel/internal/web/model/base"
	"sun-panel/internal/web/model/param/commonApiStructs"
	"sun-panel/internal/web/model/param/panelApiStructs"
	"sun-panel/internal/web/model/response"

	"github.com/gin-gonic/gin"
)

type ItemIconRouter struct {
	storage storage.RcloneStorage
}

var urlPrefix string

func NewItemIconRouter() *ItemIconRouter {
	urlPrefix = global.Config.GetValueString("base", "url_prefix")
	return &ItemIconRouter{
		storage: *global.Storage,
	}
}

func (a *ItemIconRouter) InitRouter(router *gin.RouterGroup) {
	r := router.Group("")
	r.Use(interceptor.JWTAuth)
	{
		r.POST("/panel/itemIcon/edit", a.Edit)
		r.POST("/panel/itemIcon/deletes", a.Deletes)
		r.POST("/panel/itemIcon/saveSort", a.SaveSort)
		r.POST("/panel/itemIcon/addMultiple", a.AddMultiple)
		r.POST("/panel/itemIcon/getSiteFavicon", a.GetSiteFavicon)
		r.GET("/panel/itemIcon/getListByGroupId", a.GetListByGroupId)
	}
}

func (a *ItemIconRouter) Edit(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	req := repository.ItemIcon{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	if req.ItemIconGroupId == 0 {
		response.ErrorParamFomat(c, "Group is mandatory")
		return
	}

	req.UserId = userInfo.ID
	req.IconJson = common.ToJSONString(req.Icon)

	if req.ID != 0 {
		// 修改
		updateField := []string{"IconJson", "Icon", "Title", "Url", "LanUrl", "Description", "OpenMethod", "GroupId", "UserId", "ItemIconGroupId"}
		if req.Sort != 0 {
			updateField = append(updateField, "Sort")
		}
		global.Db.Model(&repository.ItemIcon{}).Select(updateField).Where("id=?", req.ID).Updates(&req)
	} else {
		req.Sort = 9999
		// 创建
		global.Db.Create(&req)
	}

	response.SuccessData(c, req)
}

// 添加多个图标
func (a *ItemIconRouter) AddMultiple(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	// type Request
	var req []repository.ItemIcon

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	for i := 0; i < len(req); i++ {
		if req[i].ItemIconGroupId == 0 {
			response.ErrorParamFomat(c, "Group is mandatory")
			return
		}
		req[i].UserId = userInfo.ID
		// json转字符串
		if j, err := json.Marshal(req[i].Icon); err == nil {
			req[i].IconJson = string(j)
		}
	}

	global.Db.Create(&req)

	response.SuccessData(c, req)
}

func (a *ItemIconRouter) GetListByGroupId(c *gin.Context) {
	type ParamsStruct struct {
		ItemIconGroupId int `form:"itemIconGroupId" json:"itemIconGroupId"`
	}

	// 查询条件
	param := ParamsStruct{}
	if err := c.ShouldBind(&param); err != nil {
		response.ErrorParamFomat(c, err.Error())
		c.Abort()
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	var itemIcons []repository.ItemIcon

	if err := global.Db.Order("sort ,created_at").Find(&itemIcons, "item_icon_group_id = ? AND user_id=?", param.ItemIconGroupId, userInfo.ID).Error; err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	for k, v := range itemIcons {
		json.Unmarshal([]byte(v.IconJson), &itemIcons[k].Icon)
	}

	response.SuccessListData(c, itemIcons, 0)
}

func (a *ItemIconRouter) Deletes(c *gin.Context) {
	req := commonApiStructs.RequestDeleteIds[uint]{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	// Start a transaction to ensure data consistency
	err := global.Db.Transaction(func(tx *gorm.DB) error {
		// First find all items to get their icon paths
		var items []repository.ItemIcon
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
			if src, ok := icon["src"].(string); ok && strings.HasPrefix(src, urlPrefix) {
				// Extract the file path from the URL
				filePath := strings.TrimPrefix(src, urlPrefix)

				// Find and delete the file record
				var file repository.File
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
		if err := tx.Delete(&repository.ItemIcon{}, "id in ? AND user_id=?", req.Ids, userInfo.ID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	response.Success(c)
}

// GetSiteFavicon 支持获取并直接下载对方网站图标到服务器
func (a *ItemIconRouter) GetSiteFavicon(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	req := panelApiStructs.ItemIconGetSiteFaviconReq{}
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	resp := panelApiStructs.ItemIconGetSiteFaviconResp{}
	fullUrl := ""
	if iconUrl, err := favicon.GetOneFaviconURL(req.Url); err != nil {
		response.Error(c, "acquisition failed: get ico error:"+err.Error())
		return
	} else {
		fullUrl = iconUrl
	}

	parsedURL, err := url.Parse(req.Url)
	if err != nil {
		response.Error(c, "acquisition failed:"+err.Error())
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
			response.Error(c, "acquisition failed: parsed ico URL :"+err.Error())
			return
		}
		fullUrl = parsedIcoURL.Scheme + "://" + parsedIcoURL.Host + parsedIcoURL.Path
	}
	global.Logger.Debug("fullUrl:", fullUrl)

	// 下载图标
	filepath, err := favicon.DownloadImage(c.Request.Context(), fullUrl, a.storage)
	if err != nil {
		response.Error(c, "acquisition failed: download "+err.Error())
		return
	}

	// 保存到数据库
	ext := path.Ext(fullUrl)
	if ext == "" {
		ext = ".ico"
	}
	mFile := repository.File{}
	if _, err := mFile.AddFile(userInfo.ID, parsedURL.Host, ext, filepath); err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	resp.IconUrl = urlPrefix + filepath
	response.SuccessData(c, resp)
}

// SaveSort 保存排序
func (a *ItemIconRouter) SaveSort(c *gin.Context) {
	req := panelApiStructs.ItemIconSaveSortRequest{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	transactionErr := global.Db.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		for _, v := range req.SortItems {
			if err := tx.Model(&repository.ItemIcon{}).Where("user_id=? AND id=? AND item_icon_group_id=?", userInfo.ID, v.Id, req.ItemIconGroupId).Update("sort", v.Sort).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
		}

		// 返回 nil 提交事务
		return nil
	})

	if transactionErr != nil {
		response.ErrorDatabase(c, transactionErr.Error())
		return
	}

	response.Success(c)
}
