package admin

import (
	"fmt"
	"portfolio-cms/internal/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ── List Skills ──────────────────────────────────────────────────────────────

func SkillsPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var skills []model.Skill
		db.Order("sort_order ASC, category ASC").Find(&skills)

		// Flash notification
		flash := c.Query("flash")
		notifHTML := ""
		if flash == "created" {
			notifHTML = notifBanner("success", "✓ Skill berhasil dibuat.")
		} else if flash == "updated" {
			notifHTML = notifBanner("success", "✓ Skill berhasil diperbarui.")
		} else if flash == "deleted" {
			notifHTML = notifBanner("info", "Skill telah dihapus.")
		}

		// Statistik
		total := len(skills)
		pubCount := 0
		avgProficiency := 0
		for _, s := range skills {
			if s.IsPublished {
				pubCount++
			}
			avgProficiency += s.Proficiency
		}
		if total > 0 {
			avgProficiency = avgProficiency / total
		}
		stats := ""
		if total > 0 {
			stats = fmt.Sprintf(`
<div style="display:grid;grid-template-columns:repeat(4,1fr);gap:16px;margin-bottom:24px;">
  %s %s %s %s
</div>`,
				statCard("◎", "Total", fmt.Sprintf("%d", total), "#e0e7ff"),
				statCard("✓", "Published", fmt.Sprintf("%d", pubCount), "#d1fae5"),
				statCard("✗", "Draft", fmt.Sprintf("%d", total-pubCount), "#fef3c7"),
				statCard("⚡", "Avg. Proficiency", fmt.Sprintf("%d%%", avgProficiency), "#ede9fe"),
			)
		}

		// Table rows
		rows := ""
		for i, s := range skills {
			// Tentukan warna progress bar berdasarkan proficiency
			proficiencyColor := "#d1cfd0" // default gray
			if s.Proficiency >= 80 {
				proficiencyColor = "#16a34a"
			} else if s.Proficiency >= 50 {
				proficiencyColor = "#f39c12"
			} else if s.Proficiency > 0 {
				proficiencyColor = "#dc2626"
			}

			publishedBadge := yesNo(s.IsPublished)
			bg := "#ffffff"
			if i%2 == 0 {
				bg = "#faf9fa"
			}

			rows += fmt.Sprintf(`
<tr style="background:%s;transition:background 0.15s;" onmouseover="this.style.background='#f5f4f5'" onmouseout="this.style.background='%s'">
  <td style="padding:12px 16px;font-size:13px;font-weight:500;">%s</td>
  <td style="padding:12px 16px;font-size:13px;">%s</td>
  <td style="padding:12px 16px;">
    <div style="display:flex;align-items:center;gap:8px;">
      <div style="width:100px;height:6px;background:#e5e3e4;border-radius:3px;overflow:hidden;">
        <div style="width:%d%%;height:100%%;background:%s;border-radius:3px;"></div>
      </div>
      <span style="font-size:12px;color:#4d4c4d;">%d%%</span>
    </div>
  </td>
  <td style="padding:12px 16px;">%s</td>
  <td style="padding:12px 16px;display:flex;gap:8px;">
    <a href="/admin/skills/%s/edit" style="font-size:12px;font-weight:500;color:#1e1d1e;text-decoration:none;padding:4px 10px;background:#f0eff0;border-radius:4px;">Edit</a>
    <button onclick="confirmDelete('/admin/skills/%s', '%s')" style="font-size:12px;font-weight:500;color:#b91c1c;background:none;border:none;cursor:pointer;padding:4px 10px;background:#fee2e2;border-radius:4px;">Delete</button>
  </td>
</tr>`, bg, bg, s.Name, s.Category, s.Proficiency, proficiencyColor, s.Proficiency, publishedBadge, s.ID, s.ID, s.Name)
		}

		empty := ""
		if len(skills) == 0 {
			empty = fmt.Sprintf(`<tr><td colspan="5" style="text-align:center;padding:40px;">%s</td></tr>`, emptyState("◎", "Belum ada skill", "Tambahkan skill pertama kamu", "New Skill", "/admin/skills/new"))
		}

		content := fmt.Sprintf(`
%s
<div style="text-align:right;margin-bottom:24px;">
  <a href="/admin/skills/new" style="display:inline-flex;align-items:center;gap:6px;padding:8px 18px;background:#1e1d1e;color:#f5f4f5;border-radius:6px;text-decoration:none;font-weight:500;font-size:13px;">
    + New Skill
  </a>
</div>
%s
<div style="background:#ffffff;border-radius:10px;border:1px solid #e5e3e4;overflow:hidden;">
  <table style="width:100%%;border-collapse:collapse;">
    <thead><tr style="background:#f8f7f8;">
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Name</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Category</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Proficiency</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Status</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Actions</th>
    </tr></thead>
    <tbody>%s %s</tbody>
  </table>
</div>

<!-- Delete Modal -->
<div id="deleteModal" style="display:none;position:fixed;top:0;left:0;width:100%%;height:100%%;background:rgba(0,0,0,0.4);z-index:1000;align-items:center;justify-content:center;">
  <div style="background:#fff;border-radius:16px;padding:24px;max-width:400px;width:90%%;box-shadow:0 20px 60px rgba(0,0,0,0.2);">
    <h3 style="font-size:16px;font-weight:600;color:#1e1d1e;">Hapus Skill</h3>
    <p style="font-size:13px;color:#6b6a6b;margin-top:8px;">Kamu yakin ingin menghapus skill <strong id="deleteItemName"></strong>?</p>
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
function closeDeleteModal() {
  document.getElementById('deleteModal').style.display = 'none';
}
document.getElementById('confirmDeleteBtn').addEventListener('click', function() {
  if (!_deleteUrl) return;
  fetch(_deleteUrl, { method: 'DELETE' }).then(r => { if (r.ok) window.location.href = '/admin/skills?flash=deleted'; });
});
</script>`, notifHTML, stats, rows, empty)

		return c.Type("html").SendString(layout("Skills", "Skills", content))
	}
}

// ── Create / Edit Form ────────────────────────────────────────────────────────

func SkillFormPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id", "")
		isEdit := false
		var skill model.Skill

		if id != "" {
			if err := db.First(&skill, "id = ?", id).Error; err != nil {
				return c.Status(404).SendString("Skill not found")
			}
			isEdit = true
		}

		title := "Create Skill"
		formAction := "/admin/skills"
		submitLabel := "Create Skill"
		if isEdit {
			title = "Edit Skill"
			formAction = "/admin/skills/" + id
			submitLabel = "Save Changes"
		}

		content := fmt.Sprintf(`
<div style="display:flex; justify-content: flex-end; margin-bottom: 16px;">
  <a href="/admin/skills" style="display:inline-flex;align-items:center;gap:6px;font-size:13px;color:#4d4c4d;text-decoration:none;padding:6px 12px;border:1px solid #e5e3e4;border-radius:6px;">← Back to Skills</a>
</div>
<form method="POST" action="%s">
  <div style="background:#fff;border-radius:10px;border:1px solid #e5e3e4;padding:24px;margin-bottom:16px;">
    <h2 style="font-size:14px;font-weight:600;color:#1e1d1e;">Informasi Skill</h2>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:16px;">
      %s %s %s
    </div>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:16px;">
      %s %s
    </div>
    <div style="margin-top:16px;">
      %s
    </div>
  </div>
  <div style="display:flex;justify-content:flex-end;gap:12px;margin-top:16px;padding:16px 0;">
    <a href="/admin/skills" style="padding:8px 16px;border:1px solid #e5e3e4;border-radius:6px;color:#4d4c4d;text-decoration:none;">Batal</a>
    <button type="submit" style="padding:8px 20px;background:#1e1d1e;color:#f5f4f5;border:none;border-radius:6px;font-weight:500;">%s</button>
  </div>
</form>`,
			formAction,
			inputField("Name", "name", skill.Name, "text", true),
			inputField("Category", "category", skill.Category, "text", false),
			inputField("Icon URL", "icon_url", skill.IconURL, "url", false),
			inputField("Proficiency (0-100)", "proficiency", fmt.Sprintf("%d", skill.Proficiency), "number", false),
			inputField("Sort Order", "sort_order", fmt.Sprintf("%d", skill.SortOrder), "number", false),
			toggleSwitch("Published", "is_published", skill.IsPublished, "Skill bisa dilihat publik"),
			submitLabel,
		)

		return c.Type("html").SendString(layout(title, "Skills", content))
	}
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func CreateSkillHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		proficiency := 0
		if p := c.FormValue("proficiency"); p != "" {
			fmt.Sscanf(p, "%d", &proficiency)
		}
		sortOrder := 0
		if so := c.FormValue("sort_order"); so != "" {
			fmt.Sscanf(so, "%d", &sortOrder)
		}
		skill := model.Skill{
			Name:        c.FormValue("name"),
			Category:    c.FormValue("category"),
			IconURL:     c.FormValue("icon_url"),
			Proficiency: proficiency,
			SortOrder:   sortOrder,
			IsPublished: c.FormValue("is_published") == "true",
		}
		if err := db.Create(&skill).Error; err != nil {
			return c.Status(500).SendString("Failed to create skill")
		}
		return c.Redirect("/admin/skills?flash=created")
	}
}

func UpdateSkillHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var skill model.Skill
		if err := db.First(&skill, "id = ?", id).Error; err != nil {
			return c.Status(404).SendString("Skill not found")
		}
		skill.Name = c.FormValue("name")
		skill.Category = c.FormValue("category")
		skill.IconURL = c.FormValue("icon_url")
		if p := c.FormValue("proficiency"); p != "" {
			fmt.Sscanf(p, "%d", &skill.Proficiency)
		}
		if so := c.FormValue("sort_order"); so != "" {
			fmt.Sscanf(so, "%d", &skill.SortOrder)
		}
		skill.IsPublished = c.FormValue("is_published") == "true"
		db.Save(&skill)
		return c.Redirect("/admin/skills?flash=updated")
	}
}

func DeleteSkillHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		db.Delete(&model.Skill{}, "id = ?", id)
		return c.SendStatus(200)
	}
}