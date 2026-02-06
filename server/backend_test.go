package main

import (
	"context"
	"courseWork/server/storage"
	"courseWork/shared"
	"slices"
	"testing"
)

func TestSortTable(t *testing.T) {
	const dbURL = "postgres://user:password@localhost:5432/test_course_work_db"
	ctx := context.Background()
	repo, err := storage.NewAthleteRepository(ctx, dbURL)
	if err != nil {
		t.Fatalf("error to start up table: %s", err.Error())
	}

	gotAthletes, err := repo.GetAllSortedByRun100m(ctx)
	if err != nil {
		t.Fatalf("error to sort table: %s", err.Error())
	}

	var wantAthletes = []shared.Athlete{
		{3, "класно", "цу", 1, 2, 3, 4},
		{1, "привіт", "руслан", 3.4, 2, 5, 8},
		{4, "крутяк", "кербе", 10, 20, 30, 40},
		{2, "як справи", "норм", 12, 2, 12, 3},
	}

	if !slices.Equal(gotAthletes, wantAthletes) {
		t.Errorf("arrays is not equal: %v != %v", gotAthletes, wantAthletes)
	}
}
