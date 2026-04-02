package database

import (
	"github.com/tiamxu/kit/log"
	"github.com/jmoiron/sqlx"
)

// AddItem 添加项目
func AddItem(item Item) (int64, error) {
	query := `INSERT INTO item (code_id, app_name, app_group, app_type, ssh_url_to_repo, http_url_to_repo) 
			 VALUES (:code_id, :app_name, :app_group, :app_type, :ssh_url_to_repo, :http_url_to_repo)`
	
	result, err := db.NamedExec(query, item)
	if err != nil {
		log.Errorf("AddItem error: %v", err)
		return 0, err
	}
	
	return result.LastInsertId()
}

// GetAllItemData 查询所有数据
func GetAllItemData() ([]Item, error) {
	var items []Item
	err := db.Select(&items, "SELECT code_id, app_name, app_group, http_url_to_repo FROM item")
	if err != nil {
		log.Errorf("GetAllItemData error: %v", err)
		return nil, err
	}
	return items, nil
}

// SelectItemByWhereWithGroup 按组查询
func SelectItemByWhereWithGroup(group string) ([]Item, error) {
	var items []Item
	err := db.Select(&items, "SELECT code_id, app_name, app_group, http_url_to_repo FROM item WHERE app_group = ?", group)
	if err != nil {
		log.Errorf("SelectItemByWhereWithGroup error: %v", err)
		return nil, err
	}
	return items, nil
}

// SelectItemByWhereWithName 按名称和组查询
func SelectItemByWhereWithName(name, group string) ([]Item, error) {
	var items []Item
	err := db.Select(&items, "SELECT code_id, app_name, app_group, http_url_to_repo FROM item WHERE app_name = ? AND app_group = ?", name, group)
	if err != nil {
		log.Errorf("SelectItemByWhereWithName error: %v", err)
		return nil, err
	}
	return items, nil
}

// AddItems 批量添加项目（事务）
func AddItems(items []Item) error {
	return db.TransactCallback(func(tx *sqlx.Tx) error {
		for _, item := range items {
			query := `INSERT INTO item (code_id, app_name, app_group, app_type, ssh_url_to_repo, http_url_to_repo) 
					 VALUES (:code_id, :app_name, :app_group, :app_type, :ssh_url_to_repo, :http_url_to_repo)`
			
			_, err := tx.NamedExec(query, item)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
