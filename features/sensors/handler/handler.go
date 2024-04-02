package handler

import (
	"fmt"
	"mertani-golang/config"
	"mertani-golang/features/sensors"
	"mertani-golang/helper/tokens"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type sensorHandler struct {
	service   sensors.Service
	jwtConfig config.JWT
}

func NewSensorHandler(service sensors.Service, jwtConfig config.JWT) sensors.Handler {
	return &sensorHandler{
		service:   service,
		jwtConfig: jwtConfig,
	}
}

func (hdl *sensorHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		var request = new(CreateRequest)
		var response = make(map[string]any)

		token := c.Get("user")
		if token == nil {
			response["message"] = "unauthorized access"
			return c.JSON(http.StatusUnauthorized, response)
		}

		userId, err := tokens.ExtractToken(hdl.jwtConfig.Secret, token.(*jwt.Token))
		if err != nil {
			c.Logger().Error(err)

			response["message"] = "unauthorized"
			return c.JSON(http.StatusUnauthorized, response)
		}

		if err := c.Bind(request); err != nil {
			c.Logger().Error(err)

			response["message"] = "incorect input data"
			return c.JSON(http.StatusBadRequest, response)
		}

		var parseInput = new(sensors.Sensor)
		parseInput.Name = request.Name
		parseInput.Description = request.Description
		parseInput.UserId = userId

		file, err := c.FormFile("image")
		if err != nil {
			response["message"] = "error retrieving uploaded file"
			return c.JSON(http.StatusBadRequest, response)
		}
		if file != nil {
			src, err := file.Open()
			if err != nil {
				return err
			}
			defer src.Close()

			request.Image = src
		}

		parseInput.ImageRaw = request.Image

		for _, detail := range request.Details {
			parseInput.Details = append(parseInput.Details, sensors.Detail{
				Spec: detail.Spec,
			})
		}

		fmt.Println(parseInput)

		if err := hdl.service.Create(*parseInput); err != nil {
			c.Logger().Error(err)

			if strings.Contains(err.Error(), "validate") {
				response["message"] = strings.ReplaceAll(err.Error(), "validate: ", "")
				return c.JSON(http.StatusBadRequest, response)
			}

			if strings.Contains(err.Error(), "unauthorized") {
				response["message"] = "unauthorized"
				return c.JSON(http.StatusBadRequest, response)
			}

			response["message"] = "internal server error"
			return c.JSON(http.StatusInternalServerError, response)
		}

		response["message"] = "create sensor success"
		return c.JSON(http.StatusCreated, response)
	}
}
