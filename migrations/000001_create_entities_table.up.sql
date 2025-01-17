CREATE TABLE IF NOT EXISTS region (
    name        TEXT NOT NULL,
    id          INTEGER PRIMARY KEY DEFAULT 1,
    extent      GEOMETRY(Polygon, 4326) NOT NULL
);

CREATE TABLE IF NOT EXISTS building (
    rid                 INTEGER,
	name          		TEXT DEFAULT '',
	pid            		INTEGER NOT NULL PRIMARY KEY,
    geom                GEOMETRY(Polygon, 4326) NOT NULL,
	plot_num	   	    INTEGER,
	use			        TEXT DEFAULT '',
	use_type            TEXT DEFAULT '',
	state         		TEXT DEFAULT '',
    remarks        	    TEXT DEFAULT '',
    FOREIGN KEY (rid) REFERENCES region(id) 
);