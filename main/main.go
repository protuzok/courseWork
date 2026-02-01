package main

import (
	"context"
	"courseWork/backend"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const dbURL = "postgres://user:password@localhost:5432/course_work_db"

func main() {
	logFile, err := startupLogger()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Logger start up")

	defer func() {
		err = logFile.Close()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error closing log file: %v\n", err)
		}
	}()

	// ---------------------------------------------

	ctx := context.Background()

	pool, err := backend.StartupTable(ctx, dbURL)
	if err != nil {
		fmt.Println("Помилка старту таблиці бази даних: ", err)
		log.Println(err)
		os.Exit(1)
	}
	defer pool.Close()

MainLoop:
	for {
		printOptions()

		opt, err := takeOption()
		if err != nil {
			fmt.Println("Помилка вводу:", err)
			continue
		}

		switch opt {
		case 0:
			return
		case 1:
			inputRecord := takeInput("Введіть значення полів запису:")

			fields := strings.Fields(inputRecord)
			if len(fields) == 0 {
				fmt.Println("Не введено жодного поля.")
				fmt.Println("Спробуйте ще раз")
				continue
			}

			parseRun100m, err := strconv.ParseFloat(fields[2], 64)
			if err != nil {
				fmt.Println("Спробуйте ще раз")
				continue
			}

			parseRun3km, err := strconv.ParseFloat(fields[3], 64)
			if err != nil {
				fmt.Println("Спробуйте ще раз")
				continue
			}

			parsePressCnt, err := strconv.ParseInt(fields[4], 10, 64)
			if err != nil {
				fmt.Println("Спробуйте ще раз")
				continue
			}

			parseJumpDistance, err := strconv.ParseFloat(fields[5], 64)
			if err != nil {
				fmt.Println("Спробуйте ще раз")
				continue
			}

			a := backend.Athlete{
				Id:           -1,
				Name:         fields[0],
				Surname:      fields[1],
				Run100m:      float32(parseRun100m),
				Run3km:       float32(parseRun3km),
				PressCnt:     int(parsePressCnt),
				JumpDistance: float32(parseJumpDistance),
			}

			err = backend.AddField(a, pool, ctx)
			if err != nil {
				log.Println(err)
				fmt.Println("Спробуйте ще раз")
				continue
			}
		case 2:
			athletes, err := backend.SelectTable(pool, ctx)
			if err != nil {
				log.Println(err)
				fmt.Println("Спробуйте ще раз")
				continue
			}

			printTable(athletes)

		case 3:
			inputIDs := takeInput("Введіть список id для видалення:")

			fields := strings.Fields(inputIDs)
			if len(fields) == 0 {
				fmt.Println("Не введено жодного ID.")
				fmt.Println("Спробуйте ще раз")
				continue
			}

			ids := make([]int, len(fields))
			for i, v := range fields {
				parseIDs, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					fmt.Println("Спробуйте ще раз")
					continue MainLoop
				}

				ids[i] = int(parseIDs)
			}

			err = backend.DeleteFields(ids, pool, ctx)
			if err != nil {
				log.Println(err)
				fmt.Println("Спробуйте ще раз")
				continue
			}
		case 4:
			inputID := takeInput("Введіть id поля для оновлення:")

			parseId, err := strconv.ParseInt(inputID, 10, 64)
			if err != nil {
				fmt.Println("Спробуйте ще раз")
				continue
			}

			inputRecord := takeInput("Введіть оновленні значення запису:")

			fields := strings.Fields(inputRecord)
			if len(fields) == 0 {
				fmt.Println("Не введено жодного поля.")
				fmt.Println("Спробуйте ще раз")
				continue
			}

			parseRun100m, err := strconv.ParseFloat(fields[2], 64)
			if err != nil {
				fmt.Println("Спробуйте ще раз")
				continue
			}

			parseRun3km, err := strconv.ParseFloat(fields[3], 64)
			if err != nil {
				fmt.Println("Спробуйте ще раз")
				continue
			}

			parsePressCnt, err := strconv.ParseInt(fields[4], 10, 64)
			if err != nil {
				fmt.Println("Спробуйте ще раз")
				continue
			}

			parseJumpDistance, err := strconv.ParseFloat(fields[5], 64)
			if err != nil {
				fmt.Println("Спробуйте ще раз")
				continue
			}

			a := backend.Athlete{
				Id:           int(parseId),
				Name:         fields[0],
				Surname:      fields[1],
				Run100m:      float32(parseRun100m),
				Run3km:       float32(parseRun3km),
				PressCnt:     int(parsePressCnt),
				JumpDistance: float32(parseJumpDistance),
			}

			err = backend.UpdateField(a, pool, ctx)
			if err != nil {
				log.Println(err)
				fmt.Println("Спробуйте ще раз")
				continue
			}

		case 5:
			athletes, err := backend.SortTable(pool, ctx)
			if err != nil {
				log.Println(err)
				fmt.Println("Спробуйте ще раз")
				continue
			}

			printTable(athletes)

		default:
			fmt.Println("Спробуйте ще раз")
			continue
		}
	}
}

// -----------------------------------------------------------------

func startupLogger() (*os.File, error) {
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("impossible to open log file: %w", err)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	return logFile, nil
}
