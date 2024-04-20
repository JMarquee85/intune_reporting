package main

import (
	"net/http"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// Create another function to compare the WorkspaceOne stuff or use the same function/ graph?

// Line Graph that tracks enrollment numbers of iOS and Android Devices over
func renderIntuneEnrollmentGraph(w http.ResponseWriter, androidEnrollDates []time.Time, iOSEnrollDates []time.Time) error {
	var err error

	// Create a new line instance
	line := charts.NewLine()

	// Prepare the data
	androidData := make([]opts.LineData, 0)
	iOSData := make([]opts.LineData, 0)
	xAxisLabels := make([]string, 0)

	// Calculate the number of weeks to display
	numWeeks := 10 // change this to the number of weeks you want to display

	for i := 0; i < numWeeks; i++ {
		// Calculate the start and end of the week
		startOfWeek := time.Now().AddDate(0, 0, -7*i).Format("2006-01-02")
		endOfWeek := time.Now().AddDate(0, 0, -7*(i+1)).Format("2006-01-02")

		// Add the week to the X-axis labels
		xAxisLabels = append([]string{startOfWeek + " to " + endOfWeek}, xAxisLabels...)

		// Count the number of enrollments for the week
		androidCount := 0
		iOSCount := 0
		for _, date := range androidEnrollDates {
			if date.After(time.Now().AddDate(0, 0, -7*(i+1))) && date.Before(time.Now().AddDate(0, 0, -7*i)) {
				androidCount++
			}
		}
		for _, date := range iOSEnrollDates {
			if date.After(time.Now().AddDate(0, 0, -7*(i+1))) && date.Before(time.Now().AddDate(0, 0, -7*i)) {
				iOSCount++
			}
		}

		// Add the counts to the data
		androidData = append([]opts.LineData{{Value: androidCount}}, androidData...)
		iOSData = append([]opts.LineData{{Value: iOSCount}}, iOSData...)
	}

	// Set the options
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Intune Enrollments By Week",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show: true,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "category",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Type: "value",
		}),
	)

	// Add the data
	line.SetXAxis(xAxisLabels).
		AddSeries("Android", androidData).
		AddSeries("iOS", iOSData)

	// Render the chart
	page := components.NewPage()
	page.AddCharts(line)
	err = page.Render(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil

}

// Other Graph Ideas
// Stacked Bar (Ask Thea to explain)
// Breakdown by region estimates? Number of users in the groups referenced against number of users in those regions in Workspace One
// Maybe this would be good for a stacked bar.
// Would need:
// Number of users in each region group in Intune
// Number of users in each region in Workspace One
func renderBarChart(w http.ResponseWriter, title string, intuneData, workspaceOneData [3]int) error {
	// Create a new bar chart
	bar := charts.NewBar()

	// Set the options
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: title,
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show: true,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "category",
			Data: []string{"iOS", "Android", "Total"},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Type: "value",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: true,
		}),
	)

	// Add the data
	bar.AddSeries("Intune", []opts.BarData{
		{Value: intuneData[0], Name: "iOS"},
		{Value: intuneData[1], Name: "Android"},
		{Value: intuneData[2], Name: "Total"},
	}).SetSeriesOptions(
		charts.WithLabelOpts(opts.Label{
			Show: true,
		}),
	)

	bar.AddSeries("Workspace One", []opts.BarData{
		{Value: workspaceOneData[0], Name: "iOS"},
		{Value: workspaceOneData[1], Name: "Android"},
		{Value: workspaceOneData[2], Name: "Total"},
	}).SetSeriesOptions(
		charts.WithLabelOpts(opts.Label{
			Show: true,
		}),
	)

	// Render the chart
	page := components.NewPage()
	page.AddCharts(bar)
	err := page.Render(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}

// Pie Chart per region of users in Workspace One vs. users in the Intune groups
