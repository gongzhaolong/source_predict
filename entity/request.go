package entity

// 请求
type DataForAnalysisReq struct {
	Metric  string    `json:"metric"`
	Time    []int64   `json:"time"`
	DataAll []float64 `json:"data_all"`
}
