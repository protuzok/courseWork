package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// todo треба буде змінити структуру таблиці та методи для роботи з нею. І назву бази даних
const dbURL = "postgres://user:password@localhost:5432/mydb"

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
			err = addField(pool, ctx)
			if err != nil {
				log.Println(err)
				os.Exit(1)
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

		default:
			fmt.Println("Спробуйте ще раз")
			continue
		}
	}
}

func startupLogger() (*os.File, error) {
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("impossible to open log file: %w", err)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	return logFile, nil
}

func printOptions() {
	fmt.Println("")
	fmt.Println("============================")
	fmt.Println("0) Вийти з програми")
	fmt.Println("1) Додати поле в таблицю")
	fmt.Println("2) Вивести всю таблицю")
	fmt.Println("3) Видалити поле(-я) з таблиці за id")
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

// -----------------------------------------------------------------

func startupTable(ctx context.Context, dbURL string) (p *pgxpool.Pool, err error) {
	p, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("error start up data base table: %w", err)
	}

	createTableSQL := "CREATE TABLE IF NOT EXISTS users(id SERIAL PRIMARY KEY, name TEXT NOT NULL, surname TEXT NOT NULL, year INT);"
	_, err = p.Exec(ctx, createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("error start up data base table: %w", err)
	}
	return p, nil
}

func addField(pool *pgxpool.Pool, ctx context.Context) error {
	input := takeInput("Введіть значення полів запису:")

	var fields = strings.Fields(input)

	sqlArgs := make([]any, len(fields))
	for i, v := range fields {
		sqlArgs[i] = v
	}

	insertRecordSQL := "INSERT INTO users(name, surname, year) VALUES ($1, $2, $3)"

	_, err := pool.Exec(ctx, insertRecordSQL, sqlArgs...)
	if err != nil {
		return fmt.Errorf("impossible to add field: %w", err)
	}

	return nil
}

func printTable(pool *pgxpool.Pool, ctx context.Context) error {
	fmt.Println("Таблиця користувачів:")
	fmt.Printf("%-5s %-20s %-20s %-5s \n", "id", "name", "surname", "year")

	selectRecordsSQL := "SELECT * FROM users"

	rows, _ := pool.Query(ctx, selectRecordsSQL)
	defer rows.Close()

	for rows.Next() {
		var u User
		err := rows.Scan(&u.id, &u.name, &u.surname, &u.year)
		if err != nil {
			return fmt.Errorf("impossible to print table: %w", err)
		}

		fmt.Printf("%-5d %-20s %-20s %-5d \n", u.id, u.name, u.surname, u.year)
	}

	return nil
}

func deleteFields(pool *pgxpool.Pool, ctx context.Context) error {
	input := takeInput("Введіть список id для видалення:")

	fields := strings.Fields(input)
	if len(fields) == 0 {
		fmt.Println("Не введено жодного ID.")
		return nil
	}

	deleteRecordsSQL := "DELETE FROM users WHERE id = ANY($1::int[])"

	_, err := pool.Exec(ctx, deleteRecordsSQL, fields)
	if err != nil {
		return fmt.Errorf("impossible to delete fields: %w", err)
	}

	return nil
}
