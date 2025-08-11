# 🏛️ Haslaw Backend Services

Backend API untuk sistem Haslaw menggunakan Go dengan framework Gin dan MySQL database.

## 🚀 Features

- **Authentication & Authorization** - JWT based auth dengan role-based access
- **News Management** - CRUD operations untuk berita
- **Member Management** - Manajemen data member
- **File Upload** - Upload dan manajemen file
- **Rate Limiting** - Protection terhadap spam requests
- **Health Check** - Monitoring kesehatan aplikasi

## 🛠️ Tech Stack

- **Backend**: Go 1.23+ dengan Gin framework
- **Database**: MySQL 8.0 dengan GORM ORM
- **Authentication**: JWT (JSON Web Tokens)
- **Containerization**: Docker & Docker Compose
- **API Documentation**: RESTful API

## 📋 Prerequisites

- Go 1.23 atau lebih baru
- Docker & Docker Compose
- MySQL 8.0 (jika tidak menggunakan Docker)

## 🏗️ Development Setup

### 1. Clone Repository
```bash
git clone https://github.com/santaimany/haslaw-services.git
cd haslaw-services
```

### 2. Environment Configuration
```bash
cp .env.example .env
```

Edit file `.env` sesuai kebutuhan:
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=haslaw_user
DB_PASSWORD=your_password
DB_NAME=haslaw_db

JWT_SECRET=your-secret-key
JWT_REFRESH_SECRET=your-refresh-secret

PORT=8080
GIN_MODE=debug
```

### 3. Install Dependencies
```bash
go mod download
```

### 4. Database Migration
```bash
go run ./cmd/migrate
```

### 5. Seed Data (Optional)
```bash
go run ./cmd/seed
```

### 6. Run Application
```bash
go run ./cmd/api
```

Application akan berjalan di `http://localhost:8080`

## 🐳 Docker Deployment

Untuk deployment menggunakan Docker, lihat panduan lengkap di [DEPLOYMENT.md](DEPLOYMENT.md).

### Quick Start dengan Docker
```bash
# Build dan jalankan semua services
docker-compose up -d --build

# Check status
docker-compose ps

# View logs
docker-compose logs -f app
```

## 📚 API Documentation

### Authentication Endpoints
```
POST /api/v1/auth/register    - Register user baru
POST /api/v1/auth/login       - Login user
POST /api/v1/auth/refresh     - Refresh access token
POST /api/v1/auth/logout      - Logout user
```

### Admin Endpoints (Protected)
```
GET  /api/v1/admin/profile    - Get user profile
PUT  /api/v1/admin/profile    - Update user profile
```

### News Endpoints
```
GET    /api/v1/news           - List semua berita
GET    /api/v1/news/:id       - Get berita by ID
POST   /api/v1/news           - Create berita baru (Protected)
PUT    /api/v1/news/:id       - Update berita (Protected)
DELETE /api/v1/news/:id       - Delete berita (Protected)
```

### Member Endpoints (Protected)
```
GET    /api/v1/members        - List semua member
GET    /api/v1/members/:id    - Get member by ID
POST   /api/v1/members        - Create member baru
PUT    /api/v1/members/:id    - Update member
DELETE /api/v1/members/:id    - Delete member
```

### Health Check
```
GET /health                   - Application health status
```

## 🔐 Default Admin Account

Default super admin akan dibuat otomatis saat aplikasi pertama kali dijalankan:

- **Username**: `superadmin`
- **Email**: `admin@haslaw.com`
- **Password**: `SuperAdmin123!`

⚠️ **PENTING**: Ganti password default ini setelah login pertama kali!

## 🗂️ Project Structure

```
.
├── cmd/
│   ├── api/          # Main application
│   ├── migrate/      # Database migration
│   └── seed/         # Data seeding
├── internal/
│   ├── app/          # Application setup
│   ├── config/       # Configuration
│   ├── handlers/     # HTTP handlers
│   ├── middleware/   # Custom middleware
│   ├── models/       # Data models
│   ├── repository/   # Data access layer
│   ├── service/      # Business logic
│   └── utils/        # Utility functions
├── uploads/          # File uploads directory
├── docker-compose.yml
├── Dockerfile
└── .env.example
```

## 🛡️ Security Features

- **JWT Authentication** dengan access & refresh tokens
- **Password hashing** menggunakan bcrypt
- **Rate limiting** untuk mencegah spam
- **Input validation** untuk semua endpoints
- **SQL injection protection** dengan GORM
- **CORS protection**

## 📝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 🆘 Support

Jika ada pertanyaan atau issue, silakan buat issue di GitHub repository atau hubungi tim development.
