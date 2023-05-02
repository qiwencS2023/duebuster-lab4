package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type MySQLConnector struct {
	db *sql.DB
}

func (c *MySQLConnector) GetLine(request *GetLineRequest) (*Line, error) {
	table := request.Table
	pk := request.PrimaryKeyValue

	// Create SQL command
	cmd := fmt.Sprintf("SELECT * FROM %s WHERE %s = %s", table.Name, table.PrimaryKey, pk)
	log.Printf("[storage] cmd: %s\n", cmd)

	// Execute SQL command
	rows, err := c.db.Query(cmd)
	if err != nil {
		return nil, err
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Get column types
	types, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	log.Printf("columns: %v\n", columns)
	for rows.Next() {
		// Create column values
		columnValues := make(map[string]string)

		// Create column pointers
		columnPointers := make([]interface{}, len(columns))
		for i, _ := range columns {
			columnPointers[i] = new(interface{})
		}

		// Scan columns into pointers
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		// Put column values into map
		for i, column := range columns {
			columnValue := columnPointers[i].(*interface{})
			columnType := types[i].DatabaseTypeName()
			log.Printf("column: %s, type: %s\n", column, columnType)
			switch columnType {
			case "VARCHAR", "TEXT":
				columnValues[column] = string((*columnValue).([]uint8))
			case "INT":
				columnValues[column] = string((*columnValue).([]uint8))
			case "FLOAT":
				columnValues[column] = string((*columnValue).([]uint8))
			default:
				return nil, fmt.Errorf("unknown column type")
			}
		}

		// Create line
		line := &Line{
			Table:      table.Name,
			PrimaryKey: table.PrimaryKey,
			Line:       columnValues,
		}

		return line, nil
	}

	return nil, fmt.Errorf("no rows found")
}

func (connector *MySQLConnector) Connect(user, password, host, dbname string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	connector.db = db
	return nil
}

func (connector *MySQLConnector) Disconnect() error {
	if err := connector.db.Close(); err != nil {
		return err
	}
	return nil
}

func (connector *MySQLConnector) CreateTable(table *Table) error {
	// Create column definitions for SQL
	columns := ""
	for name, t := range table.Columns {
		columns += fmt.Sprintf("%s %s,", name, t)
	}
	columns = columns[:len(columns)-1] // Remove trailing comma

	// Create SQL command
	cmd := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s, PRIMARY KEY (%s))", table.Name, columns, table.PrimaryKey)

	fmt.Printf("[mysql_connector.go] cmd: %s\n", cmd)

	// Execute SQL command
	if _, err := connector.db.Exec(cmd); err != nil {
		return err
	}
	return nil
}

func (connector *MySQLConnector) DeleteTable(table *Table) error {
	cmd := fmt.Sprintf("DROP TABLE %s", table.Name)
	if _, err := connector.db.Exec(cmd); err != nil {
		return err
	}
	return nil
}

func (connector *MySQLConnector) DeleteLine(line *Line) error {
	cmd := fmt.Sprintf("DELETE FROM %s WHERE %s = '%s'", line.Table, line.PrimaryKey, line.Line[line.PrimaryKey])
	if _, err := connector.db.Exec(cmd); err != nil {
		return err
	}
	return nil
}

func (connector *MySQLConnector) InsertLine(line *Line) error {
	keys := ""
	values := ""
	for key, value := range line.Line {
		keys += key + ","
		values += fmt.Sprintf("'%s',", value)
	}
	keys = keys[:len(keys)-1]       // Remove trailing comma
	values = values[:len(values)-1] // Remove trailing comma

	cmd := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", line.Table, keys, values)
	if _, err := connector.db.Exec(cmd); err != nil {
		return err
	}
	return nil
}

func (connector *MySQLConnector) UpdateLine(line *Line) error {
	updates := ""
	for key, value := range line.Line {
		updates += fmt.Sprintf("%s = '%s',", key, value)
	}
	updates = updates[:len(updates)-1] // Remove trailing comma

	cmd := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%s'", line.Table, updates, line.PrimaryKey, line.Line[line.PrimaryKey])
	log.Printf("[mysql_connector.go] cmd: %s\n", cmd)
	if _, err := connector.db.Exec(cmd); err != nil {
		return err
	}
	return nil
}
