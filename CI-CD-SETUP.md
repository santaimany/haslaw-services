# CI/CD Setup Guide

## Overview
This setup provides automated CI/CD pipeline using GitHub Actions that will:
1. **Test** your code on every push/PR
2. **Build** Docker images and push to GitHub Container Registry
3. **Deploy** automatically to your server when pushing to main branch

## Server Setup

### 1. Prepare Your Server
Run this on your Ubuntu/Debian server:
```bash
# Download and run server setup script
curl -sSL https://raw.githubusercontent.com/santaimany/haslaw-services/main/scripts/server-setup.sh | sudo bash
```

Or manually download and run:
```bash
wget https://raw.githubusercontent.com/santaimany/haslaw-services/main/scripts/server-setup.sh
chmod +x server-setup.sh
sudo ./server-setup.sh
```

### 2. Generate SSH Key for Deployment
On your local machine:
```bash
# Generate SSH key pair
ssh-keygen -t ed25519 -f ~/.ssh/haslaw-deploy -C "haslaw-deploy"

# Copy public key to server
ssh-copy-id -i ~/.ssh/haslaw-deploy.pub deploy@your-server-ip

# Test connection
ssh -i ~/.ssh/haslaw-deploy deploy@your-server-ip
```

### 3. Setup Project on Server
```bash
# Login as deploy user
sudo su - deploy

# Clone repository
cd /opt
git clone https://github.com/santaimany/haslaw-services.git haslaw-be-services
cd haslaw-be-services

# Setup environment
cp .env.example .env.production
nano .env.production  # Edit with your production values

# Initial deployment
chmod +x scripts/deploy.sh
./scripts/deploy.sh
```

## GitHub Repository Setup

### 1. Enable GitHub Container Registry
1. Go to your GitHub repository
2. Go to **Settings > Actions > General**
3. Under "Workflow permissions", select **Read and write permissions**
4. Save changes

### 2. Add Repository Secrets
Go to **Settings > Secrets and variables > Actions** and add:

#### Required Secrets:
```
HOST=your-server-ip-address
USERNAME=deploy
SSH_KEY=contents-of-your-private-ssh-key
DEPLOY_PATH=/opt/haslaw-be-services
PORT=22
DB_PASSWORD=your-production-db-password
```

#### Optional Secrets (for notifications):
```
SLACK_WEBHOOK=your-slack-webhook-url
```

**Note:** Slack notification is optional. If you don't need it, the workflow will skip the notification step automatically.

#### For staging environment:
```
STAGING_HOST=your-staging-server-ip
STAGING_DEPLOY_PATH=/opt/haslaw-be-services-staging
STAGING_PORT=22
```

### 3. Configure SSH Key Secret
To add your SSH private key:
```bash
# Copy private key content
cat ~/.ssh/haslaw-deploy

# Copy the entire output (including -----BEGIN and -----END lines)
# Paste this as the SSH_KEY secret in GitHub
```

## Workflow Explanation

### Automatic Pipeline (ci-cd.yml)
Triggers on:
- Push to `main` or `develop` branches
- Pull requests to `main`

Steps:
1. **Test**: Runs Go tests with MySQL
2. **Build**: Creates Docker image and pushes to GitHub Container Registry
3. **Deploy**: SSH into server and deploy new version
4. **Notify**: Send deployment status (optional)

### Manual Pipeline (manual-deploy.yml)
- Can be triggered manually from GitHub Actions tab
- Allows selecting environment (production/staging)
- Allows selecting specific version/tag to deploy

## Usage

### Automatic Deployment
```bash
# Make changes to your code
git add .
git commit -m "feat: add new feature"
git push origin main

# GitHub Actions will automatically:
# 1. Run tests
# 2. Build Docker image
# 3. Deploy to server
# 4. Send notification
```

### Manual Deployment
1. Go to GitHub repository
2. Click **Actions** tab
3. Select **Deploy to Production**
4. Click **Run workflow**
5. Choose environment and version
6. Click **Run workflow**

### Monitor Deployments
- **GitHub Actions**: Check workflow runs
- **Server logs**: `docker-compose logs -f app`
- **Health check**: `curl http://your-domain.com/health`

## Security Best Practices

### 1. SSH Security
```bash
# On server, edit SSH config
sudo nano /etc/ssh/sshd_config

# Add these settings:
PermitRootLogin no
PasswordAuthentication no
PubkeyAuthentication yes
AllowUsers deploy

# Restart SSH
sudo systemctl restart sshd
```

### 2. Environment Variables
- Never commit `.env.production` to repository
- Use strong, unique passwords
- Rotate secrets regularly
- Use different secrets for staging/production

### 3. Server Security
```bash
# Setup fail2ban for SSH protection
sudo systemctl enable fail2ban
sudo systemctl start fail2ban

# Configure firewall
sudo ufw status
```

## Troubleshooting

### Common Issues

#### 1. SSH Connection Failed
```bash
# Test SSH connection
ssh -i ~/.ssh/haslaw-deploy deploy@your-server-ip

# Check SSH key format in GitHub Secrets
# Must include -----BEGIN and -----END lines
```

#### 2. Docker Permission Denied
```bash
# On server, add deploy user to docker group
sudo usermod -aG docker deploy
sudo systemctl restart docker
```

#### 3. Database Connection Failed
```bash
# Check database status
docker-compose ps
docker-compose logs db

# Reset database password
docker-compose exec db mysql -u root -p
```

#### 4. Health Check Failed
```bash
# Check application logs
docker-compose logs app

# Test health endpoint
curl -v http://localhost:8080/health
```

### Debugging Deployment
```bash
# On server, check deployment logs
journalctl -u haslaw-backend.service

# Manual deployment test
sudo su - deploy
cd /opt/haslaw-be-services
./scripts/deploy.sh
```

## Monitoring

### Application Monitoring
```bash
# View all services
docker-compose ps

# View application logs
docker-compose logs -f app

# View database logs
docker-compose logs -f db

# View nginx logs
docker-compose logs -f nginx

# Check resource usage
docker stats
```

### System Monitoring
```bash
# System resources
htop
df -h
free -h

# Network usage
sudo nethogs

# Disk I/O
sudo iotop
```

## Rollback Procedure

### Automatic Rollback
If health check fails, the workflow will automatically rollback.

### Manual Rollback
```bash
# On server
sudo su - deploy
cd /opt/haslaw-be-services

# Check available images
docker images

# Rollback to previous image
docker-compose down
docker tag ghcr.io/yourusername/haslaw-services:previous ghcr.io/yourusername/haslaw-services:latest
docker-compose up -d

# Or restore from backup
docker-compose exec -T db mysql -u haslaw_user -p haslaw_db < backup_YYYYMMDD_HHMMSS.sql
```

## Advanced Configuration

### Blue-Green Deployment
For zero-downtime deployments, you can setup blue-green deployment:

1. Run two identical environments
2. Deploy to inactive environment
3. Switch traffic after health check
4. Keep old environment as backup

### Staging Environment
Setup separate staging server for testing:

1. Create staging server with same setup
2. Add staging secrets to GitHub
3. Use manual deployment workflow
4. Test before promoting to production

### Database Migrations
Add migration step to deployment:

```yaml
# In ci-cd.yml, add before health check:
- name: Run database migrations
  run: docker-compose exec -T app ./main migrate
```

## Support

For issues with CI/CD setup:
1. Check GitHub Actions logs
2. Test SSH connection manually
3. Verify server setup script completed successfully
4. Check Docker and application logs on server
