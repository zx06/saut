package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4/request"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const CtxKeyAuth = "auth_ctx"

type Auth struct {
	secret []byte
}

func NewAuth() *Auth {
	return &Auth{
		secret: []byte("7LJtjeioQM"),
	}
}

func (a *Auth) GrantToken(uid string) (string, error) {
	claims := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&Claims{
			Issuer:    uid,
			Subject:   "saut",
			Audience:  nil,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-10 * 60 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	)
	signingString, err := claims.SigningString()
	if err != nil {
		return "", nil
	}
	return signingString, err
}

func (a *Auth) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			claims *Claims
			ok     bool
		)
		_, err := request.ParseFromRequest(c.Request, request.OAuth2Extractor, func(tk *jwt.Token) (interface{}, error) {
			if claims, ok = tk.Claims.(*Claims); !ok {
				return nil, fmt.Errorf("验证失败")
			}
			c.Set(CtxKeyAuth, claims)
			return a.secret, nil
		}, request.WithClaims(&Claims{}))
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
	}
}

func (a *Auth) GetClaims(c *gin.Context) (*Claims, error) {
	var (
		v      any
		claims *Claims
		ok     bool
	)
	if v, ok = c.Get(CtxKeyAuth); !ok {
		return nil, fmt.Errorf("获取 Claims 失败")
	}
	if claims, ok = v.(*Claims); !ok {
		return nil, fmt.Errorf("获取 Claims 失败")
	}
	return claims, nil
}
