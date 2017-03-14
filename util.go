package main

import (
	"fmt"
	"time"
)

func getDateStringComponents(t time.Time) (string, string, string) {
	/*
	* returns the mm, dd, and yyyy components in that order
	 */
	year, month, date := t.Date()

	var day_string, month_string, year_string string
	if date < 10 {
		day_string = fmt.Sprintf("0%d", date)
	} else {
		day_string = fmt.Sprintf("%d", date)
	}

	if month < 10 {
		month_string = fmt.Sprintf("0%d", month)
	} else {
		month_string = fmt.Sprintf("%d", month)
	}

	year_string = fmt.Sprintf("%d", year)

	return month_string, day_string, year_string
}
