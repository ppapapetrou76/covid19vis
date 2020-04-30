package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/go-echarts/go-echarts/charts"
)

func main() {

	err, data := readRawData()
	if err != nil {
		fmt.Println(err)
		return
	}
	keys, totalDeaths, totalCases := prepareGraphData(data)
	createGraph(keys, totalCases, totalDeaths)

}

func createGraph(keys []string, totalCases []int, totalDeaths []int) {
	graph := charts.NewLine()
	graph.SetGlobalOptions(charts.TitleOpts{
		Title: "COVID-19 stats",
	})
	graph.AddXAxis(keys).
		AddYAxis("Total Cases", totalCases).
		AddYAxis("Total Deaths", totalDeaths)

	f, err := os.Create("graph.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := graph.Render(f); err != nil {
		fmt.Println(err)
		return
	}
}

func prepareGraphData(data *CountryTimeLineData) (keys []string, totalDeaths []int, totalCases []int) {
	keys = make([]string, 0, len(data.TimeLineItems[0]))
	totalDeaths = make([]int, 0, len(data.TimeLineItems[0]))
	totalCases = make([]int, 0, len(data.TimeLineItems[0]))
	for k := range data.TimeLineItems[0] {
		parsedDate, err := time.Parse("1/2/06", k)
		if err == nil {
			keys = append(keys, parsedDate.Format("2006/01/02"))
		}

	}
	sort.Strings(keys)
	for _, k := range keys {
		parsedDate, err := time.Parse("2006/01/02", k)
		if err == nil {
			originalKey := parsedDate.Format("1/02/06")
			totalDeaths = append(totalDeaths, data.TimeLineItems[0][originalKey].TotalDeaths)
			totalCases = append(totalCases, data.TimeLineItems[0][originalKey].TotalCases)
		}

	}
	return keys, totalDeaths, totalCases
}

func readRawData() (error, *CountryTimeLineData) {
	countryCode := "GR"
	if len(os.Args[1:]) == 1 {
		countryCode = os.Args[1]
	}

	url := fmt.Sprintf("https://thevirustracker.com/free-api?countryTimeline=%s", countryCode)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var data = &CountryTimeLineData{}
	_ = json.Unmarshal(body, data)
	return err, data
}
