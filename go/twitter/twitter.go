package twitter

import (
	"errors"
	"os"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/g-hyoga/kyuko/go/model"
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

var tClient twitter.Client

func init() {
	//京田辺
	config := oauth1.NewConfig(T_CONSUMER_KEY, T_CONSUMER_SECRET)
	token := oauth1.NewToken(T_ACCESS_TOKEN, T_ACCESS_TOKEN_SECRET)
	httpClient := config.Client(oauth1.NoContext, token)
	tClient = *twitter.NewClient(httpClient)
}

// create line of template
// period:className(Instructor)
func CreateLine(kyuko model.KyukoData) (string, error) {
	if kyuko.ClassName == "" || kyuko.Instructor == "" || kyuko.Period == 0 {
		return "", errors.New("休講情報がないです")
	}

	period := strconv.Itoa(kyuko.Period)

	line := period + "限:" + kyuko.ClassName + "(" + kyuko.Instructor + ")\n"
	return line, nil
}

// convert weekday(int) to weekday(KANJI)
func ConvertWeekItos(weekday int) (string, error) {
	weekKANJI := []string{"日", "月", "火", "水", "木", "金", "土"}

	if weekday < 0 || weekday > 6 {
		return "", errors.New("存在しない曜日が入力されています")
	}

	return weekKANJI[weekday], nil
}

// create tweet template
// exsample
//
// hoge曜日の休講情報
// period:className(Instructor)
// period:className(Instructor)
// ...
//
// in 140 characters
func CreateContent(kyuko []model.KyukoData) ([]string, error) {
	var tws []string

	// create lines
	var lines []string
	for _, v := range kyuko {
		line, err := CreateLine(v)
		if err != nil {
			return tws, err
		}
		lines = append(lines, line)
	}

	// convert weekday
	weekday := kyuko[0].Weekday
	KANJIweekday, err := ConvertWeekItos(weekday)
	if err != nil {
		return tws, err
	}

	// create content
	tw := KANJIweekday + "曜日の休講情報\n"
	for _, line := range lines {
		if len([]rune(tw+line)) > 140 {
			tws = append(tws, tw)
			tw = ""
			tw = KANJIweekday + "曜日の休講情報\n"
			tw = tw + line
		} else {
			tw = tw + line
		}
	}
	tws = append(tws, tw)

	return tws, err
}

// tweet argment
func Update(text string) error {
	_, _, err := tClient.Statuses.Update(text, nil)
	if err != nil {
		return err
	}

	return nil
}
