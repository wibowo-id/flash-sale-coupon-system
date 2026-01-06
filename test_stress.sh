#!/bin/bash

# Stress Test Script for Flash Sale Coupon System
# This script tests the two evaluation scenarios

BASE_URL="http://localhost:8080"

echo "=== Flash Sale Coupon System - Stress Test ==="
echo ""

# Test 1: Flash Sale Attack (50 concurrent requests, 5 stock)
echo "Test 1: Flash Sale Attack (50 concurrent requests, 5 stock)"
echo "Creating coupon with 5 stock..."
curl -s -X POST "$BASE_URL/api/coupons" \
  -H "Content-Type: application/json" \
  -d '{"name": "FLASH_SALE", "amount": 5}' > /dev/null

echo "Sending 50 concurrent claim requests..."
success_count=0
fail_count=0

for i in {1..50}; do
  response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/coupons/claim" \
    -H "Content-Type: application/json" \
    -d "{\"user_id\": \"user_$i\", \"coupon_name\": \"FLASH_SALE\"}")
  
  http_code=$(echo "$response" | tail -n1)
  if [ "$http_code" = "200" ] || [ "$http_code" = "201" ]; then
    ((success_count++))
  else
    ((fail_count++))
  fi
done

echo "Results: $success_count successful, $fail_count failed"
echo "Checking coupon details..."
curl -s "$BASE_URL/api/coupons/FLASH_SALE" | jq '.'
echo ""

# Test 2: Double Dip Attack (10 concurrent requests from same user)
echo "Test 2: Double Dip Attack (10 concurrent requests from same user)"
echo "Creating coupon with 10 stock..."
curl -s -X POST "$BASE_URL/api/coupons" \
  -H "Content-Type: application/json" \
  -d '{"name": "DOUBLE_DIP", "amount": 10}' > /dev/null

echo "Sending 10 concurrent claim requests from same user..."
success_count=0
fail_count=0

for i in {1..10}; do
  response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/coupons/claim" \
    -H "Content-Type: application/json" \
    -d '{"user_id": "user_12345", "coupon_name": "DOUBLE_DIP"}')
  
  http_code=$(echo "$response" | tail -n1)
  if [ "$http_code" = "200" ] || [ "$http_code" = "201" ]; then
    ((success_count++))
  else
    ((fail_count++))
  fi
done

echo "Results: $success_count successful, $fail_count failed"
echo "Checking coupon details..."
curl -s "$BASE_URL/api/coupons/DOUBLE_DIP" | jq '.'
echo ""

echo "=== Stress Test Complete ==="

