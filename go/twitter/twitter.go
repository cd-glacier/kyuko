package twitter

import (
	"errors"
	"fmt"
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

// template
// className(Instructor)
func CreateLine(kyuko model.KyukoData) (string, error) {
	if kyuko.ClassName == "" || kyuko.Instructor == "" || kyuko.Period == 0 {
		return "", errors.New("休講情報がないです")
	}

	fmt.Printf("hogehoge%s\n", kyuko.Period)

	period := strconv.Itoa(kyuko.Period)

	line := period + "限:" + kyuko.ClassName + "(" + kyuko.Instructor + ")\n"
	return line, nil
}

// tweet argment
func Update(text string) error {
	_, _, err := tClient.Statuses.Update(text, nil)
	if err != nil {
		return err
	}

	return nil
}
