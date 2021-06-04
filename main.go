package main

import (
	cfg "FileStore-Server/config"
	"FileStore-Server/db/mysql"
	"FileStore-Server/route"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Printf("[GIN-debug] enter main\n")
	conf := cfg.GetConf("conf.yaml")
	mysql.Init(conf)

	router := route.Router()
	router.Run(conf.UploadServiceHost)
}
