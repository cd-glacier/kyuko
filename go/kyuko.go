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
func Exec(place int, client *goTwitter.Client) ([]model.KyukoData, error) {
	var kyukoResult []model.KyukoData
	var err error

	//今日の日付
	weekday := int(time.Now().Weekday())
	//今の時間
	nowTime := time.Now().Hour()
	// 18:00超えてたら次の日の情報にする
	if nowTime >= 18 {
		weekday += 1
	}
	// 日曜なら月曜の情報にする
	if weekday == 7 {
		weekday = 1
	}

	//第一引数:校地
	//第二引数:曜日
	url, err := scrape.SetUrl(place, weekday)
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

	kyukoData, err := scrape.Scrape(doc)
	if err != nil {
		return nil, err
	}

	var db model.DB
	err = db.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	for _, data := range kyukoData {
		_, err = db.Insert(data)
		kyukoResult = append(kyukoResult, data)
		if err != nil {
			return nil, err
		}

		canceledClass, err := model.KyukoToCanceled(data)
		if err != nil {
			return nil, err
		}

		//API用
		id, err := db.ShowCanceledClassID(canceledClass)
		if err != nil {
			return nil, err
		}

		//DBに存在しないデータなら
		if id == -1 {
			_, err = db.InsertCanceledClass(canceledClass)
			if err != nil {
				return nil, err
			}
		}

		canceledClass.ID = id
		_, err = db.AddCanceled(canceledClass.ID)
		if err != nil {
			return nil, err
		}
	}

	//tws, err := twitter.CreateContent(kyukoData)
	_, err = twitter.CreateContent(kyukoData)
	if err != nil {
		return nil, err
	}

	/*
		for _, tw := range tws {
			err := twitter.Update(client, tw)
			if err != nil {
				return nil, err
				return nil, err
			}
		}
	*/

	return kyukoData, nil
}
