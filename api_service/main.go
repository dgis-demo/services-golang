package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Point struct {
	ID        int     `json:"id"`
	Lat       float32 `json:"lat"`
	Lon       float32 `json:"lon"`
	Magnitude float32 `json:"magnitude"`
}

func getPoint(c echo.Context) error {
	id := c.Param("id")
	point := Point{}
	point.ID, _ = strconv.Atoi(id)
	return c.JSON(http.StatusOK, &point)
}

func postPoint(c echo.Context) error {
	point := &Point{}
	if err := c.Bind(point); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, point)
}

func main() {
	// TODO: move it to another repo:
	// https://earthly.dev/blog/golang-monorepo/
	app := echo.New()

	app.GET("/:id", getPoint)
	app.POST("/", postPoint)

	app.Logger.Print(app.Start(":4001"))
}
