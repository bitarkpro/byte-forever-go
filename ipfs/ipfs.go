package ipfs

import (
	"FileStore-Server/config"
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"io/ioutil"
	"mime/multipart"
)

var sh *shell.Shell

//upload file to IPFS
func UploadIPFS(file multipart.File) string {
	sh = shell.NewShell(config.IpfsUploadServiceHost)
	//hash, err := sh.Add(bufio.NewReader(file))
	hash, err := sh.Add(file)
	if err != nil {
		fmt.Println("上传ipfs时错误：", err)
		return ""
	}
	return hash
}

//download filde from IPFS
func CatIPFS(cid string) ([]byte, error) {
	sh = shell.NewShell(config.IpfsUploadServiceHost)
	read, err := sh.Cat(cid)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	body, err := ioutil.ReadAll(read)

	return body, nil
}
