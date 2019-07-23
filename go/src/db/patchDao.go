package db

import (
	"fmt"
	"github.com/lin1heart/spider/go/src/util"
)

func QueryPhotoAfterRow(rowIndex int) map[int]map[string]string {
	sqlString := fmt.Sprintf("SELECT * FROM photo LIMIT %d,50", rowIndex)
	rows, err := Mysql.Query(sqlString)
	util.CheckError(err)

	results, err := ConvertToRows(rows)
	return results
}
