package kyuko

import (
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/PuerkitoBio/goquery"
	"github.com/g-hyoga/kyuko/go/model"
)

var (
	db                   model.DB
	kyukoDoc, noKyukoDoc *goquery.Document
	testData             []model.KyukoData
	testDay              string
	testNames            []string

	T_CONSUMER_KEY        = os.Getenv("T_CONSUMER_KEY")
	T_CONSUMER_SECRET     = os.Getenv("T_CONSUMER_SECRET")
	T_ACCESS_TOKEN        = os.Getenv("T_ACCESS_TOKEN")
	T_ACCESS_TOKEN_SECRET = os.Getenv("T_ACCESS_TOKEN_SECRET")

	I_CONSUMER_KEY        = os.Getenv("I_CONSUMER_KEY")
	I_CONSUMER_SECRET     = os.Getenv("I_CONSUMER_SECRET")
	I_ACCESS_TOKEN        = os.Getenv("I_ACCESS_TOKEN")
	I_ACCESS_TOKEN_SECRET = os.Getenv("I_ACCESS_TOKEN_SECRET")
)

const (
	KYUKOFILE   = "./testdata/kyuko.html"
	NOKYUKOFILE = "./testdata/not_kyuko.html"
)

func SjisToUtf8(str string) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder()))
	if err != nil {
		return "", err
	}
	return string(ret), err
}

func EncodeTestFile(fileName string) (io.Reader, error) {
	//testfileのenocde
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	utfFile, err := SjisToUtf8(string(file))
	if err != nil {
		return nil, err
	}
	stringReader := strings.NewReader(utfFile)

	return stringReader, nil
}

func deleteTestData() {
	for _, className := range testNames {
		db.DeleteWhereDayAndClassName(testDay, className)
	}

	for id := 254; id <= 257; id++ {
		db.DeleteCanceled(id)
		db.DeleteDayWhere(id, testDay)
	}
}

func init() {
	db.Connect()
	kyukoReader, _ := EncodeTestFile(KYUKOFILE)
	kyukoDoc, _ = goquery.NewDocumentFromReader(kyukoReader)

	noKyukoReader, _ := EncodeTestFile(NOKYUKOFILE)
	noKyukoDoc, _ = goquery.NewDocumentFromReader(noKyukoReader)

	//正解データの用意
	testPeriods := []int{2, 2, 2, 5}
	testReasons := []string{"公務", "出張", "公務", ""}
	testNames = []string{"環境生理学", "電気・電子計測Ｉ－１", "応用数学ＩＩ－１", "イングリッシュ・セミナー２－７０２"}
	testInstructors := []string{"福岡義之", "松川真美", "大川領", "稲垣俊史"}
	testPlace := 2
	testDay = "2016/10/10"
	testWeekday := 1

	for i := range testPeriods {
		k := model.KyukoData{}
		k.Period = testPeriods[i]
		k.Reason = testReasons[i]
		k.ClassName = testNames[i]
		k.Instructor = testInstructors[i]
		k.Weekday = testWeekday
		k.Place = testPlace
		k.Day = testDay
		testData = append(testData, k)
	}
}

func TestExec(t *testing.T) {
	defer deleteTestData()
	var kyukoData []model.KyukoData

	//testデータを使う
	kyukoData, err := scraper(kyukoDoc)
	if err != nil {
		t.Fatalf("Error Exec: Failed scraper func \n%s", err)
	}

	//一回目のInsert
	err = manageDB(kyukoData)
	if err != nil {
		t.Fatalf("Error Exec: Failed manageDB func \n%s", err)
	}
	if !reflect.DeepEqual(kyukoData, testData) {
		t.Fatalf("Error TestExec: once \n want: %v\n got:  %v", testData, kyukoData)
	}

	//canceled_classのcanceledが1かどうか

	//二回目のInsert
	//何もInsertして欲しくない
	err = manageDB(kyukoData)
	if err != nil {
		t.Fatalf("Error Exec: Failed manageDB func \n%s", err)
	}
	if !reflect.DeepEqual(kyukoData, testData) {
		t.Fatalf("Error TestExec: once \n want: %v\n got:  %v", testData, kyukoData)
	}

	//canceled_classのcanceledが2かどうか

	// 三回目のInsert日付を変えて
	// 別の日のデータとして扱う
	// reason, dayテーブルにInsertされて
	// canceledカラムが1増えれば良い
	for _, data := range kyukoData {
		data.Day = "2016/10/12"
	}
	for _, data := range testData {
		data.Day = "2016/10/12"
	}

	err = manageDB(kyukoData)
	if err != nil {
		t.Fatalf("Error Exec: Failed manageDB func \n%s", err)
	}
	if !reflect.DeepEqual(kyukoData, testData) {
		t.Fatalf("Error TestExec: twice \n want: %v\n got:  %v", testData, kyukoData)
	}

	//canceled_classのcanceledが3かどうか
}
