package routes

import (
	"mertani-golang/features/devices"
	"mertani-golang/features/sensors"
	"mertani-golang/features/users"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type Routes struct {
	JWTKey        string
	Server        *echo.Echo
	UserHandler   users.Handler
	SensorHandler sensors.Handler
	DeviceHandler devices.Handler
}

func (router Routes) InitRouter() {
	router.UserRouter()
	router.SensorRouter()
	router.DeviceRouter()
}

func (router *Routes) UserRouter() {
	router.Server.POST("/register", router.UserHandler.Register())
	router.Server.POST("/login", router.UserHandler.Login())
}

func (router *Routes) SensorRouter() {
	router.Server.POST("/sensors", router.SensorHandler.Create(), echojwt.JWT([]byte(router.JWTKey)))
	router.Server.GET("/sensors", router.SensorHandler.GetAll())
	router.Server.DELETE("/sensors/:id", router.SensorHandler.Delete(), echojwt.JWT([]byte(router.JWTKey)))
	router.Server.PATCH("/sensors/:id", router.SensorHandler.Update(), echojwt.JWT([]byte(router.JWTKey)))
}

func (router *Routes) DeviceRouter() {
	router.Server.POST("/devices", router.DeviceHandler.Create(), echojwt.JWT([]byte(router.JWTKey)))
	router.Server.GET("/devices", router.DeviceHandler.GetAll())
	router.Server.DELETE("/devices/:id", router.DeviceHandler.Delete(), echojwt.JWT([]byte(router.JWTKey)))
	router.Server.PATCH("/devices/:id", router.DeviceHandler.Update(), echojwt.JWT([]byte(router.JWTKey)))
}
