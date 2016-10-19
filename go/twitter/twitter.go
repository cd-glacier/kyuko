package twitter

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
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

// tweet argment
func Update(text string) error {
	_, _, err := tClient.Statuses.Update(text, nil)
	if err != nil {
		return err
	}

	return nil
}
