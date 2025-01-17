package dbpg

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/franela/goblin"
)

var dbInstance *Database

func (model *Model) Insert() (bool, error) {
	// You may want to pass all fields or just the necessary ones
    // success, err := Insert(dbInstance.Conn, )
	success, err := Insert(dbInstance.Conn, model, []string{
		`rid`, `pid`, `name`, `geom`, `plot_num`,
	}, On{})

	if err != nil {
		fmt.Printf("Error inserting model: %v\n", err) // Log the error
	}

	return success, err
}

func (model *Model) Delete() (int64, error) {
	clause := WhereClause{Where: fmt.Sprintf("pid = %v", model.PiD)}
	success, err := Delete(dbInstance.Conn, model, clause)

	if err != nil {
		fmt.Printf("Error deleting model: %v\n", err)
	}
	return success, err
}

func (model *Model) QueriesByColumnNames() ([]*Model, error) {
	fields := []string{`rid`, `pid`}
    wc := WhereClause{Where: "plot_num = $1", Arguments: []any{model.PlotNumber}}
	result, err := QueriesByColumnNames(dbInstance.Conn, model, fields, wc)
	if err != nil {
		fmt.Printf("Error querying model by column names: %v\n", err)
	}
	return result, err
}

func (model *Model) Update() (bool, error) {
	field := []string{`use`, `state`}
	whereClause := WhereClause{
		Where:     fmt.Sprintf("plot_num = %d", model.PlotNumber),
		Arguments: []any{model.Use, model.State},
	}
	res, err := Update(dbInstance.Conn, model, field, whereClause)

	if err != nil {
		fmt.Printf("Error updating model: %v\n", err)
	}
	return res, err
}

func (model *Model) UpdateByExclusion() (bool, error) {
	excludedFields := []string{`rid`, `name`, `pid`, `geom`, `use`, `plot_num`, `use_type`}
	whereClause := WhereClause{
		Where:     fmt.Sprintf("pid = %v", model.PiD),
		Arguments: []any{model.State, model.Remarks},
	}
	res, err := UpdateByExclusion(dbInstance.Conn, model, excludedFields, whereClause)
	if err != nil {
		fmt.Printf("Error updating model by exclusion: %v\n", err)
	}
	return res, err
}

func initDB() {

	connStr := os.Getenv("GREENLIGHT_DB_DSN")
	var err error
	dbInstance, err = NewDatabase(connStr)
	if err != nil {
		fmt.Printf("Database connection failed: %v\n", err)
		return
	}
	if dbInstance == nil {
		fmt.Println("Database instance is nil after NewDatabase()")
		return
	}

	// _, err = dbInstance.Exec(BuildingSQLModel)
	// if err != nil {
	// 	fmt.Printf("Error executing UserSQLModel: %v\n", err)
	// 	return
	// }
	fmt.Println("Database connection successful")
}

func deInitDB() {
	if dbInstance != nil {
		dbInstance.Close()
	}
}

func TestDBPG(t *testing.T) {
	g := goblin.Goblin(t)

    geom := "POLYGON((-73.9775 40.7694, -73.9681 40.7694, -73.9681 40.7831, -73.9775 40.7831, -73.9775 40.7694))"
	initDB()
	defer deInitDB()
	var m = NewModel()

	//INSERT TEST
	g.Describe("Tests Model Insert", func() {
		g.It("user insert", func() {
			m.RiD   = 3
			m.PiD = 4
			m.Name = "Washey Hostel"
			m.PlotNumber = 321
            m.Geom = sql.NullString{String: geom, Valid: true}
			bln, err := m.Insert()
			g.Assert(bln).IsTrue()
			g.Assert(err).IsNil()
		})
	})

	// MODEL DELETE TEST
	g.Describe("Test Model Delete", func() {
		g.It("user delete", func() {
			m.PiD = 2
			intg, err := m.Delete()
			g.Assert(intg).Equal(int64(1))
			g.Assert(err).IsNil()
		})
	})

	// TEST QUERY BY COLUMN NAMES
	g.Describe("Test QueryByColumnNames", func() {
		g.It("query column names", func() {
            m.PlotNumber = 125
			res, err := m.QueriesByColumnNames()

			if len(res) > 0 {
				expected := []int{2, 3}
				actual := []int{
					res[0].RiD,
					res[0].PiD,
				}
				g.Assert(actual).Equal(expected)
			} else {
				t.Error("Expected at least one result in res")
			}

			g.Assert(err).IsNil()
		})
	})

	// TEST UPDATE
	g.Describe("Test Update", func() {
		g.It("update columns", func() {
            m.PlotNumber = 321
            m.Use = "Residential"
            m.State = "Completed"
			res, err := m.Update()
			g.Assert(err).IsNil()

			g.Assert(res).IsTrue()
		})
	})

	// TEST UPDATE BY EXCLUSION
	g.Describe("Test UpdateByExclusion", func() {
		g.It("update by exclusion", func() {
			m.PiD = 4
            m.State = "Uncompleted"
            m.Remarks = "New construction"
			res, err := m.UpdateByExclusion()
			g.Assert(err).IsNil()
			g.Assert(res).IsTrue()
		})
	})
}
