package database

import (
	"github.com/tiamxu/kit/sql"
	"github.com/tiamxu/leister/config"
)

var db *sql.DB

func Connect(dbConfig *config.DB) (err error) {
	// 转换为 kit/sql 的 Config 结构
	kConfig := &sql.Config{
		Driver:          dbConfig.Driver,
		Host:            dbConfig.Host,
		Port:            dbConfig.Port,
		Database:        dbConfig.Database,
		Username:        dbConfig.Username,
		Password:        dbConfig.Password,
		MaxOpenConns:    dbConfig.MaxOpenConns,
		MaxIdleConns:    dbConfig.MaxIdleConns,
		ConnMaxLifetime: dbConfig.ConnMaxLifetime,
		ConnMaxIdleTime: 60, // 默认值
	}

	db, err = sql.Connect(kConfig)
	return
}

// GetDB 返回数据库连接
func GetDB() *sql.DB {
	return db
}

// IsNoRows 检查是否为记录不存在错误
func IsNoRows(err error) bool {
	return sql.IsNoRows(err)
}
