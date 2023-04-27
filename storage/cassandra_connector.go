package storage

import (
	"fmt"
	"github.com/gocql/gocql"
)

type CassandraConnector struct {
	session *gocql.Session
}

func (connector *CassandraConnector) Connect(user, password, host, dbname string) error {
	cluster := gocql.NewCluster(host)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: user,
		Password: password,
	}
	cluster.Keyspace = dbname
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}
	connector.session = session
	return nil
}

func (connector *CassandraConnector) Disconnect() error {
	connector.session.Close()
	return nil
}

func (connector *CassandraConnector) CreateTable(table Table) error {
	// Create column definitions for CQL
	columns := ""
	for name, t := range table.Columns {
		columns += fmt.Sprintf("%s %s,", name, t)
	}
	columns = columns[:len(columns)-1] // Remove trailing comma

	// Create CQL command
	cmd := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s, PRIMARY KEY (%s))", table.Name, columns, table.PrimaryKey)

	// Execute CQL command
	if err := connector.session.Query(cmd).Exec(); err != nil {
		return err
	}
	return nil
}

func (connector *CassandraConnector) DeleteTable(table Table) error {
	cmd := fmt.Sprintf("DROP TABLE IF EXISTS %s", table.Name)
	if err := connector.session.Query(cmd).Exec(); err != nil {
		return err
	}
	return nil
}

func (connector *CassandraConnector) DeleteLine(line Line) error {
	cmd := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", line.Table, line.PrimaryKey)
	if err := connector.session.Query(cmd, line.Line[line.PrimaryKey]).Exec(); err != nil {
		return err
	}
	return nil
}

func (connector *CassandraConnector) InsertLine(line Line) error {
	keys := ""
	values := ""
	for key, _ := range line.Line {
		keys += key + ","
		values += "?,"
	}
	keys = keys[:len(keys)-1]       // Remove trailing comma
	values = values[:len(values)-1] // Remove trailing comma

	cmd := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", line.Table, keys, values)
	if err := connector.session.Query(cmd, line.Line).Exec(); err != nil {
		return err
	}
	return nil
}

func (connector *CassandraConnector) UpdateLine(line Line) error {
	updates := ""
	values := make([]interface{}, 0, len(line.Line))

	for key, value := range line.Line {
		updates += fmt.Sprintf("%s = ?,", key)
		values = append(values, value)
	}

	updates = updates[:len(updates)-1] // Remove trailing comma

	cmd := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?", line.Table, updates, line.PrimaryKey)
	values = append(values, line.Line[line.PrimaryKey])

	if err := connector.session.Query(cmd, values...).Exec(); err != nil {
		return err
	}
	return nil
}
