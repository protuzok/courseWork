package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"courseWork/shared"
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

func AddField(a shared.Athlete, pool *pgxpool.Pool, ctx context.Context) error {
	insertRecordSQL := `INSERT INTO athletes
    	(name, surname, run_100m, run_3km, press_сnt, jump_distance) 
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := pool.Exec(ctx, insertRecordSQL, a.Name, a.Surname, a.Run100m, a.Run3km, a.PressCnt, a.JumpDistance)
	if err != nil {
		return fmt.Errorf("impossible to add field: %w", err)
	}

	return nil
}

func SelectTable(pool *pgxpool.Pool, ctx context.Context) ([]shared.Athlete, error) {
	selectRecordsSQL := `SELECT id, name, surname, run_100m, run_3km, press_сnt, jump_distance 
		FROM athletes`

	rows, err := pool.Query(ctx, selectRecordsSQL)
	if err != nil {
		return nil, err
	}

	var athletes []shared.Athlete

	for i := 0; rows.Next(); i++ {
		athletes = append(athletes, shared.Athlete{})
		err = rows.Scan(&athletes[i].Id, &athletes[i].Name, &athletes[i].Surname, &athletes[i].Run100m, &athletes[i].Run3km, &athletes[i].PressCnt, &athletes[i].JumpDistance)
		if err != nil {
			return nil, err
		}
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

func UpdateField(a shared.Athlete, pool *pgxpool.Pool, ctx context.Context) error {
	updateRecordSQL := `UPDATE athletes 
		SET (name, surname, run_100m, run_3km, press_сnt, jump_distance) = ($1, $2, $3, $4, $5, $6) 
		WHERE id = $7`

	_, err := pool.Exec(ctx, updateRecordSQL, a.Name, a.Surname, a.Run100m, a.Run3km, a.PressCnt, a.JumpDistance, a.Id)
	if err != nil {
		return fmt.Errorf("impossible to update field: %w", err)
	}

	return nil
}

func SortByRun100m(pool *pgxpool.Pool, ctx context.Context) ([]shared.Athlete, error) {
	athletes, err := SelectTable(pool, ctx)
	if err != nil {
		return nil, fmt.Errorf("impossible to sort table: %w", err)
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

func GroupByPressAndJumpAndSortByName(pool *pgxpool.Pool, ctx context.Context) ([]shared.Athlete, error) {
	selectRecordsSQL := `SELECT * FROM athletes
		WHERE press_сnt = (SELECT MAX(press_сnt) FROM athletes)
		AND jump_distance = (SELECT MIN(jump_distance) FROM athletes)
		ORDER BY name`

	rows, err := pool.Query(ctx, selectRecordsSQL)
	if err != nil {
		return nil, fmt.Errorf("impossible to group and sort table: %w", err)
	}

	var athletes []shared.Athlete

	for i := 0; rows.Next(); i++ {
		athletes = append(athletes, shared.Athlete{})
		err = rows.Scan(&athletes[i].Id, &athletes[i].Name, &athletes[i].Surname, &athletes[i].Run100m, &athletes[i].Run3km, &athletes[i].PressCnt, &athletes[i].JumpDistance)
		if err != nil {
			return nil, fmt.Errorf("impossible to group and sort table: %w", err)
		}
	}

	return athletes, nil
}

func SelectByDeviationRun3km(pool *pgxpool.Pool, ctx context.Context) ([]shared.Athlete, error) {
	selectRecordsSQL := `SELECT * FROM athletes
		WHERE run_3km BETWEEN ((SELECT AVG(run_3km) FROM athletes) * (1 - 0.07359))
		AND ((SELECT AVG(run_3km) FROM athletes) * (1 + 0.07359))`

	rows, err := pool.Query(ctx, selectRecordsSQL)
	if err != nil {
		return nil, fmt.Errorf("impossible to select by deviation: %w", err)
	}

	var athletes []shared.Athlete

	for i := 0; rows.Next(); i++ {
		athletes = append(athletes, shared.Athlete{})
		err = rows.Scan(&athletes[i].Id, &athletes[i].Name, &athletes[i].Surname, &athletes[i].Run100m, &athletes[i].Run3km, &athletes[i].PressCnt, &athletes[i].JumpDistance)
		if err != nil {
			return nil, fmt.Errorf("impossible to select by deviation: %w", err)
		}
	}

	return athletes, nil
}

func SelectByMinPressAndGetDeviationRun100m(pool *pgxpool.Pool, ctx context.Context) ([]shared.Task4Row, error) {
	const query = `
	SELECT
		name,
		press_сnt,
		run_100m,
		run_100m / (SELECT AVG(run_100m) FROM athletes) as deviation
	FROM athletes
	WHERE press_сnt = (SELECT MIN(press_сnt) FROM athletes)`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("impossible to do task 4: %w", err)
	}

	var athletes []shared.Task4Row

	for i := 0; rows.Next(); i++ {
		athletes = append(athletes, shared.Task4Row{})
		err = rows.Scan(&athletes[i].Name, &athletes[i].PressCnt, &athletes[i].Run100m, &athletes[i].Deviation)
		if err != nil {
			return nil, fmt.Errorf("impossible to do task 4: %w", err)
		}
	}

	return athletes, nil
}

func SelectBestTotalResult(pool *pgxpool.Pool, ctx context.Context) ([]shared.Athlete, error) {
	const query = `
	WITH
		Ranks AS (
		SELECT
			*,
			RANK() OVER (ORDER BY run_100m) +
			RANK() OVER (ORDER BY run_3km) +
			RANK() OVER (ORDER BY press_сnt DESC) +
			RANK() OVER (ORDER BY jump_distance DESC) as total_rank
		FROM athletes
	)
	SELECT id, name, surname, run_100m, run_3km, press_сnt, jump_distance
	FROM Ranks
	WHERE total_rank = (SELECT MIN(total_rank) FROM Ranks)`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("impossible to select best total result: %w", err)
	}

	var athletes []shared.Athlete

	for i := 0; rows.Next(); i++ {
		athletes = append(athletes, shared.Athlete{})
		err = rows.Scan(&athletes[i].Id, &athletes[i].Name, &athletes[i].Surname, &athletes[i].Run100m, &athletes[i].Run3km, &athletes[i].PressCnt, &athletes[i].JumpDistance)
		if err != nil {
			return nil, fmt.Errorf("impossible to select best total result: %w", err)
		}
	}

	return athletes, nil
}
