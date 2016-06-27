package backend

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func createToken(u *User) *jwt.Token {
	tok := jwt.New(jwt.SigningMethodRS256)
	cl := tok.Claims.(jwt.MapClaims)
	cl["iss"] = serviceAccountEmail
	cl["sub"] = serviceAccountEmail
	cl["aud"] = "https://identitytoolkit.googleapis.com/google.identity.identitytoolkit.v1.IdentityToolkit"
	now := time.Now().Unix()
	cl["iat"] = int(now)
	cl["exp"] = int(now + (60 * 60))
	cl["uid"] = fmt.Sprintf("%d", u.ID)
	cl["claims"] = map[string]interface{}{
		"premium_account": 1,
	}
	return tok
}

func parseTokenString(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return privateKey.Public(), nil
	})
}

var TokenLookup = "header: " + echo.HeaderAuthorization

func JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}
