package iterutils_test

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/mackee/iterutils"
)

func ExampleFromTryNexter2() {
	csvText := `id,name,created_at
1,foo,2021-01-01 09:00:00
2,bar,2021-01-02 12:00:00
3,baz,2021-01-03 15:00:00
`
	csvr := csv.NewReader(strings.NewReader(csvText))
	// Skip the header row
	if _, err := csvr.Read(); err != nil {
		panic(err)
	}
	tn := iterutils.NewTryNexterWithT(csvr, func(csvr *csv.Reader) (record []string, err error) {
		return csvr.Read()
	})
	type user struct {
		ID        int
		Name      string
		CreatedAt time.Time
	}
	it := iterutils.FromTryNexter2(tn, func(t iterutils.TryNexterWithT[*csv.Reader, []string], record []string, err error) (*user, error) {
		if err != nil {
			return nil, err
		}
		if len(record) != 3 {
			line, _ := t.T().FieldPos(0)
			return nil, fmt.Errorf("expected 3 fields, got %d, lines=%d", len(record), line)
		}
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}
		createdAt, err := time.Parse(time.DateTime, record[2])
		if err != nil {
			return nil, err
		}
		return &user{
			ID:        id,
			Name:      record[1],
			CreatedAt: createdAt,
		}, nil
	})
	for u, err := range it {
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Println(u.ID, u.Name, u.CreatedAt)
	}
	// Output:
	// 1 foo 2021-01-01 09:00:00 +0000 UTC
	// 2 bar 2021-01-02 12:00:00 +0000 UTC
	// 3 baz 2021-01-03 15:00:00 +0000 UTC
}
