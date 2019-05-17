// get configuration here

package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const configFileName = "poster.json"

// conf is config struct
type conf struct {
	MPAppID     string
	MPAppSecret string
	OutputDIR   string
}

// C 配置实列
var C conf

// init 加载配置文件
func init() {
	f, err := os.Open(configFileName)
	if err != nil {
		f, _ := os.OpenFile(configFileName, os.O_CREATE|os.O_RDWR, 0666)
		b, _ := json.MarshalIndent(C, "", "    ")
		f.Write(b)
		fmt.Println("配置文件已生成，请填写完成后再次运行")
		os.Exit(0)
	}
	confBytes, _ := ioutil.ReadAll(f)
	err = json.Unmarshal(confBytes, &C)
	if err != nil {
		fmt.Println("配置文件格式无效")
		os.Exit(0)
	}

	if C.OutputDIR == "" {
		C.OutputDIR = "output/"

	}
}
