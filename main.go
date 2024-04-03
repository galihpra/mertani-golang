package main

import (
	"mertani-golang/config"
	"mertani-golang/helper/encrypt"
	"mertani-golang/routes"
	"mertani-golang/utils/cloudinary"
	"mertani-golang/utils/database"

	uh "mertani-golang/features/users/handler"
	ur "mertani-golang/features/users/repository"
	us "mertani-golang/features/users/service"

	sh "mertani-golang/features/sensors/handler"
	sr "mertani-golang/features/sensors/repository"
	ss "mertani-golang/features/sensors/service"

	dh "mertani-golang/features/devices/handler"
	dr "mertani-golang/features/devices/repository"
	ds "mertani-golang/features/devices/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	var dbConfig = new(config.DatabaseMysql)
	if err := dbConfig.LoadFromEnv(); err != nil {
		panic(err)
	}

	dbConnection, err := database.MysqlInit(*dbConfig)
	if err != nil {
		panic(err)
	}

	if err := database.MysqlMigrate(dbConnection); err != nil {
		panic(err)
	}

	var jwtConfig = new(config.JWT)
	if err := jwtConfig.LoadFromEnv(); err != nil {
		panic(err)
	}

	var cloudinaryConfig = new(config.Cloudinary)
	if err := cloudinaryConfig.LoadFromEnv(); err != nil {
		panic(err)
	}

	cloudinary, err := cloudinary.NewCloudinary(*cloudinaryConfig)
	if err != nil {
		panic(err)
	}

	enc := encrypt.New()
	userRepository := ur.NewUserRepository(dbConnection)
	userService := us.New(userRepository, enc)
	userHandler := uh.NewUserHandler(userService, *jwtConfig)

	sensorRepository := sr.NewSensorRepository(dbConnection, cloudinary)
	sensorService := ss.NewSensorService(sensorRepository)
	sensorHandler := sh.NewSensorHandler(sensorService, *jwtConfig)

	deviceRepository := dr.NewDeviceRepository(dbConnection, cloudinary)
	deviceService := ds.NewDeviceService(deviceRepository)
	deviceHandler := dh.NewDeviceHandler(deviceService, *jwtConfig)

	app := echo.New()
	app.Use(middleware.Recover())
	app.Use(middleware.CORS())

	route := routes.Routes{
		JWTKey:        jwtConfig.Secret,
		Server:        app,
		UserHandler:   userHandler,
		SensorHandler: sensorHandler,
		DeviceHandler: deviceHandler,
	}

	route.InitRouter()

	app.Logger.Fatal(app.Start(":8000"))
}
