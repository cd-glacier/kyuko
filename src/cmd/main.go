package main

import (
	"fmt"
	"os"

	"github.com/g-hyoga/kyuko/src/kyuko"
	"github.com/g-hyoga/kyuko/src/twitter"
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
	fmt.Println("runnig...")

	/*
		log := logrus.New()
		log.Formatter = new(logrus.JSONFormatter)
		f, err := os.OpenFile("../log/out.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.WithFields(logrus.Fields{
				"err": err,
			}).Fatal("error")
		}
		log.Out = f
	*/

	iClient := twitter.NewTwitterClient(I_CONSUMER_KEY, I_CONSUMER_SECRET, I_ACCESS_TOKEN, I_ACCESS_TOKEN_SECRET)
	//第一引数:校地
	//第二引数:twitter client
	kyukoData, iErr := kyuko.Exec(1, iClient)
	if iErr != nil {
		fmt.Println(iErr)
		/*
				log.WithFields(logrus.Fields{
					"err": iErr,
				}).Fatal("error")
			} else {
				log.WithFields(logrus.Fields{
					"data": kyukoData,
				}).Info("kyuko data")
		*/
	}
	fmt.Println(kyukoData)

	tClient := twitter.NewTwitterClient(T_CONSUMER_KEY, T_CONSUMER_SECRET, T_ACCESS_TOKEN, T_ACCESS_TOKEN_SECRET)
	//第一引数:校地
	//第二引数:twitter client
	kyukoData, tErr := kyuko.Exec(2, tClient)
	if tErr != nil {
		fmt.Println(iErr)
		/*
				log.WithFields(logrus.Fields{
					"err": tErr,
				}).Fatal("error")
			} else {
				log.WithFields(logrus.Fields{
					"data": kyukoData,
				}).Info("kyuko data")
		*/
	}
	fmt.Println(kyukoData)
}
