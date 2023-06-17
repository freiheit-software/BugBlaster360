package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func Connect(connectionString string) (*DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to the database")
	return &DB{conn: db}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) InsertData(tableName string, data map[string]interface{}) error {
	columns := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))

	for col, val := range data {
		columns = append(columns, col)
		values = append(values, val)
	}

	columnsStr := strings.Join(columns, ", ")
	valuePlaceholders := make([]string, len(values))
	for i := range values {
		valuePlaceholders[i] = fmt.Sprintf("$%d", i+1)
	}

	insertStatement := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columnsStr, strings.Join(valuePlaceholders, ", "))

	_, err := db.conn.Exec(insertStatement, values...)
	if err != nil {
		fmt.Println("Error While Inserting", err)
		return err
	}

	fmt.Println("Data inserted successfully")
	return nil
}

func (db *DB) GetData(tableName string) ([]map[string]interface{}, error) {
	selectStatement := fmt.Sprintf("SELECT * FROM %s", tableName)

	rows, err := db.conn.Query(selectStatement)
	if err != nil {
		fmt.Println("Error while querying data:", err)
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Println("Error while retrieving column information:", err)
		return nil, err
	}

	result := make([]map[string]interface{}, 0)

	for rows.Next() {
		columnPointers := make([]interface{}, len(columns))
		columnValues := make([]interface{}, len(columns))

		for i := range columns {
			columnPointers[i] = &columnValues[i]
		}

		err := rows.Scan(columnPointers...)
		if err != nil {
			fmt.Println("Error while scanning row values:", err)
			return nil, err
		}

		rowData := make(map[string]interface{})
		for i, colName := range columns {
			val := columnValues[i]
			if val != nil {
				rowData[colName] = val
			} else {
				rowData[colName] = nil
			}
		}

		result = append(result, rowData)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error while iterating over rows:", err)
		return nil, err
	}

	return result, nil
}

func (db *DB) GetDataByField(tableName string, fieldName string, fieldValue interface{}) ([]map[string]interface{}, error) {
	selectStatement := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", tableName, fieldName)

	rows, err := db.conn.Query(selectStatement, fieldValue)
	if err != nil {
		fmt.Println("Error while querying data:", err)
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Println("Error while retrieving column information:", err)
		return nil, err
	}

	result := make([]map[string]interface{}, 0)

	for rows.Next() {
		columnValues := make([]interface{}, len(columns))
		columnPointers := make([]interface{}, len(columns))
		for i := range columnValues {
			columnPointers[i] = &columnValues[i]
		}

		err := rows.Scan(columnPointers...)
		if err != nil {
			fmt.Println("Error while scanning row values:", err)
			return nil, err
		}

		rowData := make(map[string]interface{})
		for i, colName := range columns {
			val := columnValues[i]
			rowData[colName] = val
		}

		result = append(result, rowData)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error while iterating over rows:", err)
		return nil, err
	}

	return result, nil
}
