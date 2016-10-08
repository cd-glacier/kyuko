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

/*
//校地と曜日の情報を含んだurlを引数としてとり、休講structのsliceを返す
//urlはstaticなfileを指定しても良い(test用)
func Scrape(url string) ([]model.KyukoData, error) {
	var kyukoData []model.KyukoData
	var err error

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return kyukoData, err
	}

	//ここのエラー処理どうしよう
	doc.Find("tr.style1").Each(func(i int, s *goquery.Selection) {
		var k model.KyukoData

		k, err = ScrapePeriod(i, s)

		//classがないのでこうするしかない
		tds := s.Find("td")
		for i := range tds.Nodes {
			tmp, _ := iconv.ConvertString(tds.Eq(i).Text(), "shift-jis", "utf-8")
			//授業名の時
			if i%3 == 0 {
				k.ClassName = tmp
				//講師の時
			} else if i%3 == 1 {
				//TrimSpaceとかじゃきかない
				k.Instructor = strings.Split(tmp, " ")[0] + strings.Split(tmp, " ")[4]
			}
		}

		//休講理由
		rawReason := strings.Split(s.Find("td.style3").Text(), "&")[0]
		reason, _ := iconv.ConvertString(rawReason, "shift-jis", "utf-8")
		k.Reason = strings.Split(reason, "ﾂ")[0]

		rawPlaceDayWeek := doc.Find("tr.styleT > th").Text()
		//Place
		rawPlace := strings.Split(strings.Split(rawPlaceDayWeek, "[")[1], "]")[0]
		place, _ := iconv.ConvertString(rawPlace, "shift-jis", "utf-8")
		if place == "今出川" {
			k.Place = 1
		} else if place == "京田辺" {
			k.Place = 2
		}

		//日付と曜日取らないといけない

		kyukoData = append(kyukoData, k)
	})

	return kyukoData, err
}
*/
