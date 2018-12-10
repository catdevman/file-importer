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
)

var (
	filePath string
)

func main() {
	flag.StringVar(&filePath, "file", "staff.csv", "Choose path for file")
	flag.Parse()
	var manager types.Manager
	db, err := getMysqlDatabaseConnection("root", "password", "127.0.0.1", "pivot2_apigility")
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

	if ok {
		_ = validator.ValidateCollection()
	}

	s, _ := json.Marshal(manager.ShowData())
	if s == nil {
		fmt.Printf("%s\n", s)
	}

}

func getMysqlDatabaseConnection(dbUsername, dbPassword, dbHost, database string) (*sql.DB, error) {
	return sql.Open("mysql", buildDsnString(dbUsername, dbPassword, dbHost, database))
}

func buildDsnString(user string, password, host string, database string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?multiStatements=true", user, password, host, database)
}
