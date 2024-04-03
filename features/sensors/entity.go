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
}

type Handler interface {
	Create() echo.HandlerFunc
	GetAll() echo.HandlerFunc
	Delete() echo.HandlerFunc
	Update() echo.HandlerFunc
}

type Service interface {
	Create(newSensor Sensor) error
	GetAll() ([]Sensor, error)
	Delete(SensorId, UserId uint) error
	Update(UserId uint, updateSensor Sensor) error
}

type Repository interface {
	Create(newSensor Sensor) error
	GetAll() ([]Sensor, error)
	Delete(SensorId, UserId uint) error
	Update(UserId uint, updateSensor Sensor) error
}
