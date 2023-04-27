package storage

import (
	"database/sql"
	"fmt"
)

type MySQLConnector struct {
	db *sql.DB
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

func (connector *MySQLConnector) CreateTable(table Table) error {
	// Create column definitions for SQL
	columns := ""
	for name, t := range table.Columns {
		columns += fmt.Sprintf("%s %s,", name, t)
	}
	columns = columns[:len(columns)-1] // Remove trailing comma

	// Create SQL command
	cmd := fmt.Sprintf("CREATE TABLE %s (%s, PRIMARY KEY (%s))", table.Name, columns, table.PrimaryKey)

	// Execute SQL command
	if _, err := connector.db.Exec(cmd); err != nil {
		return err
	}
	return nil
}

func (connector *MySQLConnector) DeleteTable(table Table) error {
	cmd := fmt.Sprintf("DROP TABLE %s", table.Name)
	if _, err := connector.db.Exec(cmd); err != nil {
		return err
	}
	return nil
}

func (connector *MySQLConnector) DeleteLine(line Line) error {
	cmd := fmt.Sprintf("DELETE FROM %s WHERE %s = '%s'", line.Table, line.PrimaryKey, line.Line[line.PrimaryKey])
	if _, err := connector.db.Exec(cmd); err != nil {
		return err
	}
	return nil
}

func (connector *MySQLConnector) InsertLine(line Line) error {
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

func (connector *MySQLConnector) UpdateLine(line Line) error {
	updates := ""
	for key, value := range line.Line {
		updates += fmt.Sprintf("%s = '%s',", key, value)
	}
	updates = updates[:len(updates)-1] // Remove trailing comma

	cmd := fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%s'", line.Table, updates, line.PrimaryKey, line.Line[line.PrimaryKey])
	if _, err := connector.db.Exec(cmd); err != nil {
		return err
	}
	return nil
}
