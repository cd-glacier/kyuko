package model

import "testing"

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

	testData := KyukoData{Place: 1, Week: 1, Period: 1, Day: "2016/09/26", ClassName: "Insert Test", Instructor: "hoge man", Reason: "darui"}

	_, err = db.Insert(testData)
	if err != nil {
		t.Fatalf("insert に失敗\n%s", err)
	}

	result, err := db.DeleteWhereDayAndClassName("2016/09/26", "Insert Test")
	affectedRows, _ := result.RowsAffected()
	if err != nil || int(affectedRows) <= 0 {
		t.Fatalf("deleteに失敗\n%s", err)
	}
}

func TestSelectAll(t *testing.T) {
	db := DB{}
	var err error

	err = db.Connect()
	if err != nil {
		t.Fatal("データベースに接続できません")
	}
	defer db.Close()

	testData := KyukoData{Place: 1, Week: 1, Period: 1, Day: "2016/09/26", ClassName: "SelectAll Test", Instructor: "tsetMan", Reason: "darui"}

	_, err = db.Insert(testData)
	if err != nil {
		t.Fatalf("insert に失敗: %s", err)
	}

	_, err = db.SelectAll()
	if err != nil {
		t.Fatalf("selectAll に失敗\n%s", err)
	}

	result, err := db.DeleteWhereDayAndClassName("2016/09/26", "SelectAll Test")
	affectedRows, _ := result.RowsAffected()
	if err != nil || int(affectedRows) <= 0 {
		t.Fatalf("deleteに失敗\n%s", err)
	}

}
