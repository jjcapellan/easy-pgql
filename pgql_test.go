package pgql

import (
	"os"
	"testing"
)

func getDbConnString() string {
	host := os.Getenv("PGHOST")
	port := os.Getenv("PGPORT")
	user := os.Getenv("PGUSER")
	pass := os.Getenv("PGPASSWORD")
	name := os.Getenv("PGDATABASE")
	connStr := "host=" + host + " port=" + port + " user=" + user + " password=" + pass + " dbname=" + name + " sslmode=disable"
	return connStr
}

func delAll() {
	db, _ := openDB(getDbConnString())
	defer db.Close()
	db.Exec("DELETE FROM fixtures")
}

func fill() {
	db, _ := openDB(getDbConnString())
	defer db.Close()
	db.Exec("INSERT INTO fixtures VALUES (1,'One',10),(2,'Two',20),(3,'Three',30)")
}

func restore() {
	delAll()
	fill()
}

func TestGetUpdateStr(t *testing.T) {
	tbl := New("fixtures", getDbConnString())
	result := tbl.getUpdateStr(Data{Key: "col1", Columns: []string{"col2", "col3"}})
	wanted := "UPDATE fixtures SET col2=$1,col3=$2 WHERE col1=$3"
	if result != wanted {
		t.Errorf("Incorrect, got: %s, want: %s.", result, wanted)
	}
}

func TestGetReadStr(t *testing.T) {
	tbl := New("fixtures", getDbConnString())
	result := tbl.getReadStr(Data{Key: "col1", KeyVal: "2", Columns: []string{"col2", "col3"}})
	wanted := "SELECT col2,col3 FROM fixtures WHERE col1=$1"
	if result != wanted {
		t.Errorf("Incorrect, got: %s, want: %s.", result, wanted)
	}

	result = tbl.getReadStr(Data{Columns: []string{"*"}, OrderBy: "col2"})
	wanted = "SELECT * FROM fixtures ORDER BY col2"
	if result != wanted {
		t.Errorf("Incorrect, got: %s, want: %s.", result, wanted)
	}

	result = tbl.getReadStr(Data{Columns: []string{"*"}, OrderBy: "col2", DescOrder: true})
	wanted = "SELECT * FROM fixtures ORDER BY col2 DESC"
	if result != wanted {
		t.Errorf("Incorrect, got: %s, want: %s.", result, wanted)
	}

	result = tbl.getReadStr(Data{Columns: []string{"col2", "col3"}})
	wanted = "SELECT col2,col3 FROM fixtures"
	if result != wanted {
		t.Errorf("Incorrect, got: %s, want: %s.", result, wanted)
	}
}

func TestGetInsertStr(t *testing.T) {
	tbl := New("fixtures", getDbConnString())
	result := tbl.getInsertStr(Data{ColVals: []interface{}{4, "Four", 40}})
	wanted := "INSERT INTO fixtures VALUES ($1,$2,$3)"
	if result != wanted {
		t.Errorf("Incorrect, got: %s, want: %s.", result, wanted)
	}
}

func TestRead(t *testing.T) {
	restore()
	tbl := New("fixtures", getDbConnString())
	result, _ := tbl.Read(Data{Columns: []string{"col2", "col3"}})
	wanted := "Two"
	if result[1]["col2"] != wanted {
		t.Errorf("Incorrect, got: %s, want: %s.", result[1]["col2"], wanted)
	}
}

func TestUpdate(t *testing.T) {
	restore()
	tbl := New("fixtures", getDbConnString())
	tbl.Update(Data{Columns: []string{"col2", "col3"}, ColVals: []interface{}{"Two++", 21}, Key: "col1", KeyVal: 2})
	db, _ := openDB(getDbConnString())
	row := db.QueryRow("SELECT col2,col3 FROM fixtures WHERE col1=2")
	var str string
	var number int
	row.Scan(&str, &number)
	db.Close()
	wanted1, wanted2 := "Two++", 21
	if str != wanted1 || number != wanted2 {
		t.Errorf("Incorrect, got: %s, %d, want: %s, %d.", str, number, wanted1, wanted2)
	}
}

func TestDelete(t *testing.T) {
	restore()
	tbl := New("fixtures", getDbConnString())
	tbl.Delete(Data{Key: "col1", KeyVal: 2})
	db, _ := openDB(getDbConnString())
	row := db.QueryRow("SELECT col2 FROM fixtures WHERE col1=2")
	var str string
	row.Scan(&str)
	db.Close()
	wanted := ""
	if str != wanted {
		t.Errorf("Incorrect, got: %s, want: %s.", str, wanted)
	}
}

func TestInsert(t *testing.T) {
	restore()
	tbl := New("fixtures", getDbConnString())
	tbl.Insert(Data{ColVals: []interface{}{4, "Four", 40}})
	db, _ := openDB(getDbConnString())
	row := db.QueryRow("SELECT col2, col3 FROM fixtures WHERE col1=4")
	var str string
	var number int
	row.Scan(&str, &number)
	db.Close()
	wanted1, wanted2 := "Four", 40
	if str != wanted1 || number != wanted2 {
		t.Errorf("Incorrect, got: %s, %d, want: %s, %d.", str, number, wanted1, wanted2)
	}
}
