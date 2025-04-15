#!/bin/bash

while true; do
  curl -s http://localhost:8085/ > /dev/null
  
  # Generate random number between 200 and 500 (milliseconds)
  sleep_ms=$((RANDOM % 301 + 200))
  
  # Convert to seconds for sleep command
  sleep_sec=$(echo "scale=3; $sleep_ms/1000" | bc)
  
  echo "Hit! Sleeping for $sleep_ms milliseconds..."
  sleep $sleep_sec
done
