package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type RescueTimeJob struct {
	id string
}

func (j *RescueTimeJob) generateKey(t time.Time, user string) string {
	month, date, year := getDateStringComponents(t)

	file_key := fmt.Sprintf("%s_%s_%s", month, date, year)
	month_key := fmt.Sprintf("%s-%d", strings.ToLower(t.Month().String()[0:3]), t.Year()%100)

	return fmt.Sprintf("rescuetime/v0/%s/%s/output_%s_%d.json", user, month_key, file_key, t.Unix())
}

func (j *RescueTimeJob) getData(id string, s Storage, t time.Time) {
	log.Println("RescueTimeJob is beginning")
	month, day, year := getDateStringComponents(t)
	query_date := fmt.Sprintf("%s-%s-%s", year, month, day)

	url_string := fmt.Sprintf("https://www.rescuetime.com/anapi/data?key=%s&perspective=interval&restrict_kind=document&interval=minute&restrict_begin=%s&restrict_end=%s&format=json", RESCUETIME_KEY, query_date, query_date)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url_string, nil)
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Received error for RescueTimeJob: ", err)
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
