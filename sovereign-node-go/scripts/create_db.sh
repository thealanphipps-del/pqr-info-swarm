#!/bin/bash
gcloud compute ssh aellok@instance-20260507-075600 --zone=us-central1-c --project=model-loader-495607-m2 --command="sudo docker exec cockroachdb ./cockroach sql --insecure -e 'CREATE DATABASE IF NOT EXISTS antigravity;'"
