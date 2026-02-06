package main

import (
	"context"
	"courseWork/shared"
	"slices"
	"testing"
)

func TestSortTable(t *testing.T) {
	const dbURL = "postgres://user:password@localhost:5432/test_course_work_db"
	ctx := context.Background()

	pool, err := StartupTable(ctx, dbURL)
	if err != nil {
		t.Errorf("error to start up table: %s", err.Error())
	}

	gotAthletes, err := SortByRun100m(pool, ctx)
	if err != nil {
		t.Errorf("error to sort table: %s", err.Error())
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
