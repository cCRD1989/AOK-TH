package middleware

import (
	"ccrd/model/aokmodel"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func UserSession() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Get the Cookie off req
		tokenString, err := ctx.Cookie("Authorization")
		if err != nil {
			fmt.Println("Get the Cookie off")
			ctx.Next()
			return
		}

		// Decode/Validate
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("MY_SECRET_KEY")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check the exp
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				fmt.Println("token Exp.")
				ctx.Next()
				return
			}

			// Find the user wuth token
			jsonData, _ := json.Marshal(claims["user"])
			authedUser := aokmodel.Userlogin{}
			json.Unmarshal(jsonData, &authedUser)

			// Find the user wuth token St
			user := aokmodel.Userlogin{}
			user = user.FindUserByName(authedUser.Username)
			if user.Username == "" {
				fmt.Println("Check ID Token")
				ctx.Next()
				return
			}

			// Attach to req
			ctx.Set("user", user)

			// Continue
			ctx.Next()

		} else {

			fmt.Println("token.Claims")
			ctx.Next()
			return
		}
	}
}
