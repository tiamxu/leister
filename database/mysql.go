package database

import (
	"database/sql"
	"log"
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

// 操作数据库
func ModifyDB(sql string, args ...interface{}) (int64, error) {
	result, err := db.Exec(sql, args...)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return count, nil
}
func AddItem(item Item) (int64, error) {
	i, err := insertItem(item)
	return i, err
}
func insertItem(item Item) (int64, error) {
	sqlStr := `INSERT INTO item (code_id,app_name, app_group, app_type,ssh_url_to_repo, http_url_to_repo) VALUES (?,?,?,?,?,?);`
	return ModifyDB(sqlStr, item.CodeID, item.AppName, item.AppGroup, item.AppType, item.SSHURLToRepo, item.HTTPURLToRepo)
}
