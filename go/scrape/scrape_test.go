package scrape

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestSetUrl(t *testing.T) {
	if url, err := SetUrl(1, 1); url != "http://duet.doshisha.ac.jp/info/KK1000.jsp?katei=1&youbi=1&kouchi=1" {
		t.Fatalf("urlの生成がうまくできていないようです\n err: %s", err)
	}

	if url, err := SetUrl(2, 5); url != "http://duet.doshisha.ac.jp/info/KK1000.jsp?katei=1&youbi=5&kouchi=2" || err != nil {
		t.Fatalf("urlの生成がうまくできていないようです\n err: %s", err)
	}

	if url, err := SetUrl(3, 1); err == nil {
		t.Fatalf("存在しない校地のurlが生成されています\n created url: %s", url)
	}

	if url, err := SetUrl(1, 7); err == nil {
		t.Fatalf("日曜日のurlは必要ありません\n created url: %s", url)
	}

}

func TestScrapePeriod(t *testing.T) {
	file, err := ioutil.ReadFile("../testdata/kyuko.html")
	if err != nil {
		t.Fatalf("テストデータを開けませんでした\n%s", err)
	}
	stringReader := strings.NewReader(string(file))

	doc, err := goquery.NewDocumentFromReader(stringReader)
	if err != nil {
		t.Fatalf("テストデータを開けませんでした\n%s", err)
	}

	/*
		doc, err := goquery.NewDocument("http://duet.doshisha.ac.jp/info/KK1000.jsp?katei=1&youbi=4&kouchi=2")
		if err != nil {
			t.Fatalf("テストデータを開けませんでした\n%s", err)
		}
	*/

	periods, err := ScrapePeriod(doc)
	if err != nil {
		t.Fatal("periodをスクレイピングできませんでした\n%s", err)
	}

	testSlice := []int{2, 2, 2, 5}
	if reflect.DeepEqual(periods, testSlice) {
		t.Fatalf("取得した結果が求めるものと違ったようです\n want: %d\n got:  %d", testSlice, periods)
	}
}

func TestScrapeReason(t *testing.T) {
	file, err := ioutil.ReadFile("../testdata/kyuko.html")
	if err != nil {
		t.Fatalf("テストデータを開けませんでした\n%s", err)
	}
	stringReader := strings.NewReader(string(file))

	doc, err := goquery.NewDocumentFromReader(stringReader)
	if err != nil {
		t.Fatalf("テストデータを開けませんでした\n%s", err)
	}

	reasons, err := ScrapeReason(doc)
	if err != nil {
		t.Fatalf("reasonをスクレイピングできませんでした\n%s", err)
	}
	fmt.Printf("%s", reasons)

}

//まだできてない
func testScrape(t *testing.T) {

	//r, err := Scrape("http://duet.doshisha.ac.jp/info/KK1000.jsp?katei=1&youbi=4&kouchi=2")

	/*
		file, err := os.Open("../testdata/kyuko.html")
		if err != nil {
			t.Fatalf("テストデータを開けませんでした\n%s", err)
		}
		defer file.Close()

		r, err := Scrape("", file)
		if err != nil {
			t.Fatalf("hoge\n%s", err)
		}

		fmt.Printf("%d\nhoge\n%d", r, err)
	*/
}
