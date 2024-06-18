Temp 

MYSQL_ROOT_PASSWORD
MYSQL_ALLOW_EMPTY_PASSWORD
MYSQL_RANDOM_ROOT_PASSWORD


curl -X 'POST' \
  'http://localhost:8081/go_load.GoLoadService/CreateAccount' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "accountName": "liemldv",
  "password": "liem123456"
}'

curl -X 'POST' \
  'http://localhost:8081/go_load.GoLoadService/CreateSession' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  --cookie-jar cookie.txt \
  -d '{
  "accountName": "liemldv",
  "password": "liem123456"
}'

curl -X 'POST' \
  'http://localhost:8081/go_load.GoLoadService/GetDownloadTaskList' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  --cookie cookie.txt \
  -d '{
  "offset": 0,
  "limit": 10
}'

curl -X 'POST' \
  'http://localhost:8081/go_load.v1.GoLoadService/CreateDownloadTask' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  --cookie cookie.txt \
  -d '{
  "download_type": "1",
  "url": "https://example.com"
}'

curl -X 'POST' \
  'http://localhost:8081/go_load.v1.GoLoadService/CreateSession' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  --cookie-jar cookie.txt \
  -d '{
  "accountName": "liemldv",
  "password": "liem123456"
}'