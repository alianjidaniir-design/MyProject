package main

import (
	mysqlDataSource "MyProject/models/user/dataSourses/mySqlDS"
	"flag"
	"fmt"
	"log"
)

func main() {

	envCFG, err := mysqlDataSource.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	dsn := flag.String("dsn", envCFG.DSN, "MySQL DSN")
	table := flag.String("table", envCFG.StudentTableName, "user table name")
	flag.Parse()

	if *dsn == "" {
		fmt.Println("Error: dsn is empty")
	}

	cfg := envCFG
	cfg.DSN = *dsn
	cfg.StudentTableName = *table

	if err := mysqlDataSource.ValidateTableName(cfg.StudentTableName); err != nil {
		log.Fatal("Error validating table name:", err)
	}

}
