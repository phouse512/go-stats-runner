package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TodoistJob struct {
	id string
}

func (j *TodoistJob) generateKey(t time.Time, user string) string {
	month, date, year := getDateStringComponents(t)

	file_key := fmt.Sprintf("%s_%s_%s", month, date, year)
	month_key := fmt.Sprintf("%s-%d", strings.ToLower(t.Month().String()[0:3]), t.Year()%100)

	return fmt.Sprintf("todoist/v0/%s/%s/output_%s_%d.json", user, month_key, file_key, t.Unix())
}

func (j *TodoistJob) getData(id string, s Storage, t time.Time) {
	log.Println("TodoistJob is beginning")

	url_string := fmt.Sprintf("https://todoist.com/API/v7/backups/get")
	resp, err := http.PostForm(url_string, url.Values{"token": {TODOIST_KEY}})

	if err != nil {
		log.Printf("Received error for TodoistJob: ", err)
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
