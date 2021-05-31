package GinHandler

import (
	"FileStore-Server/common"
	"FileStore-Server/compress"
	"FileStore-Server/config"
	"FileStore-Server/db"
	"FileStore-Server/ipfs"
	"FileStore-Server/meta"
	"FileStore-Server/minio"
	"FileStore-Server/util"
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var sh *shell.Shell

type HomePagedata struct {
	FolderType int            `json:"folderType"`
	Total      int64          `json:"total"`
	Title      string         `json:"title"`
	DataList   []db.TableFile `json:"dataList"`
}
type RecentFiledata struct {
	Page     int            `json:"Page"`
	NextPage int            `json:"nextPage"`
	Total    int            `json:"total"`
	DataList []db.TableFile `json:"dataList"`
}
type Folderdata struct {
	FolderType int            `json:"folderType"`
	Title      string         `json:"title"`
	Page       int            `json:"Page"`
	NextPage   int            `json:"nextPage"`
	Total      int64          `json:"total"`
	DataList   []db.TableFile `json:"dataList"`
}

//DoUploadHandler: 处理文件上传
func DoUploadHandler(c *gin.Context) {
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] %s :insertDoUploadHandler")
	token := c.Request.Header.Get("token")
	fileType := c.Request.FormValue("fileType")
	folderType, _ := strconv.Atoi(c.Request.FormValue("folderType"))
	if folderType < 1 || folderType > 5 {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusFailed,
			"msg":  common.ErrParameter,
		})
		return
	}
	comment := c.Request.FormValue("comment")
	tags := c.Request.FormValue("tags")
	file, head, err := c.Request.FormFile("file")

	phone := IsTokenExist(token)
	if phone == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusFailed,
			"msg":  common.UserNoexist,
		})
		return
	}

	errCode := -1
	defer func() {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		if errCode < 0 {
			fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [DRROR]upload failure errCode=%d", errCode)
			c.JSON(http.StatusOK, gin.H{
				"code": common.StatusFailed,
				"msg":  common.UploadFileFailed,
			})
		} else if errCode == 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": common.StatusFailed,
				"msg":  common.OverSpace,
			})
		} else {

			c.JSON(http.StatusOK, gin.H{
				"code": common.StatusSuccess,
				"msg":  common.UploadFileSuccess,
			})

		}
	}()

	//获取表单上传的文件，并打开
	defer file.Close()

	//创建文件元信息实例
	fileMeta := meta.FileMeta{}
	fileMeta.FileName = head.Filename
	fileMeta.Ext = db.GetFileExt(head.Filename)
	fileMeta.FileType = fileType
	fileMeta.FileSize = head.Size
	fileMeta.Comment = comment
	fileMeta.Tags = tags
	fileMeta.FolderType = folderType

	//查询用户存储空间
	userStorage, _ := db.QueryUserStroage(phone)
	if err != nil {
		errCode = -2
		return
	}
	userStorage.Totalstor += fileMeta.FileSize
	if userStorage.Totalstor > config.UserTotalSpace {
		fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] \"up： %d\n", userStorage.Totalstor)
		errCode = 0
		return
	}
	sh = shell.NewShell(config.IpfsUploadServiceHost)

	//创建本地文件
	filePath := "ipfsfile/" + head.Filename
	localFile, err := os.Create(filePath)
	//复制文件信息到本地文件

	fileMeta.FileSize, err = io.Copy(bufio.NewWriter(localFile), file)

	fileMeta.Cid = ipfs.UploadIPFS(localFile)
	localFile.Close()
	if fileMeta.Cid == ""{
		fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] [ERROR] file=%s upload IPFS failed \n",head.Filename)
		errCode = -10
		return
	}
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] \"upload suc cid： %s\n", fileMeta.Cid)
	//压缩到本地
	var minioFileName string
	var fileComPressPath string
	var videoComPressPath string
	if fileType == "image" {
		fileComPressPath = compress.Picture(filePath)
		if fileComPressPath == "" {
			errCode = -3
			return
		}
	} else if fileType == "video" {
		videoComPressPath = compress.Video(filePath)
		fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] \"videoComPressPath： %s\n", videoComPressPath)
		if videoComPressPath == "" {
			errCode = -3
			return
		}
		fileComPressPath= compress.Picture(videoComPressPath)
	}
	temp := strings.Split(fileComPressPath,"/")
	if len(temp) <=1 {
		errCode = -4
		return
	}
	minioFileName = fileMeta.Cid + "_" + temp[1]

	//upload compressed picture to minio server
	uploadMinio := minio.UploadMinio(minioFileName, fileComPressPath, fileMeta.Ext)
	if !uploadMinio {
		errCode = -5
		return
	}
	fileMeta.MinioUrl = config.MinioUrl + minioFileName

	//update userfile tab
	ok := db.OnUserFileUploadFinished(phone, fileMeta.Cid, fileMeta.FileName, fileMeta.Ext,
		fileMeta.FileType, fileMeta.FileSize, fileMeta.Comment, fileMeta.FolderType, fileMeta.Tags, fileMeta.MinioUrl)
	if !ok {
		errCode = -7
		return
	}
	//update user storage
	userStorage = db.UploadStorageNum(userStorage, folderType)
	updateUserStorage := db.UpdateUserStroage(phone, userStorage)
	if !updateUserStorage {
		fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] update updateUserStorage fail\n")
		errCode = -8
		return
	}
	//delete local file
	if videoComPressPath!=""{
		err = os.Remove(videoComPressPath)
		if err != nil {
			fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
			errCode = -9
			return
		}
	}

	if fileComPressPath == filePath {
		err = os.Remove(filePath)
		if err != nil {
			fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
			errCode = -9
			return
		}
		errCode = 1
		return
	}
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug]filePath= %s\n",filePath)
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug]fileComPressPath=%s\n",fileComPressPath)
	err = os.Remove(fileComPressPath)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		errCode = -9
		return
	}
	err = os.Remove(filePath)
	if err != nil{
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		errCode = -9
		return
	}

	errCode = 1
	return
}

//foder wiew
func GetFolderHandler(c *gin.Context) {
	token := c.Request.Header.Get("token")
	folderType, err := strconv.Atoi(c.Request.FormValue("folderType"))
	if err != nil {
		return
	}
	phone := IsTokenExist(token)
	if phone == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusFailed,
			"msg":  common.UserNoexist,
		})
		return
	}
	if folderType < 1 || folderType > 5 {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusFailed,
			"msg":  common.ErrParameter,
		})
		return
	}
	Folder := [6]string{"Important moment", "My creations", "Childhood", "Young age", "Middle age", "Old age"}
	folderFile := GetFolderRsp(phone, folderType, Folder[folderType])

	c.JSON(http.StatusOK, gin.H{
		"code": common.StatusSuccess,
		"msg":  common.FileExist,
		"data": folderFile,
	})
}

func GetFolderRsp(phone string, folderType int, title string) Folderdata {
	var folderFile []db.TableFile
	var tarGetFile Folderdata
	if folderType == 0 {
		folderFile = db.ViewStarFile(phone)
	} else {
		folderFile, _ = db.GetFolderData(phone, folderType)
	}

	userStorage, err := db.QueryUserStroage(phone)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return tarGetFile
	}
	var total int64
	switch folderType {
	case 0:
		total = userStorage.StarfileNum
	case 1:
		total = userStorage.CreatNum
	case 2:
		total = userStorage.ChildNum
	case 3:
		total = userStorage.YoungNum
	case 4:
		total = userStorage.MiddleNum
	case 5:
		total = userStorage.OldNum
	default:
		fmt.Fprintf(gin.DefaultWriter, "[GIN-debug]folderType= %d is not right", folderType)
	}

	page := len(folderFile)/10 + 1
	nextPage := 0
	if page > 1 {
		nextPage = 1
	}

	tarGetFile = Folderdata{
		FolderType: folderType,
		Title:      title,
		Page:       page,
		NextPage:   nextPage,
		Total:      total,
		DataList:   folderFile,
	}
	return tarGetFile
}

func FileDetailQueryHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Request.FormValue("id"))
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return
	}
	fileDetail, err := db.GetFileInfo(id)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusFailed,
			"msg":  common.ViewSqlErr,
		})
		return
	}
	if fileDetail == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusFailed,
			"msg":  common.IdNoexist,
		})
		return
	}

	resp := util.RespMsg{
		Code: common.StatusSuccess,
		Msg:  common.FileExist,
		Data: struct {
			Id         int    `json:"id"`
			Cid        string `json:"cid"`
			FileName   string `json:"fileName"`
			Ext        string `json:"ext"`
			FileSize   int64  `json:"fileSize"`
			FileType   string `json:"FileType"`
			CreatAt    int64  `json:"creatTime"`
			Comment    string `json:"comment"`
			FolderType int    `json:"folderType"`
			Star       int    `json:"star"`
			Tags       string `json:"tags"`
			MinioUrl   string `json:"minioUrl`
		}{
			Id:         fileDetail.Id,
			Cid:        fileDetail.Cid,
			FileName:   fileDetail.FileName,
			Ext:        fileDetail.Ext,
			FileSize:   fileDetail.FileSize,
			FileType:   fileDetail.FileType,
			CreatAt:    fileDetail.CreatAt,
			Comment:    fileDetail.Comment,
			FolderType: fileDetail.FolderType,
			Star:       fileDetail.Star,
			Tags:       fileDetail.Tags,
			MinioUrl:   fileDetail.MinioUrl,
		},
	}
	c.JSON(http.StatusOK, resp)
	return
}

//DownloadHandler: download file by cid
func DownloadHandler(c *gin.Context) {
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter DownloadHandler")
	token := c.Request.Header.Get("token")
	phone := IsTokenExist(token)
	if phone == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusFailed,
			"msg":  common.UserNoexist,
		})
		return
	}

	cid := c.Request.Header.Get("cid")
	fileMeta, err := meta.GetFileMetaDB(phone, cid)

	if err != nil || fileMeta == nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"msg":  common.FileNotExist,
			"code": common.StatusFailed,
		})
		return
	}

	downLoadFile, _ := ipfs.CatIPFS(cid)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"msg":  common.DownloadFileFailed,
			"code": common.StatusFailed,
		})
		return
	}
	err = ioutil.WriteFile("ipfsfile/"+fileMeta.FileName, downLoadFile, os.ModePerm)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return
	}
	c.File("ipfsfile/" + fileMeta.FileName)

	//删除临时文件
	err = os.Remove("ipfsfile/" + fileMeta.FileName)
	if err != nil {
		fmt.Println(err)
	}
	return
}

//FileMetaUpdateHandler: update file info
func FileUpdateHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Request.FormValue("id"))
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return
	}

	commentFlag := 0
	tagsFlag := 0
	for key, valve := range c.Request.PostForm {
		if key == "comment" {
			commentFlag = 1
		}
		if key == "tags" {
			tagsFlag = 1
		}
		fmt.Printf("key:%v,valve:%v\n", key, valve)
	}
	comment := c.Request.FormValue("comment")
	tags := c.Request.FormValue("tags")
	if (commentFlag == 0) && (tagsFlag == 0) {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusFailed,
			"msg":  common.UpdateFileFail,
		})
		return
	}

	fileData := db.QueryFileData(id)

	if (commentFlag == 1) && (fileData.Comment != comment) {
		fileData.Comment = comment
	}
	if (tagsFlag == 1) && (fileData.Tags != tags) {
		fileData.Tags = tags
	}

	upData := db.UpdateFileData(id, fileData.Comment, fileData.Tags)

	if upData {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusSuccess,
			"msg":  common.UpdateFileSuc,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": common.StatusFailed,
		"msg":  common.UpdateFileFail,
	})
	return
}

//FileDeleteHandler: delete file
func FileDeleteHandler(c *gin.Context) {
	token := c.Request.Header.Get("token")
	phone := IsTokenExist(token)
	if phone == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusFailed,
			"msg":  common.UserNoexist,
		})
		return
	}
	id, err := strconv.Atoi(c.Request.FormValue("id"))
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return
	}
	fileDetail, err := db.GetFileInfo(id)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusFailed,
			"msg":  common.ViewSqlErr,
		})
		return
	}
	if fileDetail == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusFailed,
			"msg":  common.IdNoexist,
		})
		return
	}
	userStorage, _ := db.QueryUserStroage(phone)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return
	}
	//delete records in database
	if !db.RemoveFile(id) {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusFailed,
			"msg":  common.RmvFileFailed,
		})
		return
	}
	//update user info in database
	userStorageTemp := db.DeleteStorageNum(userStorage, fileDetail.FolderType, fileDetail.Star)
	userStorageTemp.Totalstor = userStorage.Totalstor - fileDetail.FileSize
	updateUserStorage := db.UpdateUserStroage(phone, userStorageTemp)
	if updateUserStorage {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusSuccess,
			"msg":  common.RmvFileSuccess,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": common.StatusFailed,
		"msg":  common.RmvFileFailed,
	})

	return
}

func GetHomePageHandler(c *gin.Context) {
	token := c.Request.Header.Get("token")
	phone := IsTokenExist(token)
	if phone == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusFailed,
			"msg":  common.UserNoexist,
		})
		return
	}
	userStorage, err := db.QueryUserStroage(phone)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return
	}

	var tarGetFiletemp [6]HomePagedata
	var tarGetFile []HomePagedata
	starFile := db.ViewStarFile(phone)
	if len(starFile) > 3 {
		starFile = starFile[:3]
	}

	tarGetFiletemp[0] = HomePagedata{
		FolderType: 0,
		Total:      userStorage.StarfileNum,
		Title:      "Important moment",
		DataList:   starFile,
	}
	temp1 := db.GetHomePageData(phone, 1)
	tarGetFiletemp[1] = HomePagedata{
		FolderType: 1,
		Total:      userStorage.CreatNum,
		Title:      "My creations",
		DataList:   temp1,
	}
	temp2 := db.GetHomePageData(phone, 2)
	tarGetFiletemp[2] = HomePagedata{
		FolderType: 2,
		Total:      userStorage.ChildNum,
		Title:      "Childhood",
		DataList:   temp2,
	}
	temp3 := db.GetHomePageData(phone, 3)
	tarGetFiletemp[3] = HomePagedata{
		FolderType: 3,
		Total:      userStorage.YoungNum,
		Title:      "Young age",
		DataList:   temp3,
	}
	temp4 := db.GetHomePageData(phone, 4)
	tarGetFiletemp[4] = HomePagedata{
		FolderType: 4,
		Total:      userStorage.MiddleNum,
		Title:      "Middle age",
		DataList:   temp4,
	}
	temp5 := db.GetHomePageData(phone, 5)
	tarGetFiletemp[5] = HomePagedata{
		FolderType: 5,
		Total:      userStorage.OldNum,
		Title:      "Old age",
		DataList:   temp5,
	}

	for i := 0; i < 6; i++ {
		tarGetFile = append(tarGetFile, tarGetFiletemp[i])
	}

	if tarGetFile == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusSuccess,
			"msg":  common.StarFileNoExist,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": common.StatusSuccess,
		"msg":  common.HomePageFile,
		"data": tarGetFile,
	})
}

func GetRecentHandler(c *gin.Context) {
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter GetRecentHandler")
	token := c.Request.Header.Get("token")
	phone := IsTokenExist(token)
	if phone == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusFailed,
			"msg":  common.UserNoexist,
		})
		return
	}
	recentFile := db.GetRecentFile(phone)
	pageTemp := len(recentFile)/10 + 1
	nextPageTemp := 0
	if pageTemp > 1 {
		nextPageTemp = 1
	}
	recentFileTemp := RecentFiledata{
		Page:     pageTemp,
		NextPage: nextPageTemp,
		Total:    len(recentFile),
		DataList: recentFile,
	}
	c.JSON(http.StatusOK, gin.H{
		"code": common.StatusSuccess,
		"msg":  common.RecentFile,
		"data": recentFileTemp,
	})

}
