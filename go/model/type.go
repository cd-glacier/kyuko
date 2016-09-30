package model

type KyukoData struct {
	ID         int    `json:id`
	Place      int    `json:place`
	Weekday    int    `json:week`
	Period     int    `json:period`
	Day        string `json:date`
	ClassName  string `josn:className`
	Instructor string `json:instructor`
	Reason     string `json:reason`
}

type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)
