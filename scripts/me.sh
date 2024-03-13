#!/bin/bash

jwt='eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhdXRoIiwiZXhwIjoxNzEwMzU1OTU1LCJhdWQiOiJ0b2RvIiwic3ViIjoiMjhjZWE0NzAtYTVkZC00Y2YyLTkzNGUtNzliNzBkZjhjNTJkIiwibmFtZSI6IkFsZXggVGhlIE1hZCIsInJvbGVzIjpbIlRPRE8iXX0.HTso_GtGtfEoiskDNAI4PDi_Kbmy2cm2XO4LHxtTYmI'

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $jwt" \
  -X GET http://localhost:8080/api/me
