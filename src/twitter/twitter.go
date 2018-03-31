package twitter

import (
	"errors"
	"log"
	"os"
	"strconv"
	"unicode/utf8"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/g-hyoga/kyuko/src/data"
)

// newTwitterClient returns a new Twitter Client
func NewTwitterClient(consumerKey, consumerSecret, token, tokenSecret string) (*twitter.Client, error) {
	if consumerKey == "" || consumerSecret == "" {
		return nil, errors.New("consumerKey or consumerSecret is missing")
	}

	if token == "" || tokenSecret == "" {
		return nil, errors.New("access token or access tokenSecret is missing")
	}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	oauthToken := oauth1.NewToken(token, tokenSecret)
	httpClient := config.Client(oauth1.NoContext, oauthToken)
	return twitter.NewClient(httpClient), nil
}

// create line of template
// period:className(Instructor)
func CreateLine(kyuko data.KyukoData) (string, error) {
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
func CreateContent(kyuko []data.KyukoData) ([]string, error) {
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
	tweetLen := 9
	for _, line := range lines {
		lineLen := utf8.RuneCountInString(line)
		if tweetLen+lineLen > 140 {
			tws = append(tws, tw)
			tw = ""
			tw = KANJIweekday + "曜日の休講情報\n"
			tweetLen = 9
		}
		tw = tw + line
		tweetLen += lineLen
	}
	tws = append(tws, tw)

	return tws, err
}

// tweet argment
func Update(client *twitter.Client, text string) error {
	if os.Getenv("KYUKO") == "PRODUCTION" {
		log.Println("tweet: ", text)
		_, _, err := client.Statuses.Update(text, nil)
		if err != nil {
			return err
		}
	} else {
		log.Println("tweet: ", text)
	}

	return nil
}
