package util

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetImgBase64ByUrl(url string) (r string, err error) {
	//获取远端图片
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("A error occurred!")
		return
	}
	defer res.Body.Close()

	// 读取获取的[]byte数据
	data, _ := ioutil.ReadAll(res.Body)

	imageBase64 := base64.StdEncoding.EncodeToString(data)
	return imageBase64, err

}
