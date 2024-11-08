package dbpg

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
)

func Delete[T ITable[T]](conn *sql.DB, model T, wc WhereClause) (int64, error) {
	var query = fmt.Sprintf(
		`DELETE FROM %v WHERE %v;`, model.TableName(), wc.Where)

	var res, err = Exec(conn, query, wc.Arguments...)
	if err != nil {
		return 0, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}