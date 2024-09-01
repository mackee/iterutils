package iterutils_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mackee/iterutils"
)

func mockDB() *sql.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	rowData := [][]driver.Value{
		{1, "foo"},
		{2, "bar"},
		{3, "baz"},
	}
	rows := mock.NewRows([]string{"id", "name"}).AddRows(rowData...)
	mock.ExpectQuery("SELECT id, name FROM users").WillReturnRows(rows)
	return db
}

func ExampleFromNexter2() {
	db := mockDB()

	ctx := context.Background()
	rows, err := db.QueryContext(ctx, "SELECT id, name FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	type user struct {
		ID   int
		Name string
	}
	users := iterutils.FromNexter2(rows, func(rows *sql.Rows) (*user, error) {
		var u user
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		return &u, nil
	})
	for users, err := range users {
		if err != nil {
			panic(err)
		}
		fmt.Println(users.ID, users.Name)
	}
	// Output:
	// 1 foo
	// 2 bar
	// 3 baz
}
