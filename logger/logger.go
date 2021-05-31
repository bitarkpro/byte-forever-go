package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func OpenLog(logRouter *gin.Engine) {

	gin.DisableConsoleColor()

	f, _ := os.Create("gin.log")


	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	gin.DefaultErrorWriter = io.MultiWriter(f, os.Stderr)

	logRouter.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	fmt.Fprintln(gin.DefaultWriter, "[GIN-debug] log system open")
	var err error
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)}

}
