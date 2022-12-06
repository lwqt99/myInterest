package tools

import (
	"encoding/pem"
	"io/ioutil"
	"os"
)

var basePath = "./gitgit/certificate/"

//生成pem格式文件
func GenPem(fileName string, content []byte) error {
	block := &pem.Block{
		Type:  "Key",
		Bytes: content,
	}
	file, err := os.Create(basePath + fileName+".pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}

func ReadPem(fileName string) ([]byte, error) {
	key,err := ioutil.ReadFile(basePath + fileName+".pem")

	return key,err
}