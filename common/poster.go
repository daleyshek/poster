package common

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"unicode"

	"github.com/golang/freetype/truetype"

	"github.com/golang/freetype"

	"github.com/anthonynsimon/bild/transform"
)

const (
	// DesignWidth 设计宽度
	DesignWidth = 640
	// DesignHeight 设计高度
	DesignHeight = 1080
	// ImageHeight 图片区域640
	ImageHeight = 840
	// BodyHeight 信息区域240
	BodyHeight = 240
)

// ResourceDIR 资源目录，包括字体
var ResourceDIR string

func init() {
	ResourceDIR = "./resources/"
}

// Style 样式
type Style struct {
	ImageURL       string     //主图http链接或文件链接
	QRCodeURL      string     //二维码http链接或文件链接
	Title          string     //标题
	Content        string     //文字内容
	OutputFileName string     //保存文件名，可选
	OutputDIR      string     //保存文件目录，可选
	BorderColor    color.RGBA // 边框颜色，可选
}

var font *truetype.Font

// DrawPoster 绘制海报
func DrawPoster(s Style) (fileName string, err error) {
	// 填充默认数据
	if s.OutputDIR == "" {
		s.OutputDIR = C.OutputDIR
	}
	if s.ImageURL == "" || s.QRCodeURL == "" || s.Title == "" || s.Content == "" {
		return "", errors.New("样式参数缺失")
	}
	if s.OutputFileName == "" {
		s.OutputFileName = GetRandomName(16)
		s.OutputFileName = s.OutputFileName + ".jpg"
	}
	var bc color.RGBA
	if s.BorderColor == bc {
		s.BorderColor = color.RGBA{255, 255, 255, 255}
	}

	// 获取画布，并涂白
	rgba := image.NewRGBA(image.Rect(0, 0, DesignWidth, DesignHeight))
	for x := 0; x < rgba.Bounds().Dx(); x++ {
		for y := 0; y < rgba.Bounds().Dy(); y++ {
			rgba.Set(x, y, color.White)
		}
	}

	// 获取网络图片并解码
	picRd, err := getResourceReader(s.ImageURL)
	if err != nil {
		log.Println("未找到图片资源")
		return "", err
	}
	pic, _, err := image.Decode(picRd)
	if err != nil {
		log.Println("图片加载失败", err)
		return "", err
	}

	// 更改尺寸以适应画布大小
	var resizedw, resizedh int
	if DesignWidth/ImageHeight < pic.Bounds().Dx()/pic.Bounds().Dy() {
		// 拉伸至高度和设计高度一致
		resizedw = int(pic.Bounds().Dx() * ImageHeight / pic.Bounds().Dy())
		resizedh = ImageHeight
	} else {
		// 拉伸至宽度和设计宽度一致
		resizedh = int(pic.Bounds().Dy() * DesignWidth / pic.Bounds().Dx())
		resizedw = DesignWidth
	}
	picResized := transform.Resize(pic, resizedw, resizedh, transform.Linear)

	// 拉伸至中心完全显示
	draw.Draw(rgba, image.Rect(0, 0, DesignWidth, ImageHeight), picResized,
		image.Point{int((picResized.Bounds().Dx() - DesignWidth) / 2), int((picResized.Bounds().Dy() - ImageHeight) / 2)},
		draw.Src)

	// 绘制边框
	border, err := GetBorder("line")
	if err != nil {
		log.Println("获取边框失败", err)
		return "", err
	}
	p = Padding{20, 20, 20, 20}
	drawBorder(rgba, border, p, s.BorderColor)

	// 绘制文字
	texts := []string{}
	line := 0
	count := 0
	for _, x := range s.Content {
		if count > 24 { // 一行最大显示的字符数量
			line++
			count = 0
		}
		if len(texts) <= line {
			texts = append(texts, string(x))
		} else {
			texts[line] = texts[line] + string(x)
		}

		//汉字+2，字母+1
		if unicode.Is(unicode.Han, x) {
			count += 2
		} else {
			count++
		}
	}
	if len(texts) > 3 {
		texts = texts[:3]
		texts = append(texts, "...")
	} else {
		for l := len(texts); l < 4; l++ {
			texts = append(texts, "")
		}
	}
	texts = append(texts, s.Title)
	fontSize := 24.0
	lineHeight := 1.6
	// 加载文字字体
	if font == nil {
		fontFile, err := os.Open(ResourceDIR + "font.ttc")
		if err != nil {
			log.Println("字体加载失败", err)
			return "", err
		}
		fontBytes, err := ioutil.ReadAll(fontFile)
		if err != nil {
			log.Println("字体读取失败", err)
			return "", err
		}
		font, err = freetype.ParseFont(fontBytes)
		if err != nil {
			log.Println("解析字体失败", err)
			return "", err
		}
	}
	c := freetype.NewContext()
	c.SetFont(font)
	c.SetFontSize(fontSize)
	// 文字区域
	c.SetClip(image.Rect(p.Left, ImageHeight+p.Top, DesignWidth-200-40, DesignHeight-p.Bottom))
	c.SetDst(rgba)
	c.SetSrc(image.Black)

	pt := freetype.Pt(p.Left, ImageHeight+p.Top+int(c.PointToFixed(fontSize)>>6))
	for _, t := range texts {
		c.DrawString(t, pt)
		pt.Y += c.PointToFixed(fontSize * lineHeight)
	}

	// 绘制二维码
	qr, err := getResourceReader(s.QRCodeURL)
	if err != nil {
		log.Println("二维码资源获取失败", err)
		return "", err
	}
	qrcode, err := jpeg.Decode(qr)
	qrcodeResized := transform.Resize(qrcode, BodyHeight-40, BodyHeight-40, transform.Linear)
	draw.Draw(rgba,
		image.Rectangle{
			image.Point{int(DesignWidth/2) + int((DesignWidth/2-qrcodeResized.Bounds().Dx())/2), ImageHeight + 20},
			rgba.Bounds().Max,
		},
		qrcodeResized,
		image.Point{0, 0},
		draw.Src)

	// 保存
sv:
	img := rgba.SubImage(rgba.Bounds())
	f, err := os.OpenFile(s.OutputDIR+s.OutputFileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Println("未找到保存目录", err)
		log.Println("即将创建保存目录")
		os.Mkdir(s.OutputDIR, 0666)
		goto sv
	}

	err = jpeg.Encode(f, img, nil)
	if err != nil {
		fmt.Println("保存图片失败", err)
		return "", err
	}
	defer f.Close()
	return s.OutputFileName, nil
}

// Border 边框样式
type Border struct {
	Width  int         //设计的边框宽度
	Length int         //设计的边框长度
	Shape  *image.RGBA //形状，黑色代表形状占用空间点
}

// GetBorder 获取边框样式
// 边框是一段矩形，黑点值代表边框实际占用
func GetBorder(style string) (b Border, err error) {
	switch style {
	case "line":
		b.Width = 6
		b.Length = 80
		b.Shape = image.NewRGBA(image.Rect(0, 0, b.Length, b.Width))
		for x := 0; x < b.Length; x++ {
			for y := 0; y < b.Width; y++ {
				// 矩形边框
				b.Shape.SetRGBA(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
	case "wave":
		b.Width = 6
		b.Length = 80
		b.Shape = image.NewRGBA(image.Rect(0, 0, b.Length, b.Width))
		for x := 0; x < b.Length; x++ {
			for y := 0; y < b.Width; y++ {
				//正弦
				if float64(y) < (float64(b.Width)/1.732)*math.Sin(float64(x)*math.Pi/float64(b.Length)) {
					b.Shape.SetRGBA(x, y, color.RGBA{0, 0, 0, 255})
				}
			}
		}
	default:
		return b, errors.New("无效边框")
	}
	return b, nil
}

// Padding 补白尺寸
type Padding struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

var p Padding

func drawBorder(pic *image.RGBA, b Border, p Padding, c color.Color) {
	// 从底部中间区域开始绘制
	react := pic.Rect
	var midx, midy int
	midx = int(react.Dx() / 2)
	midy = int(react.Dy() / 2)
	black := color.RGBA{0, 0, 0, 255}

	// 上下
	for x, bx := midx, 0; x < react.Dx()-p.Right; x, bx = x+1, bx+1 {
		for y, by := p.Bottom, 0; y < react.Dy() && by < b.Width; y, by = y+1, by+1 {
			if b.Shape.At(bx%b.Length, by) == black {
				pic.Set(x, react.Dy()-y, c)
				pic.Set(react.Dx()-x, react.Dy()-y, c) //设置相反路径
				pic.Set(x, y, c)
				pic.Set(react.Dx()-x, y, c)
			}
		}
	}

	// 左右
	for y, bx := midy, 0; y < react.Dy()-p.Bottom; y, bx = y+1, bx+1 {
		for x, by := p.Top, 0; x < react.Dx()-p.Bottom && by < b.Width; x, by = x+1, by+1 {
			if b.Shape.At(bx%b.Length, by) == black {
				pic.Set(x, react.Dy()-y, c)
				pic.Set(react.Dx()-x, react.Dy()-y, c) //设置相反路径
				pic.Set(x, y, c)
				pic.Set(react.Dx()-x, y, c)
			}
		}
	}
}

func getResourceReader(src string) (r *bytes.Reader, err error) {
	if src[0:4] == "http" {
		resp, err := http.Get(src)
		if err != nil {
			return r, err
		}
		defer resp.Body.Close()
		fileBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return r, err
		}
		r = bytes.NewReader(fileBytes)
	} else {
		fileBytes, err := ioutil.ReadFile(src)
		if err != nil {
			return nil, err
		}
		r = bytes.NewReader(fileBytes)
	}
	return r, nil
}
