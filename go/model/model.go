package model

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	db *sql.DB
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
	result, err := this.db.Exec("insert into kyuko_data values(?, ?, ?, ?, ?, ?, ?, ?)", 0, k.Place, k.Weekday, k.Period, k.Day, k.ClassName, k.Instructor, k.Reason)
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

func (this *DB) SelectAll() ([]KyukoData, error) {
	kyukoData := []KyukoData{}
	rows, err := this.db.Query("select * from kyuko_data")
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

func (this *DB) DeleteWhereDayAndClassName(day, className string) (sql.Result, error) {
	result, err := this.db.Exec("delete from kyuko_data where day = ? and class_name = ?", day, className)
	if err != nil {
		return result, err
	}
	return result, err
}
