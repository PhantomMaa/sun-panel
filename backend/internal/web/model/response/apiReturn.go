package response

import (
	"sun-panel/internal/util/i18n"

	"github.com/gin-gonic/gin"
)

func ApiReturn(ctx *gin.Context, code int, msg string, data interface{}) {
	returnData := map[string]interface{}{
		"code": code,
		"msg":  msg,
	}
	if data != nil {
		returnData["data"] = data
	}
	ctx.JSON(200, returnData)
}

func SuccessData(ctx *gin.Context, data interface{}) {
	ApiReturn(ctx, 0, "OK", data)
}

func SuccessListData(ctx *gin.Context, list interface{}, count uint) {
	ApiReturn(ctx, 0, "OK", gin.H{
		"list":  list,
		"count": count,
	})
}

func Success(ctx *gin.Context) {
	ApiReturn(ctx, 0, "OK", nil)
}

func ListData(ctx *gin.Context, list interface{}, count int64) {
	data := map[string]interface{}{
		"list":  list,
		"count": count,
	}
	ApiReturn(ctx, 0, "OK", data)
}

// 返回错误 需要个性化定义的错误|带返回数据的错误
func ErrorCode(ctx *gin.Context, code int, errMsg string, data interface{}) {
	ApiReturn(ctx, code, errMsg, data)
}

// 返回错误 普通提示错误
func Error(ctx *gin.Context, errMsg string) {
	ErrorCode(ctx, -1, errMsg, nil)
}

// 返回错误 需要个性化定义的错误|带返回数据的错误
func ErrorNoAccess(ctx *gin.Context) {
	ErrorCode(ctx, 1005, i18n.Obj.Get("common.no_access"), nil)
}

// 返回错误 参数错误
func ErrorParamFomat(ctx *gin.Context, errMsg string) {
	Error(ctx, i18n.Obj.GetAndInsert("common.api_error_param_format", "[", errMsg, "]"))
	// Error(ctx, "参数错误")
}

// // 返回错误 数据库
func ErrorDatabase(ctx *gin.Context, errMsg string) {
	ErrorByCodeAndMsg(ctx, 1200, errMsg)
}

// 返回错误 数据记录未找到
func ErrorDataNotFound(ctx *gin.Context) {
	ErrorCode(ctx, 404, "Data not found", nil)
}

func ErrorByCode(ctx *gin.Context, code int) {
	msg := "Server error"
	if v, ok := GetErrorMsgByCode(code); ok {
		msg = v
	}
	ErrorCode(ctx, code, msg, nil)
}

// 使用错误码的错误并附加错误信息
func ErrorByCodeAndMsg(ctx *gin.Context, code int, msg string) {
	defalurMsg := "Server error"
	if v, ok := GetErrorMsgByCode(code); ok {
		msg = v
	}
	ErrorCode(ctx, code, defalurMsg+"["+msg+"]", nil)
}

func GetErrorMsgByCode(code int) (string, bool) {
	if v, ok := ErrorCodeMap[code]; ok {
		return v, true
	} else {
		return "", false
	}
}

// 返回错误 需要个性化定义的错误|带返回数据的错误
// func ErrorNoAccess(ctx *gin.Context) {
// 	ErrorCode(ctx, 1005, global.Lang.Get("common.no_access"), nil)
// }

// // 返回错误 参数错误
// func ErrorParamFomat(ctx *gin.Context, errMsg string) {
// 	Error(ctx, global.Lang.GetAndInsert("common.api_error_param_format", "[", errMsg, "]"))
// }

// // 返回错误 数据库
// func ErrorDatabase(ctx *gin.Context, errMsg string) {
// 	Error(ctx, global.Lang.GetAndInsert("common.db_error", "[", errMsg, "]"))
// }
