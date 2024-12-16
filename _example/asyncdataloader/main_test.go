package asyncdataloader

import (
	"context"
	"database/sql"
	"fmt"
	"iter"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/BooleanCat/go-functional/v2/it"
	"github.com/mackee/iterutils/async"
	"github.com/samber/lo"
	"github.com/vikstrous/dataloadgen"
)

type UserDetail struct {
	Name  string
	Items []string
}

func BenchmarkUserDetails__loopQuery(b *testing.B) {
	ctx := context.Background()
	db, err := NewDB()
	if err != nil {
		b.Fatalf("failed to open database: %v", err)
	}
	for range b.N {
		ids := it.Map(it.Take(randIntGen(10000), 500), func(i int64) UserID {
			return UserID(i)
		})
		uds := it.Map2(it.Enumerate(ids), func(_ int, id UserID) (*UserDetail, error) {
			ud, err := getUserDetail(ctx, db, id)
			if err != nil {
				return nil, fmt.Errorf("failed to get user detail: %w", err)
			}
			return ud, nil
		})
		_, err = it.TryCollect(uds)
		if err != nil {
			b.Fatalf("failed to get user details: %v", err)
		}
	}
}

func getUserDetail(ctx context.Context, db *sql.DB, id UserID) (*UserDetail, error) {
	user, err := NewUsersSQL().Select().ID(id).SingleContext(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	userItems, closer := NewUserItemsSQL().Select().UserID(id).IterContext(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("failed to get user items: %w", err)
	}
	defer closer()
	itemNames := make([]string, 0, 10)
	for ui, err := range userItems {
		if err != nil {
			return nil, fmt.Errorf("failed to get user item: %w", err)
		}
		item, err := NewItemsSQL().Select().ID(ui.ItemID).SingleContext(ctx, db)
		if err != nil {
			return nil, fmt.Errorf("failed to get item: %w", err)
		}
		itemNames = append(itemNames, item.Name)
	}
	return &UserDetail{
		Name:  user.Name,
		Items: itemNames,
	}, nil
}

func BenchmarkUserDetails__dataloader(b *testing.B) {
	ctx := context.Background()
	db, err := NewDB()
	if err != nil {
		b.Fatalf("failed to open database: %v", err)
	}
	for range b.N {
		dl := newDataloader(db)
		ids := it.Map(it.Take(randIntGen(10000), 500), func(i int64) UserID {
			return UserID(i)
		})
		uds := async.Map2(it.Enumerate(ids), func(_ int, id UserID) (*UserDetail, error) {
			ud, err := getUserDetailByDataloader(ctx, dl, id)
			if err != nil {
				return nil, fmt.Errorf("failed to get user detail: %w", err)
			}
			return ud, nil
		})
		_, err = it.TryCollect(uds)
		if err != nil {
			b.Fatalf("failed to get user details: %v", err)
		}
	}
}

func getUserDetailByDataloader(ctx context.Context, dl *dataloader, id UserID) (*UserDetail, error) {
	user, err := dl.users.Load(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	userItems, err := dl.userItemsLoader.Load(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user items: %w", err)
	}
	_items := async.Map2(slices.All(userItems), func(_ int, ui *UserItems) (*Items, error) {
		return dl.itemsLoader.Load(ctx, ui.ItemID)
	})
	items, err := it.TryCollect(_items)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	itemNames := lo.Map(items, func(i *Items, _ int) string { return i.Name })
	return &UserDetail{
		Name:  user.Name,
		Items: itemNames,
	}, nil
}

type dataloader struct {
	users           *dataloadgen.Loader[UserID, *Users]
	userItemsLoader *dataloadgen.Loader[UserID, []*UserItems]
	itemsLoader     *dataloadgen.Loader[ItemID, *Items]
}

func newDataloader(db *sql.DB) *dataloader {
	usersLoader := dataloadgen.NewLoader(
		func(ctx context.Context, ids []UserID) ([]*Users, []error) {
			_users, err := NewUsersSQL().Select().IDIn(ids...).AllContext(ctx, db)
			if err != nil {
				return nil, []error{err}
			}
			users := lo.ToSlicePtr(_users)
			return sortByIndex(users, ids, func(u *Users) UserID { return u.ID }), nil
		},
	)
	userItemsLoader := dataloadgen.NewLoader(
		func(ctx context.Context, ids []UserID) ([][]*UserItems, []error) {
			_userItems, err := NewUserItemsSQL().Select().UserIDIn(ids...).AllContext(ctx, db)
			if err != nil {
				return nil, []error{err}
			}
			userItemsMap := make(map[UserID][]*UserItems, len(ids))
			for _, ui := range _userItems {
				userItemsMap[ui.UserID] = append(userItemsMap[ui.UserID], &ui)
			}
			userItems := make([][]*UserItems, 0, len(ids))
			for _, id := range ids {
				userItems = append(userItems, userItemsMap[id])
			}
			return userItems, nil
		},
	)
	itemsLoader := dataloadgen.NewLoader(
		func(ctx context.Context, ids []ItemID) ([]*Items, []error) {
			_items, err := NewItemsSQL().Select().IDIn(ids...).AllContext(ctx, db)
			if err != nil {
				return nil, []error{err}
			}
			items := lo.ToSlicePtr(_items)
			return sortByIndex(items, ids, func(i *Items) ItemID { return i.ID }), nil
		},
	)

	return &dataloader{
		users:           usersLoader,
		userItemsLoader: userItemsLoader,
		itemsLoader:     itemsLoader,
	}
}

func sortByIndex[T any, K comparable](items []T, keys []K, index func(T) K) []T {
	byIndex := make(map[K]T, len(items))
	for _, item := range items {
		byIndex[index(item)] = item
	}
	sorted := make([]T, 0, len(items))
	for _, key := range keys {
		sorted = append(sorted, byIndex[key])
	}
	return sorted
}

func randIntGen(max int64) iter.Seq[int64] {
	return func(fn func(int64) bool) {
		for {
			i := rand.Int64N(max) + 1
			if !fn(i) {
				break
			}
		}
	}
}
