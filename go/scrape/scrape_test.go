package scrape

import (
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/PuerkitoBio/goquery"
)

var kyukoDoc, noKyukoDoc *goquery.Document

const (
	KYUKOFILE   = "../testdata/kyuko.html"
	NOKYUKOFILE = "../testdata/not_kyuko.html"
)

func SjisToUtf8(str string) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder()))
	if err != nil {
		return "", err
	}
	return string(ret), err
}

func EncodeTestFile(fileName string) (io.Reader, error) {
	//testfileのenocde
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	utfFile, err := SjisToUtf8(string(file))
	if err != nil {
		return nil, err
	}
	stringReader := strings.NewReader(utfFile)

	return stringReader, nil
}

func init() {
	//休講ある
	kyukoReader, _ := EncodeTestFile(KYUKOFILE)
	kyukoDoc, _ = goquery.NewDocumentFromReader(kyukoReader)

	//休講ない
	noKyukoReader, _ := EncodeTestFile(NOKYUKOFILE)
	noKyukoDoc, _ = goquery.NewDocumentFromReader(noKyukoReader)
}

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
	periods, err := ScrapePeriod(kyukoDoc)
	if err != nil {
		t.Fatal("periodをスクレイピングできませんでした\n%s", err)
	}

	testSlice := []int{2, 2, 2, 5}
	if !reflect.DeepEqual(periods, testSlice) {
		t.Fatalf("取得した結果が求めるものと違ったようです\n want: %d\n got:  %d", testSlice, periods)
	}

	periods, err = ScrapePeriod(noKyukoDoc)
	if err != nil {
		t.Fatal("periodをスクレイピングできませんでした\n%s", err)
	}
	if len(periods) != 0 {
		t.Fatalf("取得した結果が求めるものと違ったようです\ngot:  %v", periods)
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
	reasons, err := ScrapeReason(kyukoDoc)
	if err != nil {
		t.Fatalf("reasonをスクレイピングできませんでした\n%s", err)
	}

	testSlice := []string{"公務", "出張", "公務", ""}
	if !reflect.DeepEqual(reasons, testSlice) {
		t.Fatalf("取得した結果が求めるものと違ったようです\n want: %v\n got:  %v", testSlice, reasons)
	}

	reasons, err = ScrapeReason(noKyukoDoc)
	if err != nil {
		t.Fatal("periodをスクレイピングできませんでした\n%s", err)
	}
	if len(reasons) != 0 {
		t.Fatalf("取得した結果が求めるものと違ったようです\ngot:  %v", reasons)
	}

}

func TestScrapeNameAndInstructor(t *testing.T) {
	names, instructors, err := ScrapeNameAndInstructor(kyukoDoc)
	if err != nil {
		t.Fatalf("Nameのスクレイピングに失敗したようです\n%s", err)
	}

	testSlice := []string{"環境生理学", "電気・電子計測Ｉ－１", "応用数学ＩＩ－１", "イングリッシュ・セミナー２－７０２"}
	if !reflect.DeepEqual(names, testSlice) {
		t.Fatalf("取得した結果が求めるものと違ったようです\n want: %v\n got:  %v", testSlice, names)
	}

	testSlice = []string{"福岡義之", "松川真美", "大川領", "稲垣俊史"}
	if !reflect.DeepEqual(instructors, testSlice) {
		t.Fatalf("取得した結果が求めるものと違ったようです\n want: %v\n got:  %v", testSlice, instructors)
	}

	names, instructors, err = ScrapeNameAndInstructor(noKyukoDoc)
	if err != nil {
		t.Fatal("periodをスクレイピングできませんでした\n%s", err)
	}
	if len(names) != 0 {
		t.Fatalf("取得した結果が求めるものと違ったようです\ngot:  %v", names)
	}
	if len(instructors) != 0 {
		t.Fatalf("取得した結果が求めるものと違ったようです\ngot:  %v", instructors)
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
