# Inventory API Server

## 재고 차량 조회
GET /inventory/cars

```
curl -X GET http://localhost:8888/inventory/cars
```

## 재고 차량 등록
POST /inventory/cars

```
curl -X POST -d '{"vin":"MWELKFJEWLFJDKFML","model":"b","make":"c","year":2015,"msrp":100000,"status":"Ordered","booked":"y","listed":"n"}' http://localhost:8888/inventory/cars
```

## 재고 차량 삭제
DELETE /inventory/cars/:vin

```
curl -X DELETE http://localhost:8888/inventory/cars/$vin
ex) curl -X DELETE http://localhost:8888/inventory/cars/MWELKFJEWLFJDKFML
```