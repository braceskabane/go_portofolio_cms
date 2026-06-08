package admin

import (
	"fmt"
	"portfolio-cms/internal/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ── List Experience Categories ───────────────────────────────────────────────

func ExperienceCategoriesPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var categories []model.ExperienceCategory
		db.Order("sort_order ASC, name ASC").Find(&categories)

		// Flash notification
		flash := c.Query("flash")
		notifHTML := ""
		if flash == "created" {
			notifHTML = notifBanner("success", "✓ Kategori berhasil dibuat.")
		} else if flash == "updated" {
			notifHTML = notifBanner("success", "✓ Kategori berhasil diperbarui.")
		} else if flash == "deleted" {
			notifHTML = notifBanner("info", "Kategori telah dihapus.")
		}

		// Statistik
		total := len(categories)
		pubCount := 0
		for _, c := range categories {
			if c.IsPublished {
				pubCount++
			}
		}
		stats := ""
		if total > 0 {
			stats = fmt.Sprintf(`
<div style="display:grid;grid-template-columns:repeat(3,1fr);gap:16px;margin-bottom:24px;">
  %s %s %s
</div>`,
				statCard("🏷️", "Total", fmt.Sprintf("%d", total), "#e0e7ff"),
				statCard("✓", "Published", fmt.Sprintf("%d", pubCount), "#d1fae5"),
				statCard("✗", "Draft", fmt.Sprintf("%d", total-pubCount), "#fef3c7"),
			)
		}

		// Table rows
		rows := ""
		for i, cat := range categories {
			publishedBadge := yesNo(cat.IsPublished)
			bg := "#ffffff"
			if i%2 == 0 {
				bg = "#faf9fa"
			}
			rows += fmt.Sprintf(`
<tr style="background:%s;transition:background 0.15s;" onmouseover="this.style.background='#f5f4f5'" onmouseout="this.style.background='%s'">
  <td style="padding:12px 16px;font-size:13px;font-weight:500;">%s</td>
  <td style="padding:12px 16px;font-size:12px;color:#6b6a6b;">%s</td>
  <td style="padding:12px 16px;">%s</td>
  <td style="padding:12px 16px;display:flex;gap:8px;">
    <a href="/admin/experience-categories/%s/edit" style="font-size:12px;font-weight:500;color:#1e1d1e;text-decoration:none;padding:4px 10px;background:#f0eff0;border-radius:4px;">Edit</a>
    <button onclick="confirmDelete('/admin/experience-categories/%s', '%s')" style="font-size:12px;font-weight:500;color:#b91c1c;background:none;border:none;cursor:pointer;padding:4px 10px;background:#fee2e2;border-radius:4px;">Delete</button>
  </td>
</tr>`, bg, bg, cat.Name, cat.Slug, publishedBadge, cat.ID, cat.ID, cat.Name)
		}

		empty := ""
		if len(categories) == 0 {
			empty = fmt.Sprintf(`<tr><td colspan="4" style="text-align:center;padding:40px;">%s</td></tr>`, emptyState("🏷️", "Belum ada kategori", "Tambahkan kategori untuk experience", "New Category", "/admin/experience-categories/new"))
		}

		content := fmt.Sprintf(`
%s
<div style="text-align:right;margin-bottom:24px;">
  <a href="/admin/experience-categories/new" style="display:inline-flex;align-items:center;gap:6px;padding:8px 18px;background:#1e1d1e;color:#f5f4f5;border-radius:6px;text-decoration:none;font-weight:500;font-size:13px;">
    + New Category
  </a>
</div>
%s
<div style="background:#ffffff;border-radius:10px;border:1px solid #e5e3e4;overflow:hidden;">
  <table style="width:100%%;border-collapse:collapse;">
    <thead><tr style="background:#f8f7f8;">
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Name</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Slug</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Status</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Actions</th>
    </tr></thead>
    <tbody>%s %s</tbody>
  </table>
</div>

<!-- Delete Modal -->
<div id="deleteModal" style="display:none;position:fixed;top:0;left:0;width:100%%;height:100%%;background:rgba(0,0,0,0.4);z-index:1000;align-items:center;justify-content:center;">
  <div style="background:#fff;border-radius:16px;padding:24px;max-width:400px;width:90%%;box-shadow:0 20px 60px rgba(0,0,0,0.2);">
    <h3 style="font-size:16px;font-weight:600;color:#1e1d1e;">Hapus Category</h3>
    <p style="font-size:13px;color:#6b6a6b;margin-top:8px;">Kamu yakin ingin menghapus category <strong id="deleteItemName"></strong>?</p>
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
  fetch(_deleteUrl, { method: 'DELETE' }).then(r => { if (r.ok) window.location.href = '/admin/experience-categories?flash=deleted'; });
});
</script>`, notifHTML, stats, rows, empty)

		return c.Type("html").SendString(layout("Experience Categories", "exp_cats", content))
	}
}

// ── Create / Edit Form ────────────────────────────────────────────────────────

func ExperienceCategoryFormPage(db *gorm.DB, fallbackID string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id", fallbackID)
		isEdit := false
		var category model.ExperienceCategory

		if id != "" {
			if err := db.First(&category, "id = ?", id).Error; err != nil {
				return c.Status(404).SendString("Category not found")
			}
			isEdit = true
		}

		title := "Create Experience Category"
		formAction := "/admin/experience-categories"
		submitLabel := "Create Category"
		if isEdit {
			title = "Edit Experience Category"
			formAction = "/admin/experience-categories/" + id
			submitLabel = "Save Changes"
		}

		content := fmt.Sprintf(`
<div style="display:flex; justify-content: flex-end; margin-bottom: 16px;">
  <a href="/admin/experience-categories" style="display:inline-flex;align-items:center;gap:6px;font-size:13px;color:#4d4c4d;text-decoration:none;padding:6px 12px;border:1px solid #e5e3e4;border-radius:6px;">← Back to Categories</a>
</div>
<form method="POST" action="%s">
  <div style="background:#fff;border-radius:10px;border:1px solid #e5e3e4;padding:24px;margin-bottom:16px;">
    <h2 style="font-size:14px;font-weight:600;color:#1e1d1e;">Informasi Kategori</h2>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:16px;">
      %s %s %s
    </div>
    <div style="margin-top:16px;">
      %s
    </div>
  </div>
  <div style="display:flex;justify-content:flex-end;gap:12px;margin-top:16px;padding:16px 0;">
    <a href="/admin/experience-categories" style="padding:8px 16px;border:1px solid #e5e3e4;border-radius:6px;color:#4d4c4d;text-decoration:none;">Batal</a>
    <button type="submit" style="padding:8px 20px;background:#1e1d1e;color:#f5f4f5;border:none;border-radius:6px;font-weight:500;">%s</button>
  </div>
</form>`,
			formAction,
			inputField("Name", "name", category.Name, "text", true),
			inputField("Slug", "slug", category.Slug, "text", true),
			inputField("Sort Order", "sort_order", fmt.Sprintf("%d", category.SortOrder), "number", false),
			toggleSwitch("Published", "is_published", category.IsPublished, "Kategori bisa digunakan"),
			submitLabel,
		)

		return c.Type("html").SendString(layout(title, "exp_cats", content))
	}
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func CreateExperienceCategoryHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sortOrder := 0
		if so := c.FormValue("sort_order"); so != "" {
			fmt.Sscanf(so, "%d", &sortOrder)
		}
		cat := model.ExperienceCategory{
			Name:        c.FormValue("name"),
			Slug:        c.FormValue("slug"),
			SortOrder:   sortOrder,
			IsPublished: c.FormValue("is_published") == "true",
		}
		if err := db.Create(&cat).Error; err != nil {
			return c.Status(500).SendString("Failed to create category")
		}
		return c.Redirect("/admin/experience-categories?flash=created")
	}
}

func UpdateExperienceCategoryHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var cat model.ExperienceCategory
		if err := db.First(&cat, "id = ?", id).Error; err != nil {
			return c.Status(404).SendString("Category not found")
		}
		cat.Name = c.FormValue("name")
		cat.Slug = c.FormValue("slug")
		if so := c.FormValue("sort_order"); so != "" {
			fmt.Sscanf(so, "%d", &cat.SortOrder)
		}
		cat.IsPublished = c.FormValue("is_published") == "true"
		db.Save(&cat)
		return c.Redirect("/admin/experience-categories?flash=updated")
	}
}

func DeleteExperienceCategoryHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		db.Delete(&model.ExperienceCategory{}, "id = ?", id)
		return c.SendStatus(200)
	}
}