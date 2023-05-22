package database

import "log"

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

func insertItem(item Item) (int64, error) {
	sqlStr := `INSERT INTO item (code_id,app_name, app_group, app_type,ssh_url_to_repo, http_url_to_repo) VALUES (?,?,?,?,?,?);`
	return ModifyDB(sqlStr, item.CodeID, item.AppName, item.AppGroup, item.AppType, item.SSHURLToRepo, item.HTTPURLToRepo)
}

func QueryItemWithCon(sql string) ([]Item, error) {
	sql = "select code_id,app_name,app_group, http_url_to_repo from item " + sql
	rows, err := QueryDB(sql)
	if err != nil {
		return nil, err
	}
	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.CodeID, &item.AppName, &item.AppGroup, &item.HTTPURLToRepo)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}
	return items, nil
}
