package handler

import "io"

type CreateRequest struct {
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`

	Image io.Reader

	UserId uint
}
