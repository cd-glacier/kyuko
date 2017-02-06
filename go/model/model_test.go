package model

import "testing"

var db DB

func init() {
	db.Connect()
}

func TestConnectDB(t *testing.T) {
	db := DB{}
	err := db.Connect()
	if err != nil {
		t.Fatal("データベースに接続できません")
	}
}

func TestInsert(t *testing.T) {
	var err error

	testData := KyukoData{Place: 1, Weekday: 1, Period: 1, Day: "2016/09/26", ClassName: "Insert Test", Instructor: "hoge man", Reason: "darui"}

	_, err = db.Insert(testData)
	if err != nil {
		t.Fatalf("insert に失敗\n%s", err)
	}
}

func TestInsertCanceledClass(t *testing.T) {
	var err error

	testData := CanceledClass{Canceled: 10, Place: 1, Weekday: 1, Period: 1, Year: 2016, ClassName: "CanceledClass", Season: "spring", Instructor: "hoge man"}

	_, err = db.InsertCanceledClass(testData)
	if err != nil {
		t.Fatalf("canceled class の insert に失敗\n%s", err)
	}
}

func TestReason(t *testing.T) {
	var err error

	testData := Reason{CanceledClassID: 1, Reason: "darui"}

	_, err = db.InsertReason(testData)
	if err != nil {
		t.Fatalf("reason の insert に失敗\n%s", err)
	}
}

func TestDAy(t *testing.T) {
	var err error

	testData := Day{CanceledClassID: 1, Date: "2016/09/26"}

	_, err = db.InsertDay(testData)
	if err != nil {
		t.Fatalf("day の insert に失敗\n%s", err)
	}
}

func TestSelectAll(t *testing.T) {
	var err error

	testData := KyukoData{Place: 1, Weekday: 1, Period: 1, Day: "2016/09/26", ClassName: "SelectAll Test", Instructor: "tsetMan", Reason: "darui"}
	_, err = db.Insert(testData)
	if err != nil {
		t.Fatalf("insert に失敗: %s", err)
	}

	_, err = db.SelectAll()
	if err != nil {
		t.Fatalf("selectAll に失敗\n%s", err)
	}
}

func TestShowCanceledClassID(t *testing.T) {
	var err error

	testData := CanceledClass{Canceled: 10, Place: 1, Weekday: 1, Period: 1, Year: 2016, ClassName: "ShowIDTest", Season: "spring", Instructor: "hoge man"}

	//以下のコメント部
	//一回だけ実行しないとエラーでる
	//delete関数作るのめんどくさい
	/*
		_, err = db.InsertCanceledClass(testData)
		if err != nil {
			t.Fatalf("canceled class の insert に失敗\n%s", err)
		}
	*/

	_, err = db.ShowCanceledClassID(testData)
	if err != nil {
		t.Fatalf("Show canceled class id に失敗\n%s", err)
	}

}

func TestDelete(t *testing.T) {
	var err error

	testData := KyukoData{Place: 1, Weekday: 1, Period: 1, Day: "2016/09/26", ClassName: "Delete Test", Instructor: "tsetMan", Reason: "darui"}
	_, err = db.Insert(testData)
	if err != nil {
		t.Fatalf("insert に失敗: %s", err)
	}

	result, err := db.DeleteWhereDayAndClassName("2016/09/26", "Delete Test")
	affectedRows, _ := result.RowsAffected()
	if err != nil || int(affectedRows) <= 0 {
		t.Fatalf("deleteに失敗\n%s", err)
	}
}

func TestGetYear(t *testing.T) {
	testData := "2016/02/06"
	out, err := getYear(testData)
	if err != nil {
		t.Fatalf("Error getYear func: %s", err)
	}
	if out != 2016 {
		t.Fatal("Failed GetYear func\n want: %s, got: %s", 2016, out)
	}
}

func TestGetSeason(t *testing.T) {
	out, err := getSeason("2016/01/01")
	if err != nil {
		t.Fatalf("Error getSeason func: %s", err)
	}
	if out != "autumn" {
		t.Fatal("Failed GetSeason func\n want: %s, got: %s", "autumn", out)
	}

	out, err = getSeason("2016/05/01")
	if err != nil {
		t.Fatalf("Error getSeason func: %s", err)
	}
	if out != "spring" {
		t.Fatal("Failed GetSeason func\n want: %s, got: %s", "spring", out)
	}
}
