package dto

type CategroyRequest struct {
	Name string `josn:"name" binding:"required"`
}

type CategroyResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
