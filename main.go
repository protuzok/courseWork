package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

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

	ctx := context.Background()

	pool, err := startupTable(ctx, dbURL)
	if err != nil {
		fmt.Println("Помилка старту таблиці бази даних: ", err)
		log.Println(err)
		os.Exit(1)
	}
	defer pool.Close()

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

			a := Athlete{
				-1,
				fields[0],
				fields[1],
				float32(parseRun100m),
				float32(parseRun3km),
				int(parsePressCnt),
				float32(parseJumpDistance),
			}

			err = addField(a, pool, ctx)
			if err != nil {
				log.Println(err)
				fmt.Println("Спробуйте ще раз")
				continue
			}
		case 2:
			err = printTable(pool, ctx)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}

		case 3:
			err = deleteFields(pool, ctx)
			if err != nil {
				log.Println(err)
				os.Exit(1)
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

			a := Athlete{
				int(parseId),
				fields[0],
				fields[1],
				float32(parseRun100m),
				float32(parseRun3km),
				int(parsePressCnt),
				float32(parseJumpDistance),
			}

			err = updateField(a, pool, ctx)
			if err != nil {
				log.Println(err)
				fmt.Println("Спробуйте ще раз")
				continue
			}

		default:
			fmt.Println("Спробуйте ще раз")
			continue
		}
	}
}

// -----------------------------------------------------------------

func printOptions() {
	fmt.Println("")
	fmt.Println("============================")
	fmt.Println("0) Вийти з програми")
	fmt.Println("1) Додати поле в таблицю")
	fmt.Println("2) Вивести всю таблицю")
	fmt.Println("3) Видалити поле(-я) з таблиці за id")
	fmt.Println("4) Змінити поле таблиці за id")
}

func takeOption() (int, error) {
	input := takeInput("Оберіть опцію:")

	option, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return -1, err
	}

	return int(option), nil
}

func takeInput(instruction string) (input string) {
	for {
		fmt.Println("")
		fmt.Println(instruction)
		fmt.Printf("-> ")
		scanner := bufio.NewScanner(os.Stdin)

		if scanner.Scan() {
			input = scanner.Text()

			if input == "" {
				fmt.Println("Пустий рядок. Спробуйте ще раз")
				fmt.Println("")

				continue
			}
		}

		break
	}

	return input
}
