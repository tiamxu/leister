package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Connect(dbConfig *Config) (err error) {
	db, err = sql.Open(dbConfig.Driver, dbConfig.Source())
	if err != nil {
		return
	}
	err = db.Ping()
	if err != nil {
		return
	}
	db.SetMaxOpenConns(dbConfig.MaxOpenConns)
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Second)
	return
}

func AddItem(item Item) (int64, error) {
	i, err := insertItem(item)
	return i, err
}

// 查询单行
func QueryRowDB(sql string) *sql.Row {
	return db.QueryRow(sql)
}

// 查询多行
func QueryDB(sql string) (*sql.Rows, error) {
	return db.Query(sql)
}

// func QueryItemWithName() (Item, error) {
// 	var item Item
// 	return item, nil
// }

// 查询所有数据
func GetAllItemData() ([]Item, error) {
	return QueryItemWithCon("")
}

func SelectItemByWhereWithGroup(group string, arg ...interface{}) ([]Item, error) {
	whereSql := fmt.Sprintf("where app_group='%s'", group)
	return QueryItemWithCon(whereSql)
}

func SelectItemByWhereWithName(name, group string, arg ...interface{}) ([]Item, error) {
	whereSql := fmt.Sprintf("where app_name='%s' and app_group='%s'", name, group)
	return QueryItemWithCon(whereSql)
}
