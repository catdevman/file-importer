package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/goFileImporter/file-importer/staff10col"
	"github.com/goFileImporter/file-importer/types"
	"log"
	"os"
	// "runtime/trace"
)

var (
	filePath, dbUsername, dbPassword, dbHost, dbName string
)

func main() {
	flag.StringVar(&filePath, "file", "staff.csv", "Choose path for file")
	flag.StringVar(&dbUsername, "dbUsername", os.Getenv("DB_USERNAME"), "Database Username")
	flag.StringVar(&dbPassword, "dbPassword", os.Getenv("DB_PASSWORD"), "Database Password")
	flag.StringVar(&dbHost, "dbHost", os.Getenv("DB_HOST"), "Database Host")
	flag.StringVar(&dbName, "dbName", os.Getenv("DB_NAME"), "Database Name")
	flag.Parse()
	var manager types.Manager
	db, err := getMysqlDatabaseConnection(dbUsername, dbPassword, dbHost, dbName)
	defer db.Close()
	if err != nil {
		// This is probably bad
	}
	manager = staff10col.NewStaffManager(db)
	file, err := os.Open(filePath)
	if err != nil {
		// This is probably bad
	}
	_, err = manager.LoadDataFromReader(bufio.NewReader(file))

	if err != nil {
		log.Fatal(err)
	}

	validator, ok := manager.(types.ManagerValidator)
	var errs []types.ErroredRecord
	if ok {
		errs = validator.ValidateCollection()
	}

	s, _ := json.Marshal(manager.Data())
	if s == nil {
		fmt.Printf("%s\n", s)
	}
	if errs != nil {
		fmt.Println(errs)
	}
}

func getMysqlDatabaseConnection(dbUsername, dbPassword, dbHost, database string) (*sql.DB, error) {
	return sql.Open("mysql", buildDsnString(dbUsername, dbPassword, dbHost, database))
}

func buildDsnString(user string, password, host string, database string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?multiStatements=true", user, password, host, database)
}
