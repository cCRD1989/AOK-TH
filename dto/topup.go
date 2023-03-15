package dto

type TopupRequest struct {
	Txid     string `josn:"txid" binding:"required"`
	Orderid  string `josn:"orderid" binding:"required"`
	Status   string `josn:"status" binding:"required"`
	Detail   string `josn:"detail" binding:"required"`
	Channel  string `josn:"channel" binding:"required"`
	Amount   string `josn:"amount" binding:"required"`
	Currency string `josn:"currency" binding:"required"`
	Sig      string `josn:"sig" binding:"required"`
}

type TopupResponse struct {
	Txid   string `json:"txid"`
	Status string `json:"status"`
}
