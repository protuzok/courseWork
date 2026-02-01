package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

const dbURL = "postgres://user:password@localhost:5432/course_work_db"

// -----------------------------------------------------------------

func startupTable(ctx context.Context, dbURL string) (p *pgxpool.Pool, err error) {
	p, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("error start up data base table: %w", err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS athletes(
		id SERIAL PRIMARY KEY, 
		name TEXT NOT NULL, 
		surname TEXT NOT NULL, 
		run_100m REAL NOT NULL,
		run_3km REAL NOT NULL,
		press_сnt INT NOT NULL,
		jump_distance REAL NOT NULL
	);`

	_, err = p.Exec(ctx, createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("error start up data base table: %w", err)
	}
	return p, nil
}

func addField(a Athlete, pool *pgxpool.Pool, ctx context.Context) error {
	insertRecordSQL := `INSERT INTO athletes
    	(name, surname, run_100m, run_3km, press_сnt, jump_distance) 
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := pool.Exec(ctx, insertRecordSQL, a.name, a.surname, a.run100m, a.run3km, a.pressCnt, a.jumpDistance)
	if err != nil {
		return fmt.Errorf("impossible to add field: %w", err)
	}

	return nil
}

func printTable(pool *pgxpool.Pool, ctx context.Context) error {
	fmt.Println("Таблиця атлетів:")
	fmt.Printf("%-5s %-20s %-20s %-10s %-10s %-10s %-10s \n", "id", "name", "surname", "run100m", "run3km", "pressCnt", "jumpDist")

	selectRecordsSQL := "SELECT id, name, surname, run_100m, run_3km, press_сnt, jump_distance FROM athletes"

	rows, _ := pool.Query(ctx, selectRecordsSQL)
	defer rows.Close()

	for rows.Next() {
		var a Athlete
		err := rows.Scan(&a.id, &a.name, &a.surname, &a.run100m, &a.run3km, &a.pressCnt, &a.jumpDistance)
		if err != nil {
			return fmt.Errorf("impossible to print table: %w", err)
		}

		fmt.Printf("%-5d %-20s %-20s %-10.2f %-10.2f %-10d %-10.2f \n", a.id, a.name, a.surname, a.run100m, a.run3km, a.pressCnt, a.jumpDistance)
	}

	return nil
}

func deleteFields(ids []int, pool *pgxpool.Pool, ctx context.Context) error {

	deleteRecordsSQL := "DELETE FROM athletes WHERE id = ANY($1::int[])"

	_, err := pool.Exec(ctx, deleteRecordsSQL, ids)
	if err != nil {
		return fmt.Errorf("impossible to delete fields: %w", err)
	}

	return nil
}

func updateField(a Athlete, pool *pgxpool.Pool, ctx context.Context) error {
	updateRecordSQL := `UPDATE athletes 
		SET (name, surname, run_100m, run_3km, press_сnt, jump_distance) = ($1, $2, $3, $4, $5, $6) 
		WHERE id = $7`

	_, err := pool.Exec(ctx, updateRecordSQL, a.name, a.surname, a.run100m, a.run3km, a.pressCnt, a.jumpDistance, a.id)
	if err != nil {
		return fmt.Errorf("impossible to update field: %w", err)
	}

	return nil
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
