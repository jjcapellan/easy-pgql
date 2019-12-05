package pgql

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/lib/pq" //needed
)

///////////////////////
// Types

// Table type stores needed parameters for open a table in database
type Table struct {
	Name   string
	Config string
}

// Data used in query
type Data struct {
	Key       string
	KeyVal    interface{}
	Columns   []string
	ColVals   []interface{}
	OrderBy   string
	DescOrder bool
	Limit     int
}

///////////////////////////
// Main functions

// New creates a table
func New(name, dbconfig string) Table {
	tbl := Table{Name: name, Config: dbconfig}
	return tbl
}

func openDB(config string) (*sql.DB, error) {
	db, err := sql.Open("postgres", config)
	if err != nil {
		return db, err
	}
	err = db.Ping()
	if err != nil {
		return db, err
	}
	return db, err
}

// Insert inserts one row in a table
func (t Table) Insert(d Data) error {
	db, err := openDB(t.Config)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(t.getInsertStr(d), d.ColVals...)
	if err != nil {
		return err
	}
	return err
}

// Update updates one row in a table
func (t Table) Update(d Data) error {
	db, err := openDB(t.Config)
	if err != nil {
		return err
	}
	defer db.Close()
	d.ColVals = append(d.ColVals, d.KeyVal)
	_, err = db.Exec(t.getUpdateStr(d), d.ColVals...)
	if err != nil {
		return err
	}
	return err
}

// Delete deletes a row in a table
func (t Table) Delete(d Data) error {
	db, err := openDB(t.Config)
	if err != nil {
		return err
	}
	defer db.Close()
	sqlQuery := fmt.Sprintf(`DELETE FROM %s WHERE %s=$1`, t.Name, d.Key)
	_, err = db.Exec(sqlQuery, d.KeyVal)
	if err != nil {
		return err
	}
	return err
}

// Read reads one or more rows in a table
func (t Table) Read(d Data) ([]map[string]interface{}, error) {
	db, err := openDB(t.Config)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var rows *sql.Rows
	if d.KeyVal != nil {
		rows, _ = db.Query(t.getReadStr(d), d.KeyVal)
	} else {
		rows, _ = db.Query(t.getReadStr(d))
	}
	result, err := rowsToJSON(rows)
	if err != nil {
		return result, err
	}
	return result, err
}

// GetPos returns the position (int64) of a row with the pair Key, KeyVal
func (t Table) GetPos(d Data) (int64, error) {
	db, err := openDB(t.Config)
	if err != nil {
		return 0, err
	}
	defer db.Close()
	direction := ""
	if d.DescOrder {
		direction = "DESC"
	}
	sqlQuery := fmt.Sprintf(`SELECT ROW_NUMBER FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s %s), %s FROM %s) x WHERE %s=$1`, d.OrderBy, direction, d.Key, t.Name, d.Key)
	row := db.QueryRow(sqlQuery, d.KeyVal)
	var position int64
	err = row.Scan(&position)
	if err != nil {
		return 0, err
	}
	return position, nil
}

//////////////////////////////
// sql strings helpers

func getParamsStr(n int) string {
	str := ""
	for i := 1; i < n+1; i++ {
		str += "$" + strconv.Itoa(i) + ","
	}
	str = strings.TrimRight(str, ",")
	return str
}

func (t Table) getReadStr(d Data) string {
	query := "SELECT "
	if d.Columns != nil {
		query += strings.Join(d.Columns, ",")
	} else {
		query += "*"
	}
	query += " FROM " + t.Name
	if d.Key != "" {
		query += " WHERE " + d.Key + "=$1"
	}
	if d.OrderBy != "" {
		query += " ORDER BY " + d.OrderBy
	}
	if d.DescOrder {
		query += " DESC"
	}
	if d.Limit != 0 {
		query += " LIMIT " + strconv.Itoa(d.Limit)
	}
	return query
}

func (t Table) getUpdateStr(d Data) string {
	query, sets, count := "UPDATE "+t.Name+" SET ", "", 1
	for i := 0; i < len(d.Columns); i++ {
		sets += "%s=$" + strconv.Itoa(count) + ","
		count++
	}
	sets = strings.TrimRight(sets, ",")
	query += sets + " WHERE %s=$" + strconv.Itoa(count)
	d.Columns = append(d.Columns, d.Key)
	query = fmt.Sprintf(query, arrToInterface(d.Columns)...)

	return query
}

func (t Table) getInsertStr(d Data) string {
	query := "INSERT INTO " + t.Name
	if d.Columns != nil {
		query += " (" + strings.Join(d.Columns, ",") + ")"
	}
	query += " VALUES (%s)"
	query = fmt.Sprintf(query, getParamsStr(len(d.ColVals)))
	return query
}

////////////////////////////
// Conversion data helpers

func arrToInterface(arr []string) []interface{} {
	intf := make([]interface{}, len(arr))
	for i := range arr {
		intf[i] = arr[i]
	}
	return intf
}

func rowsToJSON(rows *sql.Rows) ([]map[string]interface{}, error) {
	colNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	numCols := len(colNames)
	values := make([]interface{}, numCols)
	valuesPtr := make([]interface{}, numCols)
	container := make([]map[string]interface{}, 0)

	for rows.Next() {
		rowMap := make(map[string]interface{}, 0)
		for i := range values {
			valuesPtr[i] = &values[i]
		}
		rows.Scan(valuesPtr...)
		for i, col := range colNames {
			rowMap[col] = values[i]
		}
		container = append(container, rowMap)
	}
	return container, err
}
