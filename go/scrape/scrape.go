package scrape

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/g-hyoga/kyuko/go/model"
	"golang.org/x/text/unicode/norm"
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

/*
func SetUrl(place, week int) (string, error) {
	//url := "http://duet.doshisha.ac.jp/info/KK1000.jsp?katei=1"
	//weekに7(Sunday)はない
	if (place != 1 && place != 2) || week < 1 || week > 6 {
		return "", errors.New("place is 1 or 2, 0 < week < 7")
	} else {
		url = url + "&youbi=" + strconv.Itoa(week)
		url = url + "&kouchi=" + strconv.Itoa(place)
		return url, nil
	}
}
*/

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
	//utfBody := transform.NewReader(reader, japanese.ShiftJIS.NewDecoder())

	return reader, nil
}

func GetPlaceComponent(doc *goquery.Document, place int) *goquery.Document {
	node := doc.Find(".data").Get(place)
	return goquery.NewDocumentFromNode(node)
}

func GetKyukoTobaleLine(doc *goquery.Document, place int) [][]string {
	var lines [][]string
	GetPlaceComponent(doc, 1).Find("tr").Each(func(i int, s *goquery.Selection) {
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
	month := strings.Split(strings.Split(day, "年")[1], "月")[0]
	date := strings.Split(strings.Split(day, "日")[0], "月")[1]

	return string(year) + "/" + string(month) + "/" + string(date), nil
}

/* Duetのversion変化に伴って更新
func ScrapePeriod(doc *goquery.Document) ([]int, error) {
	var periods []int
	var err error

	//エラー処理どうにかする
	//"1講時"みたいなのが取れる
	doc.Find("tr.").Each(func(i int, s *goquery.Selection) {
		originalPeriod := s.Find("th.style2").Text()

		stringPeriod := strings.Split(originalPeriod, "講時")[0]
		stringPeriod = strings.Replace(stringPeriod, "\n", "", -1)
		period, _ := strconv.Atoi(stringPeriod)

		if period == 0 && i != 0 {
			period = periods[i-1]
		}
		periods = append(periods, period)

	})
	return periods, err
}

func ScrapeReason(doc *goquery.Document) ([]string, error) {
	var reasons []string
	var err error

	doc.Find("tr.style1").Each(func(i int, s *goquery.Selection) {
		reason := s.Find("td.style3").Text()
		reason = stringCleaner.Replace(reason)
		reasons = append(reasons, reason)
	})

	return reasons, err
}

func ScrapeNameAndInstructor(doc *goquery.Document) (names, instructors []string, err error) {
	doc.Find("tr.style1 > td").Each(func(i int, s *goquery.Selection) {
		var name, instructor string

		switch i % 3 {
		case 0:
			name = s.Text()
			name = strings.Replace(name, " ", "", -1)
			names = append(names, name)
		case 1:
			instructor = s.Text()
			instructor = strings.Replace(instructor, " ", "", -1)
			instructors = append(instructors, instructor)
		}

	})

	return names, instructors, nil
}

func ScrapeWeekday(doc *goquery.Document) (int, error) {
	weekday := doc.Find("tr.styleT > th").Text()
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
	day := doc.Find("tr.styleT > th").Text()
	day = strings.Split(day, "]")[1]
	day = strings.Split(day, "(")[0]
	day = stringCleaner.Replace(day)
	year := strings.Split(day, "年")[0]
	month := strings.Split(strings.Split(day, "年")[1], "月")[0]
	date := strings.Split(strings.Split(day, "日")[0], "月")[1]

	return string(year) + "/" + string(month) + "/" + string(date), nil
}

*/

func ScrapePlace(doc *goquery.Document) (int, error) {
	place := doc.Find("tr.styleT > th").Text()
	place = strings.Split(place, "]")[0]
	place = strings.Replace(place, "[", "", -1)

	if place == "今出川" {
		return 1, nil
	} else if place == "京田辺" {
		return 2, nil
	}

	return 0, errors.New("place not found")
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
