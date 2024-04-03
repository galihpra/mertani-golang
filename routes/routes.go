package routes

import (
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
}

func (router Routes) InitRouter() {
	router.UserRouter()
	router.SensorRouter()
}

func (router *Routes) UserRouter() {
	router.Server.POST("/register", router.UserHandler.Register())
	router.Server.POST("/login", router.UserHandler.Login())
}

func (router *Routes) SensorRouter() {
	router.Server.POST("/sensors", router.SensorHandler.Create(), echojwt.JWT([]byte(router.JWTKey)))
	router.Server.GET("/sensors", router.SensorHandler.GetAll())
}
