package kyuko

import (
	"fmt"
	"time"

	goTwitter "github.com/dghubble/go-twitter/twitter"
	"github.com/g-hyoga/kyuko/src/model"
)

func Exec(place int, client *goTwitter.Client) ([]model.KyukoData, error) {
	var kyukoData []model.KyukoData

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

	fmt.Println(nowTime)

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
