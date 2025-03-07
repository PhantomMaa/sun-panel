package response

import (
	"sun-panel/internal/language"

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

// 返回错误 验证码相关错误错误
// func ErrorVerification(ctx *gin.Context, errCode int, codeID string) {
// 	msg := ""
// 	switch errCode {
// 	case ERROR_CODE_VERIFICATION_FAIL:
// 		msg = "验证失败，请重新验证"
// 	case ERROR_CODE_VERIFICATION_MUST:
// 		msg = "需要进一步验证"
// 	}
// 	ApiReturn(ctx, errCode, msg, gin.H{
// 		"verification": commonApiStructs.VerificationResponse{
// 			CodeID:  codeID,
// 			Result:  false,
// 			Message: msg,
// 		},
// 	})
// }

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
	ErrorCode(ctx, 1005, language.Obj.Get("common.no_access"), nil)
}

// 返回错误 参数错误
func ErrorParamFomat(ctx *gin.Context, errMsg string) {
	Error(ctx, language.Obj.GetAndInsert("common.api_error_param_format", "[", errMsg, "]"))
	// Error(ctx, "参数错误")
}

// // 返回错误 数据库
func ErrorDatabase(ctx *gin.Context, errMsg string) {
	// Error(ctx, global.Lang.GetAndInsert("common.db_error", "[", errMsg, "]"))
	ErrorByCodeAndMsg(ctx, 1200, errMsg)

}

// 返回错误 数据记录未找到
func ErrorDataNotFound(ctx *gin.Context) {
	ErrorByCode(ctx, -1)
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
