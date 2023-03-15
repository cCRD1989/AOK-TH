package main

import (
	"ccrd/controller"
	"ccrd/middleware"

	"github.com/gin-gonic/gin"
)

func serveRoutes(r *gin.Engine) {

	//category
	categoryController := controller.Categroy{}
	categoryGroup := r.Group("/categorys")
	categoryGroup.GET("", categoryController.FindAll)
	categoryGroup.GET("/:id", categoryController.FindOne)
	categoryGroup.POST("", categoryController.Create)
	categoryGroup.PATCH("/:id", categoryController.Update)
	categoryGroup.DELETE("/:id", categoryController.Delete)

	//Products
	productController := controller.Products{}
	productGroup := r.Group("/products")
	productGroup.GET("", productController.FindAll)
	productGroup.GET("/:id", productController.FindOne)
	productGroup.POST("", productController.Create)
	productGroup.PATCH("/:id", productController.Update)
	productGroup.DELETE("/:id", productController.Delete)

	//orders
	orderController := controller.Order{}
	orderGroup := r.Group("/orders")
	orderGroup.GET("", orderController.FindAll)
	orderGroup.GET("/:id", orderController.FindOne)
	orderGroup.POST("", orderController.Create)

	//WeloveKhan
	userController := controller.User{}
	userGroup := r.Group("")

	userGroup.POST("/login", userController.ChaeckLogin_jwt)
	userGroup.POST("/registered", userController.Registered)

	userAll := r.Group("/users", middleware.JWTAuth())
	userAll.GET("/readallprofile", userController.ChaeckLoginAll)
	userAll.GET("/profile", userController.Profile)

	userAll.GET("/111", userController.FFFFF)

}
