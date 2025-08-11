#!/bin/bash

# Trace Extractor Script
# Usage: ./extract-trace.sh [trace_id] [output_file]

TRACE_ID=${1}
OUTPUT_FILE=${2:-"trace_${TRACE_ID}.log"}

if [ -z "$TRACE_ID" ]; then
    echo "❌ Usage: $0 <trace_id> [output_file]"
    echo "Example: $0 abc123def456 trace_abc123def456.log"
    exit 1
fi

echo "🔍 Extracting trace ID: $TRACE_ID"
echo "📄 Output file: $OUTPUT_FILE"
echo ""

# Extract all logs for specific trace ID
docker-compose logs app 2>/dev/null | grep "TRACE: $TRACE_ID" > "$OUTPUT_FILE"

if [ -s "$OUTPUT_FILE" ]; then
    LINE_COUNT=$(wc -l < "$OUTPUT_FILE")
    echo "✅ Found $LINE_COUNT log entries for trace ID: $TRACE_ID"
    echo "📋 Trace timeline:"
    echo "=================="
    
    cat "$OUTPUT_FILE" | while read line; do
        timestamp=$(echo "$line" | cut -d'|' -f1)
        action=$(echo "$line" | grep -o -E "(Started|Completed)" | head -1)
        method=$(echo "$line" | awk '{print $8}')
        endpoint=$(echo "$line" | awk '{print $9}')
        
        if echo "$line" | grep -q "Started"; then
            echo "🚀 $timestamp - $method $endpoint (Started)"
        elif echo "$line" | grep -q "Completed"; then
            duration=$(echo "$line" | grep -o '[0-9.]*[a-z]*s with' | sed 's/ with//')
            status=$(echo "$line" | grep -o 'status [0-9]*' | cut -d' ' -f2)
            case $status in
                2*) echo "✅ $timestamp - $method $endpoint (Completed in $duration - $status)" ;;
                3*) echo "🔄 $timestamp - $method $endpoint (Completed in $duration - $status)" ;;
                4*) echo "⚠️  $timestamp - $method $endpoint (Completed in $duration - $status)" ;;
                5*) echo "❌ $timestamp - $method $endpoint (Completed in $duration - $status)" ;;
                *) echo "❓ $timestamp - $method $endpoint (Completed in $duration - $status)" ;;
            esac
        else
            echo "📝 $timestamp - $(echo "$line" | cut -d']' -f2-)"
        fi
    done
    
    echo ""
    echo "💾 Full trace saved to: $OUTPUT_FILE"
    echo "🔍 View with: cat $OUTPUT_FILE"
else
    echo "❌ No logs found for trace ID: $TRACE_ID"
    rm -f "$OUTPUT_FILE"
    exit 1
fi
