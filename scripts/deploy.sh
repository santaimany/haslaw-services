#!/bin/bash

# HasLaw Backend Deployment Script
set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

PROJECT_DIR="/opt/haslaw-be-services"

echo -e "${BLUE}🚀 HasLaw Backend Deployment${NC}"
echo "=============================="

# Check if we're in the right directory
if [ ! -f "docker-compose.yml" ]; then
    if [ -d "$PROJECT_DIR" ]; then
        cd "$PROJECT_DIR"
    else
        echo -e "${RED}❌ Project directory not found${NC}"
        exit 1
    fi
fi

# Load environment variables
if [ -f ".env.production" ]; then
    export $(cat .env.production | grep -v '^#' | xargs)
    echo -e "${GREEN}✅ Loaded production environment${NC}"
else
    echo -e "${YELLOW}⚠️  No .env.production found, using defaults${NC}"
fi

# Check if services are running
if docker-compose ps | grep -q "Up"; then
    echo -e "${YELLOW}📦 Creating database backup...${NC}"
    BACKUP_FILE="backup_$(date +%Y%m%d_%H%M%S).sql"
    if docker-compose exec -T db mysqladmin ping -h localhost -u ${DB_USER:-haslaw_user} -p${DB_PASSWORD:-haslaw_password} > /dev/null 2>&1; then
        docker-compose exec -T db mysqldump -u ${DB_USER:-haslaw_user} -p${DB_PASSWORD:-haslaw_password} ${DB_NAME:-haslaw_db} > "$BACKUP_FILE"
        echo -e "${GREEN}✅ Backup created: $BACKUP_FILE${NC}"
    else
        echo -e "${YELLOW}⚠️  Could not create backup - database not accessible${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  No running services detected${NC}"
fi

# Pull latest changes
echo -e "${YELLOW}📥 Pulling latest changes...${NC}"
git pull origin main

# Pull latest Docker images
echo -e "${YELLOW}🐳 Pulling latest Docker images...${NC}"
docker-compose pull

# Start/restart services
echo -e "${YELLOW}🔄 Starting services...${NC}"
docker-compose up -d

# Wait for services to be ready
echo -e "${YELLOW}⏳ Waiting for services to be ready...${NC}"
sleep 15

# Health check with retry
echo -e "${YELLOW}🏥 Running health checks...${NC}"
HEALTH_CHECK_URL="http://localhost:${PORT:-8080}/health"

for i in {1..6}; do
    if curl -f "$HEALTH_CHECK_URL" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Health check passed!${NC}"
        HEALTH_OK=true
        break
    else
        if [ $i -eq 6 ]; then
            echo -e "${RED}❌ Health check failed after 6 attempts${NC}"
            echo -e "${YELLOW}📋 Service status:${NC}"
            docker-compose ps
            echo -e "${YELLOW}📋 Application logs:${NC}"
            docker-compose logs --tail=20 app
            exit 1
        fi
        echo -e "${YELLOW}⏳ Health check attempt $i failed, retrying in 10s...${NC}"
        sleep 10
    fi
done

# Show service status
echo -e "${YELLOW}📋 Service status:${NC}"
docker-compose ps

# Clean up old Docker images
echo -e "${YELLOW}🧹 Cleaning up old images...${NC}"
docker image prune -f

# Show final status
echo ""
echo -e "${GREEN}🎉 Deployment completed successfully!${NC}"
echo -e "${BLUE}📊 Service Information:${NC}"
echo "- Application URL: http://localhost:${PORT:-8080}"
echo "- Health Check: $HEALTH_CHECK_URL"
echo "- Database: ${DB_NAME:-haslaw_db}"
echo "- Environment: $(cat .env.production 2>/dev/null | grep GIN_MODE | cut -d'=' -f2 || echo 'production')"
echo ""
echo -e "${YELLOW}💡 Useful commands:${NC}"
echo "- View logs: docker-compose logs -f app"
echo "- Check status: docker-compose ps"
echo "- Restart: docker-compose restart app"
echo "- Stop: docker-compose down"
