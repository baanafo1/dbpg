package dbpg

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	ref "github.com/intdxdt/goreflect"
)

func Insert[T ITable[T]](conn *sql.DB, model T, insertCols []string, on On) (bool, error) {
	var fields, err = ref.Fields(model)
	if err != nil {
		return false, err
	}

	fields, colRefs, err := ref.FilterFieldReferences(fields, model)
	if err != nil {
		return false, err
	}

	var cols = make([]string, 0, len(fields))
	var values = make([]any, 0, len(fields))

	var dict = KeysToMap(insertCols, true)

	for i, field := range fields {
		if !(dict[field]) {
			continue
		}
		cols = append(cols, field)
		values = append(values, colRefs[i])
	}

	var columns = ColumnNames(cols)
	var holders = ColumnPlaceholders(cols)
	fmt.Println(values)
	for _, value := range colRefs {
		fmt.Println(*&value)
	}
	fmt.Printf("columns: %v\n holders: %v\n", columns, holders)

	var sqlStatement = fmt.Sprintf(`
		INSERT INTO %v(%v) 
		VALUES (%v);`, model.TableName(), columns, holders)

	if len(on.On) > 0 {
		fmt.Println("i entered here")
		sqlStatement = fmt.Sprintf(`
		INSERT INTO %v(%v) 
		VALUES (%v)
		ON %v;`, model.TableName(), columns, holders, on.On)
		for _, v := range on.Arguments {
			values = append(values, v)
		}
	}

	fmt.Println(sqlStatement)
	res, err := Exec(conn, sqlStatement, values...)
	if err != nil {
		fmt.Println("Error is from here", err)
		return false, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return count == 1, nil
}