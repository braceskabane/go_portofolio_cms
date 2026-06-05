package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─── Base Model ────────────────────────────────────────────────────────────────

type Base struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// ─── User (Admin) ───────────────────────────────────────────────────────────────

type User struct {
	Base
	Name     string `gorm:"size:100;not null" json:"name"`
	Email    string `gorm:"size:150;uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Role     string `gorm:"size:20;default:'admin'" json:"role"` // admin | superadmin
}

// ─── Profile (single row - owner info) ─────────────────────────────────────────

type Profile struct {
	Base
	FullName    string `gorm:"size:150;not null" json:"full_name"`
	Title       string `gorm:"size:200" json:"title"` // e.g. "Full Stack Developer"
	Bio         string `gorm:"type:text" json:"bio"`
	AvatarURL   string `gorm:"size:500" json:"avatar_url"`
	Email       string `gorm:"size:150" json:"email"`
	Phone       string `gorm:"size:20" json:"phone"`
	Location    string `gorm:"size:100" json:"location"`
	GithubURL   string `gorm:"size:300" json:"github_url"`
	LinkedinURL string `gorm:"size:300" json:"linkedin_url"`
	TwitterURL  string `gorm:"size:300" json:"twitter_url"`
	WebsiteURL  string `gorm:"size:300" json:"website_url"`
	ResumeURL   string `gorm:"size:500" json:"resume_url"`
	IsPublished bool   `gorm:"default:true" json:"is_published"`
}

// ─── Project ───────────────────────────────────────────────────────────────────

type Project struct {
	Base
	Title       string    `gorm:"size:200;not null" json:"title"`
	Slug        string    `gorm:"size:250;uniqueIndex;not null" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	Content     string    `gorm:"type:text" json:"content"` // rich text / markdown
	ThumbnailURL string   `gorm:"size:500" json:"thumbnail_url"`
	DemoURL     string    `gorm:"size:300" json:"demo_url"`
	RepoURL     string    `gorm:"size:300" json:"repo_url"`
	TechStack   string    `gorm:"type:text" json:"tech_stack"` // comma-separated or JSON string
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	IsFeatured  bool      `gorm:"default:false" json:"is_featured"`
	IsPublished bool      `gorm:"default:true" json:"is_published"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
}

// ─── Skill ─────────────────────────────────────────────────────────────────────

type Skill struct {
	Base
	Name        string `gorm:"size:100;not null" json:"name"`
	Category    string `gorm:"size:100" json:"category"` // e.g. Frontend, Backend, DevOps
	IconURL     string `gorm:"size:300" json:"icon_url"`
	Proficiency int    `gorm:"default:0" json:"proficiency"` // 0-100
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	IsPublished bool   `gorm:"default:true" json:"is_published"`
}

// ─── Experience ────────────────────────────────────────────────────────────────

type Experience struct {
	Base
	Company     string     `gorm:"size:200;not null" json:"company"`
	Position    string     `gorm:"size:200;not null" json:"position"`
	Description string     `gorm:"type:text" json:"description"`
	LogoURL     string     `gorm:"size:300" json:"logo_url"`
	CompanyURL  string     `gorm:"size:300" json:"company_url"`
	Location    string     `gorm:"size:100" json:"location"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"` // nil = present
	IsCurrent   bool       `gorm:"default:false" json:"is_current"`
	IsPublished bool       `gorm:"default:true" json:"is_published"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
}

// ─── Education ─────────────────────────────────────────────────────────────────

type Education struct {
	Base
	Institution string     `gorm:"size:200;not null" json:"institution"`
	Degree      string     `gorm:"size:200" json:"degree"`
	FieldOfStudy string    `gorm:"size:200" json:"field_of_study"`
	Description string     `gorm:"type:text" json:"description"`
	LogoURL     string     `gorm:"size:300" json:"logo_url"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	GPA         string     `gorm:"size:10" json:"gpa"`
	IsPublished bool       `gorm:"default:true" json:"is_published"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
}

// ─── Contact (form submissions from public) ────────────────────────────────────

type Contact struct {
	Base
	Name    string `gorm:"size:100;not null" json:"name"`
	Email   string `gorm:"size:150;not null" json:"email"`
	Subject string `gorm:"size:250" json:"subject"`
	Message string `gorm:"type:text;not null" json:"message"`
	IsRead  bool   `gorm:"default:false" json:"is_read"`
}
