package minio

import (
	cfg "FileStore-Server/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/net/context"
)

//Object storage
func UploadMinio(objectName string, filePath string, ext string) bool {
	fmt.Fprintln(gin.DefaultWriter, "[GIN-debug] enter UploadMinio")
	config := cfg.Conf
	ctx := context.Background()
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccesskeyID, config.AccessKeySecret, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR]creat minio err %v\n", err)
		return false
	}
	contentType := "application/" + ext
	// Upload
	info, err := minioClient.FPutObject(ctx, config.Bucket, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR]upload err %v\n", err)
		return false
	}
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] Success upload minio %s of size %d\n", objectName, info.Size)
	/*	_ = minioClient.FGetObject(ctx, bucketName, objectName, filePath, minio.GetObjectOptions{})*/
	return true
}
