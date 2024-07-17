package utils

import (
	"fmt"
	"io"
	"os"
	"strings"
	"test-pp-back/models"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
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

func DrawStockGraph(stocks []models.StockData) *charts.Kline {
	kline := charts.NewKLine()

	var adjustedStockData []models.AdjustedStockData

	for _, stock := range stocks {
		date, err := time.Parse("2006-01-02", stock.Date)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}

		formattedDate := date.Format("2006/1/2")

		newStockData := models.AdjustedStockData{
			Date: formattedDate,
			Data: [4]float64{stock.Open, stock.High, stock.Low, stock.Close},
		}

		adjustedStockData = append(adjustedStockData, newStockData)
	}

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for i := 0; i < len(adjustedStockData); i++ {
		x = append(x, adjustedStockData[i].Date)
		y = append(y, opts.KlineData{Value: adjustedStockData[i].Data})
	}

	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Stock Graph",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: opts.Bool(true),
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	// kline.SetXAxis(x).AddSeries("kline", y)

	return kline
}

func CreateGraphFile(stocks []models.StockData) error {
	page := components.NewPage()

	page.AddCharts(
		DrawStockGraph(stocks),
	)

	dynamicFileName := time.Now().Format("2006-01-02_15-04-05")
	graphFileName := fmt.Sprintf("graphs/%s.html", dynamicFileName)

	files, err := os.ReadDir("graphs")
	if err != nil {
		return err
	}

	for _, f := range files {
		os.Remove(fmt.Sprintf("graphs/%s", f.Name()))
	}

	file, err := os.Create(graphFileName)
	if err != nil {
		fmt.Println(err)

		return err
	}
	// defer file.Close()

	page.Render(io.MultiWriter(file))

	return nil
}
