![GitHub tag (latest by date)](https://img.shields.io/github/tag-date/jjcapellan/easy-pgql.svg)
![GitHub license](https://img.shields.io/github/license/jjcapellan/easy-pgql.svg)
# EASY-PGQL  
Very basic, simple, ligh and easy to use library to work with POSTGRESQL tables.  
Easy-pgql focuses on the manipulation of the records of an existing table.  
It offers basic functionality for those cases in which advanced features or high performance are not required, but ease of use.

## Tutorial  
For the tutorial we will use an existing table named employees:

| id | name | department    | age | salary
| --- |:--- | :--- | ---: | ---:
| 1 | Peter | manufacturing | 32 | 1200
| 2 | Paul | manufacturing | 27 | 1150
| 3 | Adam | manufacturing | 38 | 1260
| 4 | Alice | accounting | 41 | 1400
| 5 | Alex | manufacturing | 33 | 1250

### Creating a pgql.Table object
The type **Table** has methods to perform basic operations on a specific table in a database.  
All **Table** methods returns an error object for handle it in your way.
```go
package main

import (
    "fmt"

    "github.com/jjcapellan/easy-pgql"
)

// Gets the connection string to the database from the environment
func getDbConnString() string {
	host := os.Getenv("PGHOST")
	port := os.Getenv("PGPORT")
	user := os.Getenv("PGUSER")
	pass := os.Getenv("PGPASSWORD")
	name := os.Getenv("PGDATABASE")
    connStr := "host=" + host + " port=" + port + " user=" + user + 
    " password=" + pass + " dbname=" + name + " sslmode=disable"
	return connStr
}

func main(){

// Creates a pgql.Table object. It points to the table "mytable" from our database
mytable := pgql.New("employees", getDbConnString())

}
```  
### The pgql.Data type
All table functions use a **pgql.Data** structure as parameter. It is not necessary to define all fields (example: read all rows and columns, does not require any fields). This is its definition:
```go
type Data struct {
    Key       string         // Key field for select the row/rows
	KeyVal    interface{}    // Value of the key field
	Columns   []string       // Columns of the row we want select.
	ColVals   []interface{}  // Value for the columns we want modify
	OrderBy   string         // Do we want order results by some column? (defaul = primary key)
	DescOrder bool           // Do we want descent order? (default = ascent)
	Limit     int            // Limits the query to n rows
}
```  
**Note**: Column names with caps need to be introduced between scaped quotes "\\"Name\\""
### Insert
```go
mytable := pgql.New("employees", connStr)

// Inserts a full row in the table
query := pgql.Data{
    ColVals: []interface{}{1, "Peter", "manufacturing", 32, 1200},
    }
mytable.Insert(query)

// Inserts only the specified fields in the table
query = pgql.Data{
    Columns: []string{"id", "name","age"},
    ColVals: []interface{}{1, "Peter",32},
    }
mytable.Insert(query)
```  
### Delete
```go
mytable := pgql.New("employees", connStr)

// Deletes rows where department = "manufacturing"
query := pgql.Data{
    Key: "department",
    KeyVal: "manufacturing",
}
mytable.Delete(query)

// Deletes all rows
mytable.Delete(pgql.Data{})
```  
### Update  
```go
mytable := pgql.New("employees", connStr)

// Updates salary and age of Alice
query := pgql.Data{
    key: "id",                           // Its recommended use the primary key when possible
    KeyVal: 4,
    Columns: []string{"age","salary"},
    ColVals: []interface{}{42, 1500},
}
mytable.Update(query)

// Updates the first three fields of row with id = 5
query = pgql.Data{
    key: "id",
    KeyVal: 5,
    ColVals: []interface{}{10, "Alexis", "accounting"},
}
mytable.Update(query)
```  
### Read  
Read function returns an array of maps ([]map[string]interface{}). Each map represents one row.
```go
mytable := pgql.New("employees", connStr)

// Reads all
result, _ := mytable.Read(Data{})

// Reads first 3 results order by name in descending order
query := pgql.Data{
    OrderBy: "name",
    DescOrder: true,
    Limit: 3,
}
result, _ = mytable.Read(query)

// Reads only salary column
query = pgql.Data{
    Columns: []string{"salary"},
}
result, _ = mytable.Read(query)
```  

### Get position of selected row
The function GetPos uses 4 Data properties: *Key*, *KeyVal*, *OrderBy* and optionally *DescOrder*.
```go
mytable := pgql.New("employees", connStr)

// Returns the position (int64) of Alex by his salary
query := pgql.Data{
    Key: "name",
    KeyVal: "Alex",
    OrderBy: "salary",
    DescOrder: true,     // Bigest salary --> first position
}
result, _ := mytable.GetPos(pgql.Data)

// result == 3

``` 

## Test
Testing uses a table named fixtures with 3 columns (col1, col2, col3) of types integer, varchar(40), and integer.