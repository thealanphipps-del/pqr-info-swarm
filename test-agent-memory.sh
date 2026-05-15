#!/bin/bash
# PQR-ticketing test script for agent memory system

set -e

BASE_URL="${1:-http://localhost:8080}"
AGENT_ID="${2:-test-agent-001}"

echo "PQR Ticketing System - Agent Memory Test"
echo "=========================================="
echo "Base URL: $BASE_URL"
echo "Agent ID: $AGENT_ID"
echo

# Health check
echo "1. Health Check..."
curl -s -X GET "$BASE_URL/REST/2.0/health" | jq .
echo

# Initialize schema
echo "2. Initialize Schema..."
curl -s -X POST "$BASE_URL/REST/2.0/init" | jq .
echo

# Create a ticket
echo "3. Creating Ticket..."
TICKET_ID=$(curl -s -X POST "$BASE_URL/REST/2.0/ticket" \
  -H "Content-Type: application/json" \
  -d '{
    "Subject": "Agent Working Memory",
    "Queue": "processing",
    "Text": "Initial task content",
    "AgentID": "'$AGENT_ID'",
    "Layer": 2,
    "Intent": {
      "task": "test",
      "priority": "high"
    }
  }' | jq -r '.id')
echo "Ticket ID: $TICKET_ID"
echo

# Store memory
echo "4. Storing Agent Memory..."
curl -s -X POST "$BASE_URL/REST/2.0/agent/$AGENT_ID/memory/$TICKET_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "memory_type": "context",
    "data": {
      "status": "processing",
      "items_processed": 5,
      "items_total": 10,
      "current_item": "data_point_5"
    },
    "relevance_score": 0.95
  }' | jq .
echo

# Retrieve memory
echo "5. Retrieving Agent Memory..."
curl -s -X GET "$BASE_URL/REST/2.0/agent/$AGENT_ID/memory/$TICKET_ID?type=context" | jq .
echo

# Get ticket details
echo "6. Getting Ticket Details..."
curl -s -X GET "$BASE_URL/REST/2.0/ticket/$TICKET_ID" | jq .
echo

# Store additional memory types
echo "7. Storing Knowledge Memory..."
curl -s -X POST "$BASE_URL/REST/2.0/agent/$AGENT_ID/memory/$TICKET_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "memory_type": "knowledge",
    "data": {
      "patterns": ["pattern_a", "pattern_b"],
      "confidence": 0.87
    },
    "relevance_score": 0.85
  }' | jq .
echo

# Update ticket status
echo "8. Updating Ticket Status..."
curl -s -X PUT "$BASE_URL/REST/2.0/ticket/$TICKET_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "Status": "PROCESSING",
    "Title": "Updated: Memory Storage Test"
  }' | jq .
echo

# Get audit trail
echo "9. Getting Audit Trail..."
curl -s -X GET "$BASE_URL/REST/2.0/ticket/$TICKET_ID/audit" | jq .
echo

# Get agent context
echo "10. Getting Agent Context..."
curl -s -X GET "$BASE_URL/REST/2.0/agent/$AGENT_ID/context" | jq .
echo

echo "=========================================="
echo "Test Complete!"
echo "Ticket ID: $TICKET_ID"
echo "Use this to test further operations:"
echo "  - GET /REST/2.0/ticket/$TICKET_ID"
echo "  - GET /REST/2.0/agent/$AGENT_ID/memory/$TICKET_ID"
echo "  - PUT /REST/2.0/ticket/$TICKET_ID"


