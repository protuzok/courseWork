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
		case 3:
			deleteFields(pool, ctx)

		default:
			fmt.Println("Спробуйте ще раз")
			continue
		}
	}
}

func deleteFields(pool *pgxpool.Pool, ctx context.Context) {
	fmt.Println("Введіть список id для видалення:")
	fmt.Printf("-> ")
	scanner := bufio.NewScanner(os.Stdin)

	if scanner.Scan() {
		input := scanner.Text()

		ids := strings.Fields(input)

		if len(ids) == 0 {
			fmt.Println("Ви не ввели жодного ID")
			return
		}

		var placeholders strings.Builder
		placeholders.WriteString("(")

		for i := 1; i <= len(ids); i++ {
			var s string
			if i != len(ids) {
				s = fmt.Sprintf("$%d, ", i)
			} else {
				s = fmt.Sprintf("$%d)", i)
			}

			placeholders.WriteString(s)
		}

		var sqlIds = make([]any, len(ids))
		for i, v := range ids {
			sqlIds[i] = v
		}

		deleteRecordsSQL := "DELETE FROM users WHERE id IN " + placeholders.String()

		pool.Exec(ctx, deleteRecordsSQL, sqlIds...)
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
	fmt.Println("1) Додати поле в таблицю")
	fmt.Println("2) Вивести всю таблицю")
	fmt.Println("3) Видалити поле(-я) з таблиці за id")
}

func takeOptions(option *int) {
	fmt.Printf("-> ")
	fmt.Scan(option)
	fmt.Println("")
}

// todo написати функції для вводу (стрічка для виведення і сам ввід). Таке чат робив для ЧМ пам'ятаю
func addField(pool *pgxpool.Pool, ctx context.Context) {
	fmt.Println("Введіть значення полів запису:")
	fmt.Printf("-> ")
	scanner := bufio.NewScanner(os.Stdin)

	// todo десь тут треба додати перевірку, чи користувач взагалі шось ввів (deleteFields)
	var fields []string
	if scanner.Scan() {
		input := scanner.Text()

		fields = strings.Fields(input)
	}

	sqlArgs := make([]any, len(fields))
	for i, v := range fields {
		sqlArgs[i] = v
	}

	insertRecordSQL := "INSERT INTO users(name, surname, year) VALUES ($1, $2, $3)"

	// todo всі ці функції будуть повинні повертати err, а не просто void
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
