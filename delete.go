package dbpg

import (
	"database/sql"
	"fmt"
	//_ "github.com/lib/pq"
)

func Delete[T ITable[T]](conn *sql.DB, model T, wc WhereClause) (int64, error) {
	var query = fmt.Sprintf(
		`DELETE FROM %v WHERE %v;`, model.TableName(), wc.Where)

	fmt.Println(query)
	var res, err = Exec(conn, query, wc.Arguments...)
	fmt.Println(res.RowsAffected())
	if err != nil {
		return 0, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}
