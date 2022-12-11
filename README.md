# Evo Tech Test Task
## How to setup
You will need Docker for this app to run. To startup:
- download repository
- in the repository folder execute command   
```curl
docker compose -f "docker-compose.yml" up -d --build
```

## API Documentation

All API is accessed from localhost:8000     
To upload csv file use ```POST /upload```. CSV file must have the same structure as the one in parse/test_data/example.csv   
Example:
```curl
curl -i -X POST http://localhost:8000/upload -F "file=@example.csv" -H "Content-Type: multipart/form-data"
```
To get transaction by ID use ```GET /transaction/id```.    
Example:
```curl
curl -i -X GET http://localhost:8000/transaction/1 
```
To get transactions with filters use ```GET /transactions?filter1=value1&filter2=value2```
Available filters:
- terminal_ids: accepts comma separated list of values. Example:
```curl
curl -i -X GET http://localhost:8000/transactions?terminal_ids=3506,3511
```
- status: accepts string values (accepted, declined). Example:
```curl
curl -i -X GET http://localhost:8000/transactions?status=declined
```
- payment_type: accepts string values (cash, card). Example:
```curl
curl -i -X GET http://localhost:8000/transactions?payment_type=card
```
- date_post_from: accepts date values (date format: YYYY-MM-DDTHH:MM:SS). Filters results by leaving results, that have field date_post less or equal. Example:
```curl
curl -i -X GET http://localhost:8000/transactions?date_post_from=2022-8-24T17:15:13
```
- date_post_to: accepts date values (date format: YYYY-MM-DDTHH:MM:SS). Filters results by leaving results, that have field date_post less or equal. Example:
```curl
curl -i -X GET http://localhost:8000/transactions?date_post_to=2022-8-24T17:15:13
```
- payment_narrative: accepts string values. Filters results by searching a substring in field 'payment narrative' Example:
```curl
curl -i -X GET http://localhost:8000/transactions?payment_narrative=27122
```


Example:
```curl
curl -i -X GET "http://localhost:8000/transactions?status=declined&payment_type=card&payment_narrative=27122"
```
