package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig JWT配置
type JWTConfig struct {
	Secret  string
	Timeout int64
}

// Claims 自定义声明结构体
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWT 工具结构体
type JWT struct {
	Config JWTConfig
}

// NewJWT 创建JWT实例
func NewJWT(config JWTConfig) *JWT {
	return &JWT{
		Config: config,
	}
}

// GenerateToken 生成token
func (j *JWT) GenerateToken(userID uint, username string) (string, error) {
	// 设置token过期时间
	expireTime := time.Now().Add(time.Duration(j.Config.Timeout) * time.Second)

	// 创建声明
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "rentpro-admin",
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名token
	tokenString, err := token.SignedString([]byte(j.Config.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析token
func (j *JWT) ParseToken(tokenString string) (*Claims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证token
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
