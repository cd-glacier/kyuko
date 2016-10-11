package scrape

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/PuerkitoBio/goquery"
	"github.com/g-hyoga/kyuko/go/model"
)

//place(1: 今出川 ,2: 京田辺), week(1 ~ 6: Mon ~ Sat)を引数に持ち
//urlを生成する
func SetUrl(place, week int) (string, error) {
	url := "http://duet.doshisha.ac.jp/info/KK1000.jsp?katei=1"
	//weekに7(Sunday)はない
	if (place != 1 && place != 2) || week < 1 || week > 6 {
		return "", errors.New("place is 1 or 2, 0 < week < 7")
	} else {
		url = url + "&youbi=" + strconv.Itoa(week)
		url = url + "&kouchi=" + strconv.Itoa(place)
		return url, nil
	}
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
	utfBody := transform.NewReader(reader, japanese.ShiftJIS.NewDecoder())

	return utfBody, nil
}

func ScrapePeriod(doc *goquery.Document) ([]int, error) {
	var periods []int
	var err error

	//エラー処理どうにかする
	//"1講時"みたいなのが取れる
	doc.Find("tr.style1").Each(func(i int, s *goquery.Selection) {
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
		reason = strings.Replace(reason, "\n", "", -1)
		reason = strings.Replace(reason, " ", "", -1)
		reason = strings.Replace(reason, "\u00a0", "", -1)
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

func ScrapeDay(doc *goquery.Document) (string, error) {
	day := doc.Find("tr.styleT > th").Text()
	day = strings.Split(day, "]")[1]
	day = strings.Split(day, "(")[0]
	day = strings.Replace(day, " ", "", -1)
	day = strings.Replace(day, "\n", "", -1)
	day = strings.Replace(day, "\u00a0", "", -1)
	year := strings.Split(day, "年")[0]
	month := strings.Split(strings.Split(day, "年")[1], "月")[0]
	date := strings.Split(strings.Split(day, "日")[0], "月")[1]

	return string(year) + "/" + string(month) + "/" + string(date), nil
}

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

func ScrapeWeekday(doc *goquery.Document) (int, error) {
	weekday := doc.Find("tr.styleT > th").Text()
	weekday = strings.Split(weekday, "(")[1]
	weekday = strings.Replace(weekday, ")", "", -1)
	weekday = strings.Replace(weekday, " ", "", -1)
	weekday = strings.Replace(weekday, "\n", "", -1)
	weekday = strings.Replace(weekday, "\u00a0", "", -1)
	youbi, err := ConvertWeekStoi(weekday)

	if err != nil {
		return -1, err
	}

	return youbi, nil

}

//休講structのsliceを返す
func Scrape(doc *goquery.Document) ([]model.KyukoData, error) {
	var kyukoData []model.KyukoData
	var err error

	var periods []int
	var reasons, names, instructors []string
	var weekday, place int
	var day string

	finished := make(chan bool)
	go func() {
		periods, err = ScrapePeriod(doc)
		finished <- true
	}()
	go func() {
		reasons, err = ScrapeReason(doc)
		finished <- true
	}()
	go func() {
		names, instructors, err = ScrapeNameAndInstructor(doc)
		finished <- true
	}()
	go func() {
		weekday, err = ScrapeWeekday(doc)
		finished <- true
	}()
	go func() {
		day, err = ScrapeDay(doc)
		finished <- true
	}()
	go func() {
		place, err = ScrapePlace(doc)
		finished <- true
	}()
	for i := 1; i <= 6; i++ {
		<-finished
	}

	if err != nil {
		return nil, err
	}

	if len(periods) != len(reasons) && len(periods) != len(names) && len(periods) != len(instructors) {
		return nil, errors.New("取得できていない情報があります")
	}

	for i, _ := range periods {
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
