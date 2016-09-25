package model

import (
	"fmt"
	"testing"
)

func TestConnectDB(t *testing.T) {
	db := DB{}
	err := db.Connect()
	if err != nil {
		t.Fatal("データベースに接続できません")
	}

	err = db.Close()
	if err != nil {
		t.Fatal("Closeに失敗しました")
	}

}

func TestInsert(t *testing.T) {
	db := DB{}
	var err error

	err = db.Connect()
	if err != nil {
		t.Fatal("データベースに接続できません")
	}
	defer db.Close()

	testData := KyukoData{Place: 1, Week: 1, Period: 1, Date: "2016/09/26", ClassName: "DOSUKOI", Instructor: "hoge man", Reason: "darui"}

	result, err := db.Insert(testData)
	if err != nil {
		t.Fatalf("failed insert: %s", err)
	}

	fmt.Printf("%s", result)
}
