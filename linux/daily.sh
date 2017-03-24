#!/bin/bash

DATE=`date +%m-%d-%Y -d "1 day ago"`
echo $DATE

curl -v "http://127.0.0.1:8080/jobs/?date=$DATE"
