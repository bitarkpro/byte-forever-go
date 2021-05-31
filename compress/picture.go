package compress

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	//https://github.com/FFmpeg/FFmpeg
)
type InputArgs struct {
	OutputPath string
	LocalPath  string
	Quality    int
	Width      int
	Format     string
}
const (
	PictureWide = 500
	PictureQuality = 100
)


func Picture(filename string)  string {
	var inputArgs InputArgs

	inputArgs.LocalPath = filename
	inputArgs.Quality = PictureQuality
	inputArgs.Width = PictureWide

	pathTemp, format, top := IsPictureFormat(inputArgs.LocalPath)
	if pathTemp != "" {
		fmt.Println("start compress")
		inputArgs.OutputPath = top + "_compress." + format
		isCompress,flag:= imageCompress(
			func() (io.Reader, error) {
				return os.Open(inputArgs.LocalPath)
			},
			func() (*os.File, error) {
				return os.Open(inputArgs.LocalPath)
			},
			inputArgs.OutputPath,
			inputArgs.Quality,
			inputArgs.Width,
			format)

		if !isCompress{

			fmt.Println("compress pic fail")
			return ""

		} else {
			if flag == 0 {
				inputArgs.OutputPath = inputArgs.LocalPath
			}
			fmt.Printf("compress pic suc file=%s,flag=%d\n",inputArgs.OutputPath,flag)
			return inputArgs.OutputPath
		}

	}
return ""
}

func imageCompress(
	getReadSizeFile func() (io.Reader,error),
	getDecodeFile func() (*os.File,error),
	name string,
	Quality,
	base int,
	format string) (bool,int){
	//flag
	 compressFlag:=0
	//read file
	fileOrigin, err := getDecodeFile()
	defer fileOrigin.Close()
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] os.Open(file) fail, %v\n", err)
		return false,compressFlag
	}

	var origin image.Image
	var config image.Config
	var temp io.Reader
	/** get size */
	temp, err = getReadSizeFile()

	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] os.Open(temp),temp=%v\n",temp)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return false,compressFlag
	}
	var typeImage int64
	format = strings.ToLower(format)
	/** jpg style */
	if format=="jpg" || format =="jpeg" {
		typeImage = 1
		origin, err = jpeg.Decode(fileOrigin)
		if err != nil {
			fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR]jpeg.Decode(file_origin), %v\n", err)
			return false,compressFlag
		}
		temp, err = getReadSizeFile()
		if err != nil {
			fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR]os.Open(temp), %v\n", err)
			log.Fatal(err)
			return false,compressFlag
		}
		config,err = jpeg.DecodeConfig(temp)
		if err != nil {
			fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR]jpeg.DecodeConfig(temp), %v\n", err)
			return false,compressFlag
		}
	}else if format=="png" {
		typeImage = 0
		origin, err = png.Decode(fileOrigin)
		if err != nil {
			fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR]png.Decode(file_origin), %v\n", err)
			return false,compressFlag
		}
		temp, err = getReadSizeFile()
		if err != nil {
			fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR]os.Open(temp), %v\n", err)
			return false,compressFlag
		}
		config,err = png.DecodeConfig(temp)
		if err != nil {
			fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR]png.DecodeConfig(temp), %v\n", err)
			return false,compressFlag
		}
	}
	// compress
	if config.Width <= PictureWide{
		fmt.Printf("width=%d (<= 500) no need compress\n",config.Width)
		return true,compressFlag
	}
	width  := uint(base)
	height := uint(base*config.Height/config.Width)
	fmt.Printf("width=%d,height=%d\n",width,height)

	canvas := resize.Thumbnail(width, height, origin, resize.Lanczos3)
	fileOut, err := os.Create(name)
	defer fileOut.Close()
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		return false,compressFlag
	}
	if typeImage==0 {
		err = png.Encode(fileOut, canvas)
		if err!=nil {
			fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR]compress pic fail, %v\n", err)
			return false,compressFlag
		}
	}else{
		err = jpeg.Encode(fileOut, canvas, &jpeg.Options{Quality})
		if err!=nil {
			fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR]compress pic fail, %v\n", err)
			return false,compressFlag
		}
	}
	//set flag =1
	compressFlag = 1
	return true,compressFlag
}

/** is picture */
func IsPictureFormat(path string) (string,string,string) {
	temp := strings.Split(path,".")
	if len(temp) <=1 {
		return "","",""
	}
	mapRule := make(map[string]int64)
	mapRule["jpg"]  = 1
	mapRule["png"]  = 1
	mapRule["jpeg"] = 1
	/** other style */
	if mapRule[temp[1]] == 1  {
		fmt.Fprintf(gin.DefaultWriter, "[GIN-debug]filetype=%s\n", temp[1])
		return path,temp[1],temp[0]
	}else{
		return "","",""
	}
}


func Video(filename string) string{
	buf := GetFrame(1,filename)
	temp := strings.Split(filename,".")
	if len(temp) <=1 {
		return ""
	}
	picName:=temp[0]+".jpg"
	  err := ioutil.WriteFile(picName,buf.Bytes(), os.ModePerm)
	  if err != nil {
		  fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		  return ""
	  }
	  fmt.Fprintln(gin.DefaultWriter, "[GIN-debug] SuccesGets the first frame of the video,file=",picName)
	  return picName
  	}
func GetFrame(index int,filename string) *bytes.Buffer {
    cmd := exec.Command("ffmpeg", "-i", filename, "-vframes", strconv.Itoa(index),"-f", "singlejpeg", "-")

    buf := new(bytes.Buffer)
    cmd.Stdout = buf
    if cmd.Run() != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", cmd.Run())
    }

    return buf
}








