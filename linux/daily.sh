#!/bin/bash

DATE=`date +%m-%d-%Y`
echo $DATE

curl -v "http://127.0.0.1:8080/jobs/?date=$DATE"
