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
	//休講ある
	file, err := ioutil.ReadFile("../testdata/kyuko.html")
	if err != nil {
		t.Fatalf("テストデータを開けませんでした\n%s", err)
	}
	utfFile, err := SjisToUtf8(string(file))
	if err != nil {
		t.Fatalf("文字コードの変換に失敗しました\n%s", err)
	}
	stringReader := strings.NewReader(utfFile)

	doc, err := goquery.NewDocumentFromReader(stringReader)
	if err != nil {
		t.Fatalf("テストデータを開けませんでした\n%s", err)
	}

	//休講ない
	noKyukoFile, err := ioutil.ReadFile("../testdata/not_kyuko.html")
	if err != nil {
		t.Fatalf("テストデータを開けませんでした\n%s", err)
	}
	utfNoKyukoFile, err := SjisToUtf8(string(noKyukoFile))
	if err != nil {
		t.Fatalf("文字コードの変換に失敗しました\n%s", err)
	}

	noKyukoStringReader := strings.NewReader(string(utfNoKyukoFile))

	noKyukoDoc, err := goquery.NewDocumentFromReader(noKyukoStringReader)
	if err != nil {
		t.Fatalf("テストデータを開けませんでした\n%s", err)
	}

	//test
	periods, err := ScrapePeriod(doc)
	if err != nil {
		t.Fatal("periodをスクレイピングできませんでした\n%s", err)
	}

	testSlice := []int{2, 2, 2, 5}
	if reflect.DeepEqual(periods, testSlice) {
		t.Fatalf("取得した結果が求めるものと違ったようです\n want: %d\n got:  %d", testSlice, periods)
	}

	periods, err = ScrapePeriod(noKyukoDoc)
	if err != nil {
		t.Fatal("periodをスクレイピングできませんでした\n%s", err)
	}
	testSlice = []int{}
	if !reflect.DeepEqual(periods, testSlice) {
		t.Fatalf("取得した結果が求めるものと違ったようです\n want: %v\n got:  %v", testSlice, periods)
	}
}

func TestScrapeReason(t *testing.T) {

	/*
		//httpでやるとき
			stringReader, err := Get("http://duet.doshisha.ac.jp/info/KK1000.jsp?katei=1&youbi=2&kouchi=2")
			if err != nil {
				t.Fatal("hoge\n%v", err)
			}
	*/

	//testfileのenocde
	file, err := ioutil.ReadFile("../testdata/kyuko.html")
	if err != nil {
		t.Fatalf("テストデータを開けませんでした\n%s", err)
	}
	utfFile, err := SjisToUtf8(string(file))
	if err != nil {
		t.Fatalf("文字コードの変換に失敗しました\n%s", err)
	}
	stringReader := strings.NewReader(utfFile)

	doc, err := goquery.NewDocumentFromReader(stringReader)
	if err != nil {
		t.Fatalf("テストデータを開けませんでした\n%s", err)
	}

	reasons, err := ScrapeReason(doc)
	if err != nil {
		t.Fatalf("reasonをスクレイピングできませんでした\n%s", err)
	}

	testSlice := []string{"公務", "出張", "公務", ""}
	if !reflect.DeepEqual(reasons, testSlice) {
		t.Fatalf("取得した結果が求めるものと違ったようです\n want: %v\n got:  %v", testSlice, reasons)
	}
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
