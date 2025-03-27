package system

import (
	"fmt"
	"net/http"
	"path"
	"strings"
	"sun-panel/internal/biz/constant"
	"sun-panel/internal/global"
	"sun-panel/internal/infra/config"
	"sun-panel/internal/util"
	"sun-panel/internal/web/interceptor"
	"sun-panel/internal/web/model/base"
	"sun-panel/internal/web/model/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type FileRouter struct {
	urlPrefix string
}

func NewFileRouter() *FileRouter {
	return &FileRouter{
		urlPrefix: config.AppConfig.Base.URLPrefix,
	}
}

func (a *FileRouter) InitRouter(router *gin.RouterGroup) {
	// 公共访问组，不需要 JWT 认证
	public := router.Group("")
	{
		// S3 文件访问路由
		public.GET("/file/s3/*filepath", a.GetS3File)
	}

	// 需要 JWT 认证的私有访问组
	private := router.Group("")
	private.Use(interceptor.JWTAuth)
	{
		private.POST("/file/uploadImg", a.UploadImg)
		private.POST("/file/delete", a.Delete)
		private.GET("/file/getList", a.GetList)
	}
}

func (a *FileRouter) UploadImg(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	f, err := c.FormFile("imgfile")
	if err != nil {
		response.ErrorByCode(c, constant.CodeUploadFailed)
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

	if !util.InArray(agreeExts, fileExt) {
		response.ErrorByCode(c, constant.CodeUnsupportFileFormat)
		return
	}

	fileName := util.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String())) + fileExt

	// 打开文件以获取Reader
	src, err := f.Open()
	if err != nil {
		global.Logger.Errorf("Failed to open uploaded file: %v", err)
		response.ErrorByCode(c, constant.CodeUploadFailed)
		return
	}

	defer func() {
		if err := src.Close(); err != nil {
			global.Logger.Errorf("Failed to close file. error : %v", err)
		}
	}()

	// 使用存储接口上传文件
	filepath, err := global.Storage.Upload(c.Request.Context(), src, fileName)
	if err != nil {
		global.Logger.Errorf("Failed to upload file: %v", err)
		response.ErrorByCode(c, constant.CodeUploadFailed)
		return
	}

	// 向数据库添加记录
	_, err = global.FileRepo.AddFile(userInfo.ID, f.Filename, fileExt, filepath)
	if err != nil {
		global.Logger.Errorf("Failed to add file record to database: %v", err)
		response.ErrorByCode(c, constant.CodeUploadFailed)
		return
	}

	global.Logger.Infof("Successfully uploaded file %s to %s", f.Filename, filepath)
	response.SuccessData(c, gin.H{
		"imageUrl": a.urlPrefix + filepath,
	})
}

func (a *FileRouter) GetList(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	list, count, err := global.FileRepo.GetList(userInfo.ID)
	if err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	var data []map[string]interface{}
	for _, v := range list {
		data = append(data, map[string]interface{}{
			"src":        a.urlPrefix + v.Src,
			"fileName":   v.FileName,
			"id":         v.ID,
			"createTime": v.CreatedAt,
			"updateTime": v.UpdatedAt,
		})
	}
	response.SuccessListData(c, data, count)
}

func (a *FileRouter) Delete(c *gin.Context) {
	type RequestDeleteId struct {
		Id uint `json:"id" binding:"required"`
	}

	req := RequestDeleteId{}
	userInfo, _ := base.GetCurrentUserInfo(c)
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		response.ErrorParamFomat(c, err.Error())
		return
	}

	file, err := global.FileRepo.Get(userInfo.ID, req.Id)
	if err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	// 从存储中删除文件
	if err := global.Storage.Delete(c.Request.Context(), file.Src); err != nil {
		global.Logger.Errorf("Failed to delete file %s: %v", file.Src, err)
	}

	// 从数据库中删除记录
	if err := global.FileRepo.Delete(userInfo.ID, req.Id); err != nil {
		response.ErrorDatabase(c, err.Error())
		return
	}

	response.Success(c)
}

func (a *FileRouter) GetS3File(c *gin.Context) {
	filepath := c.Param("filepath") // 获取 /api/file/s3/ 后的所有部分
	if filepath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file path is required"})
		return
	}

	// 从存储中读取文件
	fileData, err := global.Storage.Get(c, filepath)
	if err != nil {
		global.Logger.Errorf("Failed to get file %s: %v", filepath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get file"})
		return
	}

	// 设置文件类型
	contentType := "application/octet-stream"
	ext := path.Ext(filepath)
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".svg":
		contentType = "image/svg+xml"
	}

	// 设置响应头
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "inline; filename="+path.Base(filepath))

	global.Logger.Infof("Successfully serving file: %s with content type: %s", filepath, contentType)
	// 返回文件内容
	c.Data(http.StatusOK, contentType, fileData)
}
