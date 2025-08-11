#!/bin/bash

# Trace ID Monitor Script
# Usage: ./trace-monitor.sh [trace_id] [duration]

TRACE_ID=${1:-""}
DURATION=${2:-"60"}  # Default 60 seconds

echo "🔍 Haslaw Services - Trace ID Monitor"
echo "====================================="

if [ -z "$TRACE_ID" ]; then
    echo "📡 Monitoring ALL trace IDs for ${DURATION} seconds..."
    echo "Press Ctrl+C to stop"
    echo ""
    
    # Monitor all trace IDs
    timeout ${DURATION}s docker-compose logs -f app 2>/dev/null | grep --line-buffered "TRACE:" | while read line; do
        timestamp=$(echo "$line" | cut -d'|' -f1)
        trace_part=$(echo "$line" | grep -o "TRACE: [^]]*")
        request_info=$(echo "$line" | cut -d']' -f2-)
        
        echo "🔹 $timestamp | $trace_part ]$request_info"
    done
else
    echo "🎯 Tracking specific Trace ID: $TRACE_ID"
    echo "Monitoring for ${DURATION} seconds..."
    echo "Press Ctrl+C to stop"
    echo ""
    
    # Monitor specific trace ID
    timeout ${DURATION}s docker-compose logs -f app 2>/dev/null | grep --line-buffered "TRACE: $TRACE_ID" | while read line; do
        timestamp=$(echo "$line" | cut -d'|' -f1)
        request_info=$(echo "$line" | cut -d']' -f2-)
        
        echo "✅ $timestamp | TRACE: $TRACE_ID ]$request_info"
    done
fi

echo ""
echo "🏁 Monitoring completed!"
