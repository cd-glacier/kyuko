package scrape

import (
	"errors"
	"strconv"
)

type KyukoData struct {
	Place      int
	Week       int
	Period     int
	Date       string
	ClassName  string
	Instructor string
	Reason     string
}

func SetUrl(place, week int) (string, error) {
	url := "http://duet.doshisha.ac.jp/info/KK1000.jsp?katei=1"
	if place >= 1 && place <= 2 && week >= 1 && week <= 6 {
		url = url + "&youbi=" + strconv.Itoa(week)
		url = url + "&kouchi=" + strconv.Itoa(place)
		return url, nil
	} else {
		return "", errors.New("place is 1 or 2, 0 < week < 7")
	}
}

func Scrape(url string) (KyukoData, error) {
	k := KyukoData{}
	if url == "" {
		return k, nil
	} else {
		k.Place = 1
		k.Week = 2
		k.Period = 3
		k.Date = "2016/09/24"
		k.ClassName = "建学の精神のキリスト教"
		k.Instructor = "hoge man"
		k.Reason = "めんどい"
	}
	return k, nil
}
