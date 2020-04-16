package common

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// Poster 海报生成
type Poster int

// Resp 回应
// type Resp string

// Qrcode 生成小程序二维码
func (p *Poster) Qrcode(req *QRCodeReq, r *string) error {
	fn, err := RequestQRCode(*req, "")
	if err != nil {
		return err
	}
	*r = fn
	return nil
}

// Generate 海报的生成
func (p *Poster) Generate(s *Style, r *string) error {
	fn, err := DrawPoster(*s)
	if err != nil {
		return err
	}
	*r = fn
	return nil
}

// BothReq 两个配置合并
type BothReq struct {
	QRCodeReq
	Style
}

// BothRes 返回两个文件名
type BothRes struct {
	QRCodeName string
	PosterName string
}

// Both 先生成二维码，再由二维码生成海报
func (p *Poster) Both(req *BothReq, r *BothRes) error {
	// 获取小程序码保存到本地
	qrcodeName, err := RequestQRCode((*req).QRCodeReq, "")
	if err != nil {
		log.Println("获取小程序码失败", err)
		return err
	}

	// 使用生成的二维码来生成海报
	(*req).Style.QRCodeURL = C.OutputDIR + qrcodeName
	// 生成海报
	posterName, err := DrawPoster((*req).Style)
	if err != nil {
		fmt.Println("海报生成失败", err)
		return err
	}

	(*r).QRCodeName = qrcodeName
	(*r).PosterName = posterName
	return nil
}

// ServeRPC 运行RPC server
func ServeRPC() {
	p := new(Poster)
	rpc.Register(p)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":2019")
	if err != nil {
		log.Fatal(err)
	}
	go http.Serve(l, nil)
}

// ServeJSONRPC json-rpc server over tcp
func ServeJSONRPC() {
	p := new(Poster)
	rpc.Register(p)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":2019")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			go jsonrpc.ServeConn(conn)
		}
	}()
}

// ServeJSONRPCOverHTTP json-rpc over http
func ServeJSONRPCOverHTTP() {
	server := rpc.NewServer()
	p := new(Poster)
	server.Register(p)

	l, err := net.Listen("tcp", ":2019")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCodec := jsonrpc.NewServerCodec(&HttpConn{in: r.Body, out: w})
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(200)
		err := server.ServeRequest(serverCodec)
		if err != nil {
			log.Printf("Error while serving JSON request: %v", err)
			http.Error(w, "Error while serving JSON request, details have been logged.", 500)
			return
		}
	}))
}

type HttpConn struct {
	in  io.Reader
	out io.Writer
}

func (c *HttpConn) Read(p []byte) (n int, err error)  { return c.in.Read(p) }
func (c *HttpConn) Write(d []byte) (n int, err error) { return c.out.Write(d) }
func (c *HttpConn) Close() error                      { return nil }
