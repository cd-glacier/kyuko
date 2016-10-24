package main

import (
	"fmt"
	"os"

	kyuko "github.com/g-hyoga/kyuko/go"
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

func main() {
	var err error

	iClient := twitter.NewTwitterClient(I_CONSUMER_KEY, I_CONSUMER_SECRET, I_ACCESS_TOKEN, I_ACCESS_TOKEN_SECRET)
	//第一引数:校地
	//第二引数:曜日
	//第三引数:twitter client
	err = kyuko.Exec(1, iClient)
	if err != nil {
		fmt.Printf("error!!\n%s", err)
	}

	tClient := twitter.NewTwitterClient(T_CONSUMER_KEY, T_CONSUMER_SECRET, T_ACCESS_TOKEN, T_ACCESS_TOKEN_SECRET)
	//第一引数:校地
	//第二引数:曜日
	//第三引数:twitter client
	err = kyuko.Exec(2, tClient)
	if err != nil {
		fmt.Printf("error!!\n%s", err)
	}
}
