# E-commerce Order Lifecycle

## Services

### 1. checkout-api (The Storefront)
- Language: Go (Golang)

- Role: The entry point for all user traffic. It acts as the "API Gateway."

- What it does:

    + Receives Traffic: It listens for HTTP GET requests on /ping.

    + Starts the Trace: It initiates the OpenTelemetry Trace. This is the "Root Span" that tracks the entire life of the request.

    + Forwards Request: It makes a synchronous HTTP call to the inventory-service to process the order.

- Responds: Once the inventory service replies, it sends a JSON response back to the user (e.g., "Pong from A").

### 2. inventory-service (The Logic Core)
- Language: Python (FastAPI)

- Role: The heavy lifter. It manages state, database records, and caching.

- What it does:

    + Redis (Caching): It connects to Redis and increments a simple counter (hits). This simulates checking stock or updating a high-speed cache.

    + Postgres (Database): It connects to PostgreSQL and inserts a row into the access_log table with the current timestamp. This simulates saving a permanent order record.

    + SQS (Messaging): It packages the order details into a JSON message and pushes it to an AWS SQS Queue (running locally on LocalStack). This hands off the work for "shipping" to be done later.

    + Logging: It prints a structured JSON log that includes the Trace ID. This allows Grafana/Loki to link the logs to the visual trace.

### 3. shipping-worker (The Warehouse)
- Language: Go (Golang)

- Role: Background Worker. It is completely decoupled from the user's HTTP request.

- What it does:

    + Polls Queue: It sits in an infinite loop, constantly asking SQS: "Do you have any new messages?"

    + Processes Job: When it receives a message (the order from the Inventory Service), it simulates work (e.g., printing a shipping label) by sleeping for 50ms.

    + Cleans Up: After processing, it deletes the message from the queue so it doesn't get processed twice.

## The "Data Flow" Story

When you trigger the demo (e.g., curl http://localhost:8080/ping), here is exactly what happens in milliseconds:

1. You hit Checkout API.

2. Checkout API calls Inventory Service.

3. Inventory Service:

    - ⬆️ Updates Redis counter.

    - 💾 Writes to Postgres.

    - 📨 Sends "Order #123" to SQS.

    - ✅ Returns "200 OK" to Checkout API.

4. Checkout API returns "200 OK" to You. (The HTTP transaction is now finished for you, but the work isn't done!)

5. Shipping Worker sees the message in SQS, wakes up, does the work, and finishes.

This flow generates the perfect "Distributed Trace" because it shows a mix of Sync (User waiting) and Async (Background work) operations.