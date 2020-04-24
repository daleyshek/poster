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
	"math/rand"
	"net/http"
	"os"
	"time"
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
			rgba.Set(x, y, s.BorderColor)
			// rgba.Set(x, y, color.White)
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
	if float32(DesignWidth)/float32(ImageHeight) < float32(pic.Bounds().Dx())/float32(pic.Bounds().Dy()) {
		// 拉伸至高度和设计高度一致
		resizedw = pic.Bounds().Dx() * ImageHeight / pic.Bounds().Dy()
		resizedh = ImageHeight
	} else {
		// 拉伸至宽度和设计宽度一致
		resizedh = pic.Bounds().Dy() * DesignWidth / pic.Bounds().Dx()
		resizedw = DesignWidth
	}
	picResized := transform.Resize(pic, resizedw, resizedh, transform.Linear)

	// 拉伸至中心完全显示
	draw.Draw(rgba, image.Rect(0, 0, DesignWidth, ImageHeight), picResized,
		image.Point{(picResized.Bounds().Dx() - DesignWidth) / 2, (picResized.Bounds().Dy() - ImageHeight) / 2},
		draw.Src)

	// 绘制边框
	border, err := GetBorder()
	if err != nil {
		log.Println("获取边框失败", err)
		return "", err
	}
	p = Padding{20, 20, 20, 20} // w*240/1080
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
		fontFile, err := os.Open(C.FontFilePath)
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

	// 获取二维码并绘制
	qr, err := getResourceReader(s.QRCodeURL)
	if err != nil {
		log.Println("二维码资源获取失败", err)
		return "", err
	}
	qrcode, err := jpeg.Decode(qr)
	var qrcodeResized *image.RGBA
	if s.BorderColor.R != 255 && s.BorderColor.G != 255 && s.BorderColor.B != 255 {
		// 二维码扣白
		qrcodeRGBA := fillPicWhite(qrcode, s.BorderColor)
		qrcodeResized = transform.Resize(qrcodeRGBA.SubImage(qrcodeRGBA.Bounds()), BodyHeight-40, BodyHeight-40, transform.Linear)
	} else {
		qrcodeResized = transform.Resize(qrcode, BodyHeight-40, BodyHeight-40, transform.Linear)
	}

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
func GetBorder() (b Border, err error) {
	b.Width = 6
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
	react := pic.Rect
	for x := p.Left; x < react.Dx()-p.Right; x++ {
		for y := p.Top; y < react.Dy()-p.Bottom; y++ {
			if x-p.Left <= b.Width || react.Dx()-p.Right-x <= b.Width {
				if y < react.Dy()-BodyHeight {
					pic.Set(x, y, c)
				}
			} else if y-p.Top <= b.Width {
				pic.Set(x, y, c)
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

// GetRandomName 获取随机的文件名
func GetRandomName(length int) (name string) {
	dic := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	maxL := len([]byte(dic))
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		name = name + string(dic[rand.Intn(maxL)])
	}
	return name
}

// fillPicWhite 图片扣白
func fillPicWhite(pic image.Image, c color.Color) image.RGBA {
	rgba := image.NewRGBA(pic.Bounds())
	for x := 0; x < pic.Bounds().Dx(); x++ {
		for y := 0; y < pic.Bounds().Dy(); y++ {
			r, g, b, a := pic.At(x, y).RGBA()
			if r > 64535 && g > 64535 && b > 64535 && a == 65535 {
				rgba.Set(x, y, c)
			} else {
				rgba.Set(x, y, pic.At(x, y))
			}
		}
	}
	return *rgba
}
