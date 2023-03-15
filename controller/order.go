package controller

import (
	"ccrd/db"
	"ccrd/dto"
	"ccrd/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Order struct{}

func (o Order) FindAll(ctx *gin.Context) {
	var orderts []model.Order
	db.Conn.Preload("Product").Find(&orderts)

	var result []dto.OrderResponse
	for _, order := range orderts {
		oederResult := dto.OrderResponse{
			ID:    order.ID,
			Name:  order.Name,
			Tel:   order.Tel,
			Email: order.Email,
		}
		var products []dto.OrderProductResponse
		for _, product := range order.Product {
			products = append(products, dto.OrderProductResponse{
				ID:       product.ID,
				SKU:      product.SKU,
				Name:     product.Name,
				Price:    product.Price,
				Quantity: product.Quatity,
				Image:    product.Image,
			})
		}
		oederResult.Products = products
		result = append(result, oederResult)
	}
	ctx.JSON(http.StatusOK, result)
}

func (o Order) FindOne(ctx *gin.Context) {

	id := ctx.Param("id")
	var order model.Order
	query := db.Conn.Preload("Product").First(&order, id)

	if err := query.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	result := dto.OrderResponse{
		ID:    order.ID,
		Name:  order.Name,
		Tel:   order.Tel,
		Email: order.Email,
	}
	var products []dto.OrderProductResponse
	for _, product := range order.Product {
		products = append(products, dto.OrderProductResponse{
			ID:       product.ID,
			SKU:      product.SKU,
			Name:     product.Name,
			Price:    product.Price,
			Quantity: product.Quatity,
			Image:    product.Image,
		})
	}
	result.Products = products

	ctx.JSON(http.StatusOK, result)

}

func (o Order) Create(ctx *gin.Context) {
	var form dto.OrderRequest
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"erroe": err.Error()})
		return
	}

	var order model.Order
	var orderItem []model.OrderItem

	for _, product := range form.Products {
		orderItem = append(orderItem, model.OrderItem{
			SKU:     product.SKU,
			Name:    product.Name,
			Image:   product.Image,
			Price:   product.Price,
			Quatity: product.Quantity,
		})
	}
	order.Name = form.Name
	order.Tel = form.Tel
	order.Email = form.Email
	order.Product = orderItem
	db.Conn.Create(&order)

	result := dto.OrderResponse{
		ID:    order.ID,
		Name:  order.Name,
		Tel:   order.Tel,
		Email: order.Email,
	}
	var products []dto.OrderProductResponse
	for _, product := range order.Product {
		products = append(products, dto.OrderProductResponse{
			ID:       product.ID,
			SKU:      product.SKU,
			Name:     product.Name,
			Price:    product.Price,
			Quantity: product.Quatity,
			Image:    product.Image,
		})
	}
	result.Products = products
	ctx.JSON(http.StatusOK, result)

}
