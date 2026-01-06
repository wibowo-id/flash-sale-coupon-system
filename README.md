# Flash Sale Coupon System - REST API

A scalable REST API built with Go and Gin framework for managing flash sale coupons. The system handles high concurrency, guarantees strict data consistency, and is designed to prevent race conditions during coupon claims.

## Prerequisites

Before running this application, ensure you have the following installed:

- **Docker Desktop** (or Docker Engine + Docker Compose)
  - Docker version 20.10 or higher
  - Docker Compose version 2.0 or higher

Alternatively, if running without Docker:
- **Go** 1.21 or higher
- **PostgreSQL** 15 or higher

## How to Run

### Using Docker Compose (Recommended)

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd flash-sale-coupon-system
   ```

2. Start the application and database:
   ```bash
   docker-compose up --build
   ```

   This command will:
   - Build the Go application Docker image
   - Start a PostgreSQL database container
   - Start the API server container
   - Automatically run database migrations

3. The API will be available at `http://localhost:8080`

4. **Verify the deployment:**
   ```bash
   # Check if containers are running
   docker ps
   
   # Check API health
   curl http://localhost:8080/health
   
   # View API logs
   docker logs coupon_api
   
   # View database logs
   docker logs coupon_postgres
   ```

5. To stop the services:
   ```bash
   docker-compose down
   ```

6. To stop and remove volumes (clears database):
   ```bash
   docker-compose down -v
   ```

7. To run in background (detached mode):
   ```bash
   docker-compose up -d --build
   ```

### Running Locally (Without Docker)

1. Start PostgreSQL database:
   ```bash
   # Using Docker for just the database
   docker run -d --name postgres \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=postgres \
     -e POSTGRES_DB=coupon_db \
     -p 5432:5432 \
     postgres:15-alpine
   ```

2. Create `.env` file (optional, but recommended):
   ```bash
   # Create .env file with your configuration
   cat > .env << EOF
   DB_DEFAULT=postgresql
   DB_PG_HOST=localhost
   DB_PG_DATABASE=coupon_db
   DB_PG_USERNAME=postgres
   DB_PG_PASSWORD=postgres
   DB_PG_PORT=5432
   PORT=8080
   EOF
   ```
   
   Alternatively, you can set environment variables directly:
   ```bash
   export DB_DEFAULT=postgresql
   export DB_PG_HOST=localhost
   export DB_PG_DATABASE=coupon_db
   export DB_PG_USERNAME=postgres
   export DB_PG_PASSWORD=postgres
   export DB_PG_PORT=5432
   export PORT=8080
   ```
   
   **Note:** You can also use `DATABASE_URL` directly if preferred (it takes precedence over individual DB variables).

3. Install dependencies and run:
   ```bash
   go mod download
   go run main.go
   ```
   
   **Note:** The application will automatically load `.env` file if it exists. If not found, it will use environment variables or default values.

## Quick Start Guide (After Docker Deployment)

After successfully deploying the application with Docker, follow these steps to use the API:

### 1. Verify Deployment

```bash
# Check container status
docker ps

# Test health endpoint
curl http://localhost:8080/health
```

Expected output: `{"status":"ok"}`

### 2. Available Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/coupons` | Create a new coupon |
| POST | `/api/coupons/claim` | Claim a coupon for a user |
| GET | `/api/coupons/{name}` | View coupon details |
| GET | `/health` | Health check |

### 3. API Usage Examples

#### A. Create a New Coupon
```bash
curl -X POST http://localhost:8080/api/coupons \
  -H "Content-Type: application/json" \
  -d '{"name": "PROMO_SUPER", "amount": 100}'
```

**Response:** `201 Created` (no body)

#### B. Claim a Coupon
```bash
curl -X POST http://localhost:8080/api/coupons/claim \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user_12345", "coupon_name": "PROMO_SUPER"}'
```

**Possible responses:**
- `200 OK` - Successfully claimed
- `409 Conflict` - User has already claimed this coupon
- `400 Bad Request` - Stock exhausted or invalid request
- `404 Not Found` - Coupon not found

#### C. View Coupon Details
```bash
curl http://localhost:8080/api/coupons/PROMO_SUPER
```

**Response:**
```json
{
  "name": "PROMO_SUPER",
  "amount": 100,
  "remaining_amount": 99,
  "claimed_by": ["user_12345"]
}
```

### 4. Response Codes

- `200 OK` - Request successful
- `201 Created` - Resource successfully created
- `400 Bad Request` - Invalid request or stock exhausted
- `404 Not Found` - Coupon not found
- `409 Conflict` - User has already claimed this coupon

### 5. Monitoring

```bash
# View application logs in real-time
docker logs -f coupon_api

# View database logs
docker logs -f coupon_postgres

# Check database connection
docker exec -it coupon_postgres psql -U postgres -d coupon_db -c "SELECT * FROM coupons;"

# Restart application
docker-compose restart api

# Restart all services
docker-compose restart
```

### 6. Troubleshooting

**API is not accessible:**
```bash
# Check if containers are running
docker ps

# Restart service
docker-compose restart api

# Check error logs
docker logs coupon_api
```

**Database connection error:**
```bash
# Check database container
docker logs coupon_postgres

# Test connection from API container
docker exec -it coupon_api ping postgres
```

### 7. Using with Other Tools

#### Postman / Insomnia
- Base URL: `http://localhost:8080`
- Content-Type: `application/json`

#### Python Example
```python
import requests

base_url = "http://localhost:8080"

# Create coupon
response = requests.post(
    f"{base_url}/api/coupons",
    json={"name": "PROMO_SUPER", "amount": 100}
)
print(response.status_code)

# Claim coupon
response = requests.post(
    f"{base_url}/api/coupons/claim",
    json={"user_id": "user_12345", "coupon_name": "PROMO_SUPER"}
)
print(response.status_code, response.json())

# Get coupon details
response = requests.get(f"{base_url}/api/coupons/PROMO_SUPER")
print(response.json())
```

#### JavaScript/Node.js Example
```javascript
const axios = require('axios');

const baseURL = 'http://localhost:8080';

// Create coupon
async function createCoupon() {
  const response = await axios.post(`${baseURL}/api/coupons`, {
    name: 'PROMO_SUPER',
    amount: 100
  });
  console.log('Status:', response.status);
}

// Claim coupon
async function claimCoupon() {
  try {
    const response = await axios.post(`${baseURL}/api/coupons/claim`, {
      user_id: 'user_12345',
      coupon_name: 'PROMO_SUPER'
    });
    console.log('Success:', response.status);
  } catch (error) {
    console.log('Error:', error.response.status, error.response.data);
  }
}

// Get coupon details
async function getCouponDetails() {
  const response = await axios.get(`${baseURL}/api/coupons/PROMO_SUPER`);
  console.log(response.data);
}
```

## How to Test

### Manual Testing with cURL

#### 1. Create a Coupon
```bash
curl -X POST http://localhost:8080/api/coupons \
  -H "Content-Type: application/json" \
  -d '{"name": "PROMO_SUPER", "amount": 5}'
```

Expected response: `201 Created`

#### 2. Claim a Coupon
```bash
curl -X POST http://localhost:8080/api/coupons/claim \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user_12345", "coupon_name": "PROMO_SUPER"}'
```

Expected responses:
- Success: `200 OK`
- Already claimed: `409 Conflict`
- No stock: `400 Bad Request`

#### 3. Get Coupon Details
```bash
curl http://localhost:8080/api/coupons/PROMO_SUPER
```

Expected response:
```json
{
  "name": "PROMO_SUPER",
  "amount": 5,
  "remaining_amount": 4,
  "claimed_by": ["user_12345"]
}
```

### Stress Testing Scenarios

#### Scenario 1: Flash Sale Attack (50 concurrent requests, 5 stock)
```bash
# Create coupon with 5 stock
curl -X POST http://localhost:8080/api/coupons \
  -H "Content-Type: application/json" \
  -d '{"name": "FLASH_SALE", "amount": 5}'

# Run 50 concurrent requests (using different user_ids)
for i in {1..50}; do
  curl -X POST http://localhost:8080/api/coupons/claim \
    -H "Content-Type: application/json" \
    -d "{\"user_id\": \"user_$i\", \"coupon_name\": \"FLASH_SALE\"}" &
done
wait

# Verify only 5 claims succeeded
curl http://localhost:8080/api/coupons/FLASH_SALE
```

Expected result: `remaining_amount: 0`, exactly 5 users in `claimed_by` array.

#### Scenario 2: Double Dip Attack (10 concurrent requests from same user)
```bash
# Create coupon
curl -X POST http://localhost:8080/api/coupons \
  -H "Content-Type: application/json" \
  -d '{"name": "DOUBLE_DIP", "amount": 10}'

# Run 10 concurrent requests from the same user
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/coupons/claim \
    -H "Content-Type: application/json" \
    -d '{"user_id": "user_12345", "coupon_name": "DOUBLE_DIP"}' &
done
wait

# Verify only 1 claim succeeded
curl http://localhost:8080/api/coupons/DOUBLE_DIP
```

Expected result: Exactly 1 success, 9 failures (409 Conflict), `user_12345` appears only once in `claimed_by`.

### Using Automated Testing Scripts

You can create a simple bash script for stress testing:

```bash
#!/bin/bash
# stress_test.sh

# Create coupon
curl -X POST http://localhost:8080/api/coupons \
  -H "Content-Type: application/json" \
  -d '{"name": "TEST_COUPON", "amount": 5}'

# Concurrent claims
for i in {1..50}; do
  curl -s -X POST http://localhost:8080/api/coupons/claim \
    -H "Content-Type: application/json" \
    -d "{\"user_id\": \"user_$i\", \"coupon_name\": \"TEST_COUPON\"}" \
    -w "\n%{http_code}\n" &
done
wait

# Check results
curl http://localhost:8080/api/coupons/TEST_COUPON
```

## API Endpoints

### POST /api/coupons
Creates a new coupon.

**Request Body:**
```json
{
  "name": "PROMO_SUPER",
  "amount": 100
}
```

**Response:** `201 Created`

---

### POST /api/coupons/claim
Claims a coupon for a user.

**Request Body:**
```json
{
  "user_id": "user_12345",
  "coupon_name": "PROMO_SUPER"
}
```

**Response Codes:**
- `200 OK` - Successfully claimed
- `409 Conflict` - User has already claimed this coupon
- `400 Bad Request` - Stock exhausted or invalid request
- `404 Not Found` - Coupon not found

---

### GET /api/coupons/{name}
Gets coupon details including remaining stock and list of users who claimed it.

**Response Body:**
```json
{
  "name": "PROMO_SUPER",
  "amount": 100,
  "remaining_amount": 0,
  "claimed_by": ["user_12345", "user_67890"]
}
```

---

### GET /health
Health check endpoint.

**Response:** `200 OK` with `{"status": "ok"}`

## Architecture Notes

### Database Design

The system uses **PostgreSQL** with two separate tables to maintain separation of concerns:

1. **`coupons` table**: Stores coupon metadata
   - `id` (Primary Key)
   - `name` (Unique Index)
   - `amount` (Total stock available)
   - Timestamps

2. **`claims` table**: Stores claim history
   - `id` (Primary Key)
   - `user_id` (Part of unique constraint)
   - `coupon_name` (Part of unique constraint, indexed)
   - Timestamps

**Critical Constraint**: A composite unique index on `(user_id, coupon_name)` ensures that:
- A user can claim a specific coupon only once
- The same user can claim different coupons
- Race conditions are prevented at the database level

### Concurrency & Locking Strategy

The system handles high concurrency through:

1. **Database Transactions**: All claim operations are wrapped in a PostgreSQL transaction, ensuring atomicity of:
   - Checking if user already claimed
   - Checking stock availability
   - Creating the claim record
   - Preventing race conditions

2. **Unique Constraint Enforcement**: The database-level unique constraint on `(user_id, coupon_name)` acts as a final safeguard:
   - Even if two requests pass the initial checks simultaneously
   - Only one will succeed in inserting the claim record
   - The other will receive a unique constraint violation error (409 Conflict)

3. **Transaction Isolation**: PostgreSQL's default isolation level (READ COMMITTED) ensures:
   - Each transaction sees committed data
   - Concurrent transactions are serialized for conflicting operations
   - Stock counting is accurate under high load

### Flow Diagram

```
Claim Request
    ↓
Start Transaction
    ↓
Check Coupon Exists
    ↓
Check User Already Claimed (SELECT with FOR UPDATE)
    ↓
Count Existing Claims
    ↓
Check Stock Availability
    ↓
Insert Claim Record (Unique Constraint Protection)
    ↓
Commit Transaction
    ↓
Return Success/Error
```

### Why This Design Works

1. **Separation of Concerns**: Coupon data and claim history are separate, allowing for:
   - Independent scaling
   - Easier querying and reporting
   - Clear data relationships

2. **Atomic Operations**: Database transactions ensure that all checks and inserts happen atomically, preventing:
   - Double claims
   - Stock overselling
   - Race conditions

3. **Database-Level Constraints**: The unique constraint provides a safety net that works even under extreme concurrency, as the database engine handles the serialization.

## Project Structure

```
flash-sale-coupon-system/
├── main.go                 # Application entry point
├── go.mod                  # Go module dependencies
├── Dockerfile              # Docker image definition
├── docker-compose.yml      # Docker Compose configuration
├── README.md               # This file
└── internal/
    ├── config/             # Configuration management
    │   └── config.go
    ├── database/           # Database connection and migrations
    │   └── database.go
    ├── handlers/           # HTTP request handlers
    │   └── coupon.go
    └── models/             # Data models
        └── coupon.go
```

## Environment Variables

The application supports configuration via `.env` file or environment variables. Create a `.env` file in the root directory with the following variables:

### Database Configuration (Individual Variables - Recommended)
- `DB_DEFAULT`: Database type (default: `postgresql`)
- `DB_PG_HOST`: PostgreSQL host (default: `localhost`)
- `DB_PG_DATABASE`: Database name (default: `coupon_db`)
- `DB_PG_USERNAME`: Database username (default: `postgres`)
- `DB_PG_PASSWORD`: Database password (default: empty)
- `DB_PG_PORT`: Database port (default: `5432`)

### Alternative: Direct Connection String
- `DATABASE_URL`: PostgreSQL connection string (takes precedence if set)

### Server Configuration
- `PORT`: Server port (default: `8080`)

**Priority order:**
1. `DATABASE_URL` (if set, takes precedence)
2. Individual database variables (`DB_PG_*`)
3. Default values

**Example `.env` file:**
```env
DB_DEFAULT=postgresql
DB_PG_HOST=localhost
DB_PG_DATABASE=coupon_db
DB_PG_USERNAME=postgres
DB_PG_PASSWORD=postgres
DB_PG_PORT=5432
PORT=8080
```

**Alternative using DATABASE_URL:**
```env
DATABASE_URL=host=localhost user=postgres password=postgres dbname=coupon_db port=5432 sslmode=disable
PORT=8080
```

## Troubleshooting

### Database Connection Issues
- Ensure PostgreSQL container is running: `docker ps`
- Check database logs: `docker logs coupon_postgres`
- Verify database configuration in `docker-compose.yml` or `.env` file
- Test database connection: `docker exec -it coupon_postgres psql -U postgres -d coupon_db`

### Port Already in Use
- Change the port mapping in `docker-compose.yml`
- Or stop the service using port 8080: `lsof -ti:8080 | xargs kill -9`

### Migration Issues
- Drop and recreate the database volume: `docker-compose down -v && docker-compose up`
- Check migration logs: `docker logs coupon_api | grep -i migration`

### API Not Responding
- Check if API container is running: `docker ps | grep coupon_api`
- View API logs: `docker logs coupon_api`
- Restart API service: `docker-compose restart api`
- Verify health endpoint: `curl http://localhost:8080/health`

## License

This project is created for technical assessment purposes.

