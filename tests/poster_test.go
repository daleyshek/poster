package tests

import (
	"fmt"
	"poster/common"
	"strconv"
	"testing"
)

func TestPoster(t *testing.T) {
	pics := []string{
		"https://api.mcdsh.com/storage/images/800/6WCihrJGNKzOA5cbZvIpSKZbjsckmwqA3sSVjQ7n.jpeg",
		"https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=3827272599,4144405931&fm=27&gp=0.jpg",
		"https://api.mcdsh.com/storage/images/8qjKVNsh0J9RxriNj4lD7Wgv5GLBfmrNOyWmp7IK.jpeg",
	}

	for i := range pics {
		s := common.Style{}
		s.ImageURL = pics[i]
		s.Title = "BestFriendsChina"
		s.Content = "谓读得熟，则不待解说，自晓其义也。余尝谓，读书有三到，谓心到，眼到，口到。心不在此，则眼不看仔细，心眼既不专一，却只漫浪诵读，决不能记，记亦不能久也。三到之中，心到最急。心既到矣，眼口岂不到乎？"
		s.QRCodeURL = "resources/qrcode.jpg"
		s.OutputFileName = "test_" + strconv.Itoa(i) + ".jpg"
		_, err := common.DrawPoster(s)
		if err == nil {
			fmt.Print(".")
		} else {
			t.Error(err)
		}
	}
}
