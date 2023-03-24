package main

import (
	"ccrd/controller"

	"github.com/gin-gonic/gin"
)

func serveRoutes(r *gin.Engine) {

	// frontend_user
	frontend_user := controller.Frontend{}
	frontend_userGroup := r.Group("/")
	frontend_userGroup.GET("", frontend_user.UserGetHome) //index.html
	frontend_userGroup.GET("/download/:id", frontend_user.UserGetDownload)

	// AIP Razer notify
	topup_user := controller.Topup{}
	topup_Group := r.Group("/topup")
	topup_Group.GET("", topup_user.Paytopup)

	r.GET("/topups", controller.Paytopups)
	r.GET("/topups/play", controller.Payment)

	//admin
	admin_user := controller.Admin{}
	admin_userGroup := r.Group("/admin")
	admin_userGroup.GET("", admin_user.UserGetAdmin)      //index.html
	admin_userGroup.GET("/items", admin_user.GetItemsAll) //index.html

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
