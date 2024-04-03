package handler

type GetResponse struct {
	Id          uint   `json:"device_id"`
	Name        string `json:"device_name"`
	Image       string `json:"image"`
	Description string `json:"description"`
	UserId      uint   `json:"userId"`
}
