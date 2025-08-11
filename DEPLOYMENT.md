# 🚀 Manual Deployment Guide

## Prerequisites
- Docker & Docker Compose installed on server
- Git installed on server

## 🔧 Server Setup

### 1. Clone Repository
```bash
git clone https://github.com/santaimany/haslaw-services.git
cd haslaw-services
```

### 2. Environment Configuration
Create `.env` file:
```bash
cp .env.example .env
nano .env
```

Update the values:
```env
# Database Configuration
DB_HOST=db
DB_PORT=3306
DB_USER=haslaw_user
DB_PASSWORD=haslaw_password
DB_NAME=haslaw_db

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here
JWT_REFRESH_SECRET=your-super-secret-refresh-key-here

# Server Configuration
PORT=8080
GIN_MODE=release
```

### 3. Deploy Application
```bash
# Build and start all services (without nginx)
docker-compose up -d --build

# OR if you want to use nginx reverse proxy (optional)
# First make sure nginx.conf exists, then add nginx service back to docker-compose.yml

# Check if services are running
docker-compose ps

# View logs
docker-compose logs -f app
```

### 4. Health Check
```bash
# Test if application is running (direct access)
curl http://localhost:8080/health

# Test API endpoint
curl http://localhost:8080/api/v1/auth/login
```

**Note:** Current setup runs the app directly on port 8080 without nginx. If you need nginx reverse proxy, uncomment the nginx service in docker-compose.yml.

## 🔄 Update Deployment

When you want to update the application:

```bash
# Pull latest changes
git pull origin main

# Rebuild and restart
docker-compose up -d --build

# Clean up old images (optional)
docker image prune -f
```

## 📋 Useful Commands

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f app
docker-compose logs -f db
```

### Restart Services
```bash
# Restart all
docker-compose restart

# Restart specific service
docker-compose restart app
```

### Database Management
```bash
# Access MySQL
docker-compose exec db mysql -u haslaw_user -p haslaw_db

# Backup database
docker-compose exec db mysqldump -u root -p haslaw_db > backup.sql

# Import database
docker-compose exec -T db mysql -u root -p haslaw_db < backup.sql
```

### Stop Services
```bash
# Stop all services
docker-compose down

# Stop and remove volumes (⚠️ This will delete database data)
docker-compose down -v
```

## 🔐 Production Security

1. **Change default passwords** in `.env`
2. **Use HTTPS** with reverse proxy (Nginx/Traefik)
3. **Set up firewall** to only allow necessary ports
4. **Regular backups** of database

## 🌐 Nginx Reverse Proxy (Optional)

Create `/etc/nginx/sites-available/haslaw`:
```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable site:
```bash
sudo ln -s /etc/nginx/sites-available/haslaw /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```
