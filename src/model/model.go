package model

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	db *sql.DB
}

func (db *DB) Connect() error {
	var err error
	db.db, err = sql.Open("mysql", "root:password@/kyuko")
	return err
}

func (db *DB) Close() error {
	err := db.db.Close()
	return err
}

func (db *DB) Insert(k KyukoData) (sql.Result, error) {
	result, err := db.db.Exec("INSERT INTO `kyuko_data` VALUES(?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE `reason`=?;", 0, k.Place, k.Weekday, k.Period, k.Day, k.ClassName, k.Instructor, k.Reason, k.Reason)
	return result, err
}

func (db *DB) InsertCanceledClass(c CanceledClass) (sql.Result, error) {
	sql := "INSERT INTO `canceled_class` (canceled, place, week, period, year, season, class_name, instructor) VALUES(?, ?, ?, ?, ?, ?, ?, ?);"

	result, err := db.db.Exec(sql, c.Canceled, c.Place, c.Weekday, c.Period, c.Year, c.Season, c.ClassName, c.Instructor)
	return result, err
}

func (db *DB) InsertDay(d Day) (sql.Result, error) {
	sql := "INSERT INTO `day` (canceled_class_id, day) VALUES(?, ?);"
	result, err := db.db.Exec(sql, d.CanceledClassID, d.Date)
	return result, err
}

func (db *DB) InsertReason(r Reason) (sql.Result, error) {
	sql := "INSERT INTO `reason` (canceled_class_id, reason) VALUES(?, ?);"
	result, err := db.db.Exec(sql, r.CanceledClassID, r.Reason)
	return result, err
}

func ScanAll(rows *sql.Rows) ([]KyukoData, error) {
	kyukoData := []KyukoData{}
	var err error
	for rows.Next() {
		var k KyukoData
		if err = rows.Scan(&k.ID, &k.Place, &k.Weekday, &k.Period, &k.Day, &k.ClassName, &k.Instructor, &k.Reason); err != nil {
			return kyukoData, err
		}
		kyukoData = append(kyukoData, k)
	}
	return kyukoData, err
}

func ScanCanceledClass(rows *sql.Rows) ([]CanceledClass, error) {
	canceledClass := []CanceledClass{}
	var err error
	for rows.Next() {
		var c CanceledClass
		if err = rows.Scan(&c.ID, &c.Canceled, &c.Place, &c.Weekday, &c.Period, &c.Year, &c.Season, &c.ClassName, &c.Instructor); err != nil {
			return canceledClass, err
		}
		canceledClass = append(canceledClass, c)
	}
	return canceledClass, err
}

func ScanDay(rows *sql.Rows) ([]Day, error) {
	days := []Day{}
	var err error
	for rows.Next() {
		var d Day
		if err = rows.Scan(&d.ID, &d.CanceledClassID, &d.Date); err != nil {
			return days, err
		}
		days = append(days, d)
	}
	return days, err
}

func ScanReason(rows *sql.Rows) ([]Reason, error) {
	reasons := []Reason{}
	var err error
	for rows.Next() {
		var r Reason
		if err = rows.Scan(&r.ID, &r.CanceledClassID, &r.Reason); err != nil {
			return reasons, err
		}
		reasons = append(reasons, r)
	}
	return reasons, err
}

func (db *DB) SelectAll() ([]KyukoData, error) {
	kyukoData := []KyukoData{}
	rows, err := db.db.Query("select * from kyuko_data where id in(select min(id) from kyuko_data group by class_name, day);")
	if err != nil {
		return kyukoData, err
	}
	defer rows.Close()
	kyukoData, err = ScanAll(rows)
	if err != nil {
		return kyukoData, err
	}
	return kyukoData, err
}

func (db *DB) SelectReaonFromCanceledID(canceledID int) ([]Reason, error) {
	reasons := []Reason{}
	sql := "SELECT * FROM reason WHERE canceled_id=?"
	rows, err := db.db.Query(sql, canceledID)
	if err != nil {
		return reasons, err
	}
	defer rows.Close()
	reasons, err = ScanReason(rows)
	if err != nil {
		return reasons, err
	}

	return reasons, nil
}

// DBからIDを取得します。
// 存在しない場合は-1を返します
// 一意に定まらない場合はerrorを返します
func (db *DB) ShowCanceledClassID(c CanceledClass) (int, error) {
	sql := "select * from canceled_class where class_name=? and year=? and season = ?"
	rows, err := db.db.Query(sql, c.ClassName, c.Year, c.Season)
	if err != nil {
		return -1, err
	}
	defer rows.Close()
	canceledclass := []CanceledClass{}
	canceledclass, err = ScanCanceledClass(rows)
	if err != nil {
		return -1, err
	}

	if len(canceledclass) != 1 && len(canceledclass) != 0 {
		return -1, errors.New("IDが一意に定まりませんでした")
	} else if len(canceledclass) == 0 {
		return -1, nil
	}
	return canceledclass[0].ID, nil
}

func (db *DB) ShowCanceled(id int) (int, error) {
	sql := "select * from canceled_class where id=?;"
	rows, err := db.db.Query(sql, id)
	if err != nil {
		return -1, err
	}
	defer rows.Close()
	canceledclass := []CanceledClass{}
	canceledclass, err = ScanCanceledClass(rows)
	if err != nil {
		return -1, err
	}
	if len(canceledclass) != 1 && len(canceledclass) != 0 {
		return -1, errors.New("IDが一意に定まりませんでした")
	} else if len(canceledclass) == 0 {
		return -1, nil
	}
	return canceledclass[0].Canceled, nil
}

// Dayテーブルをみて今日の日付があるのか確認
func (db *DB) IsExistToday(id int, date string) (bool, error) {
	sql := "select * from day where canceled_class_id=?"
	rows, err := db.db.Query(sql, id)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	days := []Day{}
	days, err = ScanDay(rows)
	if err != nil {
		return false, err
	}

	//dayが今日の日付があるのか
	for _, day := range days {
		if day.Date == date {
			return true, nil
		}
	}

	return false, nil
}

func (db *DB) DeleteWhereDayAndClassName(day, className string) (sql.Result, error) {
	result, err := db.db.Exec("delete from kyuko_data where day = ? and class_name = ?", day, className)
	if err != nil {
		return result, err
	}
	return result, err
}

func (db *DB) AddCanceled(id int) (sql.Result, error) {
	sql := "UPDATE canceled_class SET canceled = canceled+1 WHERE id = ?;"
	result, err := db.db.Exec(sql, id)
	if err != nil {
		return result, err
	}
	return result, err
}

func KyukoToCanceled(k KyukoData) (CanceledClass, error) {
	season, err := getSeason(k.Day)
	if err != nil {
		return CanceledClass{}, err
	}

	year, err := getYear(k.Day)
	if err != nil {
		return CanceledClass{}, err
	}

	return CanceledClass{
		Place:      k.Place,
		Weekday:    k.Weekday,
		Period:     k.Period,
		ClassName:  k.ClassName,
		Instructor: k.Instructor,
		Season:     season,
		Year:       year,
	}, nil
}

func getYear(day string) (int, error) {
	strYear := strings.Split(day, "/")[0]
	year, _ := strconv.Atoi(strYear)
	return year, nil
}

func getSeason(day string) (string, error) {
	strMonth := ""
	if strings.Contains(day, "-") {
		strMonth = strings.Split(day, "-")[1]
	} else {
		strMonth = strings.Split(day, "/")[1]
	}

	month, _ := strconv.Atoi(strMonth)

	if month >= 3 && month <= 8 {
		return "spring", nil
	} else if (month >= 9 && month <= 12) || (month >= 1 && month <= 2) {
		return "autumn", nil
	}

	return "", errors.New("Season are not uniquely determined")
}

func (db *DB) DeleteCanceled(id int) (sql.Result, error) {
	sql := "DELETE FROM canceled_class WHERE id = ?;"
	result, err := db.db.Exec(sql, id)
	if err != nil {
		return result, err
	}
	return result, err
}

func (db *DB) deleteReason(id int) (sql.Result, error) {
	sql := "DELETE FROM reason WHERE id = ?;"
	result, err := db.db.Exec(sql, id)
	if err != nil {
		return result, err
	}
	return result, err
}

func (db *DB) DeleteReasonWhere(canceledID int, reason string) (sql.Result, error) {
	sql := "DELETE FROM reason WHERE canceled_class_id = ? AND reason = ?;"
	result, err := db.db.Exec(sql, canceledID, reason)
	if err != nil {
		return result, err
	}
	return result, err
}

func (db *DB) deleteDay(id int) (sql.Result, error) {
	sql := "DELETE FROM day WHERE id = ?;"
	result, err := db.db.Exec(sql, id)
	if err != nil {
		return result, err
	}
	return result, err
}

func (db *DB) DeleteDayWhere(canceledID int, day string) (sql.Result, error) {
	sql := "DELETE FROM day WHERE canceled_class_id = ? AND day = ?;"
	result, err := db.db.Exec(sql, canceledID, day)
	if err != nil {
		return result, err
	}
	return result, err
}
