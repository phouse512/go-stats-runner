### go-stats-runner

This service is responsible for collecting and querying all of my personal data from different services. It is built with go and the aws-sdk-go library.


#### How it Works

Currently it doesn't handle many users, and it doesn't run all of the jobs concurrently (yet). Each job will query the daily data for each user and store it in S3 for future processing.


#### Operations

Run with upstart. To redeploy, pull on the server, and run `go build` to compile the binary. Restart using the following command: `sudo initctl restart go-stats-runner`
