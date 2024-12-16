package main

import (
	"context"
	"database/sql"
	"fmt"
	"iter"
	"log/slog"
	"math/rand/v2"
	"os"

	"github.com/BooleanCat/go-functional/v2/it"
	"github.com/Songmu/strrand"
	asyncdataloader "github.com/mackee/iterutils/_example/asyncdataloader"
)

func main() {
	if err := run(); err != nil {
		slog.Error("occurred error", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()
	db, err := asyncdataloader.NewDB()
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	if err := doUsers(ctx, db); err != nil {
		return fmt.Errorf("failed to do users: %w", err)
	}
	if err := doItems(ctx, db); err != nil {
		return fmt.Errorf("failed to do items: %w", err)
	}
	if err := doUserItems(ctx, db); err != nil {
		return fmt.Errorf("failed to do user_items: %w", err)
	}
	return nil
}

func doUsers(ctx context.Context, db *sql.DB) error {
	for i := range 100 {
		rnd := it.Take(randNameGen(), 100)
		bi := asyncdataloader.NewUsersSQL().BulkInsert()
		id := asyncdataloader.UserID(0) + asyncdataloader.UserID(i)*100
		for name := range rnd {
			id++
			bi.Append(
				asyncdataloader.NewUsersSQL().Insert().ValueID(id).ValueName(name),
			)
		}
		if _, err := bi.ExecContext(ctx, db); err != nil {
			return fmt.Errorf("failed to bulk insert users: %w", err)
		}
	}
	return nil
}

func doItems(ctx context.Context, db *sql.DB) error {
	for i := range 100 {
		rnd := it.Take(randNameGen(), 100)
		bi := asyncdataloader.NewItemsSQL().BulkInsert()
		id := asyncdataloader.ItemID(0) + asyncdataloader.ItemID(i)*100
		for name := range rnd {
			id++
			bi.Append(
				asyncdataloader.NewItemsSQL().Insert().ValueID(id).ValueName(name),
			)
		}
		if _, err := bi.ExecContext(ctx, db); err != nil {
			return fmt.Errorf("failed to bulk insert items: %w", err)
		}
	}
	return nil
}

func doUserItems(ctx context.Context, db *sql.DB) error {
	for i := range 100 {
		bi := asyncdataloader.NewUserItemsSQL().BulkInsert()
		userID := asyncdataloader.UserID(0) + asyncdataloader.UserID(i)*100
		for range 100 {
			userID++
			for range 10 {
				itemID := asyncdataloader.ItemID(rand.Int64N(9999)) + 1
				bi.Append(
					asyncdataloader.NewUserItemsSQL().Insert().
						ValueUserID(userID).
						ValueItemID(itemID),
				)
			}
		}
		if _, err := bi.ExecContext(ctx, db); err != nil {
			return fmt.Errorf("failed to bulk insert user_items: %w", err)
		}
	}
	return nil
}

func randNameGen() iter.Seq[string] {
	return func(fn func(string) bool) {
		for {
			s, err := strrand.RandomString("[あ-ん]{10}")
			if err != nil {
				break
			}
			if !fn(s) {
				break
			}
		}
	}
}
