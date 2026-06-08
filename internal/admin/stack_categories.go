package admin

import (
	"fmt"
	"portfolio-cms/internal/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ── List Stack Categories ─────────────────────────────────────────────────────

func StackCategoriesPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var categories []model.StackCategory
		db.Preload("Items").Order("sort_order ASC, name ASC").Find(&categories)

		// Flash notification
		flash := c.Query("flash")
		notifHTML := ""
		if flash == "created" {
			notifHTML = notifBanner("success", "✓ Stack category berhasil dibuat.")
		} else if flash == "updated" {
			notifHTML = notifBanner("success", "✓ Stack category berhasil diperbarui.")
		} else if flash == "deleted" {
			notifHTML = notifBanner("info", "Stack category telah dihapus.")
		}

		// Statistik
		total := len(categories)
		totalItems := 0
		for _, c := range categories {
			totalItems += len(c.Items)
		}
		stats := ""
		if total > 0 {
			stats = fmt.Sprintf(`
<div style="display:grid;grid-template-columns:repeat(3,1fr);gap:16px;margin-bottom:24px;">
  %s %s %s
</div>`,
				statCard("🗂️", "Categories", fmt.Sprintf("%d", total), "#e0e7ff"),
				statCard("📌", "Total Items", fmt.Sprintf("%d", totalItems), "#d1fae5"),
				statCard("📊", "Avg Items/Cat", fmt.Sprintf("%.1f", float64(totalItems)/float64(total)), "#fef3c7"),
			)
		}

		// Table rows
		rows := ""
		for i, cat := range categories {
			itemCount := len(cat.Items)
			colorIndicator := ""
			if cat.Color != "" {
				colorIndicator = fmt.Sprintf(`<span style="display:inline-block;width:12px;height:12px;border-radius:3px;background:%s;margin-right:6px;"></span>`, cat.Color)
			}
			bg := "#ffffff"
			if i%2 == 0 {
				bg = "#faf9fa"
			}
			rows += fmt.Sprintf(`
<tr style="background:%s;transition:background 0.15s;" onmouseover="this.style.background='#f5f4f5'" onmouseout="this.style.background='%s'">
  <td style="padding:12px 16px;font-size:13px;font-weight:500;">%s</td>
  <td style="padding:12px 16px;font-size:12px;color:#6b6a6b;">%s</td>
  <td style="padding:12px 16px;text-align:center;">%d</td>
  <td style="padding:12px 16px;">%s%s</td>
  <td style="padding:12px 16px;">%d</td>
  <td style="padding:12px 16px;display:flex;gap:8px;">
    <a href="/admin/stack-categories/%s/edit" style="font-size:12px;font-weight:500;color:#1e1d1e;text-decoration:none;padding:4px 10px;background:#f0eff0;border-radius:4px;">Edit</a>
    <button onclick="confirmDelete('/admin/stack-categories/%s', '%s')" style="font-size:12px;font-weight:500;color:#b91c1c;background:none;border:none;cursor:pointer;padding:4px 10px;background:#fee2e2;border-radius:4px;">Delete</button>
  </td>
</tr>`, bg, bg, cat.Name, cat.Slug, itemCount, colorIndicator, cat.Color, cat.SortOrder, cat.ID, cat.ID, cat.Name)
		}

		empty := ""
		if len(categories) == 0 {
			empty = fmt.Sprintf(`<tr><td colspan="6" style="text-align:center;padding:40px;">%s</td></tr>`, emptyState("🗂️", "Belum ada stack category", "Tambahkan kategori stack pertama kamu", "New Category", "/admin/stack-categories/new"))
		}

		content := fmt.Sprintf(`
%s
<div style="text-align:right;margin-bottom:24px;">
  <a href="/admin/stack-categories/new" style="display:inline-flex;align-items:center;gap:6px;padding:8px 18px;background:#1e1d1e;color:#f5f4f5;border-radius:6px;text-decoration:none;font-weight:500;font-size:13px;">
    + New Category
  </a>
</div>
%s
<div style="background:#ffffff;border-radius:10px;border:1px solid #e5e3e4;overflow:hidden;">
  <table style="width:100%%;border-collapse:collapse;">
    <thead><tr style="background:#f8f7f8;">
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Name</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Slug</th>
      <th style="padding:12px 16px;text-align:center;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Items</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Color</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Sort</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Actions</th>
    </tr></thead>
    <tbody>%s %s</tbody>
  </table>
</div>

<!-- Delete Modal -->
<div id="deleteModal" style="display:none;position:fixed;top:0;left:0;width:100%%;height:100%%;background:rgba(0,0,0,0.4);z-index:1000;align-items:center;justify-content:center;">
  <div style="background:#fff;border-radius:16px;padding:24px;max-width:400px;width:90%%;box-shadow:0 20px 60px rgba(0,0,0,0.2);">
    <h3 style="font-size:16px;font-weight:600;color:#1e1d1e;">Hapus Stack Category</h3>
    <p style="font-size:13px;color:#6b6a6b;margin-top:8px;">Kamu yakin ingin menghapus <strong id="deleteItemName"></strong>? Semua item di dalamnya juga akan terhapus.</p>
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
  fetch(_deleteUrl, { method: 'DELETE' }).then(r => { if (r.ok) window.location.href = '/admin/stack-categories?flash=deleted'; });
});
</script>`, notifHTML, stats, rows, empty)

		return c.Type("html").SendString(layout("Stack Categories", "stack_cats", content))
	}
}

// ── Create / Edit Form ────────────────────────────────────────────────────────

func StackCategoryFormPage(db *gorm.DB, fallbackID string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id", fallbackID)
		isEdit := false
		var category model.StackCategory

		if id != "" {
			if err := db.First(&category, "id = ?", id).Error; err != nil {
				return c.Status(404).SendString("Category not found")
			}
			isEdit = true
		}

		title := "Create Stack Category"
		formAction := "/admin/stack-categories"
		submitLabel := "Create Category"
		if isEdit {
			title = "Edit Stack Category"
			formAction = "/admin/stack-categories/" + id
			submitLabel = "Save Changes"
		}

		content := fmt.Sprintf(`
<div style="display:flex; justify-content: flex-end; margin-bottom: 16px;">
  <a href="/admin/stack-categories" style="display:inline-flex;align-items:center;gap:6px;font-size:13px;color:#4d4c4d;text-decoration:none;padding:6px 12px;border:1px solid #e5e3e4;border-radius:6px;">← Back to Categories</a>
</div>
<form method="POST" action="%s">
  <div style="background:#fff;border-radius:10px;border:1px solid #e5e3e4;padding:24px;margin-bottom:16px;">
    <h2 style="font-size:14px;font-weight:600;color:#1e1d1e;">Informasi Kategori</h2>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:16px;">
      %s %s
    </div>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:16px;">
      %s %s
    </div>
    <div style="margin-top:16px;">
      %s
    </div>
  </div>
  <div style="display:flex;justify-content:flex-end;gap:12px;margin-top:16px;padding:16px 0;">
    <a href="/admin/stack-categories" style="padding:8px 16px;border:1px solid #e5e3e4;border-radius:6px;color:#4d4c4d;text-decoration:none;">Batal</a>
    <button type="submit" style="padding:8px 20px;background:#1e1d1e;color:#f5f4f5;border:none;border-radius:6px;font-weight:500;">%s</button>
  </div>
</form>`,
			formAction,
			inputField("Name", "name", category.Name, "text", true),
			inputField("Slug", "slug", category.Slug, "text", true),
			inputField("Icon URL", "icon_url", category.IconURL, "url", false),
			inputField("Color (hex)", "color", category.Color, "text", false),
			inputField("Sort Order", "sort_order", fmt.Sprintf("%d", category.SortOrder), "number", false),
			submitLabel,
		)

		return c.Type("html").SendString(layout(title, "stack_cats", content))
	}
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func CreateStackCategoryHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sortOrder := 0
		if so := c.FormValue("sort_order"); so != "" {
			fmt.Sscanf(so, "%d", &sortOrder)
		}
		cat := model.StackCategory{
			Name:      c.FormValue("name"),
			Slug:      c.FormValue("slug"),
			IconURL:   c.FormValue("icon_url"),
			Color:     c.FormValue("color"),
			SortOrder: sortOrder,
		}
		if err := db.Create(&cat).Error; err != nil {
			return c.Status(500).SendString("Failed to create category")
		}
		return c.Redirect("/admin/stack-categories?flash=created")
	}
}

func UpdateStackCategoryHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var cat model.StackCategory
		if err := db.First(&cat, "id = ?", id).Error; err != nil {
			return c.Status(404).SendString("Category not found")
		}
		cat.Name = c.FormValue("name")
		cat.Slug = c.FormValue("slug")
		cat.IconURL = c.FormValue("icon_url")
		cat.Color = c.FormValue("color")
		if so := c.FormValue("sort_order"); so != "" {
			fmt.Sscanf(so, "%d", &cat.SortOrder)
		}
		db.Save(&cat)
		return c.Redirect("/admin/stack-categories?flash=updated")
	}
}

func DeleteStackCategoryHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		db.Delete(&model.StackCategory{}, "id = ?", id)
		return c.SendStatus(200)
	}
}