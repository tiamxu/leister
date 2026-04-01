package database

import (
	"fmt"

	"github.com/tiamxu/kit/sql"
)

var db *sql.DB

func Connect(dbConfig *Config) (err error) {
	kitConfig := &sql.Config{
		Driver:          dbConfig.Driver,
		Database:        dbConfig.Database,
		Username:        dbConfig.Username,
		Password:        dbConfig.Password,
		Host:            dbConfig.Host,
		Port:            dbConfig.Port,
		MaxIdleConns:    dbConfig.MaxIdleConns,
		MaxOpenConns:    dbConfig.MaxOpenConns,
		ConnMaxLifetime: dbConfig.ConnMaxLifetime,
	}
	db, err = sql.Connect(kitConfig)
	return
}

func AddItem(item Item) (int64, error) {
	i, err := insertItem(item)
	return i, err
}

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
