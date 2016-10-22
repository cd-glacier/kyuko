package twitter

import (
	"testing"

	"github.com/g-hyoga/kyuko/go/model"
)

var testPeriods []int
var testReasons, testNames, testInstructors []string
var testPlace, testWeekday int
var testDay string
var testData []model.KyukoData

func init() {
	testPeriods = []int{2, 2, 2, 5}
	testReasons = []string{"公務", "出張", "公務", ""}
	testNames = []string{"環境生理学", "電気・電子計測Ｉ－１", "応用数学ＩＩ－１", "イングリッシュ・セミナー２－７０２"}
	testInstructors = []string{"福岡義之", "松川真美", "大川領", "稲垣俊史"}
	testPlace = 2
	testDay = "2016/10/10"
	testWeekday = 1

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

// どうやってテストしよう
func testUpdate(t *testing.T) {
	var err error
	err = Update(tClient, "test")
	if err != nil {
		t.Fatalf("tweetに失敗しました\nerr: %s", err)
	}

	err = Update(iClient, "test")
	if err != nil {
		t.Fatalf("tweetに失敗しました\nerr: %s", err)
	}
}

func TestCreateLine(t *testing.T) {
	lines := []string{"2限:環境生理学(福岡義之)\n", "2限:電気・電子計測Ｉ－１(松川真美)\n", "2限:応用数学ＩＩ－１(大川領)\n", "5限:イングリッシュ・セミナー２－７０２(稲垣俊史)\n"}

	for i, v := range testData {
		line, err := CreateLine(v)
		if err != nil {
			t.Fatalf("tweetの作成に失敗\nerr: %s", err)
		}

		if line != lines[i] {
			t.Fatalf("lineの作成に失敗しました\nwant: %s\ngot:  %s", lines[i], line)
		}
	}

}

func TestConvertWeekItos(t *testing.T) {
	if weekday, err := ConvertWeekItos(1); weekday != "月" || err != nil {
		t.Fatalf("曜日のconvertに失敗しました\nwant: 月\ngot:  %s\nerror:%s", weekday, err)
	}

	if weekday, err := ConvertWeekItos(6); weekday != "土" || err != nil {
		t.Fatalf("曜日のconvertに失敗しました\nwant: 土\ngot:  %s\nerror:", weekday, err)
	}
	if _, err := ConvertWeekItos(7); err == nil {
		t.Fatalf("存在しない曜日でconvertできています\nerror:", err)
	}
}

func TestCreateContent(t *testing.T) {
	testContents := []string{"月曜日の休講情報\n2限:環境生理学(福岡義之)\n2限:電気・電子計測Ｉ－１(松川真美)\n2限:応用数学ＩＩ－１(大川領)\n5限:イングリッシュ・セミナー２－７０２(稲垣俊史)\n2限:環境生理学(福岡義之)\n2限:電気・電子計測Ｉ－１(松川真美)\n2限:応用数学ＩＩ－１(大川領)\n", "月曜日の休講情報\n5限:イングリッシュ・セミナー２－７０２(稲垣俊史)\n"}

	// 140文字を超えさせるためにtestDataを二回適用している
	testData = append(testData, testData...)

	contents, err := CreateContent(testData)
	if err != nil {
		t.Fatalf("CreateContentでエラー\nerr: %s", err)
	}

	for i, content := range contents {
		if content != testContents[i] {
			t.Fatalf("tweetを140文字以内に収めるの失敗\nwant: %s\ngot:  %s", testContents[i], content)
		}
	}

}
