package handler

type GetResponse struct {
	Id          uint   `json:"sensor_id"`
	Name        string `json:"sensor_name"`
	Image       string `json:"image"`
	Description string `json:"description"`
	UserId      uint   `json:"userId"`
}
