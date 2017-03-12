package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Job interface {
	getData(id string, s Storage, t time.Time)
}

type GaugesJob struct {
	id string
}

func (j *GaugesJob) generateKey(t time.Time, user string) string {
	year, month, date := t.Date()

	var day_string, month_string string
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

	file_key := fmt.Sprintf("%s_%s_%d", month_string, day_string, year)
	month_key := fmt.Sprintf("%s-%d", strings.ToLower(t.Month().String()[0:3]), year%100)

	return fmt.Sprintf("gauges/v0/%s/%s/output_%s", user, month_key, file_key)
}

func (j *GaugesJob) getData(id string, s Storage, t time.Time) {

	year, month, date := t.Date()

	var day_string, month_string string
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

	date_string := fmt.Sprintf("%d-%s-%s", year, month_string, day_string)
	url_string := fmt.Sprintf("https://secure.gaug.es/gauges/%s/traffic?date=%s", GAUGES_SITE_ID, date_string)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url_string, nil)

	req.Header.Add("X-Gauges-Token", GAUGES_SECRET)
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Received error for GaugesJob: %s", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	key := j.generateKey(t, "phouse512")

	obj_to_store := DataBlob{
		data: body,
		key:  key,
	}

	s.storeData(obj_to_store)
}
