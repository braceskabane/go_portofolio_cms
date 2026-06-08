# 📇 Portfolio CMS — Go Fiber + GORM + PostgreSQL

Backend API + lightweight admin dashboard untuk personal portfolio website.  
Mendukung manajemen proyek, pengalaman kerja, tech stack, skill, galeri/dokumen polimorfik, hingga running activity tracker.

---

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

## 📁 Project Structure (final)

```
portfolio-cms/
├── cmd/
│   └── api/
│       ├── main.go              # Entry point, wiring, server start
│       └── router.go            # All route definitions (public & admin)
├── internal/
│   ├── admin/
│   │   ├── admin.go             # Admin UI route registrations
│   │   ├── auth.go              # Admin login/logout handlers
│   │   ├── dashboard.go         # Dashboard page
│   │   ├── projects.go          # Project CRUD pages & form handlers
│   │   ├── assets.go            # Asset (gallery/doc) pages
│   │   ├── project_categories.go
│   │   ├── experience_categories.go
│   │   ├── stack_categories.go
│   │   ├── stack.go
│   │   ├── skills.go
│   │   ├── experiences.go
│   │   ├── educations.go
│   │   ├── profile.go
│   │   ├── contacts.go
|   |   ├── layout.go
│   │   └── running.go
│   ├── config/
│   │   └── config.go            # Env config loader
│   ├── database/
│   │   └── database.go          # GORM connection + AutoMigrate all models
│   ├── handler/
│   │   ├── response.go          # Standard JSON response helpers (OK, Created, etc.)
│   │   ├── validator.go         # Shared validator instance
│   │   ├── auth.go              # Auth handler (register, login, refresh, me)
│   │   ├── project.go           # Project CRUD handler
│   │   ├── asset.go             # Asset CRUD handler (polymorphic)
│   │   ├── project_category.go  # Project category CRUD handler
│   │   ├── experience_category.go
│   │   ├── stack_category.go
│   │   ├── stack_item.go
│   │   ├── skill.go             # Skill + detail handler
│   │   ├── experience.go
│   │   ├── education.go
│   │   ├── profile.go
│   │   ├── contact.go
│   │   └── running_activity.go
│   ├── middleware/
│   │   ├── jwt.go               # JWT middleware + token generation
│   │   ├── session.go           # Session middleware for admin UI
│   │   └── middleware.go        # CORS, Logger, Recover
│   ├── model/
│   │   └── model.go             # All GORM models (User, Profile, Project, Asset, Skill, …)
│   ├── service/
│   │   └── service.go           # All business logic interfaces & implementations
│   └── dto/
│       └── dto.go               # Request/response DTOs
├── docs/                        # Swagger generated docs
├── .air.toml                    # Hot reload config
├── .env.example                 # Environment variables template
├── .gitignore
├── docker-compose.yml           # Local dev with PostgreSQL
├── Dockerfile                   # Multi-stage production build
├── go.mod
├── Makefile                     # Dev commands
└── README.md
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
docker compose up -d postgres
```

### 4. Run with hot reload

```bash
# Install air (first time only)
go install github.com/air-verse/air@latest

make dev
# or directly: air
```

### 5. Access

| URL                                       | Description               |
| ----------------------------------------- | ------------------------- |
| `http://localhost:8080/health`            | Health check              |
| `http://localhost:8080/admin`             | Custom admin dashboard    |
| `http://localhost:8080/docs/index.html`   | Swagger API docs          |
| `http://localhost:8080/api/v1/public/...` | Public API (for frontend) |

---

## 🔒 Authentication

```bash
# Register (first time, creates admin user)
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

# Response includes access_token & refresh_token
# Use access_token in header:
Authorization: Bearer <access_token>
```

---

## 🌐 Public API Endpoints (for Nuxt / frontend)

### Profile

```
GET    /api/v1/public/profile
```

### Projects

```
GET    /api/v1/public/projects?page=1&limit=10&featured=true
GET    /api/v1/public/projects/:slug
```

### Skills

```
GET    /api/v1/public/skills
GET    /api/v1/public/skills/:id        → detail skill + total projects & category info
```

### Experiences

```
GET    /api/v1/public/experiences
```

### Educations

```
GET    /api/v1/public/educations
```

### Contact

```
POST   /api/v1/public/contact
```

### Tech Stack

```
GET    /api/v1/public/stack-categories?with=items
GET    /api/v1/public/stack-items?category_id=<uuid>
```

### Running Activities

```
GET    /api/v1/public/running-activities?page=1&limit=10
GET    /api/v1/public/running-activities/:id
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

Semua route di bawah memerlukan header `Authorization: Bearer <access_token>`  
Prefix: `/api/v1/admin/api`

### Projects

```
GET    /projects
POST   /projects
GET    /projects/:id
PUT    /projects/:id
DELETE /projects/:id
```

### Assets (polymorphic: photo/video/pdf/doc for project or experience)

```
GET    /assets?owner_type=project&owner_id=<uuid>
POST   /assets?owner_type=project&owner_id=<uuid>
PUT    /assets/:id
DELETE /assets/:id
```

### Project Categories

```
GET    /project-categories
POST   /project-categories
PUT    /project-categories/:id
DELETE /project-categories/:id
```

### Experience Categories

```
GET    /experience-categories
POST   /experience-categories
PUT    /experience-categories/:id
DELETE /experience-categories/:id
```

### Tech Stack Categories

```
GET    /stack-categories
GET    /stack-categories/:id
POST   /stack-categories
PUT    /stack-categories/:id
DELETE /stack-categories/:id
```

### Tech Stack Items

```
GET    /stack-items?category_id=<uuid>
GET    /stack-items/:id
POST   /stack-items
PUT    /stack-items/:id
DELETE /stack-items/:id
```

### Skills

```
GET    /skills
POST   /skills
PUT    /skills/:id
DELETE /skills/:id
```

### Experiences

```
GET    /experiences
POST   /experiences
PUT    /experiences/:id
DELETE /experiences/:id
```

### Educations

```
GET    /educations
POST   /educations
PUT    /educations/:id
DELETE /educations/:id
```

### Profile

```
GET    /profile
POST   /profile         → Upsert (create or update)
```

### Contacts

```
GET    /contacts
PATCH  /contacts/:id/read
```

### Running Activities

```
GET    /running-activities
POST   /running-activities
PUT    /running-activities/:id
DELETE /running-activities/:id
```

---

## 🧩 Data Models (Hubungan Utama)

- **User** – satu admin/superadmin.
- **Profile** – single row biodata lengkap (medsos, CV, status hire).
- **Project** – punya `Category` (ProjectCategory), `Skills` (many2many), `StackItems` (many2many), dan `Assets` (polymorphic).
- **Experience** – punya `Category` (ExperienceCategory), `Skills`, `StackItems`, dan `Assets`.
- **Skill** – bisa dimiliki banyak project & experience, memiliki proficiency 0-100.
- **StackCategory** > **StackItem** – tech stack terorganisir (Frontend, Backend, dll).
- **Asset** – polymorphic; satu tabel untuk semua media & dokumen milik project atau experience (owner_type + owner_id).
- **Education**, **Contact**, **RunningActivity** – mandiri.

---

## 📦 Deploy ke Railway

1. Push ke GitHub.
2. Di Railway → **New Project → Deploy from GitHub**.
3. Tambahkan **PostgreSQL** plugin di Railway.
4. Set environment variables dari `.env.example`.
5. Railway otomatis detect `Dockerfile` dan build.

**Environment vars penting di Railway:**

```
DATABASE_URL         → otomatis dari Railway Postgres plugin
APP_PORT=8080
JWT_SECRET=your_production_secret
ADMIN_PASSWORD=your_secure_password
CORS_ALLOWED_ORIGINS=https://yoursite.com
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
