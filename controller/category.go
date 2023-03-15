package controller

import (
	"ccrd/db"
	"ccrd/dto"
	"ccrd/model"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Categroy struct{}

func (p Categroy) FindAll(ctx *gin.Context) {
	var category []model.Category
	db.Conn.Find(&category)
	var resulf []dto.CategroyResponse
	for _, v := range category {
		resulf = append(resulf, dto.CategroyResponse{
			ID:   v.ID,
			Name: v.Name,
		})
	}
	ctx.JSON(http.StatusOK, resulf)
}

func (p Categroy) FindOne(ctx *gin.Context) {

	id := ctx.Param("id")
	var category model.Category
	if err := db.Conn.First(&category, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dto.CategroyResponse{
		ID:   category.ID,
		Name: category.Name,
	})

}

//post /category <<= {"name":"Flower"} {JSON}

func (p Categroy) Create(ctx *gin.Context) {

	var form dto.CategroyRequest
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	catrgory := model.Category{
		Name: form.Name,
	}

	if err := db.Conn.Create(&catrgory).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.CategroyResponse{
		ID:   catrgory.ID,
		Name: catrgory.Name,
	})

}

func (p Categroy) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var form dto.CategroyRequest
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var category model.Category
	if err := db.Conn.First(&category, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	category.Name = form.Name
	db.Conn.Save(&category)
	ctx.JSON(http.StatusOK, dto.CategroyResponse{
		ID:   category.ID,
		Name: category.Name,
	})

}

func (p Categroy) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	//db.Conn.Delete(&model.Category{}, id)
	db.Conn.Unscoped().Delete(&model.Category{}, id)

	ctx.JSON(http.StatusOK, gin.H{"deletedAt": time.Now()})

}
