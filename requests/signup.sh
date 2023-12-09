#!/bin/bash

body='{
  "email": "alex@gmail.com",
  "fullName": "Alex The Mad",
  "password": "perfectly-safe"
}'

curl -i \
  -d "$body" \
  -H "Content-Type: application/json" \
  -X POST http://localhost:8080/api/signup
