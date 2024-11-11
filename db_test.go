package dbpg

import (
	"fmt"
	"testing"

	_ "github.com/lib/pq"

	"github.com/franela/goblin"
)

const UserSQLModel = `
DROP TABLE IF EXISTS newtest;
CREATE TABLE IF NOT EXISTS newtest (
	id            		 INTEGER NOT NULL PRIMARY KEY,
	name          		 TEXT DEFAULT '',
	geom	   			 TEXT DEFAULT '',
	centroid			 TEXT DEFAULT '',
	number_of_floors     INTEGER DEFAULT 0,
	email         		 TEXT NOT NULL UNIQUE,
    has_paid        	 INTEGER DEFAULT 0
);
`

var dbInstance *Database

type Model struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Geom           string `json:"geom"`
	Centroid       string `json:"centroid"`
	NumberOfFloors int    `json:"number_of_floors"`
	Email          string `json:"email"`
	HasPaid        int    `json:"has_paid"`
}

func NewModel() *Model {
	return &Model{ID: -1}
}

func (model *Model) New() *Model {
	return NewModel()
}

func (model *Model) Clone() *Model {
	var o = *model
	return &o
}

func (model *Model) TableName() string {
	return "newtest"
}

func (model *Model) Insert() (bool, error) {
	// You may want to pass all fields or just the necessary ones
	success, err := Insert(dbInstance.Conn, model, []string{
		`id`, `email`, `name`, `geom`,
	}, On{})

	if err != nil {
		fmt.Printf("Error inserting model: %v\n", err) // Log the error
	}

	return success, err
}

func (model *Model) Delete() (int64, error) {
	clause := WhereClause{Where: fmt.Sprintf("id = %v", model.ID)}
	success, err := Delete(dbInstance.Conn, model, clause)

	if err != nil {
		fmt.Printf("Error deleting model: %v\n", err)
	}
	return success, err
}

func (model *Model) QueriesByColumnNames() ([]*Model, error) {
	fields := []string{`id`, `geom`, `email`}
	result, err := QueriesByColumnNames(dbInstance.Conn, model, fields)
	if err != nil {
		fmt.Printf("Error querying model by column names: %v\n", err)
	}
	return result, err
}

func (model *Model) Update() (bool, error) {
	field := []string{`name`, `geom`, `number_of_floors`}
	whereClause := WhereClause{
		Where:     fmt.Sprintf("email = '%v'", model.Email),
		Arguments: []any{"plot1", "geometry absent", 5},
	}
	res, err := Update(dbInstance.Conn, model, field, whereClause)

	if err != nil {
		fmt.Printf("Error updating model: %v\n", err)
	}
	return res, err
}

func (model *Model) UpdateByExclusion() (bool, error) {
	excludedFields := []string{`id`, `name`, `centroid`, `number_of_floors`, `email`, `has_paid`}
	whereClause := WhereClause{
		Where:     fmt.Sprintf("number_of_floors = %v", model.NumberOfFloors),
		Arguments: []any{"updated by exclusion"},
	}
	res, err := UpdateByExclusion(dbInstance.Conn, model, excludedFields, whereClause)
	if err != nil {
		fmt.Printf("Error updating model by exclusion: %v\n", err)
	}
	return res, err
}

func initDB() {
	connStr := "postgres://postgres:secret@localhost:5433/demo?sslmode=disable"
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

	//_, err = dbInstance.Exec(UserSQLModel)
	//if err != nil {
	//	fmt.Printf("Error executing UserSQLModel: %v\n", err)
	//	return
	//}
	fmt.Println("Database connection successful")
}

func deInitDB() {
	if dbInstance != nil {
		dbInstance.Close()
	}
}

func TestDBPG(t *testing.T) {
	g := goblin.Goblin(t)

	initDB()
	defer deInitDB()
	var m = NewModel()

	// INSERT TEST
	//g.Describe("Tests Model Insert", func() {
	//	g.It("user insert", func() {
	//		m.ID = 2
	//		m.Email = "secondUser@db.com"
	//		m.Name = "plot2"
	//		m.Geom = "geom missing"
	//		bln, err := m.Insert()
	//		g.Assert(bln).IsTrue()
	//		g.Assert(err).IsNil()
	//	})
	//})
	//
	//// MODEL DELETE TEST
	//g.Describe("Test Model Delete", func() {
	//	g.It("user delete", func() {
	//		m.ID = 1
	//		intg, err := m.Delete()
	//		g.Assert(intg).Equal(int64(1))
	//		g.Assert(err).IsNil()
	//	})
	//})
	//
	//// TEST QUERY BY COLUMN NAMES
	//g.Describe("Test QueryByColumnNames", func() {
	//	g.It("query column names", func() {
	//		res, err := m.QueriesByColumnNames()
	//
	//		if len(res) > 0 {
	//			expected := []string{"1", "geom"}
	//			actual := []string{
	//				fmt.Sprintf("%v", res[0].ID),
	//				res[0].Geom,
	//			}
	//			g.Assert(actual).Equal(expected)
	//		} else {
	//			t.Error("Expected at least one result in res")
	//		}
	//
	//		g.Assert(err).IsNil()
	//	})
	//})
	//
	//// TEST UPDATE
	//g.Describe("Test Update", func() {
	//	g.It("update columns", func() {
	//		m.Email = "firstUser@db.com"
	//		res, err := m.Update()
	//		g.Assert(err).IsNil()
	//
	//		g.Assert(res).IsTrue()
	//	})
	//})

	// TEST UPDATE BY EXCLUSION
	g.Describe("Test UpdateByExclusion", func() {
		g.It("update by exclusion", func() {
			m.NumberOfFloors = 5
			res, err := m.UpdateByExclusion()
			g.Assert(err).IsNil()
			g.Assert(res).IsTrue()
		})
	})
}
