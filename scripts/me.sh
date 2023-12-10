#!/bin/bash

jwt='eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhdXRoIiwiZXhwIjoxNzAyMTYxMzg4LCJhdWQiOiJ0b2RvIiwic3ViIjoidXVpZCIsIm5hbWUiOiJBbGV4IFRoZSBNYWQiLCJyb2xlcyI6WyJUT0RPIl19.DSIhbioL9esS0gsiliNl9rUFYaLaZAciVvNG7e7OxyI'

curl -i \
  -H "Accept: application/json" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $jwt" \
  -X GET http://localhost:8080/api/me
