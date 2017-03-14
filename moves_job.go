package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type MovesJob struct {
	id string
}

func (j *MovesJob) generateKey(t time.Time, user string) string {
	month, date, year := getDateStringComponents(t)

	current_ts := time.Now()

	file_key := fmt.Sprintf("%s_%s_%s", month, date, year)
	month_key := fmt.Sprintf("%s-%d", strings.ToLower(t.Month().String()[0:3]), t.Year()%100)

	return fmt.Sprintf("moves/v0/%s/%s/output_%s_%d.json", user, month_key, file_key, current_ts.Unix())
}

func (j *MovesJob) getData(id string, s Storage, t time.Time) {
	log.Println("MovesJob is beginning")
	month, day, year := getDateStringComponents(t)
	query_date := fmt.Sprintf("%s%s%s", year, month, day)

	url_string := fmt.Sprintf("https://api.moves-app.com/api/1.1/user/storyline/daily/%s", query_date)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url_string, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", MOVES_KEY))
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
