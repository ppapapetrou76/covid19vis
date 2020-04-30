package main


type CountryTimeLineData struct {
	TimeLineItems []map[string]TimeLineItem
}

type TimeLineItem struct {
	NewDailyCases int `json:"new_daily_cases"`
	NewDailyDeaths int `json:"new_daily_deaths"`
	TotalCases int `json:"total_cases"`
	TotalRecoveries int `json:"total_recoveries"`
	TotalDeaths int `json:"total_deaths"`
}
