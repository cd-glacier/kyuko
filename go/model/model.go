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
	sql := "INSERT INTO `canceled_class` (canceled, place, week, period, year, season, class_name, instructor) VALUES(?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE `canceled`=?;"

	result, err := db.db.Exec(sql, c.Canceled, c.Place, c.Weekday, c.Period, c.Year, c.Season, c.ClassName, c.Instructor, c.Canceled)
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

func (db *DB) SelectAll() ([]KyukoData, error) {
	kyukoData := []KyukoData{}
	rows, err := db.db.Query("select * from kyuko_data")
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
	if len(canceledclass) != 1 {
		return -1, errors.New("IDが一意に定まりませんでした")
	} else if len(canceledclass) == 0 {
		return -1, errors.New("DBに存在しないデータです")
	}
	return canceledclass[0].ID, nil
}

func (db *DB) DeleteWhereDayAndClassName(day, className string) (sql.Result, error) {
	result, err := db.db.Exec("delete from kyuko_data where day = ? and class_name = ?", day, className)
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
	strMonth := strings.Split(day, "/")[1]
	month, _ := strconv.Atoi(strMonth)

	if month >= 3 && month <= 8 {
		return "spring", nil
	} else if (month >= 9 && month <= 12) || (month >= 1 && month <= 2) {
		return "autumn", nil
	}

	return "", errors.New("Season are not uniquely determined")
}
