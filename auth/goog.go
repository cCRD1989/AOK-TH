package auth

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/zalando/gin-oauth2/google"
)

var redirectURL, credFile string

func init() {
	bin := path.Base(os.Args[0])
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s
================
`, bin)
		flag.PrintDefaults()
	}
	flag.StringVar(&redirectURL, "redirect", "http://127.0.0.1/auth/google", "URL to be redirected to after authorization.")
	flag.StringVar(&credFile, "cred-file", "./test-clientid.google.json", "Credential JSON file")
}
func Auth() {
	//flag.Parse()

	scopes := []string{
		"https://www.googleapis.com/auth/userinfo.profile",
	}
	secret := []byte("secret")
	sessionName := "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"

	router := gin.Default()
	// init settings for google auth
	google.Setup(redirectURL, credFile, scopes, secret)
	router.Use(google.Session(sessionName))

	router.GET("/login", google.LoginHandler)

	// protected url group
	private := router.Group("/auth")
	private.Use(google.Auth())
	private.GET("/", UserInfoHandler)
	private.GET("/api", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Hello from private for groups"})
	})

	router.Run("127.0.0.1:8081")
}

func UserInfoHandler(ctx *gin.Context) {
	name := ctx.MustGet("user").(google.User)
	ctx.JSON(http.StatusOK, gin.H{"Hello": "from private", "user": name.Name})

}
