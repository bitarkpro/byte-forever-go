package main

import (
	"FileStore-Server/config"
	"FileStore-Server/db/mysql"
	"FileStore-Server/route"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	mysql.Init()
	fmt.Printf("[GIN-debug] enter main\n")

	router := route.Router()
	router.Run(config.UploadServiceHost)
}
