package service

import (
	"errors"
	"math"
	"portfolio-cms/internal/dto"
	"portfolio-cms/internal/model"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ─── Auth Service ──────────────────────────────────────────────────────────────

type AuthService interface {
	Register(name, email, password string) (*model.User, error)
	GetByEmail(email, password string) (*model.User, error)
	GetByID(id string) (*model.User, error)
}

type authService struct{ db *gorm.DB }

func NewAuthService(db *gorm.DB) AuthService { return &authService{db: db} }

func (s *authService) Register(name, email, password string) (*model.User, error) {
	var existing model.User
	if err := s.db.Where("email = ?", email).First(&existing).Error; err == nil {
		return nil, errors.New("email already registered")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{Name: name, Email: email, Password: string(hashed)}
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) GetByEmail(email, password string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}
	return &user, nil
}

func (s *authService) GetByID(id string) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ─── Project Service ───────────────────────────────────────────────────────────

type ProjectService interface {
	List(page, limit int, featuredOnly, publishedOnly bool) ([]model.Project, *dto.PaginationMeta, error)
	GetBySlug(slug string) (*model.Project, error)
	Create(req *dto.CreateProjectRequest) (*model.Project, error)
	Update(id string, req *dto.UpdateProjectRequest) (*model.Project, error)
	Delete(id string) error
}

type projectService struct{ db *gorm.DB }

func NewProjectService(db *gorm.DB) ProjectService { return &projectService{db: db} }

func (s *projectService) List(page, limit int, featuredOnly, publishedOnly bool) ([]model.Project, *dto.PaginationMeta, error) {
	var projects []model.Project
	var total int64
	offset := (page - 1) * limit

	query := s.db.Model(&model.Project{})
	if publishedOnly {
		query = query.Where("is_published = true")
	}
	if featuredOnly {
		query = query.Where("is_featured = true")
	}
	query.Count(&total)
	err := query.Order("sort_order asc, created_at desc").
		Offset(offset).Limit(limit).Find(&projects).Error

	meta := &dto.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}
	return projects, meta, err
}

func (s *projectService) GetBySlug(slug string) (*model.Project, error) {
	var p model.Project
	err := s.db.Where("slug = ? AND is_published = true", slug).First(&p).Error
	return &p, err
}

func (s *projectService) Create(req *dto.CreateProjectRequest) (*model.Project, error) {
	p := &model.Project{
		Title: req.Title, Slug: req.Slug, Description: req.Description,
		Content: req.Content, ThumbnailURL: req.ThumbnailURL, DemoURL: req.DemoURL,
		RepoURL: req.RepoURL, TechStack: req.TechStack, IsFeatured: req.IsFeatured,
		IsPublished: req.IsPublished, SortOrder: req.SortOrder,
	}
	return p, s.db.Create(p).Error
}

func (s *projectService) Update(id string, req *dto.UpdateProjectRequest) (*model.Project, error) {
	var p model.Project
	if err := s.db.First(&p, "id = ?", id).Error; err != nil {
		return nil, err
	}
	updates := map[string]interface{}{}
	if req.Title != nil        { updates["title"] = *req.Title }
	if req.Description != nil  { updates["description"] = *req.Description }
	if req.Content != nil      { updates["content"] = *req.Content }
	if req.ThumbnailURL != nil { updates["thumbnail_url"] = *req.ThumbnailURL }
	if req.DemoURL != nil      { updates["demo_url"] = *req.DemoURL }
	if req.RepoURL != nil      { updates["repo_url"] = *req.RepoURL }
	if req.TechStack != nil    { updates["tech_stack"] = *req.TechStack }
	if req.IsFeatured != nil   { updates["is_featured"] = *req.IsFeatured }
	if req.IsPublished != nil  { updates["is_published"] = *req.IsPublished }
	if req.SortOrder != nil    { updates["sort_order"] = *req.SortOrder }
	if err := s.db.Model(&p).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *projectService) Delete(id string) error {
	return s.db.Delete(&model.Project{}, "id = ?", id).Error
}

// ─── Skill Service ─────────────────────────────────────────────────────────────

type SkillService interface {
	List(publishedOnly bool) ([]model.Skill, error)
	Create(req *dto.CreateSkillRequest) (*model.Skill, error)
	Update(id string, req *dto.UpdateSkillRequest) (*model.Skill, error)
	Delete(id string) error
}

type skillService struct{ db *gorm.DB }

func NewSkillService(db *gorm.DB) SkillService { return &skillService{db: db} }

func (s *skillService) List(publishedOnly bool) ([]model.Skill, error) {
	var skills []model.Skill
	q := s.db.Order("sort_order asc, category asc")
	if publishedOnly {
		q = q.Where("is_published = true")
	}
	return skills, q.Find(&skills).Error
}

func (s *skillService) Create(req *dto.CreateSkillRequest) (*model.Skill, error) {
	sk := &model.Skill{
		Name: req.Name, Category: req.Category, IconURL: req.IconURL,
		Proficiency: req.Proficiency, SortOrder: req.SortOrder, IsPublished: req.IsPublished,
	}
	return sk, s.db.Create(sk).Error
}

func (s *skillService) Update(id string, req *dto.UpdateSkillRequest) (*model.Skill, error) {
	var sk model.Skill
	if err := s.db.First(&sk, "id = ?", id).Error; err != nil {
		return nil, err
	}
	updates := map[string]interface{}{}
	if req.Name != nil        { updates["name"] = *req.Name }
	if req.Category != nil    { updates["category"] = *req.Category }
	if req.IconURL != nil     { updates["icon_url"] = *req.IconURL }
	if req.Proficiency != nil { updates["proficiency"] = *req.Proficiency }
	if req.SortOrder != nil   { updates["sort_order"] = *req.SortOrder }
	if req.IsPublished != nil { updates["is_published"] = *req.IsPublished }
	if err := s.db.Model(&sk).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &sk, nil
}

func (s *skillService) Delete(id string) error {
	return s.db.Delete(&model.Skill{}, "id = ?", id).Error
}

// ─── Experience Service ────────────────────────────────────────────────────────

type ExperienceService interface {
	List(publishedOnly bool) ([]model.Experience, error)
	Create(req *dto.CreateExperienceRequest) (*model.Experience, error)
	Update(id string, req *dto.UpdateExperienceRequest) (*model.Experience, error)
	Delete(id string) error
}

type experienceService struct{ db *gorm.DB }

func NewExperienceService(db *gorm.DB) ExperienceService { return &experienceService{db: db} }

func (s *experienceService) List(publishedOnly bool) ([]model.Experience, error) {
	var items []model.Experience
	q := s.db.Order("sort_order asc, start_date desc")
	if publishedOnly {
		q = q.Where("is_published = true")
	}
	return items, q.Find(&items).Error
}

func (s *experienceService) Create(req *dto.CreateExperienceRequest) (*model.Experience, error) {
	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		return nil, errors.New("invalid start_date, use RFC3339 e.g. 2023-01-01T00:00:00Z")
	}
	var endDate *time.Time
	if req.EndDate != nil {
		t, err := time.Parse(time.RFC3339, *req.EndDate)
		if err == nil {
			endDate = &t
		}
	}
	item := &model.Experience{
		Company: req.Company, Position: req.Position, Description: req.Description,
		LogoURL: req.LogoURL, CompanyURL: req.CompanyURL, Location: req.Location,
		StartDate: startDate, EndDate: endDate,
		IsCurrent: req.IsCurrent, IsPublished: req.IsPublished, SortOrder: req.SortOrder,
	}
	return item, s.db.Create(item).Error
}

func (s *experienceService) Update(id string, req *dto.UpdateExperienceRequest) (*model.Experience, error) {
	var item model.Experience
	if err := s.db.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	updates := map[string]interface{}{}
	if req.Company != nil     { updates["company"] = *req.Company }
	if req.Position != nil    { updates["position"] = *req.Position }
	if req.Description != nil { updates["description"] = *req.Description }
	if req.LogoURL != nil     { updates["logo_url"] = *req.LogoURL }
	if req.IsCurrent != nil   { updates["is_current"] = *req.IsCurrent }
	if req.IsPublished != nil { updates["is_published"] = *req.IsPublished }
	if req.SortOrder != nil   { updates["sort_order"] = *req.SortOrder }
	if err := s.db.Model(&item).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *experienceService) Delete(id string) error {
	return s.db.Delete(&model.Experience{}, "id = ?", id).Error
}

// ─── Education Service ─────────────────────────────────────────────────────────

type EducationService interface {
	List(publishedOnly bool) ([]model.Education, error)
	Create(req *dto.CreateEducationRequest) (*model.Education, error)
	Update(id string, req *dto.UpdateEducationRequest) (*model.Education, error)
	Delete(id string) error
}

type educationService struct{ db *gorm.DB }

func NewEducationService(db *gorm.DB) EducationService { return &educationService{db: db} }

func (s *educationService) List(publishedOnly bool) ([]model.Education, error) {
	var items []model.Education
	q := s.db.Order("sort_order asc, start_date desc")
	if publishedOnly {
		q = q.Where("is_published = true")
	}
	return items, q.Find(&items).Error
}

func (s *educationService) Create(req *dto.CreateEducationRequest) (*model.Education, error) {
	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		return nil, errors.New("invalid start_date, use RFC3339")
	}
	item := &model.Education{
		Institution: req.Institution, Degree: req.Degree, FieldOfStudy: req.FieldOfStudy,
		Description: req.Description, LogoURL: req.LogoURL, StartDate: startDate,
		GPA: req.GPA, IsPublished: req.IsPublished, SortOrder: req.SortOrder,
	}
	return item, s.db.Create(item).Error
}

func (s *educationService) Update(id string, req *dto.UpdateEducationRequest) (*model.Education, error) {
	var item model.Education
	if err := s.db.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	updates := map[string]interface{}{}
	if req.Institution != nil  { updates["institution"] = *req.Institution }
	if req.Degree != nil       { updates["degree"] = *req.Degree }
	if req.FieldOfStudy != nil { updates["field_of_study"] = *req.FieldOfStudy }
	if req.IsPublished != nil  { updates["is_published"] = *req.IsPublished }
	if req.SortOrder != nil    { updates["sort_order"] = *req.SortOrder }
	if err := s.db.Model(&item).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *educationService) Delete(id string) error {
	return s.db.Delete(&model.Education{}, "id = ?", id).Error
}

// ─── Profile Service ───────────────────────────────────────────────────────────

type ProfileService interface {
	Get() (*model.Profile, error)
	Upsert(req *dto.UpsertProfileRequest) (*model.Profile, error)
}

type profileService struct{ db *gorm.DB }

func NewProfileService(db *gorm.DB) ProfileService { return &profileService{db: db} }

func (s *profileService) Get() (*model.Profile, error) {
	var p model.Profile
	err := s.db.Where("is_published = true").First(&p).Error
	return &p, err
}

func (s *profileService) Upsert(req *dto.UpsertProfileRequest) (*model.Profile, error) {
	var p model.Profile
	if err := s.db.First(&p).Error; err != nil {
		p = model.Profile{
			FullName: req.FullName, Title: req.Title, Bio: req.Bio,
			AvatarURL: req.AvatarURL, Email: req.Email, Phone: req.Phone,
			Location: req.Location, GithubURL: req.GithubURL, LinkedinURL: req.LinkedinURL,
			TwitterURL: req.TwitterURL, WebsiteURL: req.WebsiteURL, ResumeURL: req.ResumeURL,
			IsPublished: req.IsPublished,
		}
		return &p, s.db.Create(&p).Error
	}
	s.db.Model(&p).Updates(map[string]interface{}{
		"full_name": req.FullName, "title": req.Title, "bio": req.Bio,
		"avatar_url": req.AvatarURL, "email": req.Email, "phone": req.Phone,
		"location": req.Location, "github_url": req.GithubURL,
		"linkedin_url": req.LinkedinURL, "twitter_url": req.TwitterURL,
		"website_url": req.WebsiteURL, "resume_url": req.ResumeURL,
		"is_published": req.IsPublished,
	})
	return &p, nil
}

// ─── Contact Service ───────────────────────────────────────────────────────────

type ContactService interface {
	Save(req *dto.SendContactRequest) error
	List() ([]model.Contact, error)
	MarkRead(id string) error
}

type contactService struct{ db *gorm.DB }

func NewContactService(db *gorm.DB) ContactService { return &contactService{db: db} }

func (s *contactService) Save(req *dto.SendContactRequest) error {
	c := &model.Contact{
		Name: req.Name, Email: req.Email,
		Subject: req.Subject, Message: req.Message,
	}
	return s.db.Create(c).Error
}

func (s *contactService) List() ([]model.Contact, error) {
	var contacts []model.Contact
	return contacts, s.db.Order("created_at desc").Find(&contacts).Error
}

func (s *contactService) MarkRead(id string) error {
	return s.db.Model(&model.Contact{}).Where("id = ?", id).Update("is_read", true).Error
}
