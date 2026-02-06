package main

import (
	"context"
	"courseWork/server/handlers"
	"courseWork/server/storage"
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v4"
)

const dbURL = "postgres://user:password@localhost:5432/course_work_db"

func main() {
	repo, err := storage.NewAthleteRepository(context.Background(), dbURL)
	if err != nil {
		fmt.Println("Помилка старту таблиці бази даних: ", err)
		log.Println(err)
		os.Exit(1)
	}
	defer repo.Close()

	handler := handlers.NewAthleteHandler(repo)

	e := echo.New()
	defer func(e *echo.Echo) {
		_ = e.Close()
	}(e)

	g := e.Group("/athlete")
	g.POST("/create", handler.Create)
	g.GET("/fetch/all", handler.FetchAll)
	g.DELETE("/delete", handler.Delete)
	g.PUT("/update", handler.Update)
	g.GET("/fetch/sorted", handler.FetchSorted)
	g.GET("/fetch/best", handler.FetchBest)
	g.GET("/fetch/best_press_min_jump", handler.FetchBestPressMinJump)
	g.GET("/fetch/deviation_run_3km", handler.FetchWithRun3kmDeviation)
	g.GET("/fetch/min_press_run_100m_stats", handler.FetchMinPressRun100mStats)

	e.Logger.Fatal(e.Start(":1323"))
}
