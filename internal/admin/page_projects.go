package admin

import (
	"fmt"
	"portfolio-cms/internal/model"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ProjectsPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var projects []model.Project
		db.Order("sort_order asc, created_at desc").Find(&projects)

		rows := ""
		for _, p := range projects {
			rows += fmt.Sprintf(`
<tr class="border-t border-gray-100 hover:bg-gray-50 transition">
  <td class="px-4 py-3 text-sm font-medium text-gray-800">%s</td>
  <td class="px-4 py-3 text-sm text-gray-500 font-mono">%s</td>
  <td class="px-4 py-3 text-sm">%s</td>
  <td class="px-4 py-3 text-sm">%s</td>
  <td class="px-4 py-3 text-sm text-gray-400">%d</td>
  <td class="px-4 py-3 text-sm text-gray-400">%s</td>
  <td class="px-4 py-3 text-sm">
    <div class="flex items-center gap-3">
      <a href="/admin/projects/%s/edit" class="text-indigo-600 hover:text-indigo-800 font-medium transition">Edit</a>
      %s
    </div>
  </td>
</tr>`,
				escapeHTML(p.Title),
				escapeHTML(p.Slug),
				yesNo(p.IsFeatured),
				yesNo(p.IsPublished),
				p.SortOrder,
				p.CreatedAt.Format("02 Jan 2006"),
				p.ID.String(),
				deleteBtn("/admin/projects/"+p.ID.String(), "Delete this project?"),
			)
		}

		table := tableWrapper(fmt.Sprintf(`
<table class="w-full text-left">
  <thead class="bg-gray-50 border-b border-gray-200">
    <tr>
      <th class="px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">Title</th>
      <th class="px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">Slug</th>
      <th class="px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">Featured</th>
      <th class="px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">Published</th>
      <th class="px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">Order</th>
      <th class="px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">Created</th>
      <th class="px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">Actions</th>
    </tr>
  </thead>
  <tbody>%s</tbody>
</table>`, rows))

		content := fmt.Sprintf(`
<div class="flex items-center justify-between mb-6">
  <div>
    <h3 class="text-lg font-semibold text-gray-800">All Projects</h3>
    <p class="text-sm text-gray-500">%d projects total</p>
  </div>
  <a href="/admin/projects/new"
   style="background:#1e1d1e;color:#f5f4f5;text-decoration:none;font-size:13px;font-weight:500;padding:8px 16px;border-radius:6px;display:inline-flex;align-items:center;gap:6px;transition:background 0.12s;"
   onmouseover="this.style.background='#30353b'"
   onmouseout="this.style.background='#1e1d1e'">
   + New Project
</a>
</div>
%s`, len(projects), table)

		flash := c.Cookies("flash")
		c.ClearCookie("flash")
		return c.Type("html").SendString(layout("Projects", "Projects", content, flash))
	}
}

func ProjectFormPage(db *gorm.DB, _ string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        id := c.Params("id")
        var p model.Project
        isEdit := id != "" && id != "new"

        if isEdit {
            if err := db.First(&p, "id = ?", id).Error; err != nil {
                return c.Redirect("/admin/projects")
            }
        }

        title := "New Project"
        action := "/admin/projects"
        saveLabel := "Create Project"
        if isEdit {
            title = "Edit Project"
            action = "/admin/projects/" + id
            saveLabel = "Update Project"
        }

        toggleRow := fmt.Sprintf(
            `<div style="display:flex; gap:24px; margin-bottom:16px;">
                <div style="flex:1;">%s</div>
                <div style="flex:1;">%s</div>
            </div>`,
            toggleField("Featured", "is_featured", p.IsFeatured),
            toggleField("Published", "is_published", p.IsPublished),
        )

        formContent := fmt.Sprintf(`
            <div style="margin-top: 12px;">
                %s
                <form method="POST" action="%s" style="display:flex; flex-direction:column; gap:12px;">
                    %s
                    %s
                    %s
                    %s
                    %s
                    %s
                    %s
                    %s
                    %s
                    %s
                    <div style="display:flex; align-items:center; gap:12px; margin-top:8px;">
                        %s
                        <a href="/admin/projects" style="font-size:13px; color:#6b6a6b; text-decoration:none; padding:8px 16px; border-radius:6px; border:1px solid #e5e3e4; transition:background 0.12s;" onmouseover="this.style.background='#f5f4f5'" onmouseout="this.style.background='transparent'">Cancel</a>
                    </div>
                </form>
            </div>`,
            card(`<p style="font-size:13px; color:#4d4c4d;">Fill in the details below. <strong>Slug</strong> must be unique (e.g. <code>my-project</code>).</p>`),
            action,
            inputField("Title *", "title", escapeHTML(p.Title), "text", true),
            inputField("Slug *", "slug", escapeHTML(p.Slug), "text", true),
            textareaField("Description", "description", escapeHTML(p.Description), 3),
            textareaField("Content (Markdown/HTML)", "content", escapeHTML(p.Content), 8),
            inputField("Thumbnail URL", "thumbnail_url", escapeHTML(p.ThumbnailURL), "url", false),
            inputField("Demo URL", "demo_url", escapeHTML(p.DemoURL), "url", false),
            inputField("Repo URL", "repo_url", escapeHTML(p.RepoURL), "url", false),
            inputField("Tech Stack (comma-separated)", "tech_stack", escapeHTML(p.TechStack), "text", false),
            inputField("Sort Order", "sort_order", fmt.Sprintf("%d", p.SortOrder), "number", false),
            toggleRow,
            btnPrimary(saveLabel),
        )

        return c.Type("html").SendString(layout(title, "Projects", formContent))
    }
}

func CreateProjectHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		p := model.Project{
			Title:        c.FormValue("title"),
			Slug:         c.FormValue("slug"),
			Description:  c.FormValue("description"),
			Content:      c.FormValue("content"),
			ThumbnailURL: c.FormValue("thumbnail_url"),
			DemoURL:      c.FormValue("demo_url"),
			RepoURL:      c.FormValue("repo_url"),
			TechStack:    c.FormValue("tech_stack"),
			IsFeatured:   c.FormValue("is_featured") == "true",
			IsPublished:  c.FormValue("is_published") == "true",
		}
		if err := db.Create(&p).Error; err != nil {
			setFlash(c, "ERR:Failed to create project: "+err.Error())
		} else {
			setFlash(c, "Project created successfully!")
		}
		return c.Redirect("/admin/projects")
	}
}

func UpdateProjectHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		updates := map[string]interface{}{
			"title":         c.FormValue("title"),
			"slug":          c.FormValue("slug"),
			"description":   c.FormValue("description"),
			"content":       c.FormValue("content"),
			"thumbnail_url": c.FormValue("thumbnail_url"),
			"demo_url":      c.FormValue("demo_url"),
			"repo_url":      c.FormValue("repo_url"),
			"tech_stack":    c.FormValue("tech_stack"),
			"is_featured":   c.FormValue("is_featured") == "true",
			"is_published":  c.FormValue("is_published") == "true",
			"updated_at":    time.Now(),
		}
		if err := db.Model(&model.Project{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			setFlash(c, "ERR:Failed to update project")
		} else {
			setFlash(c, "Project updated successfully!")
		}
		return c.Redirect("/admin/projects")
	}
}

func DeleteProjectHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		db.Delete(&model.Project{}, "id = ?", id)
		// HTMX delete — return empty row swap
		return c.SendString("")
	}
}

// setFlash stores a flash message in a short-lived cookie
func setFlash(c *fiber.Ctx, msg string) {
	c.Cookie(&fiber.Cookie{
		Name:    "flash",
		Value:   msg,
		Expires: time.Now().Add(10 * time.Second),
	})
}

// slugify is a simple slug helper (optional use)
func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
