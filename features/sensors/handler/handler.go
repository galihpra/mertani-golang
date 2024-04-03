package handler

import (
	"mertani-golang/config"
	"mertani-golang/features/sensors"
	"mertani-golang/helper/tokens"
	"net/http"
	"strconv"
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

		file, _ := c.FormFile("image")
		if file != nil {
			src, err := file.Open()
			if err != nil {
				return err
			}
			defer src.Close()

			parseInput.ImageRaw = src
		}

		if err := hdl.service.Create(*parseInput); err != nil {
			c.Logger().Error(err)

			if strings.Contains(err.Error(), "validate") {
				response["message"] = strings.ReplaceAll(err.Error(), "validate: ", "")
				return c.JSON(http.StatusBadRequest, response)
			}

			if strings.Contains(err.Error(), "duplicate") {
				response["message"] = strings.ReplaceAll(err.Error(), "duplicate: ", "")
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

func (hdl *sensorHandler) GetAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		var response = make(map[string]any)

		result, err := hdl.service.GetAll()
		if err != nil {
			c.Logger().Error(err)

			response["message"] = "internal server error"
			return c.JSON(http.StatusInternalServerError, response)
		}

		var data = make([]GetResponse, len(result))
		for i, sensor := range result {
			data[i] = GetResponse{
				Id:          sensor.Id,
				Name:        sensor.Name,
				Description: sensor.Description,
				Image:       sensor.ImageUrl,
				UserId:      sensor.UserId,
			}
		}

		response["message"] = "get all sensors success"
		response["data"] = data
		return c.JSON(http.StatusOK, response)
	}
}

func (hdl *sensorHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		var response = make(map[string]any)

		SensorId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error(err)

			response["message"] = "invalid sensor id"
		}

		token := c.Get("user")
		if token == nil {
			response["message"] = "unauthorized access"
			return c.JSON(http.StatusUnauthorized, response)
		}

		UserId, err := tokens.ExtractToken(hdl.jwtConfig.Secret, token.(*jwt.Token))
		if err != nil {
			c.Logger().Error(err)

			response["message"] = "unauthorized"
			return c.JSON(http.StatusUnauthorized, response)
		}

		if err := hdl.service.Delete(uint(SensorId), uint(UserId)); err != nil {
			c.Logger().Error(err)

			if strings.Contains(err.Error(), "not found: ") {
				response["message"] = strings.ReplaceAll(err.Error(), "not found: ", "")
				return c.JSON(http.StatusNotFound, response)
			}

			if strings.Contains(err.Error(), "not authorized: ") {
				response["message"] = strings.ReplaceAll(err.Error(), "not authorized: ", "")
				return c.JSON(http.StatusNotFound, response)
			}

			response["message"] = "internal server error"
			return c.JSON(http.StatusInternalServerError, response)
		}

		response["message"] = "delete sensor success"
		return c.JSON(http.StatusOK, response)
	}
}

func (hdl *sensorHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		var response = make(map[string]any)
		var request = new(CreateRequest)

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error(err)

			response["message"] = "invalid sensor id"
		}

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

		if c.Bind(request); err != nil {
			c.Logger().Error(err)

			response["message"] = "bad request"
			return c.JSON(http.StatusBadRequest, response)
		}

		var parseInput = new(sensors.Sensor)
		parseInput.Id = uint(id)
		parseInput.Name = request.Name
		parseInput.Description = request.Description

		file, _ := c.FormFile("image")
		if file != nil {
			src, err := file.Open()
			if err != nil {
				return err
			}
			defer src.Close()

			parseInput.ImageRaw = src
		}

		if err := hdl.service.Update(uint(userId), *parseInput); err != nil {
			c.Logger().Error(err)

			if strings.Contains(err.Error(), "validate: ") {
				response["message"] = strings.ReplaceAll(err.Error(), "validate: ", "")
				return c.JSON(http.StatusBadRequest, response)
			}

			if strings.Contains(err.Error(), "not found: ") {
				response["message"] = strings.ReplaceAll(err.Error(), "not found: ", "")
				return c.JSON(http.StatusNotFound, response)
			}

			if strings.Contains(err.Error(), "duplicate") {
				response["message"] = strings.ReplaceAll(err.Error(), "duplicate: ", "")
				return c.JSON(http.StatusBadRequest, response)
			}

			if strings.Contains(err.Error(), "not authorized: ") {
				response["message"] = strings.ReplaceAll(err.Error(), "not authorized: ", "")
				return c.JSON(http.StatusNotFound, response)
			}

			response["message"] = "internal server error"
			return c.JSON(http.StatusInternalServerError, response)
		}

		response["message"] = "update sensor success"
		return c.JSON(http.StatusOK, response)
	}
}
