package dto

type UserRequest struct {
	Username string `josn:"username" binding:"required"`
	Password string `josn:"password" binding:"required"`
}

type UserResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Status   string `json:"status"`
	JWT      string `json:"jwt"`
}

type RegiseredRequest struct {
	Username string `josn:"username" binding:"required"`
	Password string `josn:"password" binding:"required"`
	Idcode   string `josn:"idcode" binding:"required"`
}

type RegiseredResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Idcode   string `josn:"idcode"`
	Status   string `josn:"status"`
}
