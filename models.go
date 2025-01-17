package dbpg

import "database/sql"


type Model struct {
	RiD             int    			`json:"rid"`
	Name            string 			`json:"name"`
	PiD             int    			`json:"pid"`
	Geom            sql.NullString 	`json:"geom"`
	PlotNumber      int    			`json:"plot_num"`
	Use             string 			`json:"use"`
	UseType         string 			`json:"use_type"`
	State           string 			`json:"state"`
    Remarks         string 			`json:"remarks"`
}

func NewModel() *Model {
	return &Model{PiD: -1}
}

func (model *Model) New() *Model {
	return NewModel()
}

func (model *Model) Clone() *Model {
	var o = *model
	return &o
}

func (model *Model) TableName() string {
	return "building"
}