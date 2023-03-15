package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		Header := ctx.Request.Header.Get("Authorization")
		tokenString := strings.Replace(Header, "Bearer ", "", 1)

		hmacSampleSecret := []byte(os.Getenv("TOKEN_KEY"))
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return hmacSampleSecret, nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			ctx.Set("userId", claims["userId"])
		} else {

			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Status": "error", "message": err.Error()})
		}

		ctx.Next()

	}
}
