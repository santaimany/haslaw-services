#!/bin/bash

# Trace Analysis Script
# Usage: ./trace-analysis.sh [hours_back]

HOURS_BACK=${1:-"1"}  # Default 1 hour
CONTAINER_NAME="haslaw-app"

echo "📊 Haslaw Services - Trace Analysis"
echo "===================================="
echo "Analyzing logs from last ${HOURS_BACK} hour(s)..."
echo ""

# Get logs from last N hours
SINCE_TIME=$(date -d "${HOURS_BACK} hours ago" '+%Y-%m-%dT%H:%M:%S')
LOGS=$(docker logs $CONTAINER_NAME --since="$SINCE_TIME" 2>/dev/null | grep "TRACE:")

if [ -z "$LOGS" ]; then
    echo "❌ No trace logs found in the last ${HOURS_BACK} hour(s)"
    exit 1
fi

echo "📋 TRACE SUMMARY:"
echo "=================="

# Count total requests
TOTAL_REQUESTS=$(echo "$LOGS" | grep "Started" | wc -l)
echo "🔢 Total Requests: $TOTAL_REQUESTS"

# Count completed requests
COMPLETED_REQUESTS=$(echo "$LOGS" | grep "Completed" | wc -l)
echo "✅ Completed Requests: $COMPLETED_REQUESTS"

# Count requests by method
echo ""
echo "📊 Requests by HTTP Method:"
echo "$LOGS" | grep "Started" | awk '{print $8}' | sort | uniq -c | sort -nr | while read count method; do
    echo "   $method: $count requests"
done

# Count requests by endpoint
echo ""
echo "🎯 Top Endpoints:"
echo "$LOGS" | grep "Started" | awk '{print $9}' | sort | uniq -c | sort -nr | head -10 | while read count endpoint; do
    echo "   $endpoint: $count requests"
done

# Average response time for completed requests
echo ""
echo "⏱️  Response Time Analysis:"
AVG_TIME=$(echo "$LOGS" | grep "Completed" | grep -o '[0-9.]*[a-z]*s with' | sed 's/ with//' | \
    awk '{
        if ($1 ~ /ms/) { time = substr($1, 1, length($1)-2) / 1000 }
        else if ($1 ~ /µs/) { time = substr($1, 1, length($1)-2) / 1000000 }
        else if ($1 ~ /s/) { time = substr($1, 1, length($1)-1) }
        else { time = $1 }
        total += time; count++
    } END { 
        if (count > 0) printf "%.3f", total/count 
    }')

if [ ! -z "$AVG_TIME" ] && [ "$AVG_TIME" != "0.000" ]; then
    echo "   Average Response Time: ${AVG_TIME}s"
else
    echo "   Average Response Time: Unable to calculate"
fi

# Status code distribution
echo ""
echo "📈 Status Code Distribution:"
echo "$LOGS" | grep "Completed" | grep -o 'status [0-9]*' | awk '{print $2}' | sort | uniq -c | sort -nr | while read count status; do
    case $status in
        2*) echo "   ✅ $status: $count requests" ;;
        3*) echo "   🔄 $status: $count requests" ;;
        4*) echo "   ⚠️  $status: $count requests" ;;
        5*) echo "   ❌ $status: $count requests" ;;
        *) echo "   ❓ $status: $count requests" ;;
    esac
done

# Recent error traces
echo ""
echo "🚨 Recent Error Traces (4xx/5xx):"
echo "$LOGS" | grep "Completed" | grep -E "status [45][0-9][0-9]" | tail -5 | while read line; do
    trace_id=$(echo "$line" | grep -o "TRACE: [^]]*" | cut -d' ' -f2)
    status=$(echo "$line" | grep -o "status [0-9]*" | cut -d' ' -f2)
    method=$(echo "$line" | awk '{print $8}')
    endpoint=$(echo "$line" | awk '{print $9}')
    echo "   🔴 $trace_id - $method $endpoint ($status)"
done

echo ""
echo "🔍 Use './trace-monitor.sh [trace_id]' to monitor specific traces in real-time"
echo "📋 Use 'docker-compose logs app | grep \"TRACE: [trace_id]\"' to see full trace history"
