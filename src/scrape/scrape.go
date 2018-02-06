package scrape

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/text/unicode/norm"

	"github.com/PuerkitoBio/goquery"
	"github.com/g-hyoga/kyuko/go/src/model"
)

var stringCleaner *strings.Replacer

func init() {
	stringCleaner = strings.NewReplacer(" ", "", "\n", "", "\u00a0", "")
}

//place(1: 今出川 ,2: 京田辺), week(1 ~ 6: Mon ~ Sat)を引数に持ち
//urlを生成する
func SetUrl(place int, isTommorow bool) (string, error) {
	url := "https://duet.doshisha.ac.jp/kokai/html/fi/fi050/FI05001G.html"
	if isTommorow {
		url = "https://duet.doshisha.ac.jp/kokai/html/fi/fi050/FI05001G_02.html"
	}
	return url, nil
}

func Get(url string) (io.Reader, error) {
	var err error

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body := bytes.Buffer{}
	body.ReadFrom(resp.Body)

	reader := bytes.NewReader(body.Bytes())
	return reader, nil
}

func GetPlaceComponent(doc *goquery.Document, place int) *goquery.Document {
	node := doc.Find(".data").Get(place)
	return goquery.NewDocumentFromNode(node)
}

func GetKyukoTobaleLine(doc *goquery.Document, place int) [][]string {
	var lines [][]string
	GetPlaceComponent(doc, place).Find("tr").Each(func(i int, s *goquery.Selection) {
		if i != 0 {
			var elements []string
			s.Find("td").Each(func(j int, e *goquery.Selection) {
				elements = append(elements, e.Text())
			})
			lines = append(lines, elements)
		}
	})
	return lines
}

func ScrapePeriod(doc *goquery.Document, place int) ([]int, error) {
	var periods []int
	var err error

	lines := GetKyukoTobaleLine(doc, place)

	for _, line := range lines {
		stringPeriod := strings.Split(line[0], "講時")[0]
		stringPeriod = strings.Replace(stringPeriod, "\n", "", -1)
		stringPeriod = string(norm.NFKC.Bytes([]byte(stringPeriod)))
		period, _ := strconv.Atoi(stringPeriod)

		if period > 7 || period < 1 {
			err = errors.New("period is not found")
		}
		periods = append(periods, period)
	}
	return periods, err
}

func ScrapeReason(doc *goquery.Document, place int) ([]string, error) {
	var reasons []string

	lines := GetKyukoTobaleLine(doc, place)

	for _, line := range lines {
		reason := line[3]
		reasons = append(reasons, reason)
	}

	return reasons, nil
}

func ScrapeName(doc *goquery.Document, place int) ([]string, error) {
	var names []string
	lines := GetKyukoTobaleLine(doc, place)
	for _, line := range lines {
		name := strings.Replace(line[1], "\n", "", -1)
		name = strings.Replace(line[1], " ", "", -1)
		names = append(names, name)
	}
	return names, nil
}

func ScrapeInstructor(doc *goquery.Document, place int) ([]string, error) {
	var names []string
	lines := GetKyukoTobaleLine(doc, place)
	for _, line := range lines {
		name := strings.Replace(line[2], "\n", "", -1)
		name = strings.Replace(line[2], "　", "", -1)
		names = append(names, name)
	}
	return names, nil
}

func ScrapeNameAndInstructor(doc *goquery.Document, place int) (names, instructors []string, err error) {
	names, err = ScrapeName(doc, place)
	instructors, err = ScrapeInstructor(doc, place)
	return names, instructors, err
}

func ScrapeWeekday(doc *goquery.Document) (int, error) {
	weekday := doc.Find(".today").Text()
	weekday = strings.Split(weekday, "(")[1]
	weekday = strings.Replace(weekday, ")", "", -1)
	weekday = stringCleaner.Replace(weekday)
	youbi, err := ConvertWeekStoi(weekday)

	if err != nil {
		return -1, err
	}
	return youbi, nil
}

func ScrapeDay(doc *goquery.Document) (string, error) {
	day := doc.Find("#form1 > span").Text()
	day = strings.Split(day, " ")[0]
	day = stringCleaner.Replace(day)
	year := strings.Split(day, "年")[0]

	dateNode := doc.Find(".cancel_class span").First().Text()
	month := strings.Split(dateNode, "/")[0]
	if m, _ := strconv.Atoi(month); m < 10 {
		month = "0" + month
	}
	date := strings.Split(strings.Split(dateNode, "/")[1], "(")[0]

	return string(year) + "/" + string(month) + "/" + string(date), nil
}

func ConvertWeekStoi(weekday string) (int, error) {
	weekMap := map[string]int{"日": 0, "月": 1, "火": 2, "水": 3, "木": 4, "金": 5, "土": 6}

	if _, ok := weekMap[weekday]; !ok {
		return -1, errors.New("存在しない曜日が入力されています")
	}

	return weekMap[weekday], nil
}

//休講structのsliceを返す
func Scrape(doc *goquery.Document, place int) ([]model.KyukoData, error) {
	var kyukoData []model.KyukoData
	var err error

	var periods []int
	var reasons, names, instructors []string
	var weekday int
	var day string

	periods, err = ScrapePeriod(doc, place)
	reasons, err = ScrapeReason(doc, place)
	names, instructors, err = ScrapeNameAndInstructor(doc, place)
	weekday, err = ScrapeWeekday(doc)
	day, err = ScrapeDay(doc)

	if err != nil {
		return nil, err
	}

	if len(periods) != len(reasons) && len(periods) != len(names) && len(periods) != len(instructors) {
		return nil, errors.New("取得できていない情報があります")
	}

	for i := range periods {
		k := model.KyukoData{}
		k.Period = periods[i]
		k.Reason = reasons[i]
		k.ClassName = names[i]
		k.Instructor = instructors[i]
		k.Weekday = weekday
		k.Place = place
		k.Day = day
		kyukoData = append(kyukoData, k)
	}

	return kyukoData, err
}
