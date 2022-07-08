package db

import (
	"database/sql"
	"fmt"
)

const (
	host     string = "localhost"
	port     int    = 6603
	username string = "systemuser"
	password string = "password"
	database string = "ERP"
)

func InitDB() {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database)
	db, err := sql.Open("mysql", connString)
	// db, err := sql.Open("mysql", "root:6HmFbaZH2X3Z0jXjruqD@tcp(localhost:6603)/ERP")

	if err != nil {
		// Error: cant connect to DB
		fmt.Println(err.Error())
	} else if err = db.Ping(); err != nil {
		// Error: lost connect to DB
		fmt.Println(err.Error())
	}
	defer db.Close()
	var createUsersQuers string = `CREATE TABLE IF NOT EXISTS users (` +
		`RecordId int AUTO_INCREMENT PRIMARY KEY NOT NULL,` +
		`Id varchar(255) NOT NULL,` +
		`Firstname varchar(255) NOT NULL,` +
		`Lastname varchar(255) NOT NULL,` +
		`Username varchar(255) NOT NULL,` +
		`Email varchar(255) NOT NULL,` +
		`Password varchar(255) NOT NULL,` +
		`Role varchar(255) NOT NULL` +
		`);`
	res, err := db.Exec(createUsersQuers)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(res)
	}
	defer db.Close()
}

func RunSqlQueryWithReturn(query string) (*sql.Rows, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database)
	db, err := sql.Open("mysql", connString)

	if err != nil {
		// Error: cant connect to DB
		return nil, err
	} else if err = db.Ping(); err != nil {
		// Error: lost connect to DB
		return nil, err
	}
	defer db.Close()
	res, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	defer db.Close()
	return res, nil
}

func RunSqlQueryWithSingeReturn(query string) (*sql.Rows, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database)
	db, err := sql.Open("mysql", connString)

	if err != nil {
		// Error: cant connect to DB
		return nil, err
	} else if err = db.Ping(); err != nil {
		// Error: lost connect to DB
		return nil, err
	}
	defer db.Close()
	res, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	defer db.Close()
	return res, nil
}

func RunSqlQueryWithoutReturn(query string) (bool, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database)
	db, err := sql.Open("mysql", connString)

	if err != nil {
		// Error: cant connect to DB
		return false, err
	} else if err = db.Ping(); err != nil {
		// Error: lost connect to DB
		return false, err
	}
	defer db.Close()
	res, err := db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}
	fmt.Println(res)
	defer db.Close()
	return true, nil
}
