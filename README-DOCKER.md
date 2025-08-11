# Docker Deployment Guide for HasLaw Backend Services

## Prerequisites
- Docker and Docker Compose installed on your server
- At least 2GB RAM and 20GB disk space
- Open ports: 80, 443, 8080, 3306 (or configure as needed)

## Quick Start

### 1. Clone and Setup
```bash
git clone <your-repo-url>
cd haslaw-be-services
```

### 2. Environment Configuration
```bash
# Copy example environment file
cp .env.example .env.production

# Edit production environment variables
nano .env.production
```

**Important: Change these values in .env.production:**
- `JWT_SECRET`: Use a strong, random 32+ character string
- `DB_PASSWORD`: Set a secure database password
- Update other settings as needed for your server

### 3. Build and Run
```bash
# Build and start all services
docker-compose up -d

# Check if services are running
docker-compose ps

# View logs
docker-compose logs -f app
```

### 4. Initial Setup
```bash
# Run database migrations (if you have migration commands)
docker-compose exec app ./main migrate

# Create default super admin (if you have seed commands)
docker-compose exec app ./main seed
```

## Production Deployment

### SSL Certificate Setup
1. Obtain SSL certificates (Let's Encrypt recommended)
2. Place certificates in `./ssl/` directory:
   - `cert.pem` - Certificate file
   - `key.pem` - Private key file
3. Uncomment SSL configuration in `nginx.conf`
4. Update server_name in `nginx.conf` with your domain

### Environment Variables for Production
```bash
# Database
DB_HOST=db
DB_PORT=3306
DB_USER=haslaw_user
DB_PASSWORD=your_very_secure_password_here
DB_NAME=haslaw_db

# JWT (CRITICAL: Change these!)
JWT_SECRET=your-super-secure-jwt-secret-minimum-32-characters
JWT_ACCESS_TOKEN_EXPIRE=15m
JWT_REFRESH_TOKEN_EXPIRE=168h

# Server
PORT=8080
GIN_MODE=release

# Security
BCRYPT_COST=12
RATE_LIMIT_PER_MINUTE=60
```

## Service Architecture

- **app**: Go backend application (port 8080)
- **db**: MySQL 8.0 database (port 3306)
- **nginx**: Reverse proxy with SSL termination (ports 80/443)

## Useful Commands

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f app
docker-compose logs -f db
docker-compose logs -f nginx
```

### Database Access
```bash
# Connect to MySQL
docker-compose exec db mysql -u haslaw_user -p haslaw_db

# Database backup
docker-compose exec db mysqldump -u haslaw_user -p haslaw_db > backup.sql

# Database restore
docker-compose exec -T db mysql -u haslaw_user -p haslaw_db < backup.sql
```

### Application Management
```bash
# Restart application
docker-compose restart app

# Rebuild application
docker-compose build app
docker-compose up -d app

# Scale application (multiple instances)
docker-compose up -d --scale app=3
```

### Updates
```bash
# Pull latest code
git pull origin main

# Rebuild and restart
docker-compose build
docker-compose up -d
```

## Monitoring and Health Checks

### Health Check Endpoints
- Application: `http://your-domain.com/health`
- Database: Built-in Docker health checks

### View Service Status
```bash
docker-compose ps
docker stats
```

## Security Considerations

1. **Change default passwords** in production
2. **Use strong JWT secrets** (32+ characters)
3. **Enable SSL/HTTPS** for production
4. **Regular security updates**:
   ```bash
   docker-compose pull
   docker-compose up -d
   ```
5. **Firewall configuration**: Only expose necessary ports
6. **Regular backups** of database and application data

## Troubleshooting

### Common Issues

**Application won't start:**
```bash
docker-compose logs app
```

**Database connection issues:**
```bash
docker-compose logs db
docker-compose exec app ping db
```

**Permission issues with uploads:**
```bash
sudo chown -R 1001:1001 uploads/
```

**Memory issues:**
```bash
docker system prune -a
```

### Performance Tuning

**For production servers:**
1. Increase MySQL buffer pool size
2. Optimize nginx worker processes
3. Monitor resource usage with `docker stats`
4. Consider using Docker Swarm or Kubernetes for high availability

## Backup Strategy

### Database Backup
```bash
# Daily backup script
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
docker-compose exec -T db mysqldump -u haslaw_user -p haslaw_db > "backup_${DATE}.sql"
```

### File Backup
```bash
# Backup uploads directory
tar -czf uploads_backup_$(date +%Y%m%d).tar.gz uploads/
```

## Support

For issues:
1. Check logs: `docker-compose logs -f`
2. Verify environment variables
3. Ensure all services are healthy: `docker-compose ps`
4. Check resource usage: `docker stats`
