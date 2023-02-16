package db

import (
	"fmt"
	"errors"
	"os"
	"io"
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	Database *sqlx.DB
}

func CreateDB(conn string) (*DB, error){
	database, err := sqlx.Open("postgres", conn)
	if err != nil {
		return nil, errors.New("unable to connect to database")
	}
	fmt.Println("SUCCESSFULLY CONNECTED TO DB")
	newDB := &DB{database}
	return newDB, nil
}


func (DB *DB) QueryFromFile(path string) (string, error) {
	fi, err := os.Open(path)
	if err != nil {
		return "", errors.New("could not open given file")
	}

	query, err := io.ReadAll(fi)
	if err != nil {
		return "", errors.New("unable to read contents of given file")
	}
	return string(query), nil
}


// Exec is a proxy for the sql.DB method Exec
func (DB *DB) Exec(query string) (sql.Result, error) {
	res, err := DB.database.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("unable to execute given queries")
	}
	return res, nil
}

// QueryRow selects a single row and tranforms it into a struct of the given type
func (DB *DB) QueryRow(query string, dest interface{}) (interface{}, error) {
	err := DB.database.QueryRowx(query).StructScan(dest)
	if err != nil {
		return nil, errors.New("could not query row with given command")
	}
	return dest, nil
}

// Query selects multiple rows and transforms it into an slice of structs of the given type
func (DB *DB) Query(query string, queryType interface{}) (interface{}, error) {
	var objects *[]interface{}
	err := DB.database.QueryRowx(query).StructScan(objects)
	if err != nil {
		return nil, err
	}
	return objects, nil
}
