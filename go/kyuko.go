package kyuko

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/g-hyoga/kyuko/go/model"
	"github.com/g-hyoga/kyuko/go/scrape"
)

//scrapingを実行してデータベースに保存する
func Exec() error {
	var err error

	//第一引数:校地
	//第二引数:曜日
	url, err := scrape.SetUrl(1, 1)
	if err != nil {
		return err
	}

	//httpでやるとき
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
		fmt.Print("hogehogehoge\n")
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
		fmt.Printf("%s\n", data)
	}

	return nil
}
