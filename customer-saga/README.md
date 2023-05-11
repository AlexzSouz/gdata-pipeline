# Customer Orders Saga

The customer orders saga implements the microservice saga patterns, and this data pipeline service digest daily customers `.csv` files and publish them to a message broker for further processing by other systems.

## CSV Input

The partner uploads 3 files to a shared S3 bucket daily. The object keys have the form $TYPE_$DATE.csv, e.g. customers_20220130.csv.

- **Customers:** List of all customers. A customer can have multiple orders.
- **Orders:** List of all orders.
- **Items:** List of all items.
