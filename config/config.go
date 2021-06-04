package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	DriverName     string `yaml:"drivername"`
	DataSourceName string `yaml:"datasourcename"`

	Bucket          string `yaml:"bucket"`
	Location        string `yaml:"location"`
	AccesskeyID     string `yaml:"accesskeyid"`
	AccessKeySecret string `yaml:"accesskeysecret"`
	UseSSL          bool   `yaml:"usessl"`
	Endpoint        string `yaml:"endpoint"`
	MinioUrl        string `yaml:"miniourl"`

	UploadServiceHost     string `yaml:"uploadservicehost"`
	IpfsUploadServiceHost string `yaml:"ipfsuploadservicehost"`

	GrantType string `yaml:"granttype"`
	AppId     string `yaml:"appid"`
	Secret    string `yaml:"secret"`

	KeyInfo        string `yaml:"keyinfo"`
	MaxAge         int64  `yaml:"maxage"`
	UserTotalSpace int64  `yaml:"usertotalspace"`
}

var Conf = Config{}

func GetConf(filename string) Config {
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println("[GIN-debug][GIN-error]\n", err)
	}

	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		fmt.Println("[GIN-debug][GIN-error]\n", err)
	}

	fmt.Printf("[GIN-debug]load config suc:%v", Conf)
	return Conf
}
