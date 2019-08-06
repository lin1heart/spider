package db

import (
	"fmt"
	"github.com/lin1heart/spider/go/src/util"
)

func QueryPhotoAfterRow(rowIndex int, inteval int) map[int]map[string]string {
	sqlString := fmt.Sprintf("SELECT * FROM photo WHERE (url != '' AND url is not null ) AND id > %d LIMIT %d", rowIndex, inteval)
	rows, err := Mysql.Query(sqlString)
	util.CheckError(err)

	results, err := ConvertToRows(rows)
	return results
}
func UpdateUrlToEmpty(id int) {
	stmt, err := Mysql.Prepare(`UPDATE photo SET url= "" WHERE id=?`)
	util.CheckError(err)
	defer stmt.Close()

	res, err := stmt.Exec(id)
	util.CheckError(err)
	num, err := res.RowsAffected()
	fmt.Printf("UpdatePhotoUrl id %d affect row %d \n", id, num)
	util.CheckError(err)
}
