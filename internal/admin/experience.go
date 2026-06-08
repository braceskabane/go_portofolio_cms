package admin

import (
	"fmt"
	"portfolio-cms/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ── List Experiences ─────────────────────────────────────────────────────────

func ExperiencesPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var items []model.Experience
		db.Preload("Category").
			Preload("StackItems").
			Preload("Skills").
			Order("sort_order ASC, start_date DESC").
			Find(&items)

		// Flash notification
		flash := c.Query("flash")
		notifHTML := ""
		if flash == "created" {
			notifHTML = notifBanner("success", "✓ Experience berhasil dibuat.")
		} else if flash == "updated" {
			notifHTML = notifBanner("success", "✓ Experience berhasil diperbarui.")
		} else if flash == "deleted" {
			notifHTML = notifBanner("info", "Experience telah dihapus.")
		}

		// Statistik
		total := len(items)
		pubCount := 0
		currentCount := 0
		for _, e := range items {
			if e.IsPublished {
				pubCount++
			}
			if e.IsCurrent {
				currentCount++
			}
		}
		stats := ""
		if total > 0 {
			stats = fmt.Sprintf(`
<div style="display:grid;grid-template-columns:repeat(4,1fr);gap:16px;margin-bottom:24px;">
  %s %s %s %s
</div>`,
				statCard("💼", "Total", fmt.Sprintf("%d", total), "#e0e7ff"),
				statCard("✓", "Published", fmt.Sprintf("%d", pubCount), "#d1fae5"),
				statCard("✗", "Draft", fmt.Sprintf("%d", total-pubCount), "#fef3c7"),
				statCard("📍", "Current", fmt.Sprintf("%d", currentCount), "#ede9fe"),
			)
		}

		// Table rows
		rows := ""
		for i, e := range items {
			catName := "—"
			if e.Category != nil {
				catName = e.Category.Name
			}
			// Stack items
			stackList := ""
			for j, si := range e.StackItems {
				if j > 0 {
					stackList += ", "
				}
				stackList += si.Name
			}
			if stackList == "" {
				stackList = `<span style="font-size:12px;color:#9a9899;">—</span>`
			}
			endDate := "Present"
			if e.EndDate != nil {
				endDate = e.EndDate.Format("2006-01")
			}
			publishedBadge := yesNo(e.IsPublished)
			bg := "#ffffff"
			if i%2 == 0 {
				bg = "#faf9fa"
			}
			rows += fmt.Sprintf(`
<tr style="background:%s;transition:background 0.15s;" onmouseover="this.style.background='#f5f4f5'" onmouseout="this.style.background='%s'">
  <td style="padding:12px 16px;font-size:13px;font-weight:500;">%s</td>
  <td style="padding:12px 16px;font-size:13px;">%s</td>
  <td style="padding:12px 16px;font-size:12px;">%s</td>
  <td style="padding:12px 16px;font-size:12px;">%s – %s</td>
  <td style="padding:12px 16px;font-size:12px;">%s</td>
  <td style="padding:12px 16px;">%s</td>
  <td style="padding:12px 16px;display:flex;gap:8px;">
    <a href="/admin/experiences/%s/edit" style="font-size:12px;font-weight:500;color:#1e1d1e;text-decoration:none;padding:4px 10px;background:#f0eff0;border-radius:4px;">Edit</a>
    <button onclick="confirmDelete('/admin/experiences/%s', '%s')" style="font-size:12px;font-weight:500;color:#b91c1c;background:none;border:none;cursor:pointer;padding:4px 10px;background:#fee2e2;border-radius:4px;">Delete</button>
  </td>
</tr>`, bg, bg, e.Company, e.Position, catName, e.StartDate.Format("2006-01"), endDate, stackList, publishedBadge, e.ID, e.ID, e.Company)
		}

		empty := ""
		if len(items) == 0 {
			empty = fmt.Sprintf(`<tr><td colspan="7" style="text-align:center;padding:40px;">%s</td></tr>`, emptyState("💼", "Belum ada experience", "Tambahkan pengalaman kerja pertama kamu", "New Experience", "/admin/experiences/new"))
		}

		content := fmt.Sprintf(`
%s
<div style="text-align:right;margin-bottom:24px;">
  <a href="/admin/experiences/new" style="display:inline-flex;align-items:center;gap:6px;padding:8px 18px;background:#1e1d1e;color:#f5f4f5;border-radius:6px;text-decoration:none;font-weight:500;font-size:13px;">
    + New Experience
  </a>
</div>
%s
<div style="background:#ffffff;border-radius:10px;border:1px solid #e5e3e4;overflow:hidden;">
  <table style="width:100%%;border-collapse:collapse;">
    <thead><tr style="background:#f8f7f8;">
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Company</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Position</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Category</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Period</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Stack</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Status</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Actions</th>
    </tr></thead>
    <tbody>%s %s</tbody>
  </table>
</div>

<!-- Delete Modal -->
<div id="deleteModal" style="display:none;position:fixed;top:0;left:0;width:100%%;height:100%%;background:rgba(0,0,0,0.4);z-index:1000;align-items:center;justify-content:center;">
  <div style="background:#fff;border-radius:16px;padding:24px;max-width:400px;width:90%%;box-shadow:0 20px 60px rgba(0,0,0,0.2);">
    <h3 style="font-size:16px;font-weight:600;color:#1e1d1e;">Hapus Experience</h3>
    <p style="font-size:13px;color:#6b6a6b;margin-top:8px;">Kamu yakin ingin menghapus <strong id="deleteItemName"></strong>?</p>
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
  document.getElementById('deleteItemName').textContent = name;
  document.getElementById('deleteModal').style.display = 'flex';
}
function closeDeleteModal() { document.getElementById('deleteModal').style.display = 'none'; }
document.getElementById('confirmDeleteBtn').addEventListener('click', function() {
  if (!_deleteUrl) return;
  fetch(_deleteUrl, { method: 'DELETE' }).then(r => { if (r.ok) window.location.href = '/admin/experiences?flash=deleted'; });
});
</script>`, notifHTML, stats, rows, empty)

		return c.Type("html").SendString(layout("Experiences", "Experiences", content))
	}
}

// ── Create / Edit Form ────────────────────────────────────────────────────────

func ExperienceFormPage(db *gorm.DB, fallbackID string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id", fallbackID)
		isEdit := false
		var exp model.Experience

		if id != "" {
			if err := db.Preload("Skills").Preload("StackItems").Preload("Category").First(&exp, "id = ?", id).Error; err != nil {
				return c.Status(404).SendString("Experience not found")
			}
			isEdit = true
		}

		// Categories dropdown
		var categories []model.ExperienceCategory
		db.Order("sort_order ASC, name ASC").Find(&categories)
		catOpts := `<option value="">— Pilih Category —</option>`
		for _, cat := range categories {
			sel := ""
			if exp.CategoryID != nil && *exp.CategoryID == cat.ID {
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
			for _, s := range exp.Skills {
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
			for _, s := range exp.StackItems {
				if s.ID == st.ID { checked = true; break }
			}
			label := st.Name
			if st.Category.Name != "" {
				label += " (" + st.Category.Name + ")"
			}
			stackItems = append(stackItems, struct{ ID, Text string; Checked bool }{st.ID.String(), label, checked})
		}

		title := "Create Experience"
		formAction := "/admin/experiences"
		submitLabel := "Create Experience"
		if isEdit {
			title = "Edit Experience"
			formAction = "/admin/experiences/" + exp.ID.String()
			submitLabel = "Save Changes"
		}

		content := fmt.Sprintf(`
<div style="display:flex; justify-content: flex-end; margin-bottom: 16px;">
  <a href="/admin/experiences" style="display:inline-flex;align-items:center;gap:6px;font-size:13px;color:#4d4c4d;text-decoration:none;padding:6px 12px;border:1px solid #e5e3e4;border-radius:6px;">← Back to Experiences</a>
</div>
<form method="POST" action="%s">
  <!-- Basic Info -->
  <div style="background:#fff;border-radius:10px;border:1px solid #e5e3e4;padding:24px;margin-bottom:16px;">
    <h2 style="font-size:14px;font-weight:600;color:#1e1d1e;">Informasi Dasar</h2>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:16px;">
      %s %s %s %s
    </div>
    <div style="margin-top:16px;">%s</div>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:16px;">
      %s %s %s %s
    </div>
    <div style="margin-top:16px;">
      <label style="font-size:12px;font-weight:600;color:#4d4c4d;">Category</label>
      <select name="category_id" style="width:100%%;margin-top:5px;border:1px solid #d1cfd0;border-radius:6px;padding:8px;font-size:13px;">%s</select>
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
      %s
    </div>
  </div>
  <!-- Actions -->
  <div style="display:flex;justify-content:flex-end;gap:12px;margin-top:16px;padding:16px 0;">
    <a href="/admin/experiences" style="padding:8px 16px;border:1px solid #e5e3e4;border-radius:6px;color:#4d4c4d;text-decoration:none;">Batal</a>
    <button type="submit" style="padding:8px 20px;background:#1e1d1e;color:#f5f4f5;border:none;border-radius:6px;font-weight:500;">%s</button>
  </div>
</form>`,
			formAction,
			inputField("Company", "company", exp.Company, "text", true),
			inputField("Position", "position", exp.Position, "text", true),
			inputField("Location", "location", exp.Location, "text", false),
			inputField("Company URL", "company_url", exp.CompanyURL, "url", false),
			textareaField("Description", "description", exp.Description, 3),
			inputField("Logo URL", "logo_url", exp.LogoURL, "url", false),
			inputField("Start Date", "start_date", exp.StartDate.Format("2006-01-02"), "date", true),
			inputField("End Date", "end_date", dateOrEmpty(exp.EndDate), "date", false),
			inputField("Sort Order", "sort_order", fmt.Sprintf("%d", exp.SortOrder), "number", false),
			catOpts,
			checkboxGroup("Skills", "skill_ids", skillItems),
			checkboxGroup("Tech Stack", "stack_item_ids", stackItems),
			toggleSwitch("Current Position", "is_current", exp.IsCurrent, "Masih bekerja di sini"),
			toggleSwitch("Published", "is_published", exp.IsPublished, "Tampilkan di halaman publik"),
			submitLabel,
		)

		return c.Type("html").SendString(layout(title, "Experiences", content))
	}
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func CreateExperienceHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startDate, _ := time.Parse("2006-01-02", c.FormValue("start_date"))
		var endDate *time.Time
		if end := c.FormValue("end_date"); end != "" {
			t, _ := time.Parse("2006-01-02", end)
			endDate = &t
		}
		sortOrder := 0
		if so := c.FormValue("sort_order"); so != "" {
			fmt.Sscanf(so, "%d", &sortOrder)
		}

		exp := model.Experience{
			Company:     c.FormValue("company"),
			Position:    c.FormValue("position"),
			Description: c.FormValue("description"),
			LogoURL:     c.FormValue("logo_url"),
			CompanyURL:  c.FormValue("company_url"),
			Location:    c.FormValue("location"),
			StartDate:   startDate,
			EndDate:     endDate,
			IsCurrent:   c.FormValue("is_current") == "true",
			IsPublished: c.FormValue("is_published") == "true",
			SortOrder:   sortOrder,
		}
		if catID := c.FormValue("category_id"); catID != "" {
			uid, _ := uuid.Parse(catID)
			exp.CategoryID = &uid
		}

		if err := db.Create(&exp).Error; err != nil {
			return c.Status(500).SendString("Failed to create experience")
		}

		if skillIDs := getFormArray(c, "skill_ids"); len(skillIDs) > 0 {
			var skills []model.Skill
			db.Find(&skills, "id IN ?", skillIDs)
			db.Model(&exp).Association("Skills").Replace(skills)
		}
		if stackIDs := getFormArray(c, "stack_item_ids"); len(stackIDs) > 0 {
			var stacks []model.StackItem
			db.Find(&stacks, "id IN ?", stackIDs)
			db.Model(&exp).Association("StackItems").Replace(stacks)
		}

		return c.Redirect("/admin/experiences?flash=created")
	}
}

func UpdateExperienceHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var exp model.Experience
		if err := db.First(&exp, "id = ?", id).Error; err != nil {
			return c.Status(404).SendString("Experience not found")
		}

		exp.Company = c.FormValue("company")
		exp.Position = c.FormValue("position")
		exp.Description = c.FormValue("description")
		exp.LogoURL = c.FormValue("logo_url")
		exp.CompanyURL = c.FormValue("company_url")
		exp.Location = c.FormValue("location")
		exp.IsCurrent = c.FormValue("is_current") == "true"
		exp.IsPublished = c.FormValue("is_published") == "true"
		if so := c.FormValue("sort_order"); so != "" {
			fmt.Sscanf(so, "%d", &exp.SortOrder)
		}

		startDate, _ := time.Parse("2006-01-02", c.FormValue("start_date"))
		exp.StartDate = startDate
		if end := c.FormValue("end_date"); end != "" {
			t, _ := time.Parse("2006-01-02", end)
			exp.EndDate = &t
		} else {
			exp.EndDate = nil
		}
		if catID := c.FormValue("category_id"); catID != "" {
			uid, _ := uuid.Parse(catID)
			exp.CategoryID = &uid
		} else {
			exp.CategoryID = nil
		}

		db.Save(&exp)

		// skills
		skillIDs := getFormArray(c, "skill_ids")
		var skills []model.Skill
		if len(skillIDs) > 0 {
			db.Find(&skills, "id IN ?", skillIDs)
		}
		db.Model(&exp).Association("Skills").Replace(skills)

		// stacks
		stackIDs := getFormArray(c, "stack_item_ids")
		var stacks []model.StackItem
		if len(stackIDs) > 0 {
			db.Find(&stacks, "id IN ?", stackIDs)
		}
		db.Model(&exp).Association("StackItems").Replace(stacks)

		return c.Redirect("/admin/experiences?flash=updated")
	}
}

func DeleteExperienceHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var exp model.Experience
		if err := db.First(&exp, "id = ?", id).Error; err == nil {
			_ = db.Model(&exp).Association("Skills").Clear()
			_ = db.Model(&exp).Association("StackItems").Clear()
			db.Delete(&exp)
		}
		return c.SendStatus(200)
	}
}