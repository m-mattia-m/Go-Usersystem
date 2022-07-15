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

func CreateUsersTable() {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database)
	db, err := sql.Open("mysql", connString)

	if err != nil {
		fmt.Println("[DB]: Error - Can't connect to the DB \t-->\t" + err.Error())
	} else if err = db.Ping(); err != nil {
		fmt.Println("[DB]: Error - Lost connection to the DB \t-->\t" + err.Error())
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
		`Role varchar(255) NOT NULL,` +
		`Token text NULL` +
		`);`
	res, err := db.Exec(createUsersQuers)
	if err != nil {
		fmt.Println("[DB]: Error - Can't create table users \t-->\t" + err.Error())
	} else {
		fmt.Println("[DB]: Table users successfully created")
		_ = res
	}
	defer db.Close()
}

func RunSqlQueryWithReturn(query string) (*sql.Rows, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database)
	db, err := sql.Open("mysql", connString)

	if err != nil {
		fmt.Println("[DB]: Error - Can't connect to the DB \t-->\t" + err.Error())
		return nil, err
	} else if err = db.Ping(); err != nil {
		fmt.Println("[DB]: Error - Lost connection to the DB \t-->\t" + err.Error())
		return nil, err
	}
	defer db.Close()
	res, err := db.Query(query)
	if err != nil {
		fmt.Println("[DB]: Error - Can't run SQL-Query with Return \t-->\t" + err.Error())
		fmt.Println("[DB]: Query: " + query)
		return nil, err
	}

	defer db.Close()
	return res, nil
}

func RunSqlQueryWithSingeReturn(query string) (*sql.Rows, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database)
	db, err := sql.Open("mysql", connString)

	if err != nil {
		fmt.Println("[DB]: Error - Can't connect to the DB \t-->\t" + err.Error())
		return nil, err
	} else if err = db.Ping(); err != nil {
		fmt.Println("[DB]: Error - Lost connection to the DB \t-->\t" + err.Error())
		return nil, err
	}

	defer db.Close()
	res, err := db.Query(query)
	if err != nil {
		fmt.Println("[DB]: Error - Can't run SQL-Query with Single Return \t-->\t" + err.Error())
		fmt.Println("[DB]: Query: " + query)
		return nil, err
	}

	defer db.Close()
	return res, nil
}

func RunSqlQueryWithoutReturn(query string) (bool, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database)
	db, err := sql.Open("mysql", connString)

	if err != nil {
		fmt.Println("[DB]: Error - Can't connect to the DB \t-->\t" + err.Error())
		return false, err
	} else if err = db.Ping(); err != nil {
		fmt.Println("[DB]: Error - Lost connection to the DB \t-->\t" + err.Error())
		return false, err
	}
	defer db.Close()
	res, err := db.Exec(query)
	if err != nil {
		fmt.Println("[DB]: Error - Can't run SQL-Query without Return \t-->\t" + err.Error())
		fmt.Println("[DB]: Query: " + query)
		return false, err
	}
	_ = res
	defer db.Close()
	return true, nil
}
