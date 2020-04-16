// 获取小程序二维码，获取的二维码数量不限
// 微信的API需要授权，须提供参数 appID 和 appSecret

package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// MPGETAccessTokenURL 获取accessToken
const MPGETAccessTokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"

// MPAccessToken 小程序获取的access_token
type MPAccessToken struct {
	Errcode     int       `json:"errcode"`
	Errmsg      string    `json:"errmsg"`
	AccessToken string    `json:"access_token"`
	ExpiresIn   int       `json:"expires_in"`
	ExpiresAt   time.Time //过期时间，s
}

var mpToken MPAccessToken

// GetMPAccessToken 获取小程序授权的accesstoken
func GetMPAccessToken() (MPAccessToken, error) {
	// 已过期
	if mpToken.ExpiresAt.After(time.Now()) && mpToken.AccessToken != "" {
		return mpToken, nil
	}
	// 重新获取
	url := fmt.Sprintf(MPGETAccessTokenURL, C.MPAppID, C.MPAppSecret)
	resp, err := http.Get(url)
	if err != nil {
		return mpToken, err
	}
	defer resp.Body.Close()
	bd, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bd, &mpToken)
	if err != nil {
		fmt.Println("json.Unmarshal failed")
		return mpToken, err
	}

	if mpToken.Errcode == 0 && mpToken.AccessToken != "" {
		mpToken.ExpiresAt = time.Now().Add(time.Second * time.Duration(mpToken.ExpiresIn))
		return mpToken, nil
	}
	return mpToken, errors.New(mpToken.Errmsg)
}

// mpQRCodeURL 获取二维码
const mpQRCodeURL = "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s"

// QRCodeReq 请求获取二维码的数据
type QRCodeReq struct {
	// AccessToken string `json:"access_token"`
	Scene string `json:"scene"`
	Page  string `json:"page"`
	Width int    `json:"width"`
}

// RequestQRCode 获取小程序二维码
// saveName 不包含后缀，为空的话就是随机文件名
// 返回文件名和错误，文件存储在OutputDIR目录中
func RequestQRCode(req QRCodeReq, saveName string, ak string) (string, error) {
	if saveName == "" {
		saveName = GetRandomName(16)
	}
	var apiURL string
	if ak == "" {
		token, err := GetMPAccessToken()
		if err != nil {
			return "", err
		}
		apiURL = fmt.Sprintf(mpQRCodeURL, token.AccessToken)
	} else {
		apiURL = fmt.Sprintf(mpQRCodeURL, ak)
	}

	post, _ := json.Marshal(req)
	resp, err := http.Post(apiURL, "application/json", bytes.NewReader(post))
	if err != nil {
		log.Println("请求微信获取二维码接口失败", err)
		return "", err
	}

	type result struct {
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
		ContentType string `json:"contentType"`
		Buffer      []byte `json:"buffer"`
	}
	var r result
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// 返回是二进制图片，或者json错误
	if resp.Header.Get("Content-Type") == "image/jpeg" || resp.Header.Get("Content-Type") == "image/png" {
		// 保存在output目录
		outputFileName := saveName

		if resp.Header.Get("Content-Type") == "image/jpeg" {
			outputFileName = outputFileName + ".jpg"
		} else {
			outputFileName = outputFileName + ".png"
		}
	here:
		f, err := os.OpenFile(C.OutputDIR+outputFileName, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			os.Mkdir(C.OutputDIR, 0666)
			goto here
		}
		f.Write(body)
		f.Close()
		return outputFileName, nil
	}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", err
	}
	if r.ErrCode != 0 {
		return "", errors.New(r.ErrMsg)
	}
	return "", nil
}
