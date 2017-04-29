package scrape

import (
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/g-hyoga/kyuko/go/model"
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
	KYUKOFILE   = "../testdata/new_kyuko.html"
	NOKYUKOFILE = "../testdata/new_not_kyuko.html"
)

func init() {
	//休講ある
	kyukoFile, _ := ioutil.ReadFile(KYUKOFILE)
	kyukoReader := strings.NewReader(string(kyukoFile))
	kyukoDoc, _ = goquery.NewDocumentFromReader(kyukoReader)

	testPeriods = []int{4, 5}
	testReasons = []string{"公務", "公務"}
	testNames = []string{"情報処理演習", "精神病理学"}
	testInstructors = []string{"田中宏季", "佃宗紀"}
	testPlace = 2
	testDay = "2017/04/14"
	testWeekday = 4

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
	noKyukoFile, _ := ioutil.ReadFile(NOKYUKOFILE)
	noKyukoReader := strings.NewReader(string(noKyukoFile))
	noKyukoDoc, _ = goquery.NewDocumentFromReader(noKyukoReader)

	noTestPlace = 1
	noTestWeekday = 6
}

func TestSetUrl(t *testing.T) {
	if url, err := SetUrl(1, true); url != "https://duet.doshisha.ac.jp/kokai/html/fi/fi050/FI05001G_02.html" {
		t.Fatalf("明日のurlになっていません\n err: %s", err)
	}

	if url, err := SetUrl(1, false); url != "https://duet.doshisha.ac.jp/kokai/html/fi/fi050/FI05001G.html" || err != nil {
		t.Fatalf("今日のurlになっていません\n err: %s", err)
	}
}

//////////////////////////
func TestGet(t *testing.T) {

}

func TestScrapePeriod(t *testing.T) {
	periods, err := ScrapePeriod(kyukoDoc, 2)
	if err != nil {
		t.Fatalf("periodをスクレイピングできませんでした\n%s", err)
	}

	if !reflect.DeepEqual(periods, testPeriods) {
		t.Fatalf("取得した結果が求めるものと違ったようです\n want: %d\n got:  %d", testPeriods, periods)
	}

	periods, err = ScrapePeriod(noKyukoDoc, 2)
	if err != nil {
		t.Fatalf("periodをスクレイピングできませんでした\n%s", err)
	}
	if len(periods) != 0 {
		t.Fatalf("取得した結果が求めるものと違ったようです\ngot:  %v", periods)
	}
}

func TestScrapeReason(t *testing.T) {
	reasons, err := ScrapeReason(kyukoDoc, 2)
	if err != nil {
		t.Fatalf("reasonをスクレイピングできませんでした\n%s", err)
	}

	if !reflect.DeepEqual(reasons, testReasons) {
		t.Fatalf("取得した結果が求めるものと違ったようです\n want: %v\n got:  %v", testReasons, reasons)
	}

	reasons, err = ScrapeReason(noKyukoDoc, 2)
	if err != nil {
		t.Fatalf("periodをスクレイピングできませんでした\n%s", err)
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
	names, instructors, err := ScrapeNameAndInstructor(kyukoDoc, 2)
	if err != nil {
		t.Fatalf("Nameのスクレイピングに失敗したようです\n%s", err)
	}

	if !reflect.DeepEqual(names, testNames) {
		t.Fatalf("取得した結果が求めるものと違ったようです\n want: %v\n got:  %v", testNames, names)
	}

	if !reflect.DeepEqual(instructors, testInstructors) {
		t.Fatalf("取得した結果が求めるものと違ったようです\n want: %v\n got:  %v", testInstructors, instructors)
	}

	names, instructors, err = ScrapeNameAndInstructor(noKyukoDoc, 1)
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
	allData, err := Scrape(kyukoDoc, 1)
	if err != nil {
		t.Fatalf("scrapingに失敗しました\n%s", err)
	}

	if !reflect.DeepEqual(allData, testData) {
		t.Fatalf("取得した結果が求めるものと違ったようです\n want: %v\n got:  %v", testData, allData)
	}
}

func TestScrapeWithNoClass(t *testing.T) {
	result, err := Scrape(noKyukoDoc, 1)
	if err != nil {
		t.Fatalf("scrapingに失敗しました\n%s", err)
	}

	if !reflect.DeepEqual(result, noTestData) {
		t.Fatalf("取得した結果が求めるものと違ったようです\n want: %v\n got:  %v", noTestData, result)
	}
}

func BenchmarkScrape(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Scrape(kyukoDoc, 1)
	}
}

func BenchmarkSetUrl(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SetUrl(1, true)
	}
}

func BenchmarkScrapeDay(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ScrapeDay(kyukoDoc)
	}
}

func BenchmarkScrapeNameAndInstructor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ScrapeNameAndInstructor(kyukoDoc, 1)
	}
}

func BenchmarkScrapePeriod(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ScrapePeriod(kyukoDoc, 1)
	}
}

func BenchmarkScrapePlace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ScrapePeriod(kyukoDoc, 1)
	}
}

func BenchmarkScrapeReason(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ScrapeReason(kyukoDoc, 1)
	}
}

func BenchmarkScrapeWeekday(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ScrapeWeekday(kyukoDoc)
	}
}
