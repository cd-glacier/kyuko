package model

type KyukoData struct {
	ID         int    `json:"id"`
	Canceled   int    `json:canceled`
	Place      int    `json:"place"`
	Weekday    int    `json:"week"`
	Period     int    `json:"period"`
	Day        string `json:"date"`
	ClassName  string `josn:"className"`
	Instructor string `json:"instructor"`
	Reason     string `json:"reason"`
}

type CanceledClass struct {
	ID         int    `json:"id"`
	Canceled   int    `json:canceled`
	Place      int    `json:"place"`
	Weekday    int    `json:"week"`
	Period     int    `json:"period"`
	Year       int    `json:"year"`
	Season     string `json:"season"`
	ClassName  string `josn:"className"`
	Instructor string `json:"instructor"`
	Days       []Day
	Reasons    []Reason
}

type Reason struct {
	ID              int    `json:"id"`
	CanceledClassID int    `json:"canceled_class_id"`
	Reason          string `josn:"reason"`
}

type Day struct {
	ID              int    `json:"id"`
	CanceledClassID int    `json:"canceled_class_id"`
	Date            string `josn:"reason"`
}
