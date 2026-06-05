# 🗂 Portfolio CMS — Go Fiber + GORM + PostgreSQL + GoAdmin

Backend API untuk portfolio website dengan dashboard sederhana.

## 🧱 Tech Stack

| Layer      | Technology                         |
| ---------- | ---------------------------------- |
| Framework  | [Go Fiber v2](https://gofiber.io/) |
| ORM        | [GORM](https://gorm.io/)           |
| Database   | PostgreSQL 16                      |
| Auth       | JWT (golang-jwt/jwt v5)            |
| Validation | go-playground/validator v10        |
| API Docs   | Swagger (swaggo)                   |
| Hot Reload | Air                                |
| Deploy     | Railway / Render (Docker)          |

---

## 📁 Project Structure

```
portfolio-cms/
├── cmd/
│   └── api/
│       ├── main.go          # Entry point, wires everything
│       └── router.go        # All route definitions
├── internal/
│   ├── admin/
│   │   ├── admin.go         # GoAdmin setup & mount layout
form definitions per model
│   ├── config/
│   │   └── config.go        # Env config loader
│   ├── database/
│   │   └── database.go      # GORM connection + AutoMigrate
│   ├── handler/
│   │   ├── response.go      # Standard JSON response helpers
│   │   ├── validator.go     # Shared validator instance
│   │   ├── auth.go          # Auth handler
│   │   ├── project.go       # Project CRUD handler
│   │   └── handlers.go      # Skill, Experience, Education, Profile, Contact handlers
│   ├── middleware/
│   │   ├── jwt.go           # JWT middleware + token generation
│   │   └── middleware.go    # CORS, Logger, Recover
│   ├── model/
│   │   └── model.go         # All GORM models (User, Project, Skill, etc.)
│   └── service/
│       └── service.go       # Business logic layer
├── .air.toml                # Hot reload config
├── .env.example             # Environment variables template
├── .gitignore
├── docker-compose.yml       # Local dev with PostgreSQL
├── Dockerfile               # Multi-stage production build
├── go.mod
└── Makefile                 # Dev commands
```

---

## 🚀 Quick Start (Local)

### 1. Clone & install dependencies

```bash
git clone <your-repo> portfolio-cms
cd portfolio-cms
go mod tidy
```

### 2. Setup environment

```bash
cp .env.example .env
# Edit .env — set DB_PASSWORD, JWT_SECRET, etc.
```

### 3. Start PostgreSQL (Docker)

```bash
# Only postgres, tanpa build app
docker compose up -d postgres
```

### 4. Run with hot reload

```bash
# Install air dulu (sekali saja)
go install github.com/air-verse/air@latest

make dev
# atau langsung: air
```

### 5. Akses

| URL                                       | Keterangan              |
| ----------------------------------------- | ----------------------- |
| `http://localhost:8080/health`            | Health check            |
| `http://localhost:8080/admin`             | GoAdmin dashboard       |
| `http://localhost:8080/docs`              | Swagger API docs        |
| `http://localhost:8080/api/v1/public/...` | Public API (untuk Nuxt) |

---

## 🔒 Authentication

```bash
# Register (pertama kali)
POST /api/v1/auth/register
{
  "name": "Admin",
  "email": "admin@example.com",
  "password": "password123"
}

# Login
POST /api/v1/auth/login
{
  "email": "admin@example.com",
  "password": "password123"
}

# Gunakan access_token di header:
Authorization: Bearer <access_token>
```

---

## 🌐 Public API Endpoints (untuk Nuxt)

```
GET  /api/v1/public/profile          → Data profile/biodata
GET  /api/v1/public/projects         → List semua project (published)
GET  /api/v1/public/projects/:slug   → Detail project by slug
GET  /api/v1/public/skills           → List skills
GET  /api/v1/public/experiences      → List pengalaman kerja
GET  /api/v1/public/educations       → List pendidikan
POST /api/v1/public/contact          → Kirim pesan kontak
```

### Response format

```json
{
  "success": true,
  "message": "Success",
  "data": { ... },
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "total_pages": 3
  }
}
```

---

## 🛠 Admin API Endpoints (JWT required)

```
POST   /api/v1/admin/api/projects       → Create project
PUT    /api/v1/admin/api/projects/:id   → Update project
DELETE /api/v1/admin/api/projects/:id   → Delete project

# (sama untuk skills, experiences, educations)

POST   /api/v1/admin/api/profile        → Upsert profile
GET    /api/v1/admin/api/contacts       → List pesan masuk
PATCH  /api/v1/admin/api/contacts/:id/read → Mark pesan dibaca
```

---

## 📦 Deploy ke Railway

1. Push ke GitHub
2. Di Railway → **New Project → Deploy from GitHub**
3. Tambahkan **PostgreSQL** plugin di Railway
4. Set environment variables dari `.env.example`
5. Railway otomatis detect `Dockerfile` dan build

**Environment vars penting di Railway:**

```
DATABASE_URL  → otomatis dari Railway Postgres plugin
APP_PORT=8080
JWT_SECRET=your_production_secret
ADMIN_PASSWORD=your_secure_password
CORS_ALLOWED_ORIGINS=https://yoursite.com
```

---

## 🔧 Menambahkan CRUD Baru

Contoh menambahkan model `Certificate`:

**1. Tambah model** di `internal/model/model.go`:

```go
type Certificate struct {
    Base
    Title       string `gorm:"size:200" json:"title"`
    Issuer      string `gorm:"size:200" json:"issuer"`
    IssuedAt    time.Time `json:"issued_at"`
    CertURL     string `gorm:"size:300" json:"cert_url"`
    IsPublished bool   `gorm:"default:true" json:"is_published"`
}
```

**2. Tambah ke AutoMigrate** di `internal/database/database.go`:

```go
&model.Certificate{},
```

**3. Tambah service** di `internal/service/service.go`:

```go
type CertificateService interface { ... }
```

**4. Tambah handler** di `internal/handler/`:

```go
// handler/certificate.go
```

**5. Register route** di `cmd/api/router.go`

**6. Register ke GoAdmin** di `internal/admin/`:

```go
// tables.go → GetCertificateTable(...)
// generators.go → "certificates": GetCertificateTable,
```

---

## 📝 Generate Swagger Docs

```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

make swag
# atau: swag init -g cmd/api/main.go -o docs
```

---

## 🧪 Testing

```bash
make test
```
