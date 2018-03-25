package kyuko

import (
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
	goTwitter "github.com/dghubble/go-twitter/twitter"
	"github.com/g-hyoga/kyuko/src/data"
	"github.com/g-hyoga/kyuko/src/scrape"
)

func Exec(place int, client *goTwitter.Client) ([]data.KyukoData, error) {
	var kyukoData []data.KyukoData

	isTommorow := allowTommorowData()

	doc, err := readHTML(place, isTommorow)
	if err != nil {
		return kyukoData, err
	}

	kyukoData, err = scraper(doc, place)
	if err != nil {
		return kyukoData, err
	}

	return kyukoData, nil
}

func allowTommorowData() bool {
	//今の時間
	nowTime := time.Now().Hour()

	fmt.Printf("Time: %s", string(nowTime))

	// 18:00超えてたら次の日の情報にする
	if nowTime >= 18 {
		return true
	}
	//今日の曜日
	weekday := int(time.Now().Weekday())
	// 日曜なら月曜の情報にする
	if weekday >= 7 {
		return true
	}
	return false
}

func readHTML(place int, isTommorow bool) (*goquery.Document, error) {
	//第一引数:校地
	//第二引数:曜日
	url, err := scrape.SetUrl(place, isTommorow)
	if err != nil {
		return nil, err
	}
	//http
	reader, err := scrape.Get(url)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}
	return doc, err
}

func scraper(doc *goquery.Document, place int) ([]data.KyukoData, error) {
	kyukoData, err := scrape.Scrape(doc, place)
	if err != nil {
		return nil, err
	}
	return kyukoData, nil
}
