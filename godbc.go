package godbc

import (
	"database/sql"
	"fmt"
)
import _ "github.com/go-sql-driver/mysql"

// 创建DB对象;
// driver: 驱动名, 如mysql
// username/password/host/port: 连接信息
// dbName: 数据库名
// ping: 是否执行ping操作
func CreateDb(driver, username, password, host string, port int, dbName string, ping bool) (*sql.DB, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, dbName)
	db, err := sql.Open(driver, connString)
	if nil != err {
		return nil, err
	}

	if ping {
		err = db.Ping()
		if nil != err {
			return nil, err
		}
	}

	return db, nil
}

// 执行修改语句;
func ExecuteUpdate(db *sql.DB, updateSql string, args ...interface{}) (sql.Result, error) {
	stmt, err := db.Prepare(updateSql)
	if nil != err {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Exec(args...)
}

// 执行查询语句;
// processor: 行处理函数, 每读取到一行都会调用一次processor
func ExecuteQuery(db *sql.DB, querySql string, processor func(result *RowResult), args ...interface{}) error {
	stmt, err := db.Prepare(querySql)
	if nil != err {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if nil != err {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if nil != err {
		return err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range columns {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if nil != err {
			return err
		}

		result := &RowResult{values}
		processor(result)
	}

	return nil
}

// 执行查询;
// scanner: 用户传入函数, 直接操作sql.Rows对象读取数据
func ExecuteScan(db *sql.DB, querySql string, scanner func(rows *sql.Rows) error, args ...interface{}) error {
	stmt, err := db.Prepare(querySql)
	if nil != err {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if nil != err {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if nil != err {
		return err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range columns {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = scanner(rows)
		if nil != err {
			return err
		}
	}

	return nil
}
