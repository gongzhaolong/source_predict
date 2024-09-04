package function

import (
	"Source_Predict/constant"
	"Source_Predict/ecode"
	"Source_Predict/entity"
	"fmt"
	"sort"
	"strconv"
	"time"
)

func BatchReadAndAnalysis(req *entity.DataForAnalysisReq) (err *ecode.ErrCode, resp entity.ResourceResp) {
	//1、将数据按月份拆分
	detailInfo := new(entity.DataDetail)
	groupedData, groupedValues, months := SeparateByMonth(req.Time, req.DataAll)
	detailInfo.Metric = req.Metric
	for _, key := range months {
		detailInfo.Data = append(detailInfo.Data, entity.Data{MonthTime: groupedData[key], PerData: groupedValues[key]})
	}
	//2、数据周期性分析及预测
	period, predict := DataAnalysis(detailInfo.Data, detailInfo.Metric, req.DataAll)
	//3、返回响应
	return err, entity.ResourceResp{req.Metric, period, predict}
}

func SeparateByMonth(timestamps []int64, values []float64) (map[string][]int64, map[string][]float64, []string) {
	// 分组按月份
	groupedData := make(map[string][]int64)
	groupedValues := make(map[string][]float64)

	for i, ts := range timestamps {
		t := time.Unix(ts, 0).UTC()
		key := fmt.Sprintf("%d-%02d", t.Year(), t.Month())
		groupedData[key] = append(groupedData[key], ts)
		groupedValues[key] = append(groupedValues[key], values[i])
	}

	// 按时间戳顺序提取月份
	var months []string
	for ts := range timestamps {
		t := time.Unix(timestamps[ts], 0).UTC()
		key := fmt.Sprintf("%d-%02d", t.Year(), t.Month())
		if !contains(months, key) {
			months = append(months, key)
		}
	}

	// 按时间戳顺序排序月份
	sort.Slice(months, func(i, j int) bool {
		timeI, _ := time.Parse("2006-01", months[i])
		timeJ, _ := time.Parse("2006-01", months[j])
		return timeI.Before(timeJ)
	})
	return groupedData, groupedValues, months
}

// 检查切片中是否包含某个元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func DataAnalysis(Data []entity.Data, metric string, data []float64) (period uint, predict float64) {
	if Data == nil {
		return constant.Noperiod, 0.0
	}
	length, total_month, EachMonthMax := len(Data), make([]int, 0), make([]float64, 0)
	if metric == "disk_usage" { //如果是disk_usage指标，认为是无周期的，用一元线性回归拟合近十天数据
		v1 := 0.0
		if len(data) >= constant.RecentDay {
			v1 = PredictCelestial(data[len(data)-constant.RecentDay:])
		} else if len(data) >= 3 && len(data) < constant.RecentDay {
			v1 = PredictCelestial(data)
		} else if len(data) > 0 && len(data) <= 2 {
			v1 = RecentAvg(data)
		}
		return constant.Noperiod, roundToFourDecimalPlaces(v1)
	}
	//step1 周期性判断 计算每组数据中的最大值和对应的时间戳，将最大值保存在一个切片中
	for _, month_data := range Data {
		maxValue, maxIndex := month_data.PerData[0], 0
		for index, value := range month_data.PerData {
			if value > maxValue {
				maxValue = value
				maxIndex = index
			}
		}
		maxTimestamp := time.Unix(month_data.MonthTime[maxIndex], 0)

		// 检查最大值是否在大促时间范围内 （大促前1天到大促后3天）
		isPromotion := false
		year, month := maxTimestamp.Year(), maxTimestamp.Month()
		promoStart := time.Date(year, month, int(month)-constant.Celestial_Start, 0, 0, 0, 0, time.UTC)
		promoEnd := time.Date(year, month, int(month)+constant.Celestial_End, 0, 0, 0, 0, time.UTC)
		if (maxTimestamp.After(promoStart) || maxTimestamp.Equal(promoStart)) && (maxTimestamp.Before(promoEnd) || maxTimestamp.Equal(promoEnd)) {
			isPromotion = true
		}
		//判断最大值的时间戳是否在大促的区间内，如果是则total_month[i]=1，否则total_month[i]=0
		if isPromotion {
			total_month = append(total_month, 1)
			EachMonthMax = append(EachMonthMax, maxValue)
		}
	}
	//统计total_month中1的个数，如果大于70%则认为是周期性的数据，period=1 否则period=0
	sum := 0
	for _, val := range total_month {
		if val == 1 {
			sum += val
		}
	}
	//step2 数据预测
	//1、周期性数据
	//a. 根据每月的极大值用一元线性回归（y=ax+b）拟合，得到预测值v1
	v1 := 0.0
	if length != 1 && float64(sum) >= float64(length)*constant.Proportion {
		v1 = PredictCelestial(EachMonthMax)
	}
	//b. 计算出大促前十天的均值v2
	v2 := RecentAvg(data)
	//c.计算预测值 v=0.5*v1+0.5*v2
	if v1 != 0.0 {
		predict = constant.V1_weight*v1 + constant.V2_weight*v2
		return constant.Isperiod, roundToFourDecimalPlaces(predict)
	}
	//2、非周期数据
	//计算前十天平均值
	return constant.Noperiod, roundToFourDecimalPlaces(v2)

}

func PredictCelestial(data []float64) (predict float64) { //一元线性回归拟合数据
	x := make([]float64, len(data))
	for i := 0; i < len(data); i++ {
		x[i] = float64(i) + 1.0
	}
	b0, b1 := linearRegression(x, data)
	xNew := float64(len(x)) + 1.0
	predict = Predict(xNew, b0, b1)
	return
}

// Function to calculate the mean of a slice of float64
func mean(data []float64) float64 {
	sum := 0.0
	for _, value := range data {
		sum += value
	}
	return sum / float64(len(data))
}

// Function to calculate the covariance of two slices of float64
func covariance(x, y []float64) float64 {
	xMean := mean(x)
	yMean := mean(y)
	cov := 0.0
	for i := 0; i < len(x); i++ {
		cov += (x[i] - xMean) * (y[i] - yMean)
	}
	return cov / float64(len(x))
}

// Function to calculate the variance of a slice of float64
func variance(x []float64) float64 {
	xMean := mean(x)
	varSum := 0.0
	for _, value := range x {
		varSum += (value - xMean) * (value - xMean)
	}
	return varSum / float64(len(x))
}

// Linear regression function to calculate the coefficients
func linearRegression(x, y []float64) (float64, float64) {
	b1 := covariance(x, y) / variance(x)
	b0 := mean(y) - b1*mean(x)
	return b0, b1
}

// Function to make predictions based on the regression coefficients
func Predict(x float64, b0, b1 float64) float64 {
	predict := b0 + b1*x
	return predict
}

// 保留四位小数的函数
func roundToFourDecimalPlaces(num float64) float64 {
	// 使用 fmt.Sprintf 保留4位小数并返回字符串
	formattedNum := fmt.Sprintf("%.4f", num)

	// 将字符串转换回浮点数
	roundedNum, err := strconv.ParseFloat(formattedNum, 64)
	if err != nil {
		fmt.Println("转换错误:", err)
		return 0.0
	}
	return roundedNum
}

// 计算近10天均值
func RecentAvg(data []float64) float64 {
	v2 := 0.0
	recent := constant.RecentDay
	sum_recent_day := 0.0
	if len(data) >= recent {
		for i := len(data) - 1; i >= len(data)-recent; i-- {
			sum_recent_day += data[i]
		}
		v2 = sum_recent_day / float64(recent)
	} else {
		for _, v := range data {
			sum_recent_day += v
		}
		v2 = sum_recent_day / float64(len(data))
	}
	return v2
}
