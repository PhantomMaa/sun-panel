package panel

import (
	"encoding/json"
	"net/url"
	"path"
	"strings"
	"sun-panel/internal/biz/repository"
	"sun-panel/internal/global"
	"sun-panel/internal/infra/zaplog"
	"sun-panel/internal/util"
	"sun-panel/internal/util/favicon"
	"sun-panel/internal/web/interceptor"
	"sun-panel/internal/web/model/base"
	"sun-panel/internal/web/model/param/panelApi"
	"sun-panel/internal/web/model/response"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ItemIconRouter struct {
}

var urlPrefix string

func NewItemIconRouter() *ItemIconRouter {
	urlPrefix = global.Config.Base.URLPrefix
	return &ItemIconRouter{}
}

func (a *ItemIconRouter) InitRouter(router *gin.RouterGroup) {
	r := router.Group("")
	r.Use(interceptor.JWTAuth)
	{
		r.POST("/panel/itemIcon/edit", a.Edit)
		r.POST("/panel/itemIcon/delete", a.Delete)
		r.POST("/panel/itemIcon/saveSort", a.SaveSort)
		r.POST("/panel/itemIcon/addMultiple", a.AddMultiple)
		r.POST("/panel/itemIcon/getSiteFavicon", a.GetSiteFavicon)
		r.GET("/panel/itemIcon/getListByGroupId", a.GetListByGroupId)
	}
}

func (a *ItemIconRouter) Edit(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	itemIcon := repository.ItemIcon{}
	if err := c.ShouldBindBodyWith(&itemIcon, binding.JSON); err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	if itemIcon.ItemIconGroupId == 0 {
		response.ErrorParamFomat(c, "Group is mandatory")
		return
	}

	itemIcon.UserId = userInfo.ID
	itemIcon.IconJson = util.ToJSONString(itemIcon.Icon)
	if err := global.ItemIconRepo.Save(&itemIcon); err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	response.SuccessData(c, itemIcon)
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

	for i := range req {
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

	if err := global.ItemIconRepo.BatchSave(req); err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	response.SuccessData(c, req)
}

func (a *ItemIconRouter) GetListByGroupId(c *gin.Context) {
	type ParamsStruct struct {
		ItemIconGroupId uint `form:"itemIconGroupId" json:"itemIconGroupId"`
	}

	// 查询条件
	param := ParamsStruct{}
	if err := c.ShouldBind(&param); err != nil {
		response.ErrorParamFomat(c, err.Error())
		c.Abort()
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	itemIcons, err := global.ItemIconRepo.GetList(userInfo.ID, param.ItemIconGroupId)
	if err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	for k, v := range itemIcons {
		json.Unmarshal([]byte(v.IconJson), &itemIcons[k].Icon)
	}

	response.SuccessListData(c, itemIcons, 0)
}

func (a *ItemIconRouter) Delete(c *gin.Context) {
	type RequestDeleteId struct {
		Id uint `json:"id" binding:"required"`
	}

	req := RequestDeleteId{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	item, err := global.ItemIconRepo.Get(userInfo.ID, req.Id)
	if err != nil {
		return
	}

	// Delete the actual file using storage interface
	var icon map[string]any
	if err := json.Unmarshal([]byte(item.IconJson), &icon); err == nil {
		// Check if the icon has a src field indicating a file path
		if fileName, ok := icon["fileName"].(string); ok {
			if err := global.Storage.Delete(c.Request.Context(), fileName); err != nil {
				zaplog.Logger.Warnf("Failed to delete file %s: %v", fileName, err)
			}
		}
	}

	if err := global.ItemIconRepo.Delete(userInfo.ID, req.Id); err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	response.Success(c)
}

// 支持获取并直接下载对方网站图标到服务器
func (a *ItemIconRouter) GetSiteFavicon(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	req := panelApi.ItemIconGetSiteFaviconReq{}
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	resp := panelApi.ItemIconGetSiteFaviconResp{}
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

	// 如果URL以双斜杠（//）开头，则使用当前页面协议
	if strings.HasPrefix(fullUrl, "//") {
		fullUrl = protocol + "://" + fullUrl[2:]
	} else if !strings.HasPrefix(fullUrl, "http://") && !strings.HasPrefix(fullUrl, "https://") {
		// 如果URL既不以http://开头也不以https://开头，则默认为http协议
		fullUrl = "http://" + fullUrl
	}

	// 去除图标的get参数
	{
		parsedIcoURL, err := url.Parse(fullUrl)
		if err != nil {
			response.Error(c, "acquisition failed: parsed ico URL :"+err.Error())
			return
		}
		fullUrl = parsedIcoURL.Scheme + "://" + parsedIcoURL.Host + parsedIcoURL.Path
	}
	zaplog.Logger.Debug("fullUrl:", fullUrl)

	// 下载图标
	fileName, err := favicon.DownloadImage(c.Request.Context(), fullUrl)
	if err != nil {
		response.Error(c, "acquisition failed: download "+err.Error())
		return
	}

	// 保存到数据库
	ext := path.Ext(fullUrl)
	if ext == "" {
		ext = ".ico"
	}
	if _, err := global.FileRepo.AddFile(userInfo.ID, ext, fileName); err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	resp.FileName = fileName
	resp.IconUrl = urlPrefix + fileName
	response.SuccessData(c, resp)
}

// SaveSort 保存排序
func (a *ItemIconRouter) SaveSort(c *gin.Context) {
	req := panelApi.ItemIconSaveSortRequest{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	err := global.ItemIconRepo.BatchSaveSort(userInfo.ID, req.ItemIconGroupId, req.SortItems)
	if err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	response.Success(c)
}
