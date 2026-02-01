package backend

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func StartupTable(ctx context.Context, dbURL string) (p *pgxpool.Pool, err error) {
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

func AddField(a Athlete, pool *pgxpool.Pool, ctx context.Context) error {
	insertRecordSQL := `INSERT INTO athletes
    	(name, surname, run_100m, run_3km, press_сnt, jump_distance) 
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := pool.Exec(ctx, insertRecordSQL, a.Name, a.Surname, a.Run100m, a.Run3km, a.PressCnt, a.JumpDistance)
	if err != nil {
		return fmt.Errorf("impossible to add field: %w", err)
	}

	return nil
}

// todo зробити окремий метод для виводу таблиці та для отримання таблиці у вигляді масиву атлетів
func PrintTable(pool *pgxpool.Pool, ctx context.Context) error {
	fmt.Println("Таблиця атлетів:")
	fmt.Printf("%-5s %-20s %-20s %-10s %-10s %-10s %-10s \n", "id", "name", "surname", "run100m", "run3km", "pressCnt", "jumpDist")

	selectRecordsSQL := "SELECT id, name, surname, run_100m, run_3km, press_сnt, jump_distance FROM athletes"

	rows, err := pool.Query(ctx, selectRecordsSQL)
	if err != nil {
		return fmt.Errorf("impossible to print table: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var a Athlete
		err := rows.Scan(&a.Id, &a.Name, &a.Surname, &a.Run100m, &a.Run3km, &a.PressCnt, &a.JumpDistance)
		if err != nil {
			return fmt.Errorf("impossible to print table: %w", err)
		}

		fmt.Printf("%-5d %-20s %-20s %-10.2f %-10.2f %-10d %-10.2f \n", a.Id, a.Name, a.Surname, a.Run100m, a.Run3km, a.PressCnt, a.JumpDistance)
	}

	return nil
}

func SortTable(pool *pgxpool.Pool, ctx context.Context) ([]Athlete, error) {
	// todo виділити це в окремий метод для отримання даних
	selectRecordsSQL := `SELECT id, name, surname, run_100m, run_3km, press_сnt, jump_distance 
		FROM athletes`

	rows, err := pool.Query(ctx, selectRecordsSQL)
	if err != nil {
		return nil, fmt.Errorf("impossible to sort and print table: %w", err)
	}

	var athletes []Athlete

	for i := 0; rows.Next(); i++ {
		athletes = append(athletes, Athlete{})
		err = rows.Scan(&athletes[i].Id, &athletes[i].Name, &athletes[i].Surname, &athletes[i].Run100m, &athletes[i].Run3km, &athletes[i].PressCnt, &athletes[i].JumpDistance)
		if err != nil {
			return nil, fmt.Errorf("impossible to sort and print table: %w", err)
		}
	}

	// Selection sort
	n := len(athletes)
	for i := 0; i < n-1; i++ {
		minIndex := i
		for j := i + 1; j < n; j++ {
			if athletes[j].Run100m < athletes[minIndex].Run100m {
				minIndex = j
			}
		}
		athletes[i], athletes[minIndex] = athletes[minIndex], athletes[i]
	}

	return athletes, nil
}

func DeleteFields(ids []int, pool *pgxpool.Pool, ctx context.Context) error {

	deleteRecordsSQL := "DELETE FROM athletes WHERE id = ANY($1::int[])"

	_, err := pool.Exec(ctx, deleteRecordsSQL, ids)
	if err != nil {
		return fmt.Errorf("impossible to delete fields: %w", err)
	}

	return nil
}

func UpdateField(a Athlete, pool *pgxpool.Pool, ctx context.Context) error {
	updateRecordSQL := `UPDATE athletes 
		SET (name, surname, run_100m, run_3km, press_сnt, jump_distance) = ($1, $2, $3, $4, $5, $6) 
		WHERE id = $7`

	_, err := pool.Exec(ctx, updateRecordSQL, a.Name, a.Surname, a.Run100m, a.Run3km, a.PressCnt, a.JumpDistance, a.Id)
	if err != nil {
		return fmt.Errorf("impossible to update field: %w", err)
	}

	return nil
}

// -----------------------------------------------------------------
