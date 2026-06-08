package admin

import (
	"fmt"
	"portfolio-cms/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ── List Educations ──────────────────────────────────────────────────────────

func EducationsPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var items []model.Education
		db.Order("sort_order ASC, start_date DESC").Find(&items)

		// Flash notification
		flash := c.Query("flash")
		notifHTML := ""
		if flash == "created" {
			notifHTML = notifBanner("success", "✓ Education berhasil dibuat.")
		} else if flash == "updated" {
			notifHTML = notifBanner("success", "✓ Education berhasil diperbarui.")
		} else if flash == "deleted" {
			notifHTML = notifBanner("info", "Education telah dihapus.")
		}

		// Statistik
		total := len(items)
		pubCount := 0
		for _, e := range items {
			if e.IsPublished {
				pubCount++
			}
		}
		stats := ""
		if total > 0 {
			stats = fmt.Sprintf(`
<div style="display:grid;grid-template-columns:repeat(3,1fr);gap:16px;margin-bottom:24px;">
  %s %s %s
</div>`,
				statCard("🎓", "Total", fmt.Sprintf("%d", total), "#e0e7ff"),
				statCard("✓", "Published", fmt.Sprintf("%d", pubCount), "#d1fae5"),
				statCard("✗", "Draft", fmt.Sprintf("%d", total-pubCount), "#fef3c7"),
			)
		}

		// Table rows
		rows := ""
		for i, e := range items {
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
  <td style="padding:12px 16px;">%s</td>
  <td style="padding:12px 16px;display:flex;gap:8px;">
    <a href="/admin/educations/%s/edit" style="font-size:12px;font-weight:500;color:#1e1d1e;text-decoration:none;padding:4px 10px;background:#f0eff0;border-radius:4px;">Edit</a>
    <button onclick="confirmDelete('/admin/educations/%s', '%s')" style="font-size:12px;font-weight:500;color:#b91c1c;background:none;border:none;cursor:pointer;padding:4px 10px;background:#fee2e2;border-radius:4px;">Delete</button>
  </td>
</tr>`, bg, bg, e.Institution, e.Degree, e.FieldOfStudy, e.StartDate.Format("2006-01"), endDate, publishedBadge, e.ID, e.ID, e.Institution)
		}

		empty := ""
		if len(items) == 0 {
			empty = fmt.Sprintf(`<tr><td colspan="6" style="text-align:center;padding:40px;">%s</td></tr>`, emptyState("🎓", "Belum ada education", "Tambahkan riwayat pendidikan pertama kamu", "New Education", "/admin/educations/new"))
		}

		content := fmt.Sprintf(`
%s
<div style="text-align:right;margin-bottom:24px;">
  <a href="/admin/educations/new" style="display:inline-flex;align-items:center;gap:6px;padding:8px 18px;background:#1e1d1e;color:#f5f4f5;border-radius:6px;text-decoration:none;font-weight:500;font-size:13px;">
    + New Education
  </a>
</div>
%s
<div style="background:#ffffff;border-radius:10px;border:1px solid #e5e3e4;overflow:hidden;">
  <table style="width:100%%;border-collapse:collapse;">
    <thead><tr style="background:#f8f7f8;">
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Institution</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Degree</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Field of Study</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Period</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Status</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Actions</th>
    </tr></thead>
    <tbody>%s %s</tbody>
  </table>
</div>

<!-- Delete Modal -->
<div id="deleteModal" style="display:none;position:fixed;top:0;left:0;width:100%%;height:100%%;background:rgba(0,0,0,0.4);z-index:1000;align-items:center;justify-content:center;">
  <div style="background:#fff;border-radius:16px;padding:24px;max-width:400px;width:90%%;box-shadow:0 20px 60px rgba(0,0,0,0.2);">
    <h3 style="font-size:16px;font-weight:600;color:#1e1d1e;">Hapus Education</h3>
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
  fetch(_deleteUrl, { method: 'DELETE' }).then(r => { if (r.ok) window.location.href = '/admin/educations?flash=deleted'; });
});
</script>`, notifHTML, stats, rows, empty)

		return c.Type("html").SendString(layout("Educations", "Educations", content))
	}
}

// ── Create / Edit Form ────────────────────────────────────────────────────────

func EducationFormPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id", "")
		isEdit := false
		var edu model.Education

		if id != "" {
			if err := db.First(&edu, "id = ?", id).Error; err != nil {
				return c.Status(404).SendString("Education not found")
			}
			isEdit = true
		}

		title := "Create Education"
		formAction := "/admin/educations"
		submitLabel := "Create Education"
		if isEdit {
			title = "Edit Education"
			formAction = "/admin/educations/" + id
			submitLabel = "Save Changes"
		}

		content := fmt.Sprintf(`
<div style="display:flex; justify-content: flex-end; margin-bottom: 16px;">
  <a href="/admin/educations" style="display:inline-flex;align-items:center;gap:6px;font-size:13px;color:#4d4c4d;text-decoration:none;padding:6px 12px;border:1px solid #e5e3e4;border-radius:6px;">← Back to Educations</a>
</div>
<form method="POST" action="%s">
  <div style="background:#fff;border-radius:10px;border:1px solid #e5e3e4;padding:24px;margin-bottom:16px;">
    <h2 style="font-size:14px;font-weight:600;color:#1e1d1e;">Informasi Pendidikan</h2>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:16px;">
      %s %s %s %s
    </div>
    <div style="margin-top:16px;">%s</div>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:16px;">
      %s %s %s %s
    </div>
    <div style="margin-top:16px;">
      %s
    </div>
  </div>
  <div style="display:flex;justify-content:flex-end;gap:12px;margin-top:16px;padding:16px 0;">
    <a href="/admin/educations" style="padding:8px 16px;border:1px solid #e5e3e4;border-radius:6px;color:#4d4c4d;text-decoration:none;">Batal</a>
    <button type="submit" style="padding:8px 20px;background:#1e1d1e;color:#f5f4f5;border:none;border-radius:6px;font-weight:500;">%s</button>
  </div>
</form>`,
			formAction,
			inputField("Institution", "institution", edu.Institution, "text", true),
			inputField("Degree", "degree", edu.Degree, "text", false),
			inputField("Field of Study", "field_of_study", edu.FieldOfStudy, "text", false),
			inputField("Logo URL", "logo_url", edu.LogoURL, "url", false),
			textareaField("Description", "description", edu.Description, 3),
			inputField("GPA", "gpa", edu.GPA, "text", false),
			inputField("Start Date", "start_date", edu.StartDate.Format("2006-01-02"), "date", true),
			inputField("End Date", "end_date", dateOrEmpty(edu.EndDate), "date", false),
			inputField("Sort Order", "sort_order", fmt.Sprintf("%d", edu.SortOrder), "number", false),
			toggleSwitch("Published", "is_published", edu.IsPublished, "Tampilkan di halaman publik"),
			submitLabel,
		)

		return c.Type("html").SendString(layout(title, "Educations", content))
	}
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func CreateEducationHandler(db *gorm.DB) fiber.Handler {
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
		edu := model.Education{
			Institution:  c.FormValue("institution"),
			Degree:       c.FormValue("degree"),
			FieldOfStudy: c.FormValue("field_of_study"),
			Description:  c.FormValue("description"),
			LogoURL:      c.FormValue("logo_url"),
			StartDate:    startDate,
			EndDate:      endDate,
			GPA:          c.FormValue("gpa"),
			IsPublished:  c.FormValue("is_published") == "true",
			SortOrder:    sortOrder,
		}
		if err := db.Create(&edu).Error; err != nil {
			return c.Status(500).SendString("Failed to create education")
		}
		return c.Redirect("/admin/educations?flash=created")
	}
}

func UpdateEducationHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var edu model.Education
		if err := db.First(&edu, "id = ?", id).Error; err != nil {
			return c.Status(404).SendString("Education not found")
		}
		startDate, _ := time.Parse("2006-01-02", c.FormValue("start_date"))
		var endDate *time.Time
		if end := c.FormValue("end_date"); end != "" {
			t, _ := time.Parse("2006-01-02", end)
			endDate = &t
		}
		edu.Institution = c.FormValue("institution")
		edu.Degree = c.FormValue("degree")
		edu.FieldOfStudy = c.FormValue("field_of_study")
		edu.Description = c.FormValue("description")
		edu.LogoURL = c.FormValue("logo_url")
		edu.StartDate = startDate
		edu.EndDate = endDate
		edu.GPA = c.FormValue("gpa")
		edu.IsPublished = c.FormValue("is_published") == "true"
		if so := c.FormValue("sort_order"); so != "" {
			fmt.Sscanf(so, "%d", &edu.SortOrder)
		}
		db.Save(&edu)
		return c.Redirect("/admin/educations?flash=updated")
	}
}

func DeleteEducationHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		db.Delete(&model.Education{}, "id = ?", id)
		return c.SendStatus(200)
	}
}