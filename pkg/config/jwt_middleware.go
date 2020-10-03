package config

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

const (
	// KeyHeaderMemberType ---
	KeyHeaderMemberType = "x-member-type"
	// KeyHeaderMemberID ---
	KeyHeaderMemberID = "x-member-id"
)

// CustomJWTMiddleware returns echo middleware jwt with config
func CustomJWTMiddleware(pathsToSkipped ...string) echo.MiddlewareFunc {
	secret := viper.GetString("secret")
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:     []byte(secret),
		SigningMethod:  "HS256",
		SuccessHandler: successHandler,
		Skipper:        getSkipperFunc(pathsToSkipped...),
		ErrorHandler: func(err error) error {
			return echo.ErrUnauthorized
		},
	})
}

func getSkipperFunc(paths ...string) middleware.Skipper {
	return func(ctx echo.Context) bool {
		for _, p := range paths {
			if strings.Contains(ctx.Request().RequestURI, p) {
				return true
			}
		}
		return false
	}
}

func successHandler(c echo.Context) {
	accessToken := strings.Replace(c.Request().Header.Get(echo.HeaderAuthorization), "Bearer ", "", 1)
	parsedToken, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		secret := viper.GetString("secret")
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})

	if err != nil {
		return
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		h := &c.Request().Header
		h.Set(KeyHeaderMemberID, claims["unm"].(string))
		h.Set(KeyHeaderMemberType, claims["mty"].(string))
	}
	return
}
