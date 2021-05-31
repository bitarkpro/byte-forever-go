package meta

import (
	mydb "FileStore-Server/db"
)


type FileMeta struct {
	Id int
	//Phone         string
	Cid        string
	FileName   string
	Ext        string
	FileType   string
	FileSize   int64
	CreatAt    int64
	Comment    string
	FolderType int
	Star       int
	Tags       string
	MinioUrl    string
}


var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}


func GetFileMetaDB(phone string, cid string) (*FileMeta, error) {
	file, err := mydb.GetFileMeta(phone, cid)
	if err != nil {
		return &FileMeta{}, err
	}

	fmeta := FileMeta{
		Id:         file.Id,
		Cid:        file.Cid,
		FileName:   file.FileName,
		Ext:        file.Ext,
		FileSize:   file.FileSize,
		FileType:   file.FileType,
		CreatAt:    file.CreatAt,
		Comment:    file.Comment,
		FolderType: file.FolderType,
		//Star       :file.Star,
	}

	return &fmeta, nil
}
