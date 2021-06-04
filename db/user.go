package db

import (
	mydb "FileStore-Server/db/mysql"
	"fmt"
	"github.com/gin-gonic/gin"
)

//User: 用户表model

type User struct {
	Phone       string
	OpenID      string
	SessionKey  string
	Iv          string
	ForeverId   string
	Integral    int64
	Totalstor   int64
	FileNum     int64
	StarfileNum int64
	//Status       int
	CreatedAt string
	//	UpdateAt  string
}
type LoginData struct {
	OpenID string
	//UnionID     string
	SessionKey string
}

type UserStorage struct {
	Phone       string
	Totalstor   int64
	FileNum     int64
	StarfileNum int64
	Integral    int64
	CreatNum    int64
	ChildNum    int64
	YoungNum    int64
	MiddleNum   int64
	OldNum      int64
}

//UpdateToken: 刷新用户登录token
func UpdateToken(phone string, token string) bool {
	fmt.Fprintln(gin.DefaultWriter, "[GIN-debug] enter UpdateToken 2222")
	stmt, err := mydb.DBConn().Prepare(
		"replace into tbl_user_token (`phone`,`user_token`) values (?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	fmt.Fprintln(gin.DefaultWriter, "[GIN-debug] enter UpdateToken 333s")

	ret, _ := stmt.Query(phone, token)

	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] %s :enter UpdateToken 55555\n", ret)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}

func GetUserInfo(phone string) (User, error) {
	user := User{}
	stmt, err := mydb.DBConn().Prepare(
		"select phone,totalstor,integral,file_num,starfile_num,created_at from tbl_user where phone=? limit 1")
	if err != nil {
		return user, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(phone).Scan(&user.Phone, &user.Totalstor, &user.Integral, &user.FileNum, &user.StarfileNum, &user.CreatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}

func WxInsertReqData(openid string, sessionKey string) bool {
	fmt.Fprintln(gin.DefaultWriter, "[GIN-debug] enter WxInsertReqData 2222")
	stmt, err := mydb.DBConn().Prepare(
		"replace into tbl_wx_request (`openid`,`session_key`) values (?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	fmt.Fprintln(gin.DefaultWriter, "[GIN-debug] enter WxInsertReqData 333s")

	ret, _ := stmt.Query(openid, sessionKey)

	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] %s :enter WxInsertReqData 55555\n", ret)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func WxInsertUserData(phone string, openid string, sessionkey string) bool {

	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user (`openid`,`phone`,`session_key`) values (?,?,?)")
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter WxInsertUserData11\n")
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return false
	}
	defer stmt.Close()
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter WxInsertUserData22222\n")

	ret, err := stmt.Exec(openid, phone, sessionkey)
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter WxInsertUserData3333\n")
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
	}
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter WxInsertUserDatas44444\n")
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] %s :nter WxInsertUserDatas55555\n", ret)
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {

		fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter WxInsertUserDatas5555\n")
		return true
	}
	return true
}

func WxQueryReqData(openid string) (LoginData, error) {
	loginData := LoginData{}
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter WxQueryReqData\n")
	stmt, err := mydb.DBConn().Prepare(
		"select * from tbl_wx_request where openid = ? limit 1")
	defer stmt.Close()
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return loginData, err
	}
	rows, err := stmt.Query(openid)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return loginData, err
	} else if rows == nil {
		fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] openid:%s not found: \n", openid)
		return loginData, err
	}

	pRows := mydb.ParseRows(rows)
	loginData.SessionKey = string(pRows[0]["session_key"].([]byte))
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] %s :i enter Wx sign SessionKey \n", loginData.SessionKey)
	return loginData, nil

}

//update user info
func UpdateUserStroage(phone string, userStorage UserStorage) bool {
	stmt, err := mydb.DBConn().Prepare(
		"update tbl_user set totalstor=?,file_num=?,starfile_num=?,creat_num=?,child_num=?,young_num=?,middle_num=?,old_num=?,integral=? where phone=?")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	ret, _ := stmt.Exec(userStorage.Totalstor, userStorage.FileNum, userStorage.StarfileNum, userStorage.CreatNum,
		userStorage.ChildNum, userStorage.YoungNum, userStorage.MiddleNum, userStorage.OldNum, userStorage.Integral, phone)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter WxInsertUserDatas5555\n")
		return true
	}
	return false
}

func QueryUserStroage(phone string) (UserStorage, error) {
	var userStorage UserStorage //预编译sql语句
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter QueryUsedSpace\n")
	stmt, err := mydb.DBConn().Prepare(
		"select totalstor,file_num,starfile_num,creat_num,child_num,young_num,middle_num,old_num,integral from tbl_user where phone=?")
	if err != nil {
		fmt.Println(err.Error())
		return userStorage, err
	}
	defer stmt.Close()

	rows, _ := stmt.Query(phone)

	if rows.Next() {
		rows.Scan(&userStorage.Totalstor, &userStorage.FileNum, &userStorage.StarfileNum, &userStorage.CreatNum,
			&userStorage.ChildNum, &userStorage.YoungNum, &userStorage.MiddleNum, &userStorage.OldNum, &userStorage.Integral)
	}

	if err != nil || rows == nil {
		fmt.Println(err.Error())
		return userStorage, err
	}
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] enter QueryUsedSpace:%s\n", userStorage)
	return userStorage, nil
}

func UploadStorageNum(userStorage UserStorage, folderType int) UserStorage {
	userStorage.FileNum += 1
	userStorage.Integral += 5
	switch folderType {
	case 1:
		userStorage.CreatNum += 1
	case 2:
		userStorage.ChildNum += 1
	case 3:
		userStorage.YoungNum += 1
	case 4:
		userStorage.MiddleNum += 1
	case 5:
		userStorage.OldNum += 1
	default:
		fmt.Fprintf(gin.DefaultWriter, "[GIN-debug]folderType= %d is not right", folderType)
	}
	return userStorage
}
func DeleteStorageNum(userStorage UserStorage, folderType int, starFile int) UserStorage {
	if starFile == 1 {
		userStorage.StarfileNum -= 1
	}
	userStorage.FileNum -= 1
	userStorage.Integral -= 5
	switch folderType {
	case 1:
		userStorage.CreatNum -= 1
	case 2:
		userStorage.ChildNum -= 1
	case 3:
		userStorage.YoungNum -= 1
	case 4:
		userStorage.MiddleNum -= 1
	case 5:
		userStorage.OldNum -= 1
	default:
		fmt.Fprintf(gin.DefaultWriter, "[GIN-debug]folderType= %d is not right", folderType)
	}
	return userStorage
}
