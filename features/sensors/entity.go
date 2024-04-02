package sensors

import (
	"io"

	"github.com/labstack/echo/v4"
)

type Sensor struct {
	Id          uint
	Name        string
	Description string

	UserId uint

	ImageUrl string
	ImageRaw io.Reader

	Details []Detail
}

type Detail struct {
	Id   uint
	Spec string
}

type Handler interface {
	Create() echo.HandlerFunc
}

type Service interface {
	Create(newSensor Sensor) error
}

type Repository interface {
	Create(newSensor Sensor) error
}
