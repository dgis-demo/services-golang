package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Point struct {
	ID        int     `json:"id" gorm:"primarykey"`
	Lat       float32 `json:"lat"`
	Lon       float32 `json:"lon"`
	Magnitude float32 `json:"magnitude"`
}

func getPoint(c echo.Context) error {
	id := c.Param("id")
	point := Point{}
	pointID, _ := strconv.Atoi(id)

	db, _ := c.Get("db").(*gorm.DB)
	err := db.First(&point, pointID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, &point)
}

func postPoint(c echo.Context) error {
	point := Point{}
	if err := c.Bind(&point); err != nil {
		log.Println(err)
		return err
	}

	db, _ := c.Get("db").(*gorm.DB)
	db.Create(&point)
	return c.JSON(http.StatusCreated, &point)
}

func main() {
	// TODO: move it to another repo:
	// https://earthly.dev/blog/golang-monorepo/

	db, err := gorm.Open(sqlite.Open(DB_NAME), &gorm.Config{})
	if err != nil {
		log.Printf(
			"could not connect to the database %s error: %s",
			DB_NAME,
			err,
		)
		return
	}
	db.AutoMigrate(&Point{})

	app := echo.New()
	app.Use(ContextDB(db))

	app.GET("/:id", getPoint)
	app.POST("/", postPoint)

	app.Logger.Print(app.Start(":4001"))
}
