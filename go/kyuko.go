package kyuko

import (
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/g-hyoga/kyuko/go/model"
	"github.com/g-hyoga/kyuko/go/scrape"
	"github.com/g-hyoga/kyuko/go/twitter"
)

var (
	T_CONSUMER_KEY        = os.Getenv("T_CONSUMER_KEY")
	T_CONSUMER_SECRET     = os.Getenv("T_CONSUMER_SECRET")
	T_ACCESS_TOKEN        = os.Getenv("T_ACCESS_TOKEN")
	T_ACCESS_TOKEN_SECRET = os.Getenv("T_ACCESS_TOKEN_SECRET")

	I_CONSUMER_KEY        = os.Getenv("I_CONSUMER_KEY")
	I_CONSUMER_SECRET     = os.Getenv("I_CONSUMER_SECRET")
	I_ACCESS_TOKEN        = os.Getenv("I_ACCESS_TOKEN")
	I_ACCESS_TOKEN_SECRET = os.Getenv("I_ACCESS_TOKEN_SECRET")
)

// scrapingを実行してデータベースに保存する
// その後twitterに投稿
func Exec() error {
	var err error

	//第一引数:校地
	//第二引数:曜日
	url, err := scrape.SetUrl(2, 1)
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

	//iClient := NewTwitterClient(I_CONSUMER_KEY, I_CONSUMER_SECRET, I_ACCESS_TOKEN, I_ACCESS_TOKEN_SECRET)
	tClient := twitter.NewTwitterClient(T_CONSUMER_KEY, T_CONSUMER_SECRET, T_ACCESS_TOKEN, T_ACCESS_TOKEN_SECRET)

	for _, tw := range tws {
		//twitter.Update(iClient, tw)
		twitter.Update(tClient, tw)
	}

	return nil
}
