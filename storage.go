package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
)

type Storage interface {
	storeData(DataBlob)
}

type S3Storage struct {
	uploader *s3manager.Uploader
	session  *session.Session
}

func NewS3Storage(uploader *s3manager.Uploader, session *session.Session) *S3Storage {
	s := &S3Storage{
		uploader: uploader,
		session:  session,
	}

	return s
}

func (s *S3Storage) storeData(blob DataBlob) {
	log.Printf("storing message with key: %s", blob.key)

	result, err := s.uploader.Upload(&s3manager.UploadInput{
		Body:   bytes.NewReader(blob.data),
		Bucket: aws.String("go-stats-runner"),
		Key:    aws.String(blob.key),
	})

	if err != nil {
		log.Fatalln("Failed to upload to S3 with err: ", err)
	}

	log.Printf("Result: %s", result.Location)
	log.Println("successfully uploaded %s to S3", blob.key)
}

func (s *S3Storage) doesExist(key string) bool {
	params := &s3.HeadObjectInput{
		Bucket: aws.String("go-stats-runner"),
		Key:    aws.String(key),
	}

	svc := s3.New(s.session)
	resp, err := svc.HeadObject(params)

	if err != nil {
		log.Println(fmt.Sprintf("Error when checking key-existence for key: %s with error: %s", key, err.Error()))
		return false
	}

	log.Println(fmt.Sprintf("Key: %s does exist, received response: %s", key, resp))
	return true
}
