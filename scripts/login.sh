#!/bin/bash

body='{
  "email":"alex@gmail.com",
  "password":"perfectly-safe"
}'

curl -i \
  -d "$body" \
  -H "Content-Type: application/json" \
  -X POST http://localhost:8080/api/login
