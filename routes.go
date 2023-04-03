package main

import (
	"ccrd/controller"

	"github.com/gin-gonic/gin"
	"github.com/zalando/gin-oauth2/google"
)

func serveRoutes(r *gin.Engine) {

	// frontend_user
	frontend_user := controller.Frontend{}
	frontend_userGroup := r.Group("/")
	frontend_userGroup.GET("", frontend_user.UserGetHome) //index.html
	frontend_userGroup.GET("/download/:id", frontend_user.UserGetDownload)

	// auth https://console.cloud.google.com/
	auth_user := r.Group("/auth")
	auth_user.GET("/", google.LoginHandler) //google add 	ctx.Redirect(http.StatusTemporaryRedirect, GetLoginURL(state))
	auth_user.GET("/google/registered", frontend_user.Auth_google_Regis)
	auth_user.GET("/custom", frontend_user.Auth_custom)
	auth_user.GET("/customregis", frontend_user.Auth_custom_regis)
	auth_user.GET("/facebooklogin", frontend_user.Auth_facebook_login)
	auth_user.GET("/facebookCall", frontend_user.Auth_facebook_call)
	auth_user.GET("/facebookRegis", frontend_user.Auth_facebook_regis)

	private := r.Group("/auth")
	private.Use(google.Auth())
	private.GET("/google", frontend_user.Auth_google) //index.html

	// AIP Razer notify
	topup_user := controller.Topup{}
	topup_Group := r.Group("/topup")
	topup_Group.GET("", topup_user.Paytopup)                 // API notify url
	topup_Group.GET("/processingpay", topup_user.PayProcess) //redirect url

	//
	r.GET("/topups", controller.Paytopups) // เปิดหน้าเติมเงิน
	r.GET("/topups/:user", controller.UserCheck)
	r.GET("/topups/play", controller.Payment) // เมื่อลูกค้า กด ออเดอร์ เข้ามา

	//admin
	admin_user := controller.Admin{}
	admin_userGroup := r.Group("/admin")
	admin_userGroup.GET("", admin_user.UserGetAdmin)      //index.html
	admin_userGroup.GET("/items", admin_user.GetItemsAll) //index.html
	admin_userGroup.GET("/logtopup", admin_user.Logtopup)

	// //category
	// categoryController := controller.Categroy{}
	// categoryGroup := r.Group("/categorys")
	// categoryGroup.GET("", categoryController.FindAll)
	// categoryGroup.GET("/:id", categoryController.FindOne)
	// categoryGroup.POST("", categoryController.Create)
	// categoryGroup.PATCH("/:id", categoryController.Update)
	// categoryGroup.DELETE("/:id", categoryController.Delete)

	// //Products
	// productController := controller.Products{}
	// productGroup := r.Group("/products")
	// productGroup.GET("", productController.FindAll)
	// productGroup.GET("/:id", productController.FindOne)
	// productGroup.POST("", productController.Create)
	// productGroup.PATCH("/:id", productController.Update)
	// productGroup.DELETE("/:id", productController.Delete)

	// //orders
	// orderController := controller.Order{}
	// orderGroup := r.Group("/orders")
	// orderGroup.GET("", orderController.FindAll)
	// orderGroup.GET("/:id", orderController.FindOne)
	// orderGroup.POST("", orderController.Create)

	// //WeloveKhan
	// userController := controller.User{}
	// userGroup := r.Group("")

	// userGroup.POST("/login", userController.ChaeckLogin_jwt)
	// userGroup.POST("/registered", userController.Registered)

	// userAll := r.Group("/users", middleware.JWTAuth())
	// userAll.GET("/readallprofile", userController.ChaeckLoginAll)
	// userAll.GET("/profile", userController.Profile)

	// userAll.GET("/111", userController.FFFFF)

}
