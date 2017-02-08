package model

import "testing"

type TestData struct {
	Insert              KyukoData
	InsertCanceledClass CanceledClass
	SelectAll           KyukoData
	Delete              KyukoData
	ShowID              CanceledClass
	IsExistToday        KyukoData
	Add                 CanceledClass
	Reason              Reason
	Day                 Day
}

var db DB
var testData TestData

func init() {
	db.Connect()

	testData.Insert = KyukoData{Place: 1, Weekday: 1, Period: 1, Day: "2016/09/26", ClassName: "Insert Test", Instructor: "hoge man", Reason: "darui"}
	testData.InsertCanceledClass = CanceledClass{Canceled: 10, Place: 1, Weekday: 1, Period: 1, Year: 2016, ClassName: "CanceledClass", Season: "spring", Instructor: "hoge man"}
	testData.SelectAll = KyukoData{Place: 1, Weekday: 1, Period: 1, Day: "2016/09/26", ClassName: "SelectAll Test", Instructor: "tsetMan", Reason: "darui"}
	testData.ShowID = CanceledClass{Canceled: 10, Place: 1, Weekday: 1, Period: 1, Year: 2016, ClassName: "ShowIDTest", Season: "spring", Instructor: "hoge man"}
	testData.IsExistToday = KyukoData{Place: 1, Weekday: 1, Period: 1, Day: "2016/11/11", ClassName: "IsExistToday Test", Instructor: "tsetMan", Reason: "darui"}
	testData.Delete = KyukoData{Place: 1, Weekday: 1, Period: 1, Day: "2016/09/26", ClassName: "Delete Test", Instructor: "tsetMan", Reason: "darui"}
	testData.Add = CanceledClass{Canceled: 10, Place: 1, Weekday: 1, Period: 1, Year: 2016, ClassName: "ADDTest", Season: "spring", Instructor: "hoge man"}
	testData.Reason = Reason{CanceledClassID: 1, Reason: "darui"}
	testData.Day = Day{CanceledClassID: 1, Date: "2016/09/26"}
}

func deleteTestData() {
	db.DeleteWhereDayAndClassName("2016/09/26", "Insert Test")
	db.DeleteWhereDayAndClassName("2016/09/26", "SelectAll Test")
	db.DeleteWhereDayAndClassName("2016/09/26", "Delete Test")
	db.DeleteWhereDayAndClassName("2016/11/11", "IsExistToday Test")

	id, _ := db.ShowCanceledClassID(testData.InsertCanceledClass)
	db.deleteCanceled(id)
	id, _ = db.ShowCanceledClassID(testData.ShowID)
	db.deleteCanceled(id)
	id, _ = db.ShowCanceledClassID(testData.Add)
	db.deleteCanceled(id)

	db.deleteReasonWhere(testData.Reason.CanceledClassID, testData.Reason.Reason)
	db.deleteDayWhere(testData.Day.CanceledClassID, testData.Day.Date)
}

func TestConnectDB(t *testing.T) {
	db := DB{}
	err := db.Connect()
	if err != nil {
		t.Fatal("データベースに接続できません")
	}
}

func TestInsert(t *testing.T) {
	defer deleteTestData()
	var err error

	_, err = db.Insert(testData.Insert)
	if err != nil {
		t.Fatalf("insert に失敗\n%s", err)
	}
}

func TestInsertCanceledClass(t *testing.T) {
	defer deleteTestData()
	var err error

	_, err = db.InsertCanceledClass(testData.InsertCanceledClass)
	if err != nil {
		t.Fatalf("canceled class の insert に失敗\n%s", err)
	}
}

func TestReason(t *testing.T) {
	defer deleteTestData()
	var err error

	_, err = db.InsertReason(testData.Reason)
	if err != nil {
		t.Fatalf("reason の insert に失敗\n%s", err)
	}
}

func TestDay(t *testing.T) {
	defer deleteTestData()
	var err error

	_, err = db.InsertDay(testData.Day)
	if err != nil {
		t.Fatalf("day の insert に失敗\n%s", err)
	}
}

func TestSelectAll(t *testing.T) {
	defer deleteTestData()
	var err error

	_, err = db.Insert(testData.SelectAll)
	if err != nil {
		t.Fatalf("insert に失敗: %s", err)
	}

	_, err = db.SelectAll()
	if err != nil {
		t.Fatalf("selectAll に失敗\n%s", err)
	}
}

func TestShowCanceledClassID(t *testing.T) {
	defer deleteTestData()
	var err error

	_, err = db.InsertCanceledClass(testData.ShowID)
	if err != nil {
		t.Fatalf("canceled class の insert に失敗\n%s", err)
	}

	id, err := db.ShowCanceledClassID(testData.ShowID)
	if err != nil {
		t.Fatalf("Show canceled class id に失敗\n%s", err)
	}

	_, err = db.deleteCanceled(id)
	if err != nil {
		t.Fatalf("Error ShowCanceledClassID: failed DeleteCnacled func\n%s", err)
	}
}

func TestIsExistToday(t *testing.T) {
	defer deleteTestData()
	var err error

	if isExist, err := db.IsExistToday(testData.IsExistToday); isExist || err != nil {
		t.Fatalf("Faleid IsExistToday\n%s", err)
	}

	_, err = db.Insert(testData.IsExistToday)
	if err != nil {
		t.Fatalf("Error IsExistToday: failed Insert func\n%s", err)
	}

	if isExist, err := db.IsExistToday(testData.IsExistToday); !isExist || err != nil {
		t.Fatalf("Faleid IsExistToday\n%s", err)
	}

}

func TestDelete(t *testing.T) {
	var err error

	_, err = db.Insert(testData.Delete)
	if err != nil {
		t.Fatalf("insert に失敗: %s", err)
	}

	result, err := db.DeleteWhereDayAndClassName("2016/09/26", "Delete Test")
	affectedRows, _ := result.RowsAffected()
	if err != nil || int(affectedRows) <= 0 {
		t.Fatalf("deleteに失敗\n%s", err)
	}
}

func TestAddCanceled(t *testing.T) {
	defer deleteTestData()
	var err error

	_, err = db.InsertCanceledClass(testData.Add)
	if err != nil {
		t.Fatalf("Error AddCanceled Test: failed InsertCanceledClass func\n%s", err)
	}

	id, err := db.ShowCanceledClassID(testData.Add)
	if err != nil {
		t.Fatalf("Error AddCanceled Test: failed showCanceledClassID\n%s", err)
	}

	_, err = db.AddCanceled(id)
	if err != nil {
		t.Fatalf("Failed AddCanceled func\n%s", err)
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
