package mysql

import (
	cfg "FileStore-Server/config"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var Db *sql.DB

const (
	MaxOpenConns = 10000
)

func Init(conf cfg.Config) {
	fmt.Fprintln(gin.DefaultWriter, "[GIN-debug] access db init")
	Db, _ = sql.Open(conf.DriverName, conf.DataSourceName)

	if Db == nil {
		fmt.Fprintln(gin.DefaultWriter, "[GIN-debug] db is nil")
		return
	}
	Db.SetMaxOpenConns(MaxOpenConns)

	err := Db.Ping()
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR]Failed to connect to mysql, err:  %v\n", err)
		os.Exit(1)
	}
}

// DBConn: Returns the database connection object
func DBConn() *sql.DB {
	fmt.Fprintln(gin.DefaultWriter, "[GIN-debug] access DBConn")
	return Db
}

func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {

		err := rows.Scan(scanArgs...)
		checkErr(err)

		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
