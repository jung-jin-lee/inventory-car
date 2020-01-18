#!/bin/sh

PORT=8888
BASE_URI=localhost:$PORT
vin=a


# GET /inventory/cars
curl -X GET $BASE_URI/inventory/cars

# POST /inventory/cars
curl -X POST -d '{"vin":"MWELKFJEWLFJDKFML","model":"b","make":"c","year":2015,"msrp":100000,"status":"Ordered","booked":"y","listed":"n"}' $BASE_URI/inventory/cars

# DELETE /inventory/cars/:vin
curl -X DELETE $BASE_URI/inventory/cars/$vin