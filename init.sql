# Create the database and user if they don't exist
CREATE DATABASE IF NOT EXISTS haslaw_db;
CREATE USER IF NOT EXISTS 'haslaw_user'@'%' IDENTIFIED BY 'haslaw_password';
GRANT ALL PRIVILEGES ON haslaw_db.* TO 'haslaw_user'@'%';
FLUSH PRIVILEGES;
