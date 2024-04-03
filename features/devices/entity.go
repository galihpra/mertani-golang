package devices

import (
	"io"

	"github.com/labstack/echo/v4"
)

type Device struct {
	Id          uint
	Name        string
	Description string

	UserId uint

	ImageUrl string
	ImageRaw io.Reader
}

type Handler interface {
	Create() echo.HandlerFunc
	GetAll() echo.HandlerFunc
	Delete() echo.HandlerFunc
	Update() echo.HandlerFunc
}

type Service interface {
	Create(newDevice Device) error
	GetAll() ([]Device, error)
	Delete(DeviceId, UserId uint) error
	Update(UserId uint, updateDevice Device) error
}

type Repository interface {
	Create(newDevice Device) error
	GetAll() ([]Device, error)
	Delete(DeviceId, UserId uint) error
	Update(UserId uint, updateDevice Device) error
}
