## Creating a tab
curl -X POST -H "Content-Type: application/json" -d '{"table_number": 1, "waiter": "w1"}' http://localhost:8080/openTab

## Placing an order
curl -X POST -H "Content-Type: application/json" -d '{"tab_id": "2qwuWZba48SRux8AkPcFQTSdoYr", "menu_items": [1,2]}' http://localhost:8080/placeOrder

## Marking drinks as served
curl -X POST -H "Content-Type: application/json" -d '{"tab_id": "2qwuWZba48SRux8AkPcFQTSdoYr", "menu_numbers": [1,2]}' http://localhost:8080/markDrinksServed

## Closing tab
curl -X POST -H "Content-Type: application/json" -d '{"tab_id": "2qwuWZba48SRux8AkPcFQTSdoYr", "amount_paid": 3.0}' http://localhost:8080/closeTab