package admin

import (
	"fmt"
	"portfolio-cms/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ── List Projects ───────────────────────────────────────────────────────────

func ProjectsPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var projects []model.Project
		db.Preload("Category").Preload("StackItems").Preload("Skills").
			Order("sort_order ASC, created_at DESC").Find(&projects)

		// Flash notification
		flash := c.Query("flash")
		notifHTML := ""
		if flash == "created" {
			notifHTML = notifBanner("success", "✓ Project berhasil dibuat.")
		} else if flash == "updated" {
			notifHTML = notifBanner("success", "✓ Project berhasil diperbarui.")
		} else if flash == "deleted" {
			notifHTML = notifBanner("info", "Project telah dihapus.")
		}

		// Statistik
		total := len(projects)
		pubCount := 0
		for _, p := range projects {
			if p.IsPublished {
				pubCount++
			}
		}
		stats := fmt.Sprintf(`
<div style="display:grid;grid-template-columns:repeat(3,1fr);gap:16px;margin-bottom:24px;">
  %s %s %s
</div>`,
			statCard("▤", "Total", fmt.Sprintf("%d", total), "#e0e7ff"),
			statCard("✓", "Published", fmt.Sprintf("%d", pubCount), "#d1fae5"),
			statCard("✗", "Draft", fmt.Sprintf("%d", total-pubCount), "#fef3c7"),
		)

		// Table rows
		rows := ""
		for i, p := range projects {
			catName := "—"
			if p.Category != nil {
				catName = p.Category.Name
			}
			stackList := ""
			maxStack := 3
			for i, si := range p.StackItems {
				if i >= maxStack {
					remaining := len(p.StackItems) - maxStack
					stackList += fmt.Sprintf(`<span style="display:inline-block;padding:2px 8px;border-radius:999px;font-size:11px;background:#e5e3e4;color:#4d4c4d;">+%d</span>`, remaining)
					break
				}
				stackList += fmt.Sprintf(`<span style="display:inline-block;padding:2px 8px;border-radius:999px;font-size:11px;background:#f0eff0;color:#4d4c4d;">%s</span>`, si.Name)
			}
			if stackList == "" {
				stackList = `<span style="font-size:12px;color:#9a9899;">—</span>`
			}

			publishedBadge := yesNo(p.IsPublished)
			featuredBadge := ""
			if p.IsFeatured {
				featuredBadge = ` <span style="color:#f39c12;">★</span>`
			}

			rows += fmt.Sprintf(`
<tr style="background:%s;transition:background 0.15s;" onmouseover="this.style.background='#f5f4f5'" onmouseout="this.style.background='%s'">
  <td style="padding:12px 16px;font-size:13px;">
    <span style="font-weight:500;">%s</span>%s
    <br><span style="font-size:11px;color:#9a9899;">%s</span>
  </td>
  <td style="padding:12px 16px;font-size:13px;">%s</td>
  <td style="padding:12px 16px;">%s</td>
  <td style="padding:12px 16px;font-size:11px;color:#6b6a6b;">%s</td>
  <td style="padding:12px 16px;">%s</td>
  <td style="padding:12px 16px;display:flex;gap:8px;">
    <a href="/admin/projects/%s/edit" style="font-size:12px;font-weight:500;color:#1e1d1e;text-decoration:none;padding:4px 10px;background:#f0eff0;border-radius:4px;">Edit</a>
    <button onclick="confirmDelete('/admin/projects/%s', '%s')" style="font-size:12px;font-weight:500;color:#b91c1c;background:none;border:none;cursor:pointer;padding:4px 10px;background:#fee2e2;border-radius:4px;">Delete</button>
  </td>
</tr>`,
				func() string { if i%2==0 { return "#ffffff" } else { return "#faf9fa" } }(),
				func() string { if i%2==0 { return "#ffffff" } else { return "#faf9fa" } }(),
				p.Title, featuredBadge,
				p.Slug,
				catName,
				stackList,
				p.CreatedAt.Format("2006-01-02"),
				publishedBadge,
				p.ID, p.ID, p.Title,
			)
		}

		empty := ""
		if len(projects) == 0 {
			empty = fmt.Sprintf(`<tr><td colspan="6" style="text-align:center;padding:40px;">%s</td></tr>`, emptyState("📁", "Belum ada project", "Tambahkan project pertama kamu", "New Project", "/admin/projects/new"))
		}

		content := fmt.Sprintf(`
%s
<div style="text-align:right;margin-bottom:24px;">
  <a href="/admin/projects/new" style="display:inline-flex;align-items:center;gap:6px;padding:8px 18px;background:#1e1d1e;color:#f5f4f5;border-radius:6px;text-decoration:none;font-weight:500;font-size:13px;">
    + New Project
  </a>
</div>
%s
<div style="background:#ffffff;border-radius:10px;border:1px solid #e5e3e4;overflow:hidden;">
  <table style="width:100%%;border-collapse:collapse;">
    <thead><tr style="background:#f8f7f8;">
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Project</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Category</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Stack</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Created</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Status</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Actions</th>
    </tr></thead>
    <tbody>%s %s</tbody>
  </table>
</div>

<!-- Delete Modal -->
<div id="deleteModal" style="display:none;position:fixed;top:0;left:0;width:100%%;height:100%%;background:rgba(0,0,0,0.4);z-index:1000;align-items:center;justify-content:center;">
  <div style="background:#fff;border-radius:16px;padding:24px;max-width:400px;width:90%%;box-shadow:0 20px 60px rgba(0,0,0,0.2);">
    <h3 style="font-size:16px;font-weight:600;color:#1e1d1e;">Hapus Project</h3>
    <p style="font-size:13px;color:#6b6a6b;margin-top:8px;">Kamu yakin ingin menghapus project <strong id="deleteProjectName"></strong>?</p>
    <div style="display:flex;gap:12px;margin-top:20px;">
      <button onclick="closeDeleteModal()" style="flex:1;padding:8px;background:#f0eff0;border:none;border-radius:6px;font-weight:500;">Batal</button>
      <button id="confirmDeleteBtn" style="flex:1;padding:8px;background:#b91c1c;color:#fff;border:none;border-radius:6px;font-weight:500;">Ya, Hapus</button>
    </div>
  </div>
</div>

<script>
let _deleteUrl = '';
function confirmDelete(url, name) {
  _deleteUrl = url;
  document.getElementById('deleteProjectName').textContent = name;
  document.getElementById('deleteModal').style.display = 'flex';
}
function closeDeleteModal() {
  document.getElementById('deleteModal').style.display = 'none';
}
document.getElementById('confirmDeleteBtn').addEventListener('click', function() {
  if (!_deleteUrl) return;
  fetch(_deleteUrl, { method: 'DELETE' }).then(r => { if (r.ok) window.location.href = '/admin/projects?flash=deleted'; });
});
</script>`, notifHTML, stats, rows, empty)

		return c.Type("html").SendString(layout("Projects", "Projects", content))
	}
}

// ── Form Create/Edit ────────────────────────────────────────────────────────
func ProjectFormPage(db *gorm.DB, fallbackID string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id", fallbackID)
		isEdit := false
		var project model.Project

		if id != "" {
			if err := db.Preload("Skills").Preload("StackItems").Preload("Category").First(&project, "id = ?", id).Error; err != nil {
				return c.Status(404).SendString("Project not found")
			}
			isEdit = true
		}

		// Categories
		var categories []model.ProjectCategory
		db.Order("sort_order ASC, name ASC").Find(&categories)
		catOpts := `<option value="">— Pilih Category —</option>`
		for _, cat := range categories {
			sel := ""
			if project.CategoryID != nil && *project.CategoryID == cat.ID {
				sel = "selected"
			}
			catOpts += fmt.Sprintf(`<option value="%s" %s>%s</option>`, cat.ID, sel, cat.Name)
		}

		// Skills checkboxes
		var skills []model.Skill
		db.Find(&skills)
		var skillItems []struct{ ID, Text string; Checked bool }
		for _, sk := range skills {
			checked := false
			for _, s := range project.Skills {
				if s.ID == sk.ID { checked = true; break }
			}
			skillItems = append(skillItems, struct{ ID, Text string; Checked bool }{sk.ID.String(), sk.Name, checked})
		}

		// Stack items checkboxes
		var stacks []model.StackItem
		db.Preload("Category").Find(&stacks)
		var stackItems []struct{ ID, Text string; Checked bool }
		for _, st := range stacks {
			checked := false
			for _, s := range project.StackItems {
				if s.ID == st.ID { checked = true; break }
			}
			label := st.Name
			if st.Category.Name != "" {
				label += " (" + st.Category.Name + ")"
			}
			stackItems = append(stackItems, struct{ ID, Text string; Checked bool }{st.ID.String(), label, checked})
		}

		title := "Create Project"
		formAction := "/admin/projects"
		submitLabel := "Create Project"
		if isEdit {
			title = "Edit Project"
			formAction = "/admin/projects/" + project.ID.String()
			submitLabel = "Save Changes"
		}

		content := fmt.Sprintf(`
<div style="display:flex; justify-content: flex-end; margin-bottom: 16px;">
  <a href="/admin/projects" style="display:inline-flex;align-items:center;gap:6px;font-size:13px;color:#4d4c4d;text-decoration:none;padding:6px 12px;border:1px solid #e5e3e4;border-radius:6px;">← Back to Projects</a>
</div>
<form method="POST" action="%s">
  <!-- Basic Info -->
  <div style="background:#fff;border-radius:10px;border:1px solid #e5e3e4;padding:24px;margin-bottom:16px;">
    <h2 style="font-size:14px;font-weight:600;color:#1e1d1e;">Informasi Dasar</h2>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:16px;">
      %s %s
      <div style="grid-column:span 2;">%s</div>
      <div style="grid-column:span 2;">
        <label style="font-size:12px;font-weight:600;color:#4d4c4d;">Category</label>
        <select name="category_id" style="width:100%%;margin-top:5px;border:1px solid #d1cfd0;border-radius:6px;padding:8px;font-size:13px;">%s</select>
      </div>
    </div>
  </div>
  <!-- Links -->
  <div style="background:#fff;border-radius:10px;border:1px solid #e5e3e4;padding:24px;margin-bottom:16px;">
    <h2 style="font-size:14px;font-weight:600;color:#1e1d1e;">Links & Thumbnail</h2>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:16px;">
      %s %s %s %s
    </div>
  </div>
  <!-- Case Study -->
  <div style="background:#fff;border-radius:10px;border:1px solid #e5e3e4;padding:24px;margin-bottom:16px;">
    <h2 style="font-size:14px;font-weight:600;color:#1e1d1e;">Case Study</h2>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:16px;">
      %s %s %s %s
      <div style="grid-column:span 2;">%s</div>
    </div>
  </div>
  <!-- Skills & Stack -->
  <div style="background:#fff;border-radius:10px;border:1px solid #e5e3e4;padding:24px;margin-bottom:16px;">
    <h2 style="font-size:14px;font-weight:600;color:#1e1d1e;">Skills & Tech Stack</h2>
    <div style="margin-top:16px;">%s</div>
    <div style="margin-top:16px;">%s</div>
  </div>
  <!-- Settings -->
  <div style="background:#fff;border-radius:10px;border:1px solid #e5e3e4;padding:24px;margin-bottom:16px;">
    <h2 style="font-size:14px;font-weight:600;color:#1e1d1e;">Pengaturan</h2>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:16px;">
      %s %s
      %s %s
    </div>
  </div>
  <!-- Sticky Footer -->
  <div style="display:flex;justify-content:flex-end;gap:12px;margin-top:16px;padding:16px 0;">
    <a href="/admin/projects" style="padding:8px 16px;border:1px solid #e5e3e4;border-radius:6px;color:#4d4c4d;text-decoration:none;">Batal</a>
    <button type="submit" style="padding:8px 20px;background:#1e1d1e;color:#f5f4f5;border:none;border-radius:6px;font-weight:500;">%s</button>
  </div>
</form>`,
			formAction,
			inputField("Title", "title", project.Title, "text", true),
			inputField("Slug", "slug", project.Slug, "text", true),
			textareaField("Description", "description", project.Description, 3),
			catOpts,
			inputField("Demo URL", "demo_url", project.DemoURL, "url", false),
			inputField("Repo URL", "repo_url", project.RepoURL, "url", false),
			inputField("Doc URL", "doc_url", project.DocURL, "url", false),
			inputField("Thumbnail URL", "thumbnail_url", project.ThumbnailURL, "url", false),
			textareaField("Problem", "problem", project.Problem, 2),
			textareaField("Solution", "solution", project.Solution, 2),
			textareaField("My Role", "my_role", project.MyRole, 2),
			textareaField("Impact", "impact", project.Impact, 2),
			textareaField("Content (Markdown)", "content", project.Content, 6),
			checkboxGroup("Skills", "skill_ids", skillItems),
			checkboxGroup("Tech Stack", "stack_item_ids", stackItems),
			inputField("Start Date", "start_date", dateOrEmpty(project.StartDate), "date", false),
			inputField("End Date", "end_date", dateOrEmpty(project.EndDate), "date", false),
			toggleSwitch("Featured", "is_featured", project.IsFeatured, "Tampilkan di halaman utama"),
			toggleSwitch("Published", "is_published", project.IsPublished, "Project bisa dilihat publik"),
			submitLabel,
		)
		return c.Type("html").SendString(layout(title, "Projects", content))
	}
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func CreateProjectHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		project := model.Project{
			Title:        c.FormValue("title"),
			Slug:         c.FormValue("slug"),
			Description:  c.FormValue("description"),
			Content:      c.FormValue("content"),
			ThumbnailURL: c.FormValue("thumbnail_url"),
			DemoURL:      c.FormValue("demo_url"),
			RepoURL:      c.FormValue("repo_url"),
			DocURL:       c.FormValue("doc_url"),
			Problem:      c.FormValue("problem"),
			Solution:     c.FormValue("solution"),
			MyRole:       c.FormValue("my_role"),
			Impact:       c.FormValue("impact"),
			IsFeatured:   c.FormValue("is_featured") == "true",
			IsPublished:  c.FormValue("is_published") == "true",
			SortOrder:    0,
		}

		if start := c.FormValue("start_date"); start != "" {
			t, _ := time.Parse("2006-01-02", start)
			project.StartDate = &t
		}
		if end := c.FormValue("end_date"); end != "" {
			t, _ := time.Parse("2006-01-02", end)
			project.EndDate = &t
		}
		if catID := c.FormValue("category_id"); catID != "" {
			uid, err := uuid.Parse(catID)
			if err == nil {
				project.CategoryID = &uid
			}
		}

		if err := db.Create(&project).Error; err != nil {
			return c.Status(500).SendString("Failed to create project")
		}

		if skillIDs := getFormArray(c, "skill_ids"); len(skillIDs) > 0 {
			var skills []model.Skill
			db.Find(&skills, "id IN ?", skillIDs)
			db.Model(&project).Association("Skills").Replace(skills)
		}
		if stackIDs := getFormArray(c, "stack_item_ids"); len(stackIDs) > 0 {
			var stacks []model.StackItem
			db.Find(&stacks, "id IN ?", stackIDs)
			db.Model(&project).Association("StackItems").Replace(stacks)
		}

		return c.Redirect("/admin/projects?flash=created")
	}
}

func UpdateProjectHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var project model.Project
		if err := db.First(&project, "id = ?", id).Error; err != nil {
			return c.Status(404).SendString("Project not found")
		}

		project.Title = c.FormValue("title")
		project.Slug = c.FormValue("slug")
		project.Description = c.FormValue("description")
		project.Content = c.FormValue("content")
		project.ThumbnailURL = c.FormValue("thumbnail_url")
		project.DemoURL = c.FormValue("demo_url")
		project.RepoURL = c.FormValue("repo_url")
		project.DocURL = c.FormValue("doc_url")
		project.Problem = c.FormValue("problem")
		project.Solution = c.FormValue("solution")
		project.MyRole = c.FormValue("my_role")
		project.Impact = c.FormValue("impact")
		project.IsFeatured = c.FormValue("is_featured") == "true"
		project.IsPublished = c.FormValue("is_published") == "true"

		if start := c.FormValue("start_date"); start != "" {
			t, _ := time.Parse("2006-01-02", start)
			project.StartDate = &t
		} else {
			project.StartDate = nil
		}
		if end := c.FormValue("end_date"); end != "" {
			t, _ := time.Parse("2006-01-02", end)
			project.EndDate = &t
		} else {
			project.EndDate = nil
		}
		if catID := c.FormValue("category_id"); catID != "" {
			uid, _ := uuid.Parse(catID)
			project.CategoryID = &uid
		} else {
			project.CategoryID = nil
		}

		db.Save(&project)

		skillIDs := getFormArray(c, "skill_ids")
		var skills []model.Skill
		if len(skillIDs) > 0 {
			db.Find(&skills, "id IN ?", skillIDs)
		}
		db.Model(&project).Association("Skills").Replace(skills)

		stackIDs := getFormArray(c, "stack_item_ids")
		var stacks []model.StackItem
		if len(stackIDs) > 0 {
			db.Find(&stacks, "id IN ?", stackIDs)
		}
		db.Model(&project).Association("StackItems").Replace(stacks)

		return c.Redirect("/admin/projects?flash=updated")
	}
}

func DeleteProjectHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var p model.Project
		if err := db.First(&p, "id = ?", id).Error; err == nil {
			_ = db.Model(&p).Association("Skills").Clear()
			_ = db.Model(&p).Association("StackItems").Clear()
			db.Delete(&p)
		}
		return c.SendStatus(200)
	}
}

// ── UI Component Helpers ──────────────────────────────────────────────────────

func notifBanner(kind, message string) string {
	bg := "bg-emerald-50 border-emerald-200 text-emerald-800"
	icon := `<svg class="w-4 h-4 text-emerald-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`
	if kind == "info" {
		bg = "bg-blue-50 border-blue-200 text-blue-800"
		icon = `<svg class="w-4 h-4 text-blue-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`
	}
	return fmt.Sprintf(`
<div id="notif-banner" class="flex items-center gap-3 px-4 py-3 mb-5 rounded-xl border %s text-sm font-medium transition-all duration-300" style="transition: opacity 0.3s, transform 0.3s;">
  %s
  <span>%s</span>
  <button onclick="this.parentElement.remove()" class="ml-auto opacity-60 hover:opacity-100 transition-opacity">
    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
  </button>
</div>`, bg, icon, message)
}

func formInput(label, name, value, inputType string, required bool, placeholder string) string {
	req := ""
	reqMark := ""
	if required {
		req = "required"
		reqMark = `<span class="text-red-500 ml-0.5">*</span>`
	}
	return fmt.Sprintf(`
<div>
  <label class="block text-xs font-semibold text-gray-600 uppercase tracking-wide mb-1.5">%s%s</label>
  <input type="%s" name="%s" value="%s" placeholder="%s" %s
         class="w-full px-3 py-2 text-sm bg-white border border-gray-200 rounded-lg text-gray-800 placeholder-gray-400
                focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-400 transition-all"/>
</div>`, label, reqMark, inputType, name, value, placeholder, req)
}

func formTextarea(label, name, value string, rows int, placeholder string) string {
	return fmt.Sprintf(`
<div>
  <label class="block text-xs font-semibold text-gray-600 uppercase tracking-wide mb-1.5">%s</label>
  <textarea name="%s" rows="%d" placeholder="%s"
            class="w-full px-3 py-2 text-sm bg-white border border-gray-200 rounded-lg text-gray-800 placeholder-gray-400
                   focus:outline-none focus:ring-2 focus:ring-gray-900/10 focus:border-gray-400 transition-all resize-y">%s</textarea>
</div>`, label, name, rows, placeholder, value)
}

func toggleSwitch(label, name string, checked bool, description string) string {
	checkedAttr := ""
	if checked {
		checkedAttr = "checked"
	}
	return fmt.Sprintf(`
<div class="flex items-start justify-between gap-4 p-4 rounded-xl border border-gray-200 bg-gray-50/50">
  <div>
    <p class="text-sm font-medium text-gray-800">%s</p>
    <p class="text-xs text-gray-500 mt-0.5">%s</p>
  </div>
  <label class="relative inline-flex items-center cursor-pointer flex-shrink-0 mt-0.5">
    <input type="hidden" name="%s" value="false">
    <input type="checkbox" name="%s" value="true" %s class="sr-only peer">
    <div class="w-10 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-gray-400 rounded-full peer
                peer-checked:after:translate-x-4 after:content-[''] after:absolute after:top-[2px] after:left-[2px]
                after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all
                peer-checked:bg-gray-900 transition-colors"></div>
  </label>
</div>`, label, description, name, name, checkedAttr)
}

func buildCheckboxGroup(skills []model.Skill, selected []model.Skill, fieldName string) string {
	if len(skills) == 0 {
		return `<span class="text-sm text-gray-400 italic">Belum ada skill tersedia</span>`
	}
	result := ""
	for _, sk := range skills {
		checked := ""
		for _, ps := range selected {
			if ps.ID == sk.ID {
				checked = "checked"
				break
			}
		}
		checkedClass := ""
		if checked != "" {
			checkedClass = "border-gray-800 bg-gray-800 text-white"
		} else {
			checkedClass = "border-gray-200 bg-white text-gray-700 hover:border-gray-400"
		}
		result += fmt.Sprintf(`
<label class="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg border cursor-pointer text-xs font-medium transition-all %s select-none">
  <input type="checkbox" name="%s" value="%s" %s class="sr-only peer">
  <span class="peer-checked:text-white">%s</span>
</label>`, checkedClass, fieldName, sk.ID, checked, sk.Name)
	}
	return result
}

func buildStackCheckboxGroup(items []model.StackItem, selected []model.StackItem) string {
	if len(items) == 0 {
		return `<span class="text-sm text-gray-400 italic">Belum ada stack item tersedia</span>`
	}
	result := ""
	for _, st := range items {
		checked := ""
		for _, ps := range selected {
			if ps.ID == st.ID {
				checked = "checked"
				break
			}
		}
		catLabel := ""
		if st.Category.Name != "" {
			catLabel = fmt.Sprintf(`<span class="opacity-60 font-normal">· %s</span>`, st.Category.Name)
		}
		checkedClass := ""
		if checked != "" {
			checkedClass = "border-violet-600 bg-violet-600 text-white"
		} else {
			checkedClass = "border-gray-200 bg-white text-gray-700 hover:border-violet-400"
		}
		result += fmt.Sprintf(`
<label class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg border cursor-pointer text-xs font-medium transition-all %s select-none">
  <input type="checkbox" name="stack_item_ids" value="%s" %s class="sr-only">
  %s %s
</label>`, checkedClass, st.ID, checked, st.Name, catLabel)
	}
	return result
}

// ── Misc helpers ──────────────────────────────────────────────────────────────

func dateOrEmpty(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

func getFormArray(c *fiber.Ctx, key string) []string {
	var values []string
	c.Request().PostArgs().VisitAll(func(k, v []byte) {
		if string(k) == key {
			values = append(values, string(v))
		}
	})
	return values
}