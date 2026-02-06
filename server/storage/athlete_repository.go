package storage

import (
	"context"
	"courseWork/shared"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AthleteRepository struct {
	db *pgxpool.Pool
}

func NewAthleteRepository(ctx context.Context, dbURL string) (*AthleteRepository, error) {
	p, err := pgxpool.New(ctx, dbURL)
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
		p.Close()
		return nil, fmt.Errorf("error creating table: %w", err)
	}

	return &AthleteRepository{db: p}, nil
}

func (r *AthleteRepository) Close() {
	r.db.Close()
}

func (r *AthleteRepository) Create(ctx context.Context, a shared.Athlete) error {
	insertRecordSQL := `INSERT INTO athletes
    	(name, surname, run_100m, run_3km, press_сnt, jump_distance) 
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(ctx, insertRecordSQL, a.Name, a.Surname, a.Run100m, a.Run3km, a.PressCnt, a.JumpDistance)
	if err != nil {
		return fmt.Errorf("impossible to add field: %w", err)
	}

	return nil
}

func (r *AthleteRepository) GetAll(ctx context.Context) ([]shared.Athlete, error) {
	selectRecordsSQL := `SELECT id, name, surname, run_100m, run_3km, press_сnt, jump_distance 
		FROM athletes`

	rows, err := r.db.Query(ctx, selectRecordsSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var athletes []shared.Athlete

	for rows.Next() {
		var a shared.Athlete
		err = rows.Scan(&a.Id, &a.Name, &a.Surname, &a.Run100m, &a.Run3km, &a.PressCnt, &a.JumpDistance)
		if err != nil {
			return nil, err
		}
		athletes = append(athletes, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return athletes, nil
}

func (r *AthleteRepository) Delete(ctx context.Context, ids []int) error {
	deleteRecordsSQL := "DELETE FROM athletes WHERE id = ANY($1::int[])"

	_, err := r.db.Exec(ctx, deleteRecordsSQL, ids)
	if err != nil {
		return fmt.Errorf("impossible to delete fields: %w", err)
	}

	return nil
}

func (r *AthleteRepository) Update(ctx context.Context, a shared.Athlete) error {
	updateRecordSQL := `UPDATE athletes 
		SET (name, surname, run_100m, run_3km, press_сnt, jump_distance) = ($1, $2, $3, $4, $5, $6) 
		WHERE id = $7`

	_, err := r.db.Exec(ctx, updateRecordSQL, a.Name, a.Surname, a.Run100m, a.Run3km, a.PressCnt, a.JumpDistance, a.Id)
	if err != nil {
		return fmt.Errorf("impossible to update field: %w", err)
	}

	return nil
}

func (r *AthleteRepository) GetAllSortedByRun100m(ctx context.Context) ([]shared.Athlete, error) {
	athletes, err := r.GetAll(ctx)
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

func (r *AthleteRepository) GetBestPressMinJump(ctx context.Context) ([]shared.Athlete, error) {
	selectRecordsSQL := `SELECT id, name, surname, run_100m, run_3km, press_сnt, jump_distance FROM athletes
		WHERE press_сnt = (SELECT MAX(press_сnt) FROM athletes)
		AND jump_distance = (SELECT MIN(jump_distance) FROM athletes)
		ORDER BY name`

	rows, err := r.db.Query(ctx, selectRecordsSQL)
	if err != nil {
		return nil, fmt.Errorf("impossible to group and sort table: %w", err)
	}
	defer rows.Close()

	var athletes []shared.Athlete
	for rows.Next() {
		var a shared.Athlete
		err = rows.Scan(&a.Id, &a.Name, &a.Surname, &a.Run100m, &a.Run3km, &a.PressCnt, &a.JumpDistance)
		if err != nil {
			return nil, fmt.Errorf("impossible to group and sort table: %w", err)
		}
		athletes = append(athletes, a)
	}

	return athletes, nil
}

func (r *AthleteRepository) GetWithRun3kmDeviation(ctx context.Context) ([]shared.Athlete, error) {
	selectRecordsSQL := `SELECT id, name, surname, run_100m, run_3km, press_сnt, jump_distance FROM athletes
		WHERE run_3km BETWEEN ((SELECT AVG(run_3km) FROM athletes) * (1 - 0.07359))
		AND ((SELECT AVG(run_3km) FROM athletes) * (1 + 0.07359))`

	rows, err := r.db.Query(ctx, selectRecordsSQL)
	if err != nil {
		return nil, fmt.Errorf("impossible to select by deviation: %w", err)
	}
	defer rows.Close()

	var athletes []shared.Athlete
	for rows.Next() {
		var a shared.Athlete
		err = rows.Scan(&a.Id, &a.Name, &a.Surname, &a.Run100m, &a.Run3km, &a.PressCnt, &a.JumpDistance)
		if err != nil {
			return nil, fmt.Errorf("impossible to select by deviation: %w", err)
		}
		athletes = append(athletes, a)
	}

	return athletes, nil
}

func (r *AthleteRepository) GetMinPressRun100mStats(ctx context.Context) ([]shared.Task4Row, error) {
	const query = `
	SELECT
		name,
		press_сnt,
		run_100m,
		run_100m / (SELECT AVG(run_100m) FROM athletes) as deviation
	FROM athletes
	WHERE press_сnt = (SELECT MIN(press_сnt) FROM athletes)`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("impossible to do task 4: %w", err)
	}
	defer rows.Close()

	var athletes []shared.Task4Row
	for rows.Next() {
		var a shared.Task4Row
		err = rows.Scan(&a.Name, &a.PressCnt, &a.Run100m, &a.Deviation)
		if err != nil {
			return nil, fmt.Errorf("impossible to do task 4: %w", err)
		}
		athletes = append(athletes, a)
	}

	return athletes, nil
}

func (r *AthleteRepository) GetBestOverallAthlete(ctx context.Context) ([]shared.Athlete, error) {
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

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("impossible to select best total result: %w", err)
	}
	defer rows.Close()

	var athletes []shared.Athlete
	for rows.Next() {
		var a shared.Athlete
		err = rows.Scan(&a.Id, &a.Name, &a.Surname, &a.Run100m, &a.Run3km, &a.PressCnt, &a.JumpDistance)
		if err != nil {
			return nil, fmt.Errorf("impossible to select best total result: %w", err)
		}
		athletes = append(athletes, a)
	}

	return athletes, nil
}
