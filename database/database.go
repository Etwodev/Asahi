package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	"xorm.io/xorm/names"

	c "github.com/SpeedSlime/Asahi/config"
)

var engine 	*xorm.Engine

func Connect() error {
	db, err := xorm.NewEngine(c.DRV(), c.DSN())
	if err != nil {
		return fmt.Errorf("Connect: an error has occured: %s", err)
	}
	db.SetMapper(names.SameMapper{})
	engine = db
	return nil
}

func Create(tables ...interface{}) error {
	err := engine.CreateTables(tables...)
	if err != nil {
		return fmt.Errorf("Create: failed creating table: %w", err)
	}
	return nil
}

func Delete(table, record interface{}) error {
	_, err := engine.Delete(record)
	if err != nil {
		return fmt.Errorf("Delete: failed deleting record: %w", err)
	}
	return nil
}

func Update(record interface{}, cols ...string) error {
	db := engine.AllCols()
	if cols != nil {
		db = engine.Cols(cols...)
	}
	_, err := db.Update(record)
	if err != nil {
		return fmt.Errorf("Update: failed updating record: %w", err)
	}
	return nil
}

func Select(record interface{}) (bool, error) {
	has, err := engine.Get(record)
	if err != nil {
		return has, fmt.Errorf("Select: failed selecting record: %w", err)
	}
	return has, nil
}

func Exists(record interface{}) bool {
	has, _ := engine.Exist(record)
	return has
}

func Find(records []interface{}, cond interface{}) (bool, error) {
	has := Exists(cond)
	if !has {
		return has, nil
	}
	err := engine.Find(records, cond)
	if err != nil {
		return has, fmt.Errorf("Find: failed finding record: %w", err)
	}
	return has, nil
}

func Insert(record ...interface{}) error {
	_, err := engine.Insert(record...)
	if err != nil {
		return fmt.Errorf("Insert: failed inserting item: %w", err)
	}
	return nil
}

func Query(sql string) ([]map[string][]byte, error) {
	results, err := engine.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("Query: failed querying item: %w", err)
	}
	return results, nil
}