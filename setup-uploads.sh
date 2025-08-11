#!/bin/bash

# Script untuk setup permissions untuk bind mount
# Jalankan script ini di server sebelum docker-compose up

echo "🔧 Setting up uploads directory permissions..."

# Buat direktori uploads jika belum ada
mkdir -p uploads/news uploads/members

# Set ownership ke user yang sama dengan container (UID 1001)
# Jika tidak ada user dengan UID 1001, gunakan current user
if id -u 1001 >/dev/null 2>&1; then
    echo "Setting ownership to UID 1001..."
    sudo chown -R 1001:1001 uploads/
else
    echo "Setting ownership to current user..."
    sudo chown -R $USER:$USER uploads/
fi

# Set permissions
chmod -R 755 uploads/

echo "✅ Permissions setup completed!"
echo "Directory structure:"
ls -la uploads/
