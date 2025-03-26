package interceptor

import (
	"errors"
	"sun-panel/internal/infra/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 自定义JWT claims结构
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

var (
	secretKey       []byte
	expireHours     = 72                          // 默认72小时
	ErrInvalidToken = errors.New("invalid token") // 自定义错误
)

// InitJWT 初始化JWT配置
func InitJWT() error {
	// 从配置文件读取JWT配置
	secret := config.AppConfig.JWT.Secret
	if secret == "" {
		secret = "sun-panel-default-jwt-secret-key" // 默认密钥
	}

	expire := config.AppConfig.JWT.Expire
	if expire <= 0 {
		expire = 72 // 默认72小时
	}

	// 初始化JWT密钥
	InitSecret(secret)
	SetExpire(expire)

	return nil
}

// InitSecret 初始化JWT密钥
func InitSecret(secret string) {
	secretKey = []byte(secret)
}

// SetExpire 设置Token过期时间（小时）
func SetExpire(hours int) {
	expireHours = hours
}

// GenerateToken 生成JWT Token
func GenerateToken(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   "user_access",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
