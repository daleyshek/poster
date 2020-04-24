package common

import (
	"encoding/json"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// ServeHTTP 简单http接口
func ServeHTTP() {
	http.HandleFunc("/", postHandle)
	http.ListenAndServe(":"+C.ListenPort, nil)
}

type httpReq struct {
	AccessToken string `form:"access_token" json:"access_token"`
	Scene       string `form:"scene" json:"scene"`
	Page        string `form:"page" json:"page"`
	Width       string `form:"width" json:"width"`
	Title       string `form:"title" json:"title"`         //标题
	Content     string `form:"content" json:"content"`     //文字内容
	ImageURL    string `form:"image_url" json:"image_url"` //主图http链接或文件链接
	qRCodeURL   string //二维码http链接或文件链接
	OutputDIR   string `form:"output_dir" json:"output_dir"`     //保存文件目录，可选
	BorderColor string `form:"border_color" json:"border_color"` //边框颜色，可选
}

type httpResp struct {
	Poster string `json:"poster"` //海报图文件名或链接
	QRCode string `json:"qrcode"` //二维码文件名或链接
}

func postHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Println("请使用POST方法")
		writeResponse(w, map[string]interface{}{"error": 1, "message": "请使用POST方法"})
		return
	}
	var rq httpReq
	switch r.Header.Get("Content-Type") {
	case "application/json":
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("请求参数异常", err)
			writeResponse(w, map[string]interface{}{"error": 1, "request": rq, "message": "请求参数异常：缺少参数"})
			return
		}
		err = json.Unmarshal(data, &rq)
		if err != nil {
			log.Println("请求参数异常", err)
			writeResponse(w, map[string]interface{}{"error": 1, "request": rq, "message": "请求参数异常:不是有效的JSON格式"})
			return
		}
	default:
		rq.AccessToken = r.PostFormValue("access_token")
		rq.Scene = r.PostFormValue("scene")
		rq.Page = r.PostFormValue("page")
		rq.Width = r.PostFormValue("width")
		rq.Title = r.PostFormValue("title")
		rq.Content = r.PostFormValue("content")
		rq.ImageURL = r.PostFormValue("image_url")
		rq.OutputDIR = r.PostFormValue("output_dir")
		rq.BorderColor = r.PostFormValue("border_color")
	}

	if rq.AccessToken == "" || rq.Scene == "" || rq.Page == "" || rq.Title == "" || rq.Content == "" {
		log.Println("POST参数无效")
		writeResponse(w, map[string]interface{}{"error": 1, "request": rq, "message": "POST参数无效"})
		return
	}

	// 获取小程序码保存到本地
	_width, err := strconv.Atoi(rq.Width)
	if err != nil {
		log.Println("width 必须是有效的数字")
		writeResponse(w, map[string]interface{}{"error": 1, "request": rq, "message": "width 必须是有效的数字"})
		return
	}

	qrcodeName, err := RequestQRCode(QRCodeReq{
		Scene: rq.Scene,
		Page:  rq.Page,
		Width: _width,
	}, "", rq.AccessToken)

	if err != nil {
		log.Println("获取小程序码失败", err)
		writeResponse(w, map[string]interface{}{"error": 1, "request": rq, "message": "获取小程序码失败"})
		return
	}

	// 使用生成的二维码来生成海报
	rq.qRCodeURL = C.OutputDIR + qrcodeName
	// 生成海报
	var br, bg, bb uint8
	br, bg, bb = 255, 255, 255
	if rq.BorderColor != "" {
		if rq.BorderColor[0:1] == "#" {
			var _r, _g, _b uint64
			var rHex, gHex, bHex string
			switch {
			case len(rq.BorderColor) == 4:
				rHex = "0x" + rq.BorderColor[1:2]
				gHex = "0x" + rq.BorderColor[2:3]
				bHex = "0x" + rq.BorderColor[3:4]
				_r, _ = strconv.ParseUint(rHex, 0, 8)
				_g, _ = strconv.ParseUint(gHex, 0, 8)
				_b, _ = strconv.ParseUint(bHex, 0, 8)
				br = uint8(_r * 16)
				bg = uint8(_g * 16)
				bb = uint8(_b * 16)
			case len(rq.BorderColor) == 7:
				rHex = "0x" + rq.BorderColor[1:3]
				gHex = "0x" + rq.BorderColor[3:5]
				bHex = "0x" + rq.BorderColor[5:7]
				_r, _ = strconv.ParseUint(rHex, 0, 8)
				_g, _ = strconv.ParseUint(gHex, 0, 8)
				_b, _ = strconv.ParseUint(bHex, 0, 8)
				br = uint8(_r)
				bg = uint8(_g)
				bb = uint8(_b)
			default:
				log.Println("颜色值无效，支持格式如下：#ffffff，#fff", err)
				writeResponse(w, map[string]interface{}{"error": 1, "request": rq, "message": "颜色值无效，支持格式如下：#ffffff，#fff"})
				return
			}
		}
	}

	bdC := color.RGBA{br, bg, bb, 255}
	posterName, err := DrawPoster(Style{
		ImageURL:       rq.ImageURL,
		QRCodeURL:      rq.qRCodeURL,
		Title:          rq.Title,
		Content:        rq.Content,
		OutputFileName: rq.OutputDIR,
		OutputDIR:      rq.OutputDIR,
		BorderColor:    bdC,
	})
	if err != nil {
		log.Println("海报生成失败", err)
		writeResponse(w, map[string]interface{}{"error": 1, "request": rq, "message": "海报生成失败"})
		return
	}
	writeResponse(w, map[string]interface{}{"error": 0, "result": map[string]string{"poster": posterName, "qrcode": qrcodeName}})
}

func writeResponse(w http.ResponseWriter, d map[string]interface{}) {
	data, _ := json.MarshalIndent(d, "", "  ")
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}
