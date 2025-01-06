## Get active table numbers
curl -H "Content-Type: application/json" http://localhost:8081/activeTableNumbers

## Get tab Id for table
curl -H "Content-Type: application/json" http://localhost:8081/tabIdForTable?table_number=1

## Get tab status for table
curl -H "Content-Type: application/json" http://localhost:8081/tabForTable?table_number=1

## Get invoice for table
curl -H "Content-Type: application/json" http://localhost:8081/invoiceForTable?table_number=1

## Get TODO list for waiter
curl -H "Content-Type: application/json" http://localhost:8081/todoListForWaiter?waiter=w1