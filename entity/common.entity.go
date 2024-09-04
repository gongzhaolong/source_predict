package entity

type DataDetail struct {
	Metric string `json:"metric"`
	Data   []Data `json:"data"`
}
type Data struct {
	MonthTime []int64   `json:"month_time"`
	PerData   []float64 `json:"month_data"`
}
