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
	Mapp         mappConf `json:"mapp"`
	OutputDIR    string   `json:"output_dir"`    // 图片保存目录
	ListenPort   string   `json:"listen_port"`   // 监听端口
	FontFilePath string   `json:"fontfile_path"` // 字体路径
}

// mappConf 小程序配置
type mappConf struct {
	AppID     string `json:"app_id"`     // 小程序app_id
	AppSecret string `json:"app_secret"` // 小程序app_secret
}

// C 配置实列
var C conf

// Init 加载配置文件
func Init() {
	f, err := os.Open(configFileName)
	if err != nil {
		f, _ := os.OpenFile(configFileName, os.O_CREATE|os.O_RDWR, 0666)
		b, _ := json.MarshalIndent(conf{
			OutputDIR:    "output/",
			ListenPort:   "2020",
			FontFilePath: "./resources/font.ttc",
			Mapp: mappConf{
				AppID:     "小程序app_id,可选",
				AppSecret: "小程序app_secret,可选",
			},
		}, "", "    ")
		f.Write(b)
		fmt.Println("配置文件已生成。除非测试，没必要填写小程序app_id和app_secret，使用自己的业务逻辑获取access_token传值代替。")
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
	} else {
		if C.OutputDIR[len(C.OutputDIR)-1:] != "/" {
			C.OutputDIR = C.OutputDIR + "/"
		}
	}

}
