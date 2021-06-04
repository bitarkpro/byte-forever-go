package GinHandler

import (
	"FileStore-Server/aes_128_cdc"
	"FileStore-Server/common"
	cfg "FileStore-Server/config"
	dblayer "FileStore-Server/db"
	mydb "FileStore-Server/db/mysql"
	"FileStore-Server/util"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
	//"encoding/json"
)

type User struct {
	Phone      string
	OpenID     string
	SessionKey string
	Iv         string
	//Status       int
	Totalstor int64
	CreatedAt string
	//	UpdateAt  string
}
type LoginData struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
}

//phone data
type PhoneData struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
	Watermark
}
type Watermark struct {
	Appid     string `json:"appid"`
	Timestamp string `json:"timestamp"`
}

//UserInfoHandler:query user info
func UserInfoHandler(c *gin.Context) {

	token := c.Request.FormValue("token")
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] myToken is:%s\n", token)
	phone := IsTokenExist(token)
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] phone is:%s\n", phone)
	if phone == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusTokenInvalid,
			"msg":  common.TokenInvalid,
		})

		return
	}

	user, err := dblayer.GetUserInfo(phone)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusFailed,
			"msg":  common.ViewUserInfoFail,
		})
		return
	}
	config := cfg.Conf

	resp := util.RespMsg{
		Code: common.StatusSuccess,
		Msg:  common.ViewUserInfoSuc,
		Data: struct {
			Phone       string `json:"phone"`
			ForeverId   string `json:"foreverid"`
			Integral    int64  `json:"integral"`
			TotalSpace  int64  `json:"totalspace"`
			Totalstor   int64  `json:"totalstor"`
			FileNum     int64  `json:"filenum"`
			StarFileNum int64  `json:"starfilenum"`
		}{
			Phone:       user.Phone,
			ForeverId:   "暂无",
			Integral:    user.Integral,
			TotalSpace:  config.UserTotalSpace,
			Totalstor:   user.Totalstor,
			FileNum:     user.FileNum,
			StarFileNum: user.StarfileNum,
		},
	}
	c.JSON(http.StatusOK, resp)
}

//Wechat login：
func WxDoGetCodeHandler(c *gin.Context) {
	fmt.Fprintln(gin.DefaultWriter, "[GIN-debug] wxuser get code")
	code := c.Request.FormValue("code")
	if code == "" {
		resp := util.RespMsg{
			Code: common.StatusFailed,
			Msg:  common.CodeInvalid,
		}
		c.Data(http.StatusOK, "application/json", resp.JSONBytes())

		return

	}
	fmt.Fprintln(gin.DefaultWriter, "[GIN-debug] code=", code)
	config := cfg.Conf
	loginData, err := WxGetOpenIdHandler(code, config.AppId, config.Secret)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
	}
	if loginData.OpenID == "" {
		resp := util.RespMsg{
			Code: common.StatusFailed,
			Msg:  common.CodeInvalid,
		}
		c.Data(http.StatusOK, "application/json", resp.JSONBytes())
		return
	}

	fmt.Fprintln(gin.DefaultWriter, "[GIN-debug] openid= ", loginData.OpenID)
	fmt.Fprintln(gin.DefaultWriter, "[GIN-debug] SessionKey= ", loginData.SessionKey)

	//return openid
	resp := util.RespMsg{
		Code: common.StatusSuccess,
		Msg:  "openid",
		Data: struct {
			Openid string `json:"openid"`
		}{
			Openid: loginData.OpenID,
		},
	}
	c.JSON(http.StatusOK, resp)
	suc := dblayer.WxInsertReqData(loginData.OpenID, loginData.SessionKey)
	if suc {
		fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] %s :insert reqdata success\n", loginData.OpenID)
	} else {
		fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] %s :insert reqdata failed\n", loginData.OpenID)

	}
	return
}

func WxDoSignInHandler(c *gin.Context) {

	user := dblayer.User{}
	user.OpenID = c.Request.FormValue("openid")
	user.Iv = c.Request.FormValue("iv")

	//query sessionkey and unionid
	loginData, err := dblayer.WxQueryReqData(user.OpenID)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return
	}

	user.SessionKey = loginData.SessionKey
	//get enc phone
	encPhone := c.Request.FormValue("encryptedData")
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug]encPhone: %s \n", encPhone)
	user.Phone = DecPhoneNum(encPhone, user.Iv, user.SessionKey)
	if user.Phone == "" {
		c.JSON(http.StatusOK, gin.H{
			"msg":  common.UserLoginfail,
			"code": common.StatusFailed,
		})
		return
	}

	userChecked := dblayer.WxInsertUserData(user.Phone, user.OpenID, user.SessionKey)

	if !userChecked {
		c.JSON(http.StatusOK, gin.H{
			"msg":  common.UserLoginfail,
			"code": common.StatusFailed,
		})
		return
	}

	//2.generate Token
	token := GenToken(user.Phone)
	upRes := dblayer.UpdateToken(user.Phone, token)
	if !upRes {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusFailed,
			"msg":  common.TokenUpdateFail,
		})
		return
	}

	//3.resp message
	resp := util.RespMsg{
		Code: common.StatusSuccess,
		Msg:  common.UserLoginSuc,
		Data: struct {
			Phone string `json:"phone"`
			Token string `json:"token"`
		}{
			Phone: user.Phone,
			Token: token,
		},
	}

	c.JSON(http.StatusOK, resp)
}

//get openid ，sessionkey
func WxGetOpenIdHandler(code string, appid string, secret string) (LoginData, error) {
	config := cfg.Conf
	var params = "?appid=" + appid + "&secret=" + secret + "&js_code=" + code +
		"&grant_type=" + config.GrantType
	var url = "https://api.weixin.qq.com/sns/jscode2session" + params
	//fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] params：%s\n",params)
	loginData := LoginData{}
	response, err := http.Get(url)
	if err != nil {
		return loginData, err
	}
	data, err := ioutil.ReadAll(response.Body)
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] data：%s\n", data)
	if err != nil {
		return loginData, err
	}

	err = json.Unmarshal(data, &loginData)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return loginData, err
	}
	return loginData, nil
}

func DecPhoneNum(encryptedData string, session_key string, iv string) string {

	decPhoneData := aes_128_cdc.PswDecrypt(encryptedData, session_key, iv)
	phoneData := PhoneData{}
	err := json.Unmarshal([]byte(decPhoneData), &phoneData)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return ""
	}
	return phoneData.PhoneNumber
}

//GenToken
func GenToken(phone string) string {
	config := cfg.Conf
	claims := jwt.MapClaims{
		"username": phone,
		"exp":      time.Now().Add(time.Duration(config.MaxAge) * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.KeyInfo))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("token: %v\n", tokenString)

	ret, err := ParseToken(tokenString)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("userinfo: %v\n", ret)
	return tokenString
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	config := cfg.Conf
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(config.KeyInfo), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func IsTokenExist(token string) string {
	stmt, err := mydb.DBConn().Prepare(
		"select phone from tbl_user_token where user_token = ? limit 1")
	defer stmt.Close()
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return ""
	}

	rows, err := stmt.Query(token)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return ""
	} else if rows == nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] phone:%s not found\n", token)
		return ""
	}

	pRows := mydb.ParseRows(rows)

	if len(pRows) > 0 {
		return string(pRows[0]["phone"].([]byte))
	}

	return ""
}
