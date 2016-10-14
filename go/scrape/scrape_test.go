package scrape

import (
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/g-hyoga/kyuko/go/model"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var kyukoDoc, noKyukoDoc *goquery.Document
var testPeriods []int
var testReasons, testNames, testInstructors []string
var testPlace, testWeekday int
var testDay string
var testData []model.KyukoData
var noTestPlace, noTestWeekday int
var noTestDay string
var noTestData []model.KyukoData

const (
	KYUKOFILE   = "../testdata/kyuko.html"
	NOKYUKOFILE = "../testdata/not_kyuko.html"
)

func init() {
	//休講ある
	kyukoReader, _ := EncodeTestFile(KYUKOFILE)
	kyukoDoc, _ = goquery.NewDocumentFromReader(kyukoReader)

	testPeriods = []int{2, 2, 2, 5}
	testReasons = []string{"公務", "出張", "公務", ""}
	testNames = []string{"環境生理学", "電気・電子計測Ｉ－１", "応用数学ＩＩ－１", "イングリッシュ・セミナー２－７０２"}
	testInstructors = []string{"福岡義之", "松川真美", "大川領", "稲垣俊史"}
	testPlace = 2
	testDay = "2016/10/10"
	testWeekday = 1

	for i := range testPeriods {
		k := model.KyukoData{}
		k.Period = testPeriods[i]
		k.Reason = testReasons[i]
		k.ClassName = testNames[i]
		k.Instructor = testInstructors[i]
		k.Weekday = testWeekday
		k.Place = testPlace
		k.Day = testDay
		testData = append(testData, k)
	}

	//休講ない
	noKyukoReader, _ := EncodeTestFile(NOKYUKOFILE)
	noKyukoDoc, _ = goquery.NewDocumentFromReader(noKyukoReader)

	noTestPlace = 1
	noTestWeekday = 6
}

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

//////////////////////////
func TestGet(t *testing.T) {

}

func TestScrapePeriod(t *testing.T) {
	periods, err := ScrapePeriod(kyukoDoc)
	if err != nil {
		t.Fatal("periodをスクレイピングできませんでした\n%s", err)
	}

	if !reflect.DeepEqual(periods, testPeriods) {
		t.Fatalf("取得した結果が求めるものと違ったようです\n want: %d\n got:  %d", testPeriods, periods)
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

	if !reflect.DeepEqual(reasons, testReasons) {
		t.Fatalf("取得した結果が求めるものと違ったようです\n want: %v\n got:  %v", testReasons, reasons)
	}

	reasons, err = ScrapeReason(noKyukoDoc)
	if err != nil {
		t.Fatal("periodをスクレイピングできませんでした\n%s", err)
	}
	if len(reasons) != 0 {
		t.Fatalf("取得した結果が求めるものと違ったようです\ngot:  %v", reasons)
	}

}

func TestScrapeDay(t *testing.T) {
	day, err := ScrapeDay(kyukoDoc)
	if err != nil {
		t.Fatalf("日付を取得できませんでした\n%s", err)
	}

	if day != testDay {
		t.Fatalf("取得した結果が求めるものと違ったようです\nwant: %v\ngot:  %v", testDay, day)
	}

}

func TestScrapePlace(t *testing.T) {
	place, err := ScrapePlace(kyukoDoc)
	if err != nil {
		t.Fatalf("placeを取得できませんでした\n%s", err)
	}

	if place != testPlace {
		t.Fatalf("取得した結果が求めるものと違ったようです\nwant: %v\ngot:  %v", testPlace, place)
	}

	place, err = ScrapePlace(noKyukoDoc)
	if err != nil {
		t.Fatalf("placeを取得できませんでした\n%s", err)
	}

	if place != noTestPlace {
		t.Fatalf("取得した結果が求めるものと違ったようです\nwant: %v\ngot:  %v", noTestPlace, place)
	}

}

func TestScrapeWeeday(t *testing.T) {
	weekday, err := ScrapeWeekday(kyukoDoc)
	if err != nil {
		t.Fatalf("曜日をスクレイピングできませんでした\n%s", err)
	}

	if weekday != testWeekday {
		t.Fatalf("取得した結果が求めるものと違ったようです\nwant: %v\ngot:  %v", testWeekday, weekday)
	}

	weekday, err = ScrapeWeekday(noKyukoDoc)
	if err != nil {
		t.Fatalf("曜日をスクレイピングできませんでした\n%s", err)
	}

	if weekday != noTestWeekday {
		t.Fatalf("取得した結果が求めるものと違ったようです\nwant: %v\ngot:  %v", noTestWeekday, weekday)
	}

}

func TestScrapeNameAndInstructor(t *testing.T) {
	names, instructors, err := ScrapeNameAndInstructor(kyukoDoc)
	if err != nil {
		t.Fatalf("Nameのスクレイピングに失敗したようです\n%s", err)
	}

	if !reflect.DeepEqual(names, testNames) {
		t.Fatalf("取得した結果が求めるものと違ったようです\n want: %v\n got:  %v", testNames, names)
	}

	if !reflect.DeepEqual(instructors, testInstructors) {
		t.Fatalf("取得した結果が求めるものと違ったようです\n want: %v\n got:  %v", testInstructors, instructors)
	}

	names, instructors, err = ScrapeNameAndInstructor(noKyukoDoc)
	if err != nil {
		t.Fatalf("periodをスクレイピングできませんでした\n%s", err)
	}
	if len(names) != 0 {
		t.Fatalf("取得した結果が求めるものと違ったようです\ngot:  %v", names)
	}
	if len(instructors) != 0 {
		t.Fatalf("取得した結果が求めるものと違ったようです\ngot:  %v", instructors)
	}

}

func TestScrape(t *testing.T) {
	allData, err := Scrape(kyukoDoc)
	if err != nil {
		t.Fatal("scrapingに失敗しました\n%s", err)
	}

	if !reflect.DeepEqual(allData, testData) {
		t.Fatalf("取得した結果が求めるものと違ったようです\n want: %v\n got:  %v", testData, allData)
	}

	allData, err = Scrape(noKyukoDoc)
	if err != nil {
		t.Fatal("scrapingに失敗しました\n%s", err)
	}
}
