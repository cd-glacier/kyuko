package main

import (
	"os"

	"github.com/Sirupsen/logrus"
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
	log := logrus.New()
	log.Formatter = new(logrus.JSONFormatter)

	f, _ := os.OpenFile("./log/test.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.Out = f

	log.WithFields(logrus.Fields{}).Info("running")

	iClient := twitter.NewTwitterClient(I_CONSUMER_KEY, I_CONSUMER_SECRET, I_ACCESS_TOKEN, I_ACCESS_TOKEN_SECRET)
	//第一引数:校地
	//第二引数:twitter client
	iErr := kyuko.Exec(1, iClient)
	if iErr != nil {
		log.WithFields(logrus.Fields{
			"err": iErr,
		}).Fatal("error")
	}

	tClient := twitter.NewTwitterClient(T_CONSUMER_KEY, T_CONSUMER_SECRET, T_ACCESS_TOKEN, T_ACCESS_TOKEN_SECRET)
	//第一引数:校地
	//第二引数:twitter client
	tErr := kyuko.Exec(2, tClient)
	if tErr != nil {
		log.WithFields(logrus.Fields{
			"err": tErr,
		}).Fatal("error")

	}
}
