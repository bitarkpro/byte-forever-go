package GinHandler

import (
	"FileStore-Server/common"
	"FileStore-Server/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	//"encoding/json"
)

func AddStarFileHandler(c *gin.Context) {
	token := c.Request.Header.Get("token")
	phone := IsTokenExist(token)
	if phone == "" {
		c.JSON(http.StatusOK, gin.H{
			"msg":  common.UserNoexist,
			"code": common.StatusFailed,
		})
		return
	}

	id, err := strconv.Atoi(c.Request.FormValue("id"))
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return
	}
	userfile, _ := db.GetFileInfo(id)
	//fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter DownloadHandler:%s\n", userfile)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR]1111111 %v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"msg":  common.ErrResult,
			"code": common.StatusFailed,
		})
		return
	}
	if userfile == nil {
		c.JSON(http.StatusOK, gin.H{
			"msg":  common.FileNotExist,
			"code": common.StatusFailed,
		})
		return
	}
	starfileOk := db.UpdateStarDb(id, common.AddStarFile)
	if !starfileOk {
		c.JSON(http.StatusOK, gin.H{
			"msg":  common.AddStarfileFail,
			"code": common.StatusFailed,
		})
		return
	}
	//update user info
	fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR]111111123333 %v\n", err)
	userStorage, _ := db.QueryUserStroage(phone)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR]111111122 %v\n", err)
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return
	}

	userStorage.StarfileNum += 1
	updateUserStorage := db.UpdateUserStroage(phone, userStorage)
	if !updateUserStorage {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  common.AddStarfileSuccess,
		"code": common.StatusSuccess,
	})

	return
}

func CancelStarFileHandler(c *gin.Context) {
	token := c.Request.Header.Get("token")
	phone := IsTokenExist(token)
	if phone == "" {
		c.JSON(http.StatusOK, gin.H{
			"msg":  common.UserNoexist,
			"code": common.StatusFailed,
		})
		return
	}

	id, err := strconv.Atoi(c.Request.FormValue("id"))
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return
	}
	userfile, err := db.GetFileInfo(id)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"msg":  common.ErrResult,
			"code": common.StatusFailed,
		})
		return
	}
	if userfile == nil {
		c.JSON(http.StatusOK, gin.H{
			"msg":  common.FileNotExist,
			"code": common.StatusFailed,
		})
		return
	}
	starfileOk := db.UpdateStarDb(id, common.CancelStarFile)
	if !starfileOk {
		c.JSON(http.StatusOK, gin.H{
			"msg":  common.CancelStarfileFail,
			"code": common.StatusFailed,
		})
		return
	}
	//update user info
	userStorage, _ := db.QueryUserStroage(phone)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return
	}

	userStorage.StarfileNum -= 1
	updateUserStorage := db.UpdateUserStroage(phone, userStorage)
	if !updateUserStorage {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  common.CancelStarfileSuccess,
		"code": common.StatusSuccess,
	})

	return
}
