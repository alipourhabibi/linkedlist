package v2

import (
	"context"
	"fmt"
	"linkedlist/linkedlist"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ListEntity struct {
	Index uint `json:"index" param:"index"`
	Value int  `json:"value" param:"value"`
}

type server struct {
	list *linkedlist.LinkedList
}

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func V2() (*echo.Echo, error) {
	e := echo.New()

	logger := slog.Default()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	e.Validator = &customValidator{validator: validator.New()}

	s := &server{}
	l := linkedlist.NewLinkedList()
	s.list = l
	e.POST("/numbers/:index/:value", s.Insert)
	e.DELETE("/numbers/:index", s.Remove)
	e.GET("/numbers/value/:value", s.Find)
	e.GET("/numbers/index/:index", s.Get)

	return e, nil
}

func (s *server) Insert(c echo.Context) error {
	data := ListEntity{}

	if err := c.Bind(&data); err != nil {
		return err
	}
	if err := c.Validate(&data); err != nil {
		return err
	}

	ok := s.list.Insert(data.Index, data.Value)
	if !ok {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid index")
	}
	c.JSON(http.StatusCreated, data)
	return nil
}

func (s *server) Remove(c echo.Context) error {
	indexStr := c.Param("index")
	index, err := strconv.ParseUint(indexStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid index")
	}

	ok := s.list.Remove(uint(index))
	if !ok {
		return echo.NewHTTPError(echo.ErrNotFound.Code, "Index not found")
	}

	c.NoContent(http.StatusOK)
	return nil

}

func (s *server) Find(c echo.Context) error {
	valueStr := c.Param("value")
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid value")
	}

	index, ok := s.list.Find(value)
	if !ok {
		return echo.NewHTTPError(echo.ErrNotFound.Code, "Value not found")
	}

	data := ListEntity{
		Index: index,
		Value: value,
	}
	c.JSON(http.StatusOK, data)
	return nil
}

func (s *server) Get(c echo.Context) error {

	indexStr := c.Param("index")
	index, err := strconv.ParseUint(indexStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, "Invalid index")
	}

	value, ok := s.list.Get(uint(index))
	if !ok {
		return echo.NewHTTPError(echo.ErrNotFound.Code, "Index not found")
	}
	data := ListEntity{
		Index: uint(index),
		Value: value,
	}

	c.JSON(http.StatusOK, data)
	return nil
}
