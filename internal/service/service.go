package service

import (
	"errors"
	"math"
	"portfolio-cms/internal/dto"
	"portfolio-cms/internal/model"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ─── Pagination Meta ─────────────────────────────────────────────────────────

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

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
	List(page, limit int, featuredOnly, publishedOnly bool) ([]model.Project, *PaginationMeta, error)
	GetBySlug(slug string) (*model.Project, error)
	GetByID(id string) (*model.Project, error)
	Create(req *dto.CreateProjectRequest) (*model.Project, error)
	Update(id string, req *dto.UpdateProjectRequest) (*model.Project, error)
	Delete(id string) error
}

type projectService struct{ db *gorm.DB }

func NewProjectService(db *gorm.DB) ProjectService { return &projectService{db: db} }

func (s *projectService) preloadQuery() *gorm.DB {
	return s.db.
		Preload("StackItems.Category").
		Preload("Skills").
		Preload("Category"). // project category
		Preload("Assets", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order asc")
		})
}

func (s *projectService) List(page, limit int, featuredOnly, publishedOnly bool) ([]model.Project, *PaginationMeta, error) {
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

	meta := &PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}
	return projects, meta, err
}

func (s *projectService) GetBySlug(slug string) (*model.Project, error) {
	var p model.Project
	err := s.preloadQuery().Where("slug = ? AND is_published = true", slug).First(&p).Error
	return &p, err
}

func (s *projectService) GetByID(id string) (*model.Project, error) {
	var p model.Project
	err := s.preloadQuery().First(&p, "id = ?", id).Error
	return &p, err
}

func (s *projectService) Create(req *dto.CreateProjectRequest) (*model.Project, error) {
	p := &model.Project{
		Title:        req.Title,
		Slug:         req.Slug,
		Description:  req.Description,
		Content:      req.Content,
		ThumbnailURL: req.ThumbnailURL,
		DemoURL:      req.DemoURL,
		RepoURL:      req.RepoURL,
		DocURL:       req.DocURL,
		Problem:      req.Problem,
		Solution:     req.Solution,
		MyRole:       req.MyRole,
		Impact:       req.Impact,
		IsFeatured:   req.IsFeatured,
		IsPublished:  req.IsPublished,
		SortOrder:    req.SortOrder,
	}

	if req.CategoryID != nil {
		p.CategoryID = req.CategoryID
	}

	if req.StartDate != nil {
		t, err := time.Parse(time.RFC3339, *req.StartDate)
		if err != nil {
			return nil, errors.New("invalid start_date, use RFC3339")
		}
		p.StartDate = &t
	}
	if req.EndDate != nil {
		t, err := time.Parse(time.RFC3339, *req.EndDate)
		if err != nil {
			return nil, errors.New("invalid end_date, use RFC3339")
		}
		p.EndDate = &t
	}

	if err := s.db.Create(p).Error; err != nil {
		return nil, err
	}

	// attach stack items
	if len(req.StackItemIDs) > 0 {
		var items []model.StackItem
		if err := s.db.Find(&items, "id IN ?", req.StackItemIDs).Error; err != nil {
			return nil, err
		}
		if err := s.db.Model(p).Association("StackItems").Replace(items); err != nil {
			return nil, err
		}
	}

	// attach skills
	if len(req.SkillIDs) > 0 {
		var skills []model.Skill
		if err := s.db.Find(&skills, "id IN ?", req.SkillIDs).Error; err != nil {
			return nil, err
		}
		if err := s.db.Model(p).Association("Skills").Replace(skills); err != nil {
			return nil, err
		}
	}

	return s.GetByID(p.ID.String())
}

func (s *projectService) Update(id string, req *dto.UpdateProjectRequest) (*model.Project, error) {
	var p model.Project
	if err := s.db.First(&p, "id = ?", id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Title != nil        { updates["title"] = *req.Title }
	if req.Slug != nil         { updates["slug"] = *req.Slug }
	if req.Description != nil  { updates["description"] = *req.Description }
	if req.Content != nil      { updates["content"] = *req.Content }
	if req.ThumbnailURL != nil { updates["thumbnail_url"] = *req.ThumbnailURL }
	if req.DemoURL != nil      { updates["demo_url"] = *req.DemoURL }
	if req.RepoURL != nil      { updates["repo_url"] = *req.RepoURL }
	if req.DocURL != nil       { updates["doc_url"] = *req.DocURL }
	if req.Problem != nil      { updates["problem"] = *req.Problem }
	if req.Solution != nil     { updates["solution"] = *req.Solution }
	if req.MyRole != nil       { updates["my_role"] = *req.MyRole }
	if req.Impact != nil       { updates["impact"] = *req.Impact }
	if req.IsFeatured != nil   { updates["is_featured"] = *req.IsFeatured }
	if req.IsPublished != nil  { updates["is_published"] = *req.IsPublished }
	if req.SortOrder != nil    { updates["sort_order"] = *req.SortOrder }
	if req.CategoryID != nil   { updates["category_id"] = *req.CategoryID }

	if req.StartDate != nil {
		t, err := time.Parse(time.RFC3339, *req.StartDate)
		if err != nil {
			return nil, errors.New("invalid start_date, use RFC3339")
		}
		updates["start_date"] = t
	}
	if req.EndDate != nil {
		t, err := time.Parse(time.RFC3339, *req.EndDate)
		if err != nil {
			return nil, errors.New("invalid end_date, use RFC3339")
		}
		updates["end_date"] = t
	}

	if len(updates) > 0 {
		if err := s.db.Model(&p).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	// update stack items
	if req.StackItemIDs != nil {
		var items []model.StackItem
		if len(*req.StackItemIDs) > 0 {
			if err := s.db.Find(&items, "id IN ?", *req.StackItemIDs).Error; err != nil {
				return nil, err
			}
		}
		if err := s.db.Model(&p).Association("StackItems").Replace(items); err != nil {
			return nil, err
		}
	}

	// update skills
	if req.SkillIDs != nil {
		var skills []model.Skill
		if len(*req.SkillIDs) > 0 {
			if err := s.db.Find(&skills, "id IN ?", *req.SkillIDs).Error; err != nil {
				return nil, err
			}
		}
		if err := s.db.Model(&p).Association("Skills").Replace(skills); err != nil {
			return nil, err
		}
	}

	return s.GetByID(id)
}

func (s *projectService) Delete(id string) error {
	var p model.Project
	if err := s.db.First(&p, "id = ?", id).Error; err != nil {
		return err
	}
	_ = s.db.Model(&p).Association("StackItems").Clear()
	_ = s.db.Model(&p).Association("Skills").Clear()
	// Assets akan terhapus otomatis oleh CASCADE di database (sebaiknya atur constraint di migrasi)
	return s.db.Delete(&model.Project{}, "id = ?", id).Error
}

// ─── Asset Service (polymorphic: project / experience) ─────────────────────────

type AssetService interface {
	Create(ownerType, ownerID string, req *dto.CreateAssetRequest) (*model.Asset, error)
	List(ownerType, ownerID string) ([]model.Asset, error)
	Update(id string, req *dto.UpdateAssetRequest) (*model.Asset, error)
	Delete(id string) error
}

type assetService struct{ db *gorm.DB }

func NewAssetService(db *gorm.DB) AssetService { return &assetService{db: db} }

func (s *assetService) Create(ownerType string, ownerID string, req *dto.CreateAssetRequest) (*model.Asset, error) {
	ownerUUID, err := uuid.Parse(ownerID)
	if err != nil {
		return nil, errors.New("invalid owner_id")
	}
	if ownerType != "project" && ownerType != "experience" {
		return nil, errors.New("owner_type must be 'project' or 'experience'")
	}
	a := &model.Asset{
		OwnerType: ownerType,
		OwnerID:   ownerUUID,
		Type:      model.AssetType(req.Type),
		URL:       req.URL,
		Title:     req.Title,
		Caption:   req.Caption,
		SortOrder: req.SortOrder,
	}
	return a, s.db.Create(a).Error
}

func (s *assetService) List(ownerType string, ownerID string) ([]model.Asset, error) {
	var assets []model.Asset
	ownerUUID, err := uuid.Parse(ownerID)
	if err != nil {
		return nil, errors.New("invalid owner_id")
	}
	err = s.db.Where("owner_type = ? AND owner_id = ?", ownerType, ownerUUID).
		Order("sort_order asc").Find(&assets).Error
	return assets, err
}

func (s *assetService) Update(id string, req *dto.UpdateAssetRequest) (*model.Asset, error) {
	var a model.Asset
	if err := s.db.First(&a, "id = ?", id).Error; err != nil {
		return nil, err
	}
	updates := map[string]interface{}{}
	if req.Type != nil      { updates["type"] = *req.Type }
	if req.URL != nil       { updates["url"] = *req.URL }
	if req.Title != nil     { updates["title"] = *req.Title }
	if req.Caption != nil   { updates["caption"] = *req.Caption }
	if req.SortOrder != nil { updates["sort_order"] = *req.SortOrder }
	if err := s.db.Model(&a).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *assetService) Delete(id string) error {
	return s.db.Delete(&model.Asset{}, "id = ?", id).Error
}

// ─── Project Category Service ─────────────────────────────────────────────────

type ProjectCategoryService interface {
	List() ([]model.ProjectCategory, error)
	Create(req *dto.CreateCategoryRequest) (*model.ProjectCategory, error)
	Update(id string, req *dto.UpdateCategoryRequest) (*model.ProjectCategory, error)
	Delete(id string) error
}

type projectCategoryService struct{ db *gorm.DB }

func NewProjectCategoryService(db *gorm.DB) ProjectCategoryService {
	return &projectCategoryService{db: db}
}

func (s *projectCategoryService) List() ([]model.ProjectCategory, error) {
	var cats []model.ProjectCategory
	err := s.db.Order("sort_order asc, name asc").Find(&cats).Error
	return cats, err
}

func (s *projectCategoryService) Create(req *dto.CreateCategoryRequest) (*model.ProjectCategory, error) {
	// cek slug
	var existing model.ProjectCategory
	if err := s.db.Where("slug = ?", req.Slug).First(&existing).Error; err == nil {
		return nil, errors.New("slug already exists")
	}
	cat := &model.ProjectCategory{
		Name:        req.Name,
		Slug:        req.Slug,
		SortOrder:   req.SortOrder,
		IsPublished: req.IsPublished,
	}
	return cat, s.db.Create(cat).Error
}

func (s *projectCategoryService) Update(id string, req *dto.UpdateCategoryRequest) (*model.ProjectCategory, error) {
	var cat model.ProjectCategory
	if err := s.db.First(&cat, "id = ?", id).Error; err != nil {
		return nil, err
	}
	updates := map[string]interface{}{}
	if req.Name != nil        { updates["name"] = *req.Name }
	if req.Slug != nil        { updates["slug"] = *req.Slug }
	if req.SortOrder != nil   { updates["sort_order"] = *req.SortOrder }
	if req.IsPublished != nil { updates["is_published"] = *req.IsPublished }
	if err := s.db.Model(&cat).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

func (s *projectCategoryService) Delete(id string) error {
	return s.db.Delete(&model.ProjectCategory{}, "id = ?", id).Error
}

// ─── Experience Category Service ──────────────────────────────────────────────

type ExperienceCategoryService interface {
	List() ([]model.ExperienceCategory, error)
	Create(req *dto.CreateCategoryRequest) (*model.ExperienceCategory, error)
	Update(id string, req *dto.UpdateCategoryRequest) (*model.ExperienceCategory, error)
	Delete(id string) error
}

type experienceCategoryService struct{ db *gorm.DB }

func NewExperienceCategoryService(db *gorm.DB) ExperienceCategoryService {
	return &experienceCategoryService{db: db}
}

func (s *experienceCategoryService) List() ([]model.ExperienceCategory, error) {
	var cats []model.ExperienceCategory
	err := s.db.Order("sort_order asc, name asc").Find(&cats).Error
	return cats, err
}

func (s *experienceCategoryService) Create(req *dto.CreateCategoryRequest) (*model.ExperienceCategory, error) {
	var existing model.ExperienceCategory
	if err := s.db.Where("slug = ?", req.Slug).First(&existing).Error; err == nil {
		return nil, errors.New("slug already exists")
	}
	cat := &model.ExperienceCategory{
		Name:        req.Name,
		Slug:        req.Slug,
		SortOrder:   req.SortOrder,
		IsPublished: req.IsPublished,
	}
	return cat, s.db.Create(cat).Error
}

func (s *experienceCategoryService) Update(id string, req *dto.UpdateCategoryRequest) (*model.ExperienceCategory, error) {
	var cat model.ExperienceCategory
	if err := s.db.First(&cat, "id = ?", id).Error; err != nil {
		return nil, err
	}
	updates := map[string]interface{}{}
	if req.Name != nil        { updates["name"] = *req.Name }
	if req.Slug != nil        { updates["slug"] = *req.Slug }
	if req.SortOrder != nil   { updates["sort_order"] = *req.SortOrder }
	if req.IsPublished != nil { updates["is_published"] = *req.IsPublished }
	if err := s.db.Model(&cat).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

func (s *experienceCategoryService) Delete(id string) error {
	return s.db.Delete(&model.ExperienceCategory{}, "id = ?", id).Error
}

// ─── Stack Category Service ────────────────────────────────────────────────────

type StackCategoryService interface {
	List(withItems bool) ([]model.StackCategory, error)
	GetByID(id string) (*model.StackCategory, error)
	Create(req *dto.CreateStackCategoryRequest) (*model.StackCategory, error)
	Update(id string, req *dto.UpdateStackCategoryRequest) (*model.StackCategory, error)
	Delete(id string) error
}

type stackCategoryService struct{ db *gorm.DB }

func NewStackCategoryService(db *gorm.DB) StackCategoryService {
	return &stackCategoryService{db: db}
}

func (s *stackCategoryService) List(withItems bool) ([]model.StackCategory, error) {
	var categories []model.StackCategory
	q := s.db.Order("sort_order asc, name asc")
	if withItems {
		q = q.Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_published = true").Order("sort_order asc, name asc")
		})
	}
	return categories, q.Find(&categories).Error
}

func (s *stackCategoryService) GetByID(id string) (*model.StackCategory, error) {
	var cat model.StackCategory
	err := s.db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order asc, name asc")
	}).First(&cat, "id = ?", id).Error
	return &cat, err
}

func (s *stackCategoryService) Create(req *dto.CreateStackCategoryRequest) (*model.StackCategory, error) {
	var existing model.StackCategory
	if err := s.db.Where("slug = ?", req.Slug).First(&existing).Error; err == nil {
		return nil, errors.New("slug already exists")
	}
	cat := &model.StackCategory{
		Name:      req.Name,
		Slug:      req.Slug,
		IconURL:   req.IconURL,
		Color:     req.Color,
		SortOrder: req.SortOrder,
	}
	return cat, s.db.Create(cat).Error
}

func (s *stackCategoryService) Update(id string, req *dto.UpdateStackCategoryRequest) (*model.StackCategory, error) {
	var cat model.StackCategory
	if err := s.db.First(&cat, "id = ?", id).Error; err != nil {
		return nil, err
	}
	updates := map[string]interface{}{}
	if req.Name != nil      { updates["name"] = *req.Name }
	if req.Slug != nil      { updates["slug"] = *req.Slug }
	if req.IconURL != nil   { updates["icon_url"] = *req.IconURL }
	if req.Color != nil     { updates["color"] = *req.Color }
	if req.SortOrder != nil { updates["sort_order"] = *req.SortOrder }
	if err := s.db.Model(&cat).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

func (s *stackCategoryService) Delete(id string) error {
	return s.db.Delete(&model.StackCategory{}, "id = ?", id).Error
}

// ─── Stack Item Service ────────────────────────────────────────────────────────

type StackItemService interface {
	List(publishedOnly bool, categoryID string) ([]model.StackItem, error)
	GetByID(id string) (*model.StackItem, error)
	Create(req *dto.CreateStackItemRequest) (*model.StackItem, error)
	Update(id string, req *dto.UpdateStackItemRequest) (*model.StackItem, error)
	Delete(id string) error
}

type stackItemService struct{ db *gorm.DB }

func NewStackItemService(db *gorm.DB) StackItemService {
	return &stackItemService{db: db}
}

func (s *stackItemService) List(publishedOnly bool, categoryID string) ([]model.StackItem, error) {
	var items []model.StackItem
	q := s.db.Preload("Category").Order("sort_order asc, name asc")
	if publishedOnly {
		q = q.Where("is_published = true")
	}
	if categoryID != "" {
		catUUID, err := uuid.Parse(categoryID)
		if err != nil {
			return nil, errors.New("invalid category_id")
		}
		q = q.Where("category_id = ?", catUUID)
	}
	return items, q.Find(&items).Error
}

func (s *stackItemService) GetByID(id string) (*model.StackItem, error) {
	var item model.StackItem
	err := s.db.Preload("Category").First(&item, "id = ?", id).Error
	return &item, err
}

func (s *stackItemService) Create(req *dto.CreateStackItemRequest) (*model.StackItem, error) {
	var cat model.StackCategory
	if err := s.db.First(&cat, "id = ?", req.CategoryID).Error; err != nil {
		return nil, errors.New("category not found")
	}
	var existing model.StackItem
	if err := s.db.Where("slug = ?", req.Slug).First(&existing).Error; err == nil {
		return nil, errors.New("slug already exists")
	}
	item := &model.StackItem{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Slug:        req.Slug,
		IconURL:     req.IconURL,
		SortOrder:   req.SortOrder,
		IsPublished: req.IsPublished,
	}
	if err := s.db.Create(item).Error; err != nil {
		return nil, err
	}
	return s.GetByID(item.ID.String())
}

func (s *stackItemService) Update(id string, req *dto.UpdateStackItemRequest) (*model.StackItem, error) {
	var item model.StackItem
	if err := s.db.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	updates := map[string]interface{}{}
	if req.CategoryID != nil  { updates["category_id"] = *req.CategoryID }
	if req.Name != nil        { updates["name"] = *req.Name }
	if req.Slug != nil        { updates["slug"] = *req.Slug }
	if req.IconURL != nil     { updates["icon_url"] = *req.IconURL }
	if req.SortOrder != nil   { updates["sort_order"] = *req.SortOrder }
	if req.IsPublished != nil { updates["is_published"] = *req.IsPublished }
	if err := s.db.Model(&item).Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.GetByID(id)
}

func (s *stackItemService) Delete(id string) error {
	return s.db.Delete(&model.StackItem{}, "id = ?", id).Error
}

// ─── Skill Service ─────────────────────────────────────────────────────────────

type SkillDetailResponse struct {
	model.Skill
	TotalProjects        int             `json:"total_projects"`
	TotalItemsInCategory int             `json:"total_items_in_category"`
	Projects             []model.Project `json:"projects,omitempty"`
}

type SkillService interface {
	List(publishedOnly bool) ([]model.Skill, error)
	GetDetail(id string) (*SkillDetailResponse, error)
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

func (s *skillService) GetDetail(id string) (*SkillDetailResponse, error) {
	var skill model.Skill
	if err := s.db.Preload("Projects").First(&skill, "id = ?", id).Error; err != nil {
		return nil, err
	}
	totalProjects := len(skill.Projects)
	var totalInCategory int64
	s.db.Model(&model.Skill{}).Where("category = ? AND id != ?", skill.Category, id).Count(&totalInCategory)
	return &SkillDetailResponse{
		Skill:                skill,
		TotalProjects:        totalProjects,
		TotalItemsInCategory: int(totalInCategory),
		Projects:             skill.Projects,
	}, nil
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
	GetByID(id string) (*model.Experience, error)
	Create(req *dto.CreateExperienceRequest) (*model.Experience, error)
	Update(id string, req *dto.UpdateExperienceRequest) (*model.Experience, error)
	Delete(id string) error
}

type experienceService struct{ db *gorm.DB }

func NewExperienceService(db *gorm.DB) ExperienceService { return &experienceService{db: db} }

func (s *experienceService) preloadQuery() *gorm.DB {
	return s.db.
		Preload("StackItems").
		Preload("Skills").
		Preload("Category").
		Preload("Assets", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order asc")
		})
}

func (s *experienceService) List(publishedOnly bool) ([]model.Experience, error) {
	var items []model.Experience
	q := s.preloadQuery().Order("sort_order asc, start_date desc")
	if publishedOnly {
		q = q.Where("is_published = true")
	}
	return items, q.Find(&items).Error
}

func (s *experienceService) GetByID(id string) (*model.Experience, error) {
	var item model.Experience
	err := s.preloadQuery().First(&item, "id = ?", id).Error
	return &item, err
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
		Company:     req.Company,
		Position:    req.Position,
		Description: req.Description,
		LogoURL:     req.LogoURL,
		CompanyURL:  req.CompanyURL,
		Location:    req.Location,
		StartDate:   startDate,
		EndDate:     endDate,
		IsCurrent:   req.IsCurrent,
		IsPublished: req.IsPublished,
		SortOrder:   req.SortOrder,
	}

	if req.CategoryID != nil {
		item.CategoryID = req.CategoryID
	}

	if err := s.db.Create(item).Error; err != nil {
		return nil, err
	}

	if len(req.StackItemIDs) > 0 {
		var stacks []model.StackItem
		if err := s.db.Find(&stacks, "id IN ?", req.StackItemIDs).Error; err != nil {
			return nil, err
		}
		if err := s.db.Model(item).Association("StackItems").Replace(stacks); err != nil {
			return nil, err
		}
	}

	if len(req.SkillIDs) > 0 {
		var skills []model.Skill
		if err := s.db.Find(&skills, "id IN ?", req.SkillIDs).Error; err != nil {
			return nil, err
		}
		if err := s.db.Model(item).Association("Skills").Replace(skills); err != nil {
			return nil, err
		}
	}

	return s.GetByID(item.ID.String())
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
	if req.CompanyURL != nil  { updates["company_url"] = *req.CompanyURL }
	if req.Location != nil    { updates["location"] = *req.Location }
	if req.IsCurrent != nil   { updates["is_current"] = *req.IsCurrent }
	if req.IsPublished != nil { updates["is_published"] = *req.IsPublished }
	if req.SortOrder != nil   { updates["sort_order"] = *req.SortOrder }
	if req.CategoryID != nil  { updates["category_id"] = *req.CategoryID }

	if err := s.db.Model(&item).Updates(updates).Error; err != nil {
		return nil, err
	}

	if req.StackItemIDs != nil {
		var stacks []model.StackItem
		if len(*req.StackItemIDs) > 0 {
			if err := s.db.Find(&stacks, "id IN ?", *req.StackItemIDs).Error; err != nil {
				return nil, err
			}
		}
		if err := s.db.Model(&item).Association("StackItems").Replace(stacks); err != nil {
			return nil, err
		}
	}

	if req.SkillIDs != nil {
		var skills []model.Skill
		if len(*req.SkillIDs) > 0 {
			if err := s.db.Find(&skills, "id IN ?", *req.SkillIDs).Error; err != nil {
				return nil, err
			}
		}
		if err := s.db.Model(&item).Association("Skills").Replace(skills); err != nil {
			return nil, err
		}
	}

	return s.GetByID(id)
}

func (s *experienceService) Delete(id string) error {
	var item model.Experience
	if err := s.db.First(&item, "id = ?", id).Error; err != nil {
		return err
	}
	_ = s.db.Model(&item).Association("StackItems").Clear()
	_ = s.db.Model(&item).Association("Skills").Clear()
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
			FullName:         req.FullName,
			Title:            req.Title,
			About:            req.About,
			Bio:              req.Bio,
			AvatarURL:        req.AvatarURL,
			Location:         req.Location,
			Email:            req.Email,
			Phone:            req.Phone,
			GithubURL:        req.GithubURL,
			LinkedinURL:      req.LinkedinURL,
			TwitterURL:       req.TwitterURL,
			InstagramURL:     req.InstagramURL,
			TiktokURL:        req.TiktokURL,
			StravaURL:        req.StravaURL,
			WebsiteURL:       req.WebsiteURL,
			ResumeURL:        req.ResumeURL,
			YearsExperience:  req.YearsExperience,
			AvailableForHire: req.AvailableForHire,
			IsPublished:      req.IsPublished,
		}
		return &p, s.db.Create(&p).Error
	}
	s.db.Model(&p).Updates(map[string]interface{}{
		"full_name":          req.FullName,
		"title":              req.Title,
		"about":              req.About,
		"bio":                req.Bio,
		"avatar_url":         req.AvatarURL,
		"location":           req.Location,
		"email":              req.Email,
		"phone":              req.Phone,
		"github_url":         req.GithubURL,
		"linkedin_url":       req.LinkedinURL,
		"twitter_url":        req.TwitterURL,
		"instagram_url":      req.InstagramURL,
		"tiktok_url":         req.TiktokURL,
		"strava_url":         req.StravaURL,
		"website_url":        req.WebsiteURL,
		"resume_url":         req.ResumeURL,
		"years_experience":   req.YearsExperience,
		"available_for_hire": req.AvailableForHire,
		"is_published":       req.IsPublished,
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

// ─── Running Activity Service ──────────────────────────────────────────────────

type RunningActivityService interface {
	List(publishedOnly bool, page, limit int) ([]model.RunningActivity, *PaginationMeta, error)
	GetByID(id string) (*model.RunningActivity, error)
	Create(req *dto.CreateRunningActivityRequest) (*model.RunningActivity, error)
	Update(id string, req *dto.UpdateRunningActivityRequest) (*model.RunningActivity, error)
	Delete(id string) error
}

type runningActivityService struct{ db *gorm.DB }

func NewRunningActivityService(db *gorm.DB) RunningActivityService {
	return &runningActivityService{db: db}
}

func (s *runningActivityService) List(publishedOnly bool, page, limit int) ([]model.RunningActivity, *PaginationMeta, error) {
	var activities []model.RunningActivity
	var total int64
	offset := (page - 1) * limit

	q := s.db.Model(&model.RunningActivity{})
	if publishedOnly {
		q = q.Where("is_published = true")
	}
	q.Count(&total)
	err := q.Order("date desc").Offset(offset).Limit(limit).Find(&activities).Error

	meta := &PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}
	return activities, meta, err
}

func (s *runningActivityService) GetByID(id string) (*model.RunningActivity, error) {
	var a model.RunningActivity
	err := s.db.First(&a, "id = ?", id).Error
	return &a, err
}

func (s *runningActivityService) Create(req *dto.CreateRunningActivityRequest) (*model.RunningActivity, error) {
	a := &model.RunningActivity{
		Title:          req.Title,
		Notes:          req.Notes,
		MapImageURL:    req.MapImageURL,
		DurationSec:    req.DurationSec,
		DistanceMeters: req.DistanceMeters,
		TotalCalories:  req.TotalCalories,
		ActiveCalories: req.ActiveCalories,
		AvgPaceSec:     req.AvgPaceSec,
		AvgSpeedKph:    req.AvgSpeedKph,
		AvgCadence:     req.AvgCadence,
		AvgStride:      req.AvgStride,
		Steps:          req.Steps,
		AvgHeartRate:   req.AvgHeartRate,
		IsPublished:    req.IsPublished,
	}
	if req.Date != nil {
		t, err := time.Parse(time.RFC3339, *req.Date)
		if err != nil {
			return nil, errors.New("invalid date format, use RFC3339")
		}
		a.Date = t
	}
	return a, s.db.Create(a).Error
}

func (s *runningActivityService) Update(id string, req *dto.UpdateRunningActivityRequest) (*model.RunningActivity, error) {
	var a model.RunningActivity
	if err := s.db.First(&a, "id = ?", id).Error; err != nil {
		return nil, err
	}
	updates := map[string]interface{}{}
	if req.Title != nil          { updates["title"] = *req.Title }
	if req.Notes != nil          { updates["notes"] = *req.Notes }
	if req.MapImageURL != nil    { updates["map_image_url"] = *req.MapImageURL }
	if req.DurationSec != nil    { updates["duration_sec"] = *req.DurationSec }
	if req.DistanceMeters != nil { updates["distance_meters"] = *req.DistanceMeters }
	if req.TotalCalories != nil  { updates["total_calories"] = *req.TotalCalories }
	if req.ActiveCalories != nil { updates["active_calories"] = *req.ActiveCalories }
	if req.AvgPaceSec != nil     { updates["avg_pace_sec"] = *req.AvgPaceSec }
	if req.AvgSpeedKph != nil    { updates["avg_speed_kph"] = *req.AvgSpeedKph }
	if req.AvgCadence != nil     { updates["avg_cadence"] = *req.AvgCadence }
	if req.AvgStride != nil      { updates["avg_stride"] = *req.AvgStride }
	if req.Steps != nil          { updates["steps"] = *req.Steps }
	if req.AvgHeartRate != nil   { updates["avg_heart_rate"] = *req.AvgHeartRate }
	if req.IsPublished != nil    { updates["is_published"] = *req.IsPublished }
	if req.Date != nil {
		t, err := time.Parse(time.RFC3339, *req.Date)
		if err == nil {
			updates["date"] = t
		}
	}
	if err := s.db.Model(&a).Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.GetByID(id)
}

func (s *runningActivityService) Delete(id string) error {
	return s.db.Delete(&model.RunningActivity{}, "id = ?", id).Error
}