Berikut README yang sudah diupdate dengan fitur running analysis:

```markdown
# 📇 Portfolio CMS — Go Fiber + GORM + PostgreSQL

Backend API + lightweight admin dashboard untuk personal portfolio website.  
Mendukung manajemen proyek, pengalaman kerja, tech stack, skill, galeri/dokumen polimorfik, running activity tracker, hingga **AI-powered running analysis** dengan integrasi Google Calendar.

---

## 🧱 Tech Stack

| Layer       | Technology                          |
| ----------- | ----------------------------------- |
| Framework   | [Go Fiber v2](https://gofiber.io/)  |
| ORM         | [GORM](https://gorm.io/)            |
| Database    | PostgreSQL 16                       |
| Auth        | JWT (golang-jwt/jwt v5)             |
| Validation  | go-playground/validator v10         |
| API Docs    | Swagger (swaggo)                    |
| AI Analysis | Google Gemini Flash 2.5 (free tier) |
| Calendar    | Google Calendar API v3 (OAuth 2.0)  |
| Hot Reload  | Air                                 |
| Deploy      | Railway / Render (Docker)           |

---

## 📁 Project Structure
```

portfolio-cms/
├── cmd/
│ └── api/
│ ├── main.go
│ └── router.go
├── internal/
│ ├── admin/
│ │ ├── admin.go
│ │ ├── auth.go
│ │ ├── config.go # Global vars (googleClientID, Init())
│ │ ├── dashboard.go
│ │ ├── projects.go
│ │ ├── assets.go
│ │ ├── project_categories.go
│ │ ├── experience_categories.go
│ │ ├── stack_categories.go
│ │ ├── stack.go
│ │ ├── skills.go
│ │ ├── experiences.go
│ │ ├── educations.go
│ │ ├── profile.go
│ │ ├── contacts.go
│ │ ├── layout.go
│ │ ├── running.go
│ │ └── running_analysis.go # AI analysis page + Calendar integration
│ ├── config/
│ │ └── config.go
│ ├── database/
│ │ └── database.go
│ ├── handler/
│ │ ├── response.go
│ │ ├── validator.go
│ │ ├── auth.go
│ │ ├── project.go
│ │ ├── asset.go
│ │ ├── project_category.go
│ │ ├── experience_category.go
│ │ ├── stack_category.go
│ │ ├── stack_item.go
│ │ ├── skill.go
│ │ ├── experience.go
│ │ ├── education.go
│ │ ├── profile.go
│ │ ├── contact.go
│ │ ├── running_activity.go
│ │ └── running_analysis.go # API handler untuk analysis & calendar sync
│ ├── middleware/
│ │ ├── jwt.go
│ │ ├── session.go
│ │ └── middleware.go
│ ├── model/
│ │ └── model.go
│ ├── service/
│ │ ├── service.go
│ │ ├── running_analysis_service.go # Gemini prompt builder + analisis
│ │ ├── running_analysis_prompts.go # System prompt ilmiah pelatihan lari
│ │ └── calendar_service.go # Google Calendar API integration
│ └── dto/
│ ├── dto.go
│ └── running_analysis.go # DTO untuk analisis & calendar events
├── docs/
├── .air.toml
├── .env.example
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── Makefile
└── README.md

````

---

## 🚀 Quick Start (Local)

### 1. Clone & install dependencies

```bash
git clone <your-repo> portfolio-cms
cd portfolio-cms
go mod tidy
````

### 2. Setup environment

```bash
cp .env.example .env
# Edit .env sesuai kebutuhan
```

### 3. Start PostgreSQL (Docker)

```bash
docker compose up -d postgres
```

### 4. Run with hot reload

```bash
go install github.com/air-verse/air@latest
make dev
```

### 5. Access

| URL                                            | Description               |
| ---------------------------------------------- | ------------------------- |
| `http://localhost:8080/health`                 | Health check              |
| `http://localhost:8080/admin`                  | Admin dashboard           |
| `http://localhost:8080/admin/running-analysis` | AI Running Analysis page  |
| `http://localhost:8080/docs/index.html`        | Swagger API docs          |
| `http://localhost:8080/api/v1/public/...`      | Public API (for frontend) |

---

## 🔑 Environment Variables

```env
# App
APP_NAME=portfolio-cms
APP_ENV=development
APP_PORT=8080
APP_URL=http://localhost:8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=portfolio_cms
DB_SSLMODE=disable
DB_TIMEZONE=Asia/Jakarta

# JWT
JWT_SECRET=change_me_in_production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# Admin
ADMIN_PATH=/admin
ADMIN_TITLE=Portfolio CMS
ADMIN_USERNAME=admin
ADMIN_PASSWORD=admin123

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000

# Gemini AI — untuk extract screenshot & analisis lari
GEMINI_API_KEY=AIzaSy...
GEMINI_MODEL=gemini-2.5-flash

# Google OAuth — untuk integrasi Google Calendar di admin dashboard
# Buat di: console.cloud.google.com → APIs & Services → Credentials → OAuth 2.0 Client ID
GOOGLE_CLIENT_ID=xxxxx.apps.googleusercontent.com
```

---

## 🔒 Authentication

```bash
# Register
POST /api/v1/auth/register
{ "name": "Admin", "email": "admin@example.com", "password": "password123" }

# Login
POST /api/v1/auth/login
{ "email": "admin@example.com", "password": "password123" }

# Gunakan access_token di header:
Authorization: Bearer <access_token>
```

---

## 🌐 Public API Endpoints

### Profile

```
GET  /api/v1/public/profile
```

### Projects

```
GET  /api/v1/public/projects?page=1&limit=10&featured=true
GET  /api/v1/public/projects/:slug
```

### Skills

```
GET  /api/v1/public/skills
GET  /api/v1/public/skills/:id
```

### Experiences & Educations

```
GET  /api/v1/public/experiences
GET  /api/v1/public/educations
```

### Tech Stack

```
GET  /api/v1/public/stack-categories?with=items
GET  /api/v1/public/stack-items?category_id=<uuid>
```

### Running Activities

```
GET  /api/v1/public/running-activities?page=1&limit=10
GET  /api/v1/public/running-activities/:id
```

### Contact

```
POST /api/v1/public/contact
```

### Response format

```json
{
  "success": true,
  "message": "Success",
  "data": {},
  "meta": { "page": 1, "limit": 10, "total": 25, "total_pages": 3 }
}
```

---

## 🛠 Admin API Endpoints (JWT required)

Prefix: `/api/v1/admin/api`  
Header: `Authorization: Bearer <access_token>`

### Projects

```
GET    /projects
POST   /projects
GET    /projects/:id
PUT    /projects/:id
DELETE /projects/:id
```

### Assets (polymorphic: photo/video/pdf/doc)

```
GET    /assets?owner_type=project&owner_id=<uuid>
POST   /assets
PUT    /assets/:id
DELETE /assets/:id
```

### Categories

```
GET/POST/PUT/DELETE  /project-categories
GET/POST/PUT/DELETE  /experience-categories
GET/POST/PUT/DELETE  /stack-categories
GET/POST/PUT/DELETE  /stack-items
```

### Skills / Experiences / Educations

```
GET/POST/PUT/DELETE  /skills
GET/POST/PUT/DELETE  /experiences
GET/POST/PUT/DELETE  /educations
```

### Profile & Contacts

```
GET/POST  /profile
GET       /contacts
PATCH     /contacts/:id/read
```

### Running Activities

```
GET    /running-activities
POST   /running-activities
POST   /running-activities/screenshot   ← upload gambar Huawei Health, Gemini extract otomatis
PUT    /running-activities/:id
DELETE /running-activities/:id
```

### 🧠 Running Analysis (AI)

```
POST   /running-analysis/generate       ← generate analisis dari semua data lari via Gemini
POST   /running-analysis/sync-calendar  ← sync jadwal latihan ke Google Calendar
```

---

## 🧠 Fitur AI Running Analysis

Halaman `/admin/running-analysis` menganalisis seluruh histori sesi lari menggunakan **Google Gemini Flash** berdasarkan prinsip ilmiah pelatihan lari.

### Cara Kerja

```
1. Admin klik "Generate Analisis"
2. Backend ambil semua sesi lari dari database
3. Build prompt: data historis + prinsip ilmiah (Banister, Maffetone, Jack Daniels)
4. Kirim ke Gemini Flash → dapat response JSON terstruktur
5. Tampilkan hasil di dashboard
6. Admin klik "Add to Google Calendar" → OAuth popup → sync jadwal ke kalender
```

### Hasil Analisis

| Komponen               | Deskripsi                                                            |
| ---------------------- | -------------------------------------------------------------------- |
| **Laporan Pelatih**    | Narasi Bahasa Indonesia seperti laporan coach pribadi                |
| **Kondisi Kebugaran**  | Level, fatigue score, aerobic base, trend, CTL/ATL/TSB               |
| **Zona Pace Personal** | Easy / Aerobic / Tempo / Threshold berdasarkan data aktual           |
| **Rencana Mingguan**   | 7 hari latihan dengan tipe sesi, target pace, HR zone, dan rationale |
| **Peringatan**         | Deteksi overtraining, injury risk, recovery needed                   |
| **Jadwal Kalender**    | Event siap sync ke Google Calendar dengan color coding               |

### Prinsip Ilmiah yang Digunakan

- **Banister Impulse-Response Model** — ATL, CTL, TSB untuk training load
- **Maffetone Method** — zona aerobic berbasis HR
- **Jack Daniels VDOT** — pace zones dari performa aktual
- **10% Rule** — progressive overload yang aman
- **Cadence Analysis** — deteksi overstriding dari data cadence

### Google Calendar Integration

Integrasi calendar berjalan **sepenuhnya di frontend** tanpa menyimpan token di server:

```
1. Klik "Add to Google Calendar"
2. Cek sessionStorage — ada token valid? → langsung sync
3. Tidak ada / expired → tampilkan OAuth popup Google
4. Dapat access_token → simpan ke sessionStorage (1 jam)
5. Insert semua event via Google Calendar API langsung dari browser
6. Tampilkan link ke masing-masing event yang berhasil dibuat
```

**Color coding event:**

- 🟢 Sage — Easy run
- 🟠 Tangerine — Tempo / Threshold
- 🔵 Blueberry — Long run

### Setup Google OAuth untuk Calendar

1. Buka [console.cloud.google.com](https://console.cloud.google.com)
2. **APIs & Services → Enable APIs** → aktifkan **Google Calendar API**
3. **Credentials → Create Credentials → OAuth 2.0 Client ID**
4. Application type: **Web application**
5. Authorized JavaScript origins: `http://localhost:8080` (tambah domain produksi juga)
6. Copy Client ID → set `GOOGLE_CLIENT_ID` di `.env`
7. **OAuth consent screen → Test users** → tambahkan email kamu

---

## 📸 Running Activity dari Screenshot

Upload screenshot Huawei Health → Gemini otomatis extract dan simpan ke database:

```bash
POST /api/v1/admin/api/running-activities/screenshot
Content-Type: multipart/form-data
Body: screenshot=<file.jpg>

# Response:
{
  "success": true,
  "data": {
    "distance_meters": 5200,
    "duration_sec": 2120,
    "avg_pace_sec": 405,
    "avg_heart_rate": 158,
    "avg_cadence": 172,
    ...
  }
}
```

---

## 🧩 Data Models

- **User** — admin/superadmin
- **Profile** — single row biodata (medsos, CV, status hire)
- **Project** — punya Category, Skills (m2m), StackItems (m2m), Assets (polymorphic)
- **Experience** — punya Category, Skills, StackItems, Assets
- **Skill** — many2many dengan Project & Experience
- **StackCategory → StackItem** — tech stack terorganisir
- **Asset** — polymorphic (owner_type + owner_id) untuk Project & Experience
- **RunningActivity** — sesi lari dengan 10+ metrik (distance, pace, HR, cadence, stride, dll)
- **Education**, **Contact** — mandiri

---

## 📦 Deploy ke Railway

1. Push ke GitHub
2. Railway → **New Project → Deploy from GitHub**
3. Tambahkan **PostgreSQL** plugin
4. Set environment variables (lihat bagian Environment Variables di atas)
5. Railway otomatis detect `Dockerfile` dan build

**Tambahan untuk produksi:**

- Set `GOOGLE_CLIENT_ID` di Railway environment
- Tambahkan domain produksi ke **Authorized JavaScript origins** di Google Cloud Console
- Tambahkan akun email ke **Test users** di OAuth consent screen (atau publish ke Production)

---

## 📝 Generate Swagger Docs

```bash
go install github.com/swaggo/swag/cmd/swag@latest
make swag
```

---

## 🧪 Testing

```bash
make test
```

```

```
