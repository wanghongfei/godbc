package godbc

import (
	"database/sql"
	"errors"
	"strconv"
	"time"
)


// 表中的一行数据
type RowResult struct {
	rowData		[]sql.RawBytes
}

var InvalidIndexError error = errors.New("invalid index")

// 获取第index列的值, 以string类型输出, index从0开始
func (res *RowResult) GetString(index int) (string, error) {
	if err := res.checkIndex(index); err != nil {
		return "", err
	}

	return string(res.rowData[index]), nil
}

func (res *RowResult) GetInt(index int) (int, error) {
	str, err := res.GetString(index)
	if nil != err {
		return 0, err
	}

	return strconv.Atoi(str)
}

func (res *RowResult) GetInt64(index int) (int64, error) {
	str, err := res.GetString(index)
	if nil != err {
		return 0, err
	}

	return strconv.ParseInt(str, 10, 64)
}

// pattern: 日期字符串的格式, 如2006-01-02 15:04:05
func (res *RowResult) GetTime(index int, pattern string) (*time.Time, error) {
	str, err := res.GetString(index)
	if nil != err {
		return nil, err
	}

	t, err := time.Parse(pattern, str)
	if nil != err {
		return nil, err
	}

	return &t, nil
}

func (res *RowResult) checkIndex(index int) error {
	if index < 0 || index >= len(res.rowData) {
		return InvalidIndexError
	}

	return nil
}

