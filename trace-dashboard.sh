#!/bin/bash

# Trace Monitoring Dashboard
# Usage: ./trace-dashboard.sh

clear
echo "🚀 Haslaw API Trace Monitoring Dashboard"
echo "========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to get current stats
get_current_stats() {
    local last_hour=$(date -d '1 hour ago' '+%Y-%m-%d %H:%M:%S')
    
    echo -e "${BLUE}📊 Current System Stats (Last Hour)${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    # Total requests
    local total_requests=$(docker logs haslaw-app 2>&1 | grep "TRACE:" | grep -c "Request:")
    echo -e "Total Requests: ${GREEN}$total_requests${NC}"
    
    # Active traces (last 10 minutes)
    local active_traces=$(docker logs haslaw-app 2>&1 | grep "TRACE:" | tail -1000 | grep -c "Request:")
    echo -e "Recent Activity: ${GREEN}$active_traces${NC} requests (last 1000 logs)"
    
    # Error rate
    local errors=$(docker logs haslaw-app 2>&1 | grep "TRACE:" | tail -1000 | grep -c "ERROR")
    local error_rate=0
    if [ $active_traces -gt 0 ]; then
        error_rate=$((errors * 100 / active_traces))
    fi
    
    if [ $error_rate -gt 10 ]; then
        echo -e "Error Rate: ${RED}$error_rate%${NC} ($errors errors)"
    elif [ $error_rate -gt 5 ]; then
        echo -e "Error Rate: ${YELLOW}$error_rate%${NC} ($errors errors)"
    else
        echo -e "Error Rate: ${GREEN}$error_rate%${NC} ($errors errors)"
    fi
    
    # Response time analysis
    echo -e "\n${BLUE}⏱️  Response Time Analysis${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    docker logs haslaw-app 2>&1 | grep "Response:" | tail -100 | \
    sed -n 's/.*Duration: \([0-9.]*[a-z]*\).*/\1/p' | \
    awk '
    BEGIN { total=0; count=0; fast=0; slow=0 }
    {
        # Convert to milliseconds
        if ($1 ~ /ms$/) {
            time = substr($1, 1, length($1)-2)
        } else if ($1 ~ /s$/) {
            time = substr($1, 1, length($1)-1) * 1000
        } else if ($1 ~ /µs$/) {
            time = substr($1, 1, length($1)-2) / 1000
        } else {
            time = $1
        }
        
        total += time
        count++
        
        if (time < 100) fast++
        else if (time > 1000) slow++
    }
    END {
        if (count > 0) {
            avg = total/count
            printf "Average Response Time: %.2fms\n", avg
            printf "Fast Responses (<100ms): %d (%.1f%%)\n", fast, (fast/count)*100
            printf "Slow Responses (>1s): %d (%.1f%%)\n", slow, (slow/count)*100
        }
    }'
}

# Function to show top endpoints
show_top_endpoints() {
    echo -e "\n${BLUE}🔥 Top Endpoints (Last 500 requests)${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    docker logs haslaw-app 2>&1 | grep "Request:" | tail -500 | \
    sed -n 's/.*Method: \([A-Z]*\), Path: \([^,]*\).*/\1 \2/p' | \
    sort | uniq -c | sort -rn | head -10 | \
    awk '{printf "%-4s %-8s %s\n", $1"x", $2, $3}'
}

# Function to show recent errors
show_recent_errors() {
    echo -e "\n${BLUE}❌ Recent Errors (Last 20)${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    docker logs haslaw-app 2>&1 | grep -E "(ERROR|WARN)" | grep "TRACE:" | tail -20 | \
    while IFS= read -r line; do
        if echo "$line" | grep -q "ERROR"; then
            echo -e "${RED}$line${NC}"
        else
            echo -e "${YELLOW}$line${NC}"
        fi
    done
}

# Function to show status codes distribution
show_status_codes() {
    echo -e "\n${BLUE}📈 HTTP Status Codes (Last 200 responses)${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    docker logs haslaw-app 2>&1 | grep "Response:" | tail -200 | \
    sed -n 's/.*Status: \([0-9]*\).*/\1/p' | \
    sort | uniq -c | sort -rn | \
    awk '{
        if ($2 >= 200 && $2 < 300) color="\033[0;32m"      # Green for 2xx
        else if ($2 >= 300 && $2 < 400) color="\033[0;33m" # Yellow for 3xx  
        else if ($2 >= 400 && $2 < 500) color="\033[0;31m" # Red for 4xx
        else if ($2 >= 500) color="\033[1;31m"             # Bright red for 5xx
        else color="\033[0m"                               # Default
        
        printf "%s%-4s %s\033[0m\n", color, $1"x", $2
    }'
}

# Function to show live trace
show_live_trace() {
    echo -e "\n${BLUE}🔴 Live Trace Monitor${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "Press Ctrl+C to return to dashboard..."
    echo ""
    docker logs -f haslaw-app 2>&1 | grep --line-buffered "TRACE:" | \
    while IFS= read -r line; do
        if echo "$line" | grep -q "Request:"; then
            echo -e "${GREEN}➤ $line${NC}"
        elif echo "$line" | grep -q "Response:"; then
            echo -e "${BLUE}➤ $line${NC}"
        elif echo "$line" | grep -q "ERROR"; then
            echo -e "${RED}➤ $line${NC}"
        else
            echo -e "➤ $line"
        fi
    done
}

# Function to search trace by ID
search_trace() {
    echo -e "\n${BLUE}🔍 Search Trace by ID${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    read -p "Enter Trace ID: " trace_id
    
    if [ -n "$trace_id" ]; then
        echo -e "\n${GREEN}Found traces for ID: $trace_id${NC}"
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        docker logs haslaw-app 2>&1 | grep "TRACE: $trace_id" | \
        while IFS= read -r line; do
            if echo "$line" | grep -q "Request:"; then
                echo -e "${GREEN}$line${NC}"
            elif echo "$line" | grep -q "Response:"; then
                echo -e "${BLUE}$line${NC}"
            elif echo "$line" | grep -q "ERROR"; then
                echo -e "${RED}$line${NC}"
            else
                echo "$line"
            fi
        done
    fi
    
    echo ""
    read -p "Press Enter to continue..."
}

# Main menu loop
while true; do
    clear
    echo "🚀 Haslaw API Trace Monitoring Dashboard"
    echo "========================================="
    echo ""
    echo "1) 📊 Current Stats"
    echo "2) 🔥 Top Endpoints" 
    echo "3) ❌ Recent Errors"
    echo "4) 📈 Status Codes"
    echo "5) 🔴 Live Monitor"
    echo "6) 🔍 Search Trace"
    echo "7) 🔄 Refresh"
    echo "8) ❌ Exit"
    echo ""
    read -p "Select option (1-8): " choice
    
    case $choice in
        1) get_current_stats; echo ""; read -p "Press Enter to continue..." ;;
        2) show_top_endpoints; echo ""; read -p "Press Enter to continue..." ;;
        3) show_recent_errors; echo ""; read -p "Press Enter to continue..." ;;
        4) show_status_codes; echo ""; read -p "Press Enter to continue..." ;;
        5) show_live_trace ;;
        6) search_trace ;;
        7) continue ;;
        8) echo "Goodbye! 👋"; exit 0 ;;
        *) echo "Invalid option"; sleep 1 ;;
    esac
done
