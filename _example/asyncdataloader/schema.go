package asyncdataloader

import "github.com/mackee/go-genddl/index"

//go:generate go run github.com/mackee/go-sqlla/v2/cmd/sqlla
//go:generate go run github.com/mackee/go-genddl/cmd/genddl -outpath=./mysql.sql -driver=mysql -outerforeignkey

type UserID int64

//sqlla:table users
//genddl:table users
type Users struct {
	ID   UserID `db:"id,autoincrement,primarykey"`
	Name string `db:"name"`
}

type ItemID int64

//sqlla:table items
//genddl:table items
type Items struct {
	ID   ItemID `db:"id,autoincrement,primarykey"`
	Name string `db:"name"`
}

type UserItemID int64

//sqlla:table user_items
//genddl:table user_items
type UserItems struct {
	ID     UserItemID `db:"id,autoincrement,primarykey"`
	UserID UserID     `db:"user_id"`
	ItemID ItemID     `db:"item_id"`
}

func (u UserItems) _schemaIndex(methods index.Methods) []index.Definition {
	return []index.Definition{
		methods.ForeignKey(u.UserID, Users{}.ID, index.ForeignKeyDeleteCascade),
		methods.ForeignKey(u.ItemID, Items{}.ID, index.ForeignKeyDeleteCascade),
	}
}
