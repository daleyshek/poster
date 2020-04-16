package tests

import (
	"net/rpc"
	"net/rpc/jsonrpc"
	"testing"

	"github.com/daleyshek/poster/common"
)

func testRPC(t *testing.T) {
	common.ServeRPC()

	client, err := rpc.DialHTTP("tcp", "127.0.0.1:2019")
	if err != nil {
		t.Fatal("client 连接失败", err)
	}
	s := common.Style{}
	s.Title = "BestFriendsChina"
	s.Content = "凡读书......须要读得字字响亮，不可误一字，不可少一字，不可多一字，不可倒一字，不可牵强暗记，只是要多诵数遍，自然上口，久远不忘。古人云，“读书百遍，其义自见”。谓读得熟，则不待解说，自晓其义也。余尝谓，读书有三到，谓心到，眼到，口到。心不在此，则眼不看仔细，心眼既不专一，却只漫浪诵读，决不能记，记亦不能久也。三到之中，心到最急。心既到矣，眼口岂不到乎？"
	s.ImageURL = "https://api.mcdsh.com/storage/images/800/g45rDYNgaI3z0ZACqzkI0iysuXIz4omyBhZSGBUM.jpeg"
	s.QRCodeURL = "../resources/qrcode.jpg"

	var reply string
	err = client.Call("Poster.Generate", &s, &reply)
	if err != nil {
		t.Fatal("调用失败", err)
	} else {
		t.Log("生成的图片为", reply)
	}
}

func TestJSONRPC(t *testing.T) {
	common.ServeJSONRPC()

	client, err := jsonrpc.Dial("tcp", "127.0.0.1:2019")
	if err != nil {
		t.Fatal("client 连接失败", err)
	}
	s := common.Style{}
	s.Title = "BestFriendsChina"
	s.Content = "凡读书......须要读得字字响亮，不可误一字，不可少一字，不可多一字，不可倒一字，不可牵强暗记，只是要多诵数遍，自然上口，久远不忘。古人云，“读书百遍，其义自见”。谓读得熟，则不待解说，自晓其义也。余尝谓，读书有三到，谓心到，眼到，口到。心不在此，则眼不看仔细，心眼既不专一，却只漫浪诵读，决不能记，记亦不能久也。三到之中，心到最急。心既到矣，眼口岂不到乎？"
	s.ImageURL = "https://api.mcdsh.com/storage/images/800/g45rDYNgaI3z0ZACqzkI0iysuXIz4omyBhZSGBUM.jpeg"
	s.QRCodeURL = "../resources/qrcode.jpg"
	var reply string
	err = client.Call("Poster.Generate", &s, &reply)
	if err != nil {
		t.Fatal("调用失败", err)
	} else {
		t.Log("生成的图片为", reply)
	}
}

func init() {
	common.OutPutDIR = "../output/"
	common.ResourceDIR = "../resources/"
}
