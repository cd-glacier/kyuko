package kyuko

import (
	"errors"
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
	testDay2             string
	testNames            []string
	testReasons          []string

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

func init() {
	db.Connect()
	kyukoReader, _ := EncodeTestFile(KYUKOFILE)
	kyukoDoc, _ = goquery.NewDocumentFromReader(kyukoReader)

	noKyukoReader, _ := EncodeTestFile(NOKYUKOFILE)
	noKyukoDoc, _ = goquery.NewDocumentFromReader(noKyukoReader)

	//正解データの用意
	testPeriods := []int{2, 2, 2, 5}
	testReasons = []string{"公務", "出張", "公務", ""}
	testNames = []string{"環境生理学", "電気・電子計測Ｉ－１", "応用数学ＩＩ－１", "イングリッシュ・セミナー２－７０２"}
	testInstructors := []string{"福岡義之", "松川真美", "大川領", "稲垣俊史"}
	testPlace := 2
	testDay = "2016/10/10"
	testDay2 = "2016/10/12"
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
	for i, className := range testNames {
		db.DeleteWhereDayAndClassName(testDay, className)

		id, _ := db.ShowCanceledClassID(model.CanceledClass{ClassName: className, Year: 2016, Season: "autumn"})

		db.DeleteCanceled(id)
		db.DeleteReasonWhere(id, testReasons[i])
		db.DeleteDayWhere(id, testDay)
		db.DeleteDayWhere(id, testDay2)
	}
}

func showCanceled(c model.CanceledClass) (int, error) {
	id, err := db.ShowCanceledClassID(c)
	if err != nil {
		return -1, err
	}
	canceled, err := db.ShowCanceled(id)
	if err != nil {
		return -1, err
	}
	return canceled, nil
}

func reproducer(op int) error {
	var kyukoData []model.KyukoData
	kyukoData, err := scraper(kyukoDoc)
	if err != nil {
		return err
	}

	// 3回目だけ違うデータを用意
	if op == 3 {
		for i, _ := range kyukoData {
			kyukoData[i].Day = testDay2
		}
		for i, _ := range testData {
			testData[i].Day = testDay2
		}
	}

	//Insert
	err = manageDB(kyukoData)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(kyukoData, testData) {
		return errors.New("Error exec: not equal answer data")
	}
	return nil
}

func TestExec(t *testing.T) {
	defer deleteTestData()

	//一回目のテスト
	err := reproducer(1)
	if err != nil {
		t.Fatalf("error testexec: failed reproducer func\n%s", err)
	}

	//canceled_classのcanceledが1かどうか
	for i, _ := range testData {
		c, err := model.KyukoToCanceled(testData[i])
		if err != nil {
			t.Fatalf("Error TestExec: Failed KyukoToCanceled func\n%s", err)
		}
		canceled, err := showCanceled(c)
		if err != nil {
			t.Fatalf("Error TestExec: Failed showCanceled func\n%s", err)
		}
		if canceled != 1 {
			t.Fatalf("Error TestExec: canceled in DB\n want: %s\n got:  %s\n", 1, canceled)
		}
	}

	//二回目のテスト
	//重複データを与えているので何もInsertして欲しくない
	err = reproducer(2)
	if err != nil {
		t.Fatalf("error testexec: failed reproducer func\n%s", err)
	}

	//canceled_classのcanceledが1かどうか
	for i, _ := range testData {
		c, err := model.KyukoToCanceled(testData[i])
		if err != nil {
			t.Fatalf("Error TestExec: Failed KyukoToCanceled func\n%s", err)
		}
		canceled, err := showCanceled(c)
		if err != nil {
			t.Fatalf("Error TestExec: Failed isCorrectCanceled func\n%s", err)
		}
		if canceled != 1 {
			t.Fatalf("Error TestExec: canceled in DB\n want: %s\n got:  %s\n", 1, canceled)
		}
	}

	// 三回目のInsert日付を変えて
	// 別の日のデータとして扱う
	// reason, dayテーブルにInsertされて
	// canceledカラムが1増えれば良い
	err = reproducer(3)
	if err != nil {
		t.Fatalf("error testexec: failed reproducer func\n%s", err)
	}

	//canceled_classのcanceledが2かどうか
	for i, _ := range testData {
		c, err := model.KyukoToCanceled(testData[i])
		if err != nil {
			t.Fatalf("Error TestExec: Failed KyukoToCanceled func\n%s", err)
		}
		canceled, err := showCanceled(c)
		if err != nil {
			t.Fatalf("Error TestExec: Failed isCorrectCanceled func\n%s", err)
		}
		if canceled != 2 {
			t.Fatalf("Error TestExec: canceled in DB\n want: %s\n got:  %s\n", 2, canceled)
		}
	}

}
