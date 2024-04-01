package main

import (
	"mertani-golang/config"
	"mertani-golang/helper/encrypt"
	"mertani-golang/routes"
	"mertani-golang/utils/database"

	uh "mertani-golang/features/users/handler"
	ur "mertani-golang/features/users/repository"
	us "mertani-golang/features/users/service"

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

	enc := encrypt.New()
	userRepository := ur.NewUserRepository(dbConnection)
	userService := us.New(userRepository, enc)
	userHandler := uh.NewUserHandler(userService, *jwtConfig)

	app := echo.New()
	app.Use(middleware.Recover())
	app.Use(middleware.CORS())

	route := routes.Routes{
		JWTKey:      jwtConfig.Secret,
		Server:      app,
		UserHandler: userHandler,
	}

	route.InitRouter()

	app.Logger.Fatal(app.Start(":8000"))
}
