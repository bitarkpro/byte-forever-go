package db

import (
	mydb "FileStore-Server/db/mysql"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

//UserFile:
type UserFile struct {
	Id         int    `json:"id"`
	Phone      string `json:"phone"`
	Cid        string `json:"cid"`
	FileName   string `json:"fileName"`
	Ext        string `json:"ext"`
	FileSize   int64  `json:"fileSize"`
	FileType   string `json:"fileType"`
	CreatAt    string `json:"creatTime"`
	Comment    string `json:"comment"`
	FolderType int    `json:"folderType"`
	Star       int    `json:"star"`
}

func OnUserFileUploadFinished(phone string, cid string, fileName string, ext string, fileType string,
	fileSize int64, comment string, folderType int, tags string, minioUrl string) bool {
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter OnUserFileUploadFinished\n")
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user_file (`phone`,`cid`,`file_name`,`ext`,`file_type`,`file_size`," +
			"`created_at`,`comment`,`folder_type`,`tags`,`minio_url`) values (?,?,?,?,?,?,?,?,?,?,?)")

	if err != nil {

		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	ret, _ := stmt.Exec(phone, cid, fileName, ext, fileType, fileSize,
		time.Now().Unix(), comment, folderType, tags ,minioUrl)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("tbl_user_file has been update. ")
		}
		return true
	}
	return true
}
