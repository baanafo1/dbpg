CREATE TABLE IF NOT EXISTS region (
    name        TEXT NOT NULL,
    id          INTEGER PRIMARY KEY DEFAULT 1,
    extent      GEOMETRY(Polygon, 4326) NOT NULL
);

INSERT INTO region (name, ID, EXTENT) VALUES
    ('North America', 1, ST_GeomFromText('POLYGON((-170 75, -50 75, -50 20, -170 20, -170 75))', 4326)),
    ('South America', 2, ST_GeomFromText('POLYGON((-90 12, -35 12, -35 -55, -90 -55, -90 12))', 4326)),
    ('Europe', 3, ST_GeomFromText('POLYGON((-10 35, 40 35, 40 70, -10 70, -10 35))', 4326)),
    ('Africa', 4, ST_GeomFromText('POLYGON((-20 35, 50 35, 50 -35, -20 -35, -20 35))', 4326)),
    ('Asia', 5, ST_GeomFromText('POLYGON((-180 10, -70 10, -70 75, 180 75, 180 10, -180 10))', 4326)),
    ('Australia', 6, ST_GeomFromText('POLYGON((110 -10, 160 -10, 160 -45, 110 -45, 110 -10))', 4326)),
    ('Antarctica', 7, ST_GeomFromText('POLYGON((-180 -60, -50 -60, -50 -90, 180 -90, 180 -60, -180 -60))', 4326));

    
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