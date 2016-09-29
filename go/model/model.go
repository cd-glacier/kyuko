package model

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	db *sql.DB
}

type KyukoData struct {
	ID         int    `json:id`
	Place      int    `json:place`
	Week       int    `json:week`
	Period     int    `json:period`
	Day        string `json:date`
	ClassName  string `josn:className`
	Instructor string `json:instructor`
	Reason     string `json:reason`
}

func (this *DB) Connect() error {
	var err error
	this.db, err = sql.Open("mysql", "root:password@/kyuko")
	return err
}

func (this *DB) Close() error {
	err := this.db.Close()
	return err
}

func (this *DB) Insert(k KyukoData) (sql.Result, error) {
	result, err := this.db.Exec("insert into kyuko_data values(?, ?, ?, ?, ?, ?, ?, ?)", 0, k.Place, k.Week, k.Period, k.Day, k.ClassName, k.Instructor, k.Reason)
	return result, err
}

func ScanAll(rows *sql.Rows) ([]KyukoData, error) {
	kyukoSlice := []KyukoData{}
	var err error
	for rows.Next() {
		var k KyukoData
		if err = rows.Scan(&k.ID, &k.Place, &k.Week, &k.Period, &k.Day, &k.ClassName, &k.Instructor, &k.Reason); err != nil {
			return kyukoSlice, err
		}
		kyukoSlice = append(kyukoSlice, k)
	}
	return kyukoSlice, err
}

func (this *DB) SelectAll() ([]KyukoData, error) {
	kyukoSlice := []KyukoData{}
	rows, err := this.db.Query("select * from kyuko_data")
	if err != nil {
		return kyukoSlice, err
	}
	defer rows.Close()
	kyukoSlice, err = ScanAll(rows)
	if err != nil {
		return kyukoSlice, err
	}
	return kyukoSlice, err
}

func (this *DB) DeleteWhereDayAndClassName(day, className string) (sql.Result, error) {
	result, err := this.db.Exec("delete from kyuko_data where day = ? and class_name = ?", day, className)
	if err != nil {
		return result, err
	}
	return result, err
}
