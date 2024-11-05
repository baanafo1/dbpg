package dbpg

import (
	"fmt"
	"testing"

	_ "github.com/lib/pq"

	"github.com/franela/goblin"
)

const UserSQLModel = `
DROP TABLE IF EXISTS test;
CREATE TABLE IF NOT EXISTS test (
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
	ID             int           `json:"id"`
	Name           string		 `json:"name"`
	Geom           string		 `json:"geom"`
	Centriod       string		 `json:"centriod"`
	NumberOfFloors string	     `json:"number_of_floors"`
	Email          string		 `json:"email"`
	HasPaid        int			 `json:"has_paid"`
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
	return "test"
}

func (model *Model) Insert() (bool, error) {
    // You may want to pass all fields or just the necessary ones
    success, err := Insert(dbInstance.Conn, model, []string{
        `id`, `email`, `name`,`geom`,
    }, On{})

    if err != nil {
        fmt.Printf("Error inserting model: %v\n", err) // Log the error
    }

    return success, err
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
	
    _, err = dbInstance.Exec(UserSQLModel)
    if err != nil {
		fmt.Printf("Error executing UserSQLModel: %v\n", err)
    }
	// fmt.Println(&dbInstance)
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

	g.Describe("Tests Model Insert", func() {
		g.It("user insert", func() {
			var m = NewModel()
			m.ID = 1
		    m.Email = "email@db.com"
			m.Name = "test"
			m.Geom = "geom"
			// m.Address = "123 db street"
			bln, err := m.Insert()
			g.Assert(bln).IsTrue()
			g.Assert(err).IsNil()
		})
	})
}