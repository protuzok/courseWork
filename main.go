package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

const dbURL = "postgres://user:password@localhost:5432/mydb"

func main() {
	ctx := context.Background()

	pool, done := startTable(ctx)
	if !done {
		return
	}
	defer pool.Close()

	for true {
		printOptions()

		var option int
		takeOptions(&option)

		switch option {
		case 0:
			return
		case 1:
			addField(pool, ctx)
		case 2:
			readTable(pool, ctx)
		default:
			fmt.Println("Спробуйте ще раз")
			continue
		}
	}
}

func startTable(ctx context.Context) (p *pgxpool.Pool, done bool) {
	p, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, false
	}

	createTableSQL := "CREATE TABLE IF NOT EXISTS users(id SERIAL PRIMARY KEY, name TEXT NOT NULL, surname TEXT NOT NULL, year INT);"
	_, err = p.Exec(ctx, createTableSQL)
	if err != nil {
		return nil, false
	}
	return p, true
}

func printOptions() {
	fmt.Println("")
	fmt.Println("============================")
	fmt.Println("Оберіть опцію:")
	fmt.Println("0) Вийти з програми")
	fmt.Println("1) Додати поле")
	fmt.Println("2) Прочитати таблицю")
}

func takeOptions(option *int) {
	fmt.Printf("-> ")
	fmt.Scan(option)
	fmt.Println("")
}

func addField(pool *pgxpool.Pool, ctx context.Context) {
	fmt.Println("Введіть значення полів запису:")
	fmt.Printf("-> ")
	scanner := bufio.NewScanner(os.Stdin)

	var fields []string
	if scanner.Scan() {
		input := scanner.Text()

		fields = strings.Fields(input)
	}

	sqlArgs := make([]any, len(fields), len(fields))
	for i, v := range fields {
		sqlArgs[i] = v
	}

	insertRecordSQL := "INSERT INTO users(name, surname, year) VALUES ($1, $2, $3)"

	pool.Exec(ctx, insertRecordSQL, sqlArgs...)
}

func readTable(pool *pgxpool.Pool, ctx context.Context) {
	fmt.Println("Таблиця користувачів:")
	fmt.Printf("%-5s %-20s %-20s %-5s \n", "id", "name", "surname", "year")

	selectRecordsSQL := "SELECT * FROM users"

	rows, _ := pool.Query(ctx, selectRecordsSQL)
	defer rows.Close()

	for rows.Next() {
		var u User
		rows.Scan(&u.id, &u.name, &u.surname, &u.year)

		fmt.Printf("%-5d %-20s %-20s %-5d \n", u.id, u.name, u.surname, u.year)
	}
}
