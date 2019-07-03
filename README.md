# GODBC

对Go标准库中的sql接口进行简单封装，模拟Java的JDBC使用方式。

example:

```go
func TestExecuteQuery(t *testing.T) {
	// 创建DB对象
	db, err := CreateDb("mysql", "root", "", "127.0.0.1", 3306, "gosql", true)
	if nil != err {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// 执行查询
	ExecuteQuery(db, "SELECT * FROM engine_task ORDER BY task_id DESC limit ?", func(result *RowResult) {
		// 获取第8列, 以string返回
		val, _ := result.GetString(7)
		fmt.Println(val)

		id, _ := result.GetInt(0)
		fmt.Println(id)

		time, _ := result.GetTime(4, "2006-01-02 15:04:05")
		fmt.Println(time)

		time, _ = result.GetTime(5, "2006-01-02 15:04:05")
		fmt.Println(time)
	}, 2)
	
}

func TestExecuteUpdate(t *testing.T) {
	db, err := CreateDb("mysql", "root", "", "127.0.0.1", 3306, "gosql", true)
	if nil != err {
		t.Fatal(err)
		return
	}
	defer db.Close()

	res, err := ExecuteUpdate(db, "UPDATE engine_task SET output_path = ? WHERE task_id = ?", "test_path", 3704)
	if nil != err {
		t.Fatal(err)
		return
	}

	fmt.Println(res.RowsAffected())
}

func TestExecuteScan(t *testing.T) {
	db, err := CreateDb("mysql", "root", "", "127.0.0.1", 3306, "gosql", true)
	if nil != err {
		t.Fatal(err)
		return
	}
	defer db.Close()

	ExecuteScan(db, "select task_id, status FROM engine_task ORDER BY task_id DESC limit ?", func(rows *sql.Rows) error {
		id := 0
		status := ""
		rows.Scan(&id, &status)

		fmt.Printf("%d, %s\n", id, status)

		return nil
	}, 3)
}
```

