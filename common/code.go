package common

type ErrorCode int32

const (
	StatusSuccess      = 200
	StatusFailed       = 500
	StatusTokenInvalid = -1

	AddStarFile           = 1
	CancelStarFile        = 0
	CodeInvalid           = "invalid code"
	TokenInvalid          = "invalid token"
	TokenUpdateFail       = "UPdate token failed"
	UserNoexist           = "User does not exist, please login or register first"
	UserLoginSuc          = "Log in successfully"
	UserLoginfail         = "Log in failed"
	IdNoexist             = "ID error, query failed"
	ViewUserInfoSuc       = "Query user information successful"
	ViewUserInfoFail      = "Failed to query user information"
	AddStarfileSuccess    = "Added star file successfully"
	AddStarfileFail       = "Failed to add star file"
	CancelStarfileSuccess = "The star file was cancelled successfully"
	CancelStarfileFail    = "Failed to cancel star file"
	StarFileNoExist       = "No star file"

	FileNotExist      = "File information does not exist"
	FileExist         = "Query file successful"
	ViewSqlErr        = "Query database failed"
	RmvFileSuccess    = "File deleted successfully"
	RmvFileFailed     = "Failed to delete file"
	UploadFileSuccess = "Upload the file successfully"
	UploadFileFailed  = "Failed to upload file"

	DownloadFileFailed = "Failed to download file"
	OverSpace          = "The space is insufficient"
	HomePageFile       = "Home page files"
	RecentFile         = "Recently uploaded file"
	ErrResult          = "Encounters an error"
	ErrParameter       = "Parameter is wrong"

	UpdateFileSuc  = "Update success"
	UpdateFileFail = "Update fail"
)
