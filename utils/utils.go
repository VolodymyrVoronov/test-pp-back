package utils

import (
	"strings"
	"test-pp-back/models"
	"time"

	"gonum.org/v1/gonum/floats"
)

func ParseDate(dateStr string) (float64, error) {
	if strings.Contains(dateStr, "/") {
		date, err := time.Parse("02/01/2006", dateStr)
		if err != nil {
			return 0, err
		}

		return float64(date.Unix()), nil
	} else if strings.Contains(dateStr, "-") {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return 0, err
		}

		return float64(date.Unix()), nil
	}

	return 0, nil
}

func FormatDate(timestamp float64) string {
	dartTime := time.Unix(int64(timestamp), 0)

	return dartTime.Format("2006-01-02")
}

func ExtractData(stocks []models.StockData) ([]float64, []float64, error) {
	var dates []float64
	var closingPrices []float64

	for _, stock := range stocks {
		date, err := ParseDate(stock.Date)
		if err != nil {
			return nil, nil, err
		}

		dates = append(dates, date)
		closingPrices = append(closingPrices, stock.Close)
	}

	return dates, closingPrices, nil
}

func LinearRegression(x, y []float64) (slope, intercept float64) {
	n := float64(len(x))
	sumX := floats.Sum(x)
	sumY := floats.Sum(y)
	sumXY := floats.Dot(x, y)
	sumX2 := floats.Dot(x, x)

	slope = (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	intercept = (sumY - slope*sumX) / n

	return slope, intercept
}

func PredictPricesForGivenDays(slope, intercept, lastDate float64, days int) []models.Prediction {
	var predictions []models.Prediction

	for i := 0; i < days; i++ {
		nextDate := lastDate + float64(i*86400)
		predictedPrice := slope*nextDate + intercept
		predictions = append(predictions, models.Prediction{
			Date:  FormatDate(nextDate),
			Price: predictedPrice,
		})
	}

	return predictions
}
