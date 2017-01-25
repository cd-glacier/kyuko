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

	testData := KyukoData{Place: 1, Canceled: 1, Weekday: 1, Period: 1, Day: "2016/09/26", ClassName: "Insert Test", Instructor: "hoge man", Reason: "darui"}

	_, err = db.InsertCanceledClass(testData)
	if err != nil {
		t.Fatalf("insertCanceledClass に失敗\n%s", err)
	}
}

func TestUpdateCanceledClass(t *testing.T) {
	var err error

	_, err = db.UpdateCanceledClass(10, 10)
	if err != nil {
		t.Fatalf("update に失敗\n%s", err)
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

func TestSelectAllCanceledClass(t *testing.T) {
	var err error

	testData := KyukoData{Place: 1, Weekday: 1, Period: 1, Day: "2016/09/26", ClassName: "SelectAll Test", Instructor: "tsetMan", Reason: "darui"}
	_, err = db.Insert(testData)
	if err != nil {
		t.Fatalf("insert に失敗: %s", err)
	}

	_, err := db.SelectAllCanceledClass()
	if err != nil {
		t.Fatalf("selectAll に失敗\n%s", err)
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
