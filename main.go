package main

import (
	"flag"
	"fmt"
	"os"
	"poster/common"
)

func main() {
	// common.ServeJSONRPCOverHTTP()
	// common.ServeJSONRPC()
	cli()
	common.ServeJSONRPCOverHTTP()
	select {}
}

func cli() {
	cliMode := flag.Bool("cli", false, "是否使用cli mode")
	imageURL := flag.String("imageURL", "", "图片资源")
	title := flag.String("title", "", "标题")
	content := flag.String("content", "", "内容")

	scene := flag.String("scene", "", "场景值")
	page := flag.String("page", "", "小程序页面")
	width := flag.Int("width", 480, "小程序码的尺寸")
	flag.Parse()

	if *cliMode == false {
		return
	}

	qr := common.QRCodeReq{Scene: *scene, Page: *page, Width: *width}
	// 获取小程序码保存到本地
	qrcodeName, err := common.RequestQRCode(qr, "")
	if err != nil {
		fmt.Println("获取小程序码失败", err)
		os.Exit(1)
	}

	// 生成海报
	style := common.Style{
		ImageURL:  *imageURL,
		QRCodeURL: common.C.OutputDIR + qrcodeName,
		Title:     *title,
		Content:   *content,
	}
	postName, err := common.DrawPoster(style)
	if err != nil {
		fmt.Println("海报生成失败", err)
		os.Exit(1)
	}
	fmt.Println("运行成功")
	fmt.Println("小程序码文件", qrcodeName)
	fmt.Println("海报文件", postName)
	os.Exit(0)
}
