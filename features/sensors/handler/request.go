package handler

import "io"

type CreateRequest struct {
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`

	Image io.Reader

	UserId uint

	Details []Detail
}

type Detail struct {
	Spec string `json:"spec" form:"spec"`
}
