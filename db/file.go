package db

import (
	"FileStore-Server/config"
	mydb "FileStore-Server/db/mysql"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"path"
	"time"
)

// TableFile :
type TableFile struct {
	Id int `json:"id"`
	//Phone         string   `json:"phone"`
	Cid        string `json:"cid"`
	FileName   string `json:"fileName"`
	Ext        string `json:"ext"`
	FileSize   int64  `json:"fileSize"`
	FileType   string `json:"fileType"`
	CreatAt    int64  `json:"creatTime"`
	Comment    string `json:"comment"`
	FolderType int    `json:"folderType"`
	Star       int    `json:"star"`
	Tags       string `json:"tags"`
	MinioUrl   string `json:"minioUrl"`
}
type FileData struct {
	Comment string
	Tags    string
}


func GetFileMeta(phone string, cid string) (*TableFile, error) {
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter GetFileMeta\n")
	//,comment,star,folder_type
	stmt, err := mydb.DBConn().Prepare(
		"select id,cid,file_name,ext,file_type,file_size,created_at,comment,folder_type from tbl_user_file where phone=? and cid=?")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()
	tfile := TableFile{}
	err = stmt.QueryRow(phone, cid).Scan(&tfile.Id, &tfile.Cid, &tfile.FileName, &tfile.Ext, &tfile.FileType, &tfile.FileSize, &tfile.CreatAt, &tfile.Comment, &tfile.FolderType)
	if err != nil {
		if err == sql.ErrNoRows {
			// 查不到对应记录， 返回参数及错误均为nil
			return nil, nil
		} else {
			fmt.Println(err.Error())
			return nil, err
		}
	}
	return &tfile, nil
}

func GetFileInfo(id int) (*TableFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select id,cid,file_name,ext,file_type,file_size,created_at,comment,folder_type,star,tags,minio_url from tbl_user_file where id=?")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	tfile := TableFile{}
	err = stmt.QueryRow(id).Scan(&tfile.Id, &tfile.Cid, &tfile.FileName, &tfile.Ext, &tfile.FileType, &tfile.FileSize,
		&tfile.CreatAt, &tfile.Comment, &tfile.FolderType, &tfile.Star, &tfile.Tags, &tfile.MinioUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			// 查不到对应记录， 返回参数及错误均为nil
			return nil, nil
		} else {
			fmt.Println(err.Error())
			return nil, err
		}
	}

	return &tfile, nil
}

func GetFileExt(fileName string) string {
	fileSuffix := path.Ext(fileName)
	fileExt := fileSuffix[1:]
	return fileExt
}

func UpdateStarDb(id int, star int) bool {
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter UpdateStarDb star:%d\n", star)
	stmt, err := mydb.DBConn().Prepare(
		"update tbl_user_file set star=? where id=?")

	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	ret, _ := stmt.Exec(star, id)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if ret == nil {
		fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter AddStar ret in nil\n")
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter AddStar\n")
		return true
	}
	return false
}

func ViewStarFile(phone string) []TableFile {
	stmt, err := mydb.DBConn().Prepare(
		"select id,cid,file_name,ext,file_type,file_size,created_at,comment,folder_type,star,tags,minio_url from tbl_user_file where phone=? and star=1 order by created_at desc")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer stmt.Close()

	var tfile TableFile
	var i = 0
	var starfle [10]TableFile
	results, _ := stmt.Query(phone)
	for results.Next() {
		err = results.Scan(&tfile.Id, &tfile.Cid, &tfile.FileName, &tfile.Ext, &tfile.FileType, &tfile.FileSize,
			&tfile.CreatAt, &tfile.Comment, &tfile.FolderType, &tfile.Star, &tfile.Tags,&tfile.MinioUrl)
		if err != nil {
			panic(err.Error())
		}
		starfle[i] = tfile
		i++
		if i > 10-1 {
			fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] ViewStarFile break\n")
			break
		}
	}
	// Intercept valid data
	newstarfle := starfle[:i]
	fmt.Fprintf(gin.DefaultWriter, "ViewStarFile[GIN-debug] %s\n", newstarfle)
	return newstarfle
}

// query data by folder type
func GetFolderData(phone string, foldertype int) ([]TableFile, error) {
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter GetFolderData\n")
	stmt, err := mydb.DBConn().Prepare(
		"select id,cid,file_name,ext,file_type,file_size,created_at,comment,folder_type,star,tags,minio_url from tbl_user_file where phone=? and folder_type=? order by created_at desc")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	var tfile TableFile
	var i = 0
	var homePageFile [10]TableFile
	results, _ := stmt.Query(phone, foldertype)
	for results.Next() {
		err = results.Scan(&tfile.Id, &tfile.Cid, &tfile.FileName, &tfile.Ext, &tfile.FileType, &tfile.FileSize,
			&tfile.CreatAt, &tfile.Comment, &tfile.FolderType, &tfile.Star, &tfile.Tags,&tfile.MinioUrl)
		if err != nil {
			panic(err.Error())
		}
		homePageFile[i] = tfile
		i++
		if i > 10-1 {
			fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] GetFolderData break\n")
			break
		}
	}

	newhomePageFile := homePageFile[:i]
	return newhomePageFile, nil
}

func RemoveFile(id int) bool {
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter RemoveFile")
	stmt, err := mydb.DBConn().Prepare(
		"delete from tbl_user_file where id=?")
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return false
	}
	defer stmt.Close()

	res, _ := stmt.Exec(id)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return false
	}
	num, err := res.RowsAffected()
	if num < 1 {
		fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] rmv tab_user_file faild\n")
		return false
	}
	return true
}

func GetHomePageData(phone string, foldertype int) []TableFile {
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter GetHomePageData\n")

	stmt, err := mydb.DBConn().Prepare(
		"select id,cid,file_name,ext,file_type,file_size,created_at,comment,folder_type,star,tags,minio_url from tbl_user_file where phone=? and folder_type=? order by created_at desc")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer stmt.Close()

	var tfile TableFile
	var i = 0
	var homePageFile [3]TableFile
	results, _ := stmt.Query(phone, foldertype)
	for results.Next() {
		err = results.Scan(&tfile.Id, &tfile.Cid, &tfile.FileName, &tfile.Ext, &tfile.FileType, &tfile.FileSize, &tfile.CreatAt, &tfile.Comment, &tfile.FolderType, &tfile.Star, &tfile.Tags,&tfile.MinioUrl)
		if err != nil {
			panic(err.Error())
		}
		homePageFile[i] = tfile
		i++
		if i > 3-1 {
			fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] GetHomePageData break\n")
			break
		}

	}
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter :%d\n", i)

	newhomePageFile := homePageFile[:i]
	return newhomePageFile
}

func GetRecentFile(phone string) []TableFile {
	stmt, err := mydb.DBConn().Prepare(
		"select id,cid,file_name,ext,file_type,file_size,created_at,comment,folder_type,star,tags,minio_url from tbl_user_file where phone=? and created_at>? order by created_at desc")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer stmt.Close()

	var tfile TableFile
	var i = 0
	var starfle [50]TableFile
	createdAt := time.Now().Unix() - config.MaxAge
	results, _ := stmt.Query(phone, createdAt)
	for results.Next() {
		err = results.Scan(&tfile.Id, &tfile.Cid, &tfile.FileName, &tfile.Ext, &tfile.FileType, &tfile.FileSize,
			&tfile.CreatAt, &tfile.Comment, &tfile.FolderType, &tfile.Star, &tfile.Tags,&tfile.MinioUrl)
		if err != nil {
			panic(err.Error())
		}
		starfle[i] = tfile
		i++
		if i > 50-1 {
			fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] GetRecentFile break\n")
			break
		}
	}

	newstarfle := starfle[:i]
	fmt.Fprintf(gin.DefaultWriter, "ViewStarFile[GIN-debug] %s\n", newstarfle)
	return newstarfle
}

func UpdateFileData(id int, comment string, tags string) bool {
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter UpdateFileData\n")
	stmt, err := mydb.DBConn().Prepare(
		"update tbl_user_file set comment=?,tags=? where id=?")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	ret, _ := stmt.Exec(comment, tags, id)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if ret == nil {
		fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter UpdateFileData ret in nil\n")
		return false
	}
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter UpdateFileData ret:%s\n", ret)
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

func QueryFileData(id int) FileData {
	var fileData FileData

	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter QueryFileData\n")
	stmt, err := mydb.DBConn().Prepare(
		"select comment,tags from tbl_user_file where id=?")
	if err != nil {
		fmt.Println(err.Error())
		return fileData
	}
	defer stmt.Close()

	rows, _ := stmt.Query(id)
	if rows.Next() {
		rows.Scan(&fileData.Comment, &fileData.Tags)
	}
	if err != nil || rows == nil {
		fmt.Println(err.Error())
	}
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] query coment:%s and tags:%s\n", fileData.Comment, fileData.Tags)
	return fileData
}
