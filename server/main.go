package main

import (
	"context"
	"courseWork/shared"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

const dbURL = "postgres://user:password@localhost:5432/course_work_db"

func main() {
	pool, err := StartupTable(context.Background(), dbURL)
	if err != nil {
		fmt.Println("Помилка старту таблиці бази даних: ", err)
		log.Println(err)
		os.Exit(1)
	}
	defer pool.Close()

	e := echo.New()
	defer func(e *echo.Echo) {
		_ = e.Close()
	}(e)

	e.POST("/athlete/create", func(c echo.Context) error {
		a := new(shared.Athlete)
		err := c.Bind(a)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		err = AddField(*a, pool, c.Request().Context())
		if err != nil {
			return err
		}

		return c.NoContent(http.StatusOK)
	})

	e.GET("/athlete/fetch/all", func(c echo.Context) error {
		athletes, err := SelectTable(pool, c.Request().Context())
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, athletes)
	})

	e.PUT("/athlete/update", func(c echo.Context) error {
		a := &shared.Athlete{}
		err := c.Bind(a)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = UpdateField(*a, pool, c.Request().Context())
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusOK)
	})

	e.GET("/athlete/fetch/best", func(c echo.Context) error {
		athletes, err := SelectBestTotalResult(pool, c.Request().Context())
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, athletes)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
