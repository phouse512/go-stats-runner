package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"net/http"
	"time"
)

var storageQueue chan Storage
var jobs []Job

func init() {
	log.Printf("initializing server")

	storageQueue = make(chan Storage, 5)

	session, _ := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewSharedCredentials("", "go-stats-runner"),
	})
	uploader := s3manager.NewUploader(session)

	for i := 0; i < 5; i++ {
		var new_storage Storage
		new_storage = NewS3Storage(uploader, session)
		storageQueue <- new_storage
	}

	gauges_job := new(GaugesJob)
	rescuetime_job := new(RescueTimeJob)
	todoist_job := new(TodoistJob)
	jobs = append(jobs, gauges_job)
	jobs = append(jobs, rescuetime_job)
	jobs = append(jobs, todoist_job)
}

func main() {
	http.HandleFunc("/jobs/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request: ", r.URL.Path[1:])

	query_params := r.URL.Query()

	var date_obj time.Time
	var err error
	if date, exists := query_params["date"]; exists {
		date_obj, err = time.Parse("01-02-2006", date[0])

		if err != nil {
			log.Println("Improper date format given: ", err.Error())
			return
		}
	} else {
		log.Println("date param not attached, exiting")
		return
	}

	for _, job := range jobs {
		storage := getStorageService()
		job.getData("string", storage, date_obj)

		queueStorage(storage)
	}
}

func getStorageService() Storage {
	select {
	case store := <-storageQueue:
		return store
	}
}

func queueStorage(s Storage) {
	select {
	case storageQueue <- s:

	case <-time.After(1 * time.Second):
		log.Printf("Couldn't find a home")
	}
}
