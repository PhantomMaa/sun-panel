package jwt

import (
	"sun-panel/global"
	"sun-panel/lib/jwt"
)

// InitJWT 初始化JWT配置
func InitJWT() error {
	// 从配置文件读取JWT配置
	secret := global.Config.GetValueString("jwt", "secret")
	if secret == "" {
		secret = "sun-panel-default-jwt-secret-key" // 默认密钥
	}

	expire := global.Config.GetValueInt("jwt", "expire")
	if expire <= 0 {
		expire = 72 // 默认72小时
	}

	// 初始化JWT密钥
	jwt.InitSecret(secret)
	jwt.SetExpire(expire)
	
	return nil
}
