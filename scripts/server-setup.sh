#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔧 HasLaw Backend Services - Server Setup Script${NC}"
echo "=================================================="

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}❌ Please run as root (use sudo)${NC}"
    exit 1
fi

# Update system
echo -e "${YELLOW}📦 Updating system packages...${NC}"
apt update && apt upgrade -y

# Install Docker
echo -e "${YELLOW}🐳 Installing Docker...${NC}"
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh get-docker.sh
    usermod -aG docker $SUDO_USER
    rm get-docker.sh
else
    echo -e "${GREEN}✅ Docker already installed${NC}"
fi

# Install Docker Compose
echo -e "${YELLOW}🐳 Installing Docker Compose...${NC}"
if ! command -v docker-compose &> /dev/null; then
    curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
else
    echo -e "${GREEN}✅ Docker Compose already installed${NC}"
fi

# Install Git
echo -e "${YELLOW}📥 Installing Git...${NC}"
if ! command -v git &> /dev/null; then
    apt install -y git
else
    echo -e "${GREEN}✅ Git already installed${NC}"
fi

# Install other useful tools
echo -e "${YELLOW}🛠️ Installing additional tools...${NC}"
apt install -y curl wget htop nano ufw fail2ban

# Setup firewall
echo -e "${YELLOW}🔥 Configuring firewall...${NC}"
ufw --force enable
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow 80
ufw allow 443
ufw allow 8080

# Create deployment user
echo -e "${YELLOW}👤 Creating deployment user...${NC}"
if ! id "deploy" &>/dev/null; then
    useradd -m -s /bin/bash deploy
    usermod -aG docker deploy
    echo -e "${GREEN}✅ User 'deploy' created${NC}"
else
    echo -e "${GREEN}✅ User 'deploy' already exists${NC}"
fi

# Setup SSH key for deployment user
echo -e "${YELLOW}🔑 Setting up SSH for deployment...${NC}"
sudo -u deploy mkdir -p /home/deploy/.ssh
sudo -u deploy touch /home/deploy/.ssh/authorized_keys
chmod 700 /home/deploy/.ssh
chmod 600 /home/deploy/.ssh/authorized_keys
chown -R deploy:deploy /home/deploy/.ssh

echo -e "${BLUE}📋 SSH Public Key Setup:${NC}"
echo "Add your GitHub Actions public key to: /home/deploy/.ssh/authorized_keys"
echo "Or run: echo 'YOUR_PUBLIC_KEY' >> /home/deploy/.ssh/authorized_keys"

# Create deployment directory
echo -e "${YELLOW}📁 Creating deployment directory...${NC}"
mkdir -p /opt/haslaw-be-services
chown deploy:deploy /opt/haslaw-be-services

# Setup log rotation
echo -e "${YELLOW}📝 Setting up log rotation...${NC}"
cat > /etc/logrotate.d/haslaw-docker << EOF
/var/lib/docker/containers/*/*.log {
    rotate 7
    daily
    compress
    size 100M
    missingok
    delaycompress
    copytruncate
}
EOF

# Setup automatic security updates
echo -e "${YELLOW}🔒 Setting up automatic security updates...${NC}"
apt install -y unattended-upgrades
echo 'Unattended-Upgrade::Automatic-Reboot "false";' >> /etc/apt/apt.conf.d/50unattended-upgrades

# Setup system monitoring
echo -e "${YELLOW}📊 Installing monitoring tools...${NC}"
apt install -y htop iotop nethogs

# Create deployment script
echo -e "${YELLOW}📜 Creating deployment helper script...${NC}"
cat > /home/deploy/deploy.sh << 'EOF'
#!/bin/bash

# HasLaw Backend Deployment Script
set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

PROJECT_DIR="/opt/haslaw-be-services"

cd $PROJECT_DIR

echo -e "${YELLOW}🚀 Starting deployment...${NC}"

# Backup database
echo -e "${YELLOW}💾 Creating database backup...${NC}"
docker-compose exec -T db mysqldump -u haslaw_user -phaslaw_password haslaw_db > backup_$(date +%Y%m%d_%H%M%S).sql

# Pull latest changes
echo -e "${YELLOW}📥 Pulling latest changes...${NC}"
git pull origin main

# Update containers
echo -e "${YELLOW}🐳 Updating containers...${NC}"
docker-compose pull
docker-compose up -d

# Health check
echo -e "${YELLOW}🏥 Running health check...${NC}"
sleep 10
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Deployment successful!${NC}"
else
    echo -e "${RED}❌ Health check failed!${NC}"
    exit 1
fi

# Cleanup
echo -e "${YELLOW}🧹 Cleaning up...${NC}"
docker image prune -f

echo -e "${GREEN}🎉 Deployment completed!${NC}"
EOF

chmod +x /home/deploy/deploy.sh
chown deploy:deploy /home/deploy/deploy.sh

# Setup systemd service for auto-start
echo -e "${YELLOW}⚙️ Setting up systemd service...${NC}"
cat > /etc/systemd/system/haslaw-backend.service << EOF
[Unit]
Description=HasLaw Backend Services
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/haslaw-be-services
ExecStart=/usr/local/bin/docker-compose up -d
ExecStop=/usr/local/bin/docker-compose down
TimeoutStartSec=0
User=deploy
Group=deploy

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable haslaw-backend.service

echo -e "${GREEN}✅ Server setup completed!${NC}"
echo ""
echo -e "${BLUE}📋 Next Steps:${NC}"
echo "1. Add your SSH public key to /home/deploy/.ssh/authorized_keys"
echo "2. Clone your repository to /opt/haslaw-be-services"
echo "3. Configure your GitHub Secrets:"
echo "   - HOST: your-server-ip"
echo "   - USERNAME: deploy"
echo "   - SSH_KEY: your-private-ssh-key"
echo "   - DEPLOY_PATH: /opt/haslaw-be-services"
echo "   - PORT: 22 (or your SSH port)"
echo "   - DB_PASSWORD: your-database-password"
echo "4. Setup your .env.production file"
echo "5. Test deployment with: sudo -u deploy /home/deploy/deploy.sh"
echo ""
echo -e "${YELLOW}🔐 Security Reminders:${NC}"
echo "- Change default passwords"
echo "- Setup SSL certificates"
echo "- Configure domain name in nginx.conf"
echo "- Setup regular backups"
echo ""
echo -e "${GREEN}🎉 Server is ready for deployment!${NC}"
