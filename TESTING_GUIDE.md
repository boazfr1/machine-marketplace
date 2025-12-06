# Machine Marketplace - Complete Testing Guide

## Overview
This guide will help you test the complete Linux machine marketplace flow from login to purchase with Kafka integration.

## Prerequisites

1. **Install Kafka Go dependency:**
   ```bash
   cd server
   go get github.com/segmentio/kafka-go
   go mod tidy
   ```

2. **Generate database code (if needed):**
   ```bash
   cd server
   sqlc generate
   ```

## Starting the Application

### Option 1: Docker Compose (Recommended)
```bash
# From project root
docker-compose up --build

# Wait for all services to start:
# - Client: http://localhost:80
# - Auth Service: http://localhost:3001
# - Order Service: http://localhost:3002
# - PostgreSQL: localhost:5432
# - Kafka: localhost:9092
# - Zookeeper: localhost:2181
```

### Option 2: Local Development
```bash
# Terminal 1 - Start database and Kafka
docker-compose up postgres kafka zookeeper

# Terminal 2 - Start auth service
cd server
make run-auth

# Terminal 3 - Start order service
cd server
make run-order

# Terminal 4 - Start frontend
cd client
npm run dev
```

## Complete User Flow Test

### 1. **Test Authentication**

**Signup:**
1. Navigate to http://localhost (or http://localhost:5173 for local dev)
2. Click "Sign Up"
3. Fill in:
   - Name: testuser
   - Email: testuser@example.com
   - Password: password123
4. Click "Sign Up"
5. ✅ Should redirect to login page

**Login:**
1. Go to Login page
2. Use credentials:
   - Email: alice@test.com (or your created user)
   - Password: password1
3. Click "Log In"
4. ✅ Should redirect to marketplace

### 2. **Test Machine Browsing & Filtering**

**View All Machines:**
1. After login, you're on the marketplace feed
2. ✅ Should see 10 available Linux machines displayed in cards
3. ✅ Each card shows: Name, CPU, RAM, GPU, Storage, Owner ID

**Test Filters:**
1. Use the slider filters:
   - Min CPU: Set to 4 cores
   - ✅ Should see only machines with 4+ CPU cores

2. Filter by GPU:
   - Min GPU: Set to 1
   - ✅ Should see only machines with GPU

3. Filter by RAM:
   - Min RAM: Set to 16 GB
   - ✅ Should see only machines with 16+ GB RAM

4. Test Category Filter:
   - Select "With GPU" from dropdown
   - ✅ Should see only GPU machines

5. Clear all filters:
   - Click "Clear filters" button
   - ✅ Should see all machines again

**Search:**
1. Use search bar: Type "Linux-AI"
2. ✅ Should filter to matching machine

### 3. **Test Machine Details & Purchase**

**View Details:**
1. Click on any available machine card
2. ✅ Should navigate to /machine/{id}
3. ✅ Should see:
   - Full specifications including GPU
   - Pricing options
   - Owner information
   - Purchase button

**Purchase a Machine:**
1. On machine details page
2. Select rental duration (click one of the pricing tiers):
   - Daily (24h)
   - 3 Days (72h) ← Select this
   - Weekly (168h)
   - Monthly (720h)
3. Click "Purchase for 72h" button
4. ✅ Button should change to "Processing Purchase..."
5. ✅ After success: "✓ Purchase Successful!"
6. ✅ Should auto-redirect to /purchased-machines after 2 seconds

### 4. **Verify Kafka Integration**

**Check Kafka Received the Purchase Event:**

```bash
# Method 1: Using Docker
docker exec -it machine-marketplace-kafka-1 /bin/bash
kafka-console-consumer --bootstrap-server localhost:9092 \
  --topic machine-purchases --from-beginning

# Method 2: Using kafka-go-cli (if installed)
kafka-console-consumer --bootstrap-server localhost:9092 \
  --topic machine-purchases --from-beginning
```

**Expected Output:**
```json
{
  "machine_id": 3,
  "buyer_id": 1,
  "deal_expiration": "2025-12-09T21:00:00Z",
  "purchase_time": "2025-12-06T21:00:00Z",
  "machine_name": "Linux-AI-3"
}
```

### 5. **Test Database Updates**

**Verify Purchase in Database:**
```bash
# Connect to PostgreSQL
docker exec -it machine-marketplace-postgres-1 psql -U postgres -d machine_market

# Check if machine is now purchased
SELECT id, name, buyer_id, owner_id FROM machines WHERE id = 3;

# Expected: buyer_id should be set (not NULL)
```

### 6. **Test Purchased Machines Page**

1. Navigate to "My Purchases" or /purchased-machines
2. ✅ Should see the machine you just purchased
3. ✅ Should show: Name, CPU, RAM, GPU, Storage, Owner

## API Testing (Optional)

### Test Filtering API
```bash
# Get all machines
curl http://localhost:3002/api/v1/order \
  -H "Cookie: jwt=YOUR_JWT_TOKEN" \
  | jq

# Filter by CPU
curl "http://localhost:3002/api/v1/order?cpu=4" \
  -H "Cookie: jwt=YOUR_JWT_TOKEN" \
  | jq

# Filter by GPU
curl "http://localhost:3002/api/v1/order?gpu=1" \
  -H "Cookie: jwt=YOUR_JWT_TOKEN" \
  | jq

# Filter by multiple criteria
curl "http://localhost:3002/api/v1/order?cpu=8&ram=16&gpu=2" \
  -H "Cookie: jwt=YOUR_JWT_TOKEN" \
  | jq
```

### Test Purchase API
```bash
# First, login to get JWT cookie
curl -X POST http://localhost:3001/api/v1/user/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@test.com","password":"password1"}' \
  -c cookies.txt

# Purchase a machine
curl -X POST http://localhost:3002/api/v1/order/buy \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "machine_id": 2,
    "deal_duration_hours": 48
  }' | jq
```

## Expected Results Summary

✅ **Authentication**: Users can signup, login, and logout
✅ **Marketplace**: Displays all available machines with specs including GPU
✅ **Filtering**: Can filter by CPU, RAM, GPU (min values)
✅ **Search**: Can search machines by name
✅ **Machine Details**: Shows full specs and purchase option
✅ **Purchase Flow**: Can select duration and purchase
✅ **Kafka Integration**: Purchase events sent to Kafka topic
✅ **Database Updates**: Machines marked as purchased (buyer_id set)
✅ **UI Feedback**: Loading states, success messages, error handling

## Troubleshooting

**Frontend not connecting:**
- Check CORS settings in server
- Verify services are running on correct ports
- Check browser console for errors

**Kafka not receiving messages:**
- Check if Kafka and Zookeeper are healthy: `docker ps`
- View Kafka logs: `docker logs machine-marketplace-kafka-1`
- Ensure KAFKA_BROKER env variable is set correctly

**Database errors:**
- Restart PostgreSQL: `docker-compose restart postgres`
- Check schema is created: Connect to DB and run `\dt`
- Verify seed data exists: `SELECT * FROM machines;`

**Purchase fails:**
- Check machine is available (buyer_id IS NULL)
- Verify JWT token is valid
- Check server logs for errors

## Clean Slate Testing

To start fresh:
```bash
# Stop all services
docker-compose down

# Remove volumes (clears database)
docker-compose down -v

# Rebuild and restart
docker-compose up --build
```

This will recreate the database with fresh seed data (10 available machines).
