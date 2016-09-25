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
	Date       string `json:date`
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
	result, err := this.db.Exec("insert into kyuko_data values(?, ?, ?, ?, ?, ?, ?, ?)", 0, k.Place, k.Week, k.Period, k.Date, k.ClassName, k.Instructor, k.Reason)
	return result, err
}

func (this.*DB)
