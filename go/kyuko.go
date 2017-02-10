package kyuko

import (
	"time"

	"github.com/PuerkitoBio/goquery"
	goTwitter "github.com/dghubble/go-twitter/twitter"
	"github.com/g-hyoga/kyuko/go/model"
	"github.com/g-hyoga/kyuko/go/scrape"
	"github.com/g-hyoga/kyuko/go/twitter"
)

func weekdayToday() int {
	//今日の曜日
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
	return weekday
}

func scraper(place, weekday int) ([]model.KyukoData, error) {
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
	return kyukoData, nil
}

func ManageDB(kyukoData []model.KyukoData) error {
	var db model.DB
	err = db.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	for _, data := range kyukoData {
		_, err = db.Insert(data)
		kyukoResult = append(kyukoResult, data)
		if err != nil {
			return err
		}

		canceledClass, err := model.KyukoToCanceled(data)
		if err != nil {
			return err
		}

		//挿入するデータが存在するのか確認
		id, err := db.ShowCanceledClassID(canceledClass)
		if err != nil {
			return err
		}

		//DBに存在するデータかつ今日のデータでないなら
		if isExist, _ := db.IsExistToday(data); id != -1 && !isExist {
			canceledClass.ID = id
			_, err = db.AddCanceled(canceledClass.ID)
			if err != nil {
				return err
			}
		}

		_, err = db.InsertCanceledClass(canceledClass)
		if err != nil {
			return err
		}
	}
	return nil
}

func manageTwitter(kyukoData model.KyukoData) error {
	tws, err := twitter.CreateContent(kyukoData)
	//_, err = twitter.CreateContent(kyukoData)
	if err != nil {
		return err
	}

	for _, tw := range tws {
		err := twitter.Update(client, tw)
		if err != nil {
			return err
		}
	}

	return nil
}

// scrapingを実行してデータベースに保存する
// その後twitterに投稿
func Exec(place int, client *goTwitter.Client) ([]model.KyukoData, error) {
	var kyukoResult []model.KyukoData

	weekday := weekdayToday()

	kyukoData, err := scraper(place, weekday)
	if err != nil {
		return kyukoResult, err
	}

	err := manageDB(kyuokData)
	if err != nil {
		return kyukoResult, err
	}

	err := manageTwitter(kyukoData)
	if err != nil {
		return kyukoResult, err
	}

	return kyukoData, nil
}
