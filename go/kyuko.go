package kyuko

import (
	"time"

	"github.com/PuerkitoBio/goquery"
	goTwitter "github.com/dghubble/go-twitter/twitter"
	"github.com/g-hyoga/kyuko/go/model"
	"github.com/g-hyoga/kyuko/go/scrape"
	"github.com/g-hyoga/kyuko/go/twitter"
)

// scrapingを実行してデータベースに保存する
// その後twitterに投稿
func Exec(place int, client *goTwitter.Client) error {
	var err error

	//今日の日付
	weekday := int(time.Now().Weekday())
	//今の時間
	nowTime := time.Now().Hour()
	// 18:00超えてたら次の日の情報にする
	if nowTime >= 18 {
		weekday += 1
	}

	//第一引数:校地
	//第二引数:曜日
	url, err := scrape.SetUrl(place, weekday)
	if err != nil {
		return err
	}

	//http
	reader, err := scrape.Get(url)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return err
	}

	kyukoData, err := scrape.Scrape(doc)
	if err != nil {
		return err
	}

	var db model.DB
	err = db.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	for _, data := range kyukoData {
		_, err = db.Insert(data)
		if err != nil {
			return err
		}
	}

	tws, err := twitter.CreateContent(kyukoData)
	if err != nil {
		return err
	}

	for _, tw := range tws {
		twitter.Update(client, tw)
	}

	return nil
}
