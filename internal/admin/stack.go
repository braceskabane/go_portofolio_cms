package admin

import (
	"fmt"
	"portfolio-cms/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ── List Stack Items ─────────────────────────────────────────────────────────

func StackItemsPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var items []model.StackItem
		db.Preload("Category").Order("sort_order ASC, name ASC").Find(&items)

		// Flash notification
		flash := c.Query("flash")
		notifHTML := ""
		if flash == "created" {
			notifHTML = notifBanner("success", "✓ Stack item berhasil dibuat.")
		} else if flash == "updated" {
			notifHTML = notifBanner("success", "✓ Stack item berhasil diperbarui.")
		} else if flash == "deleted" {
			notifHTML = notifBanner("info", "Stack item telah dihapus.")
		}

		// Statistik
		total := len(items)
		pubCount := 0
		for _, it := range items {
			if it.IsPublished {
				pubCount++
			}
		}
		stats := ""
		if total > 0 {
			stats = fmt.Sprintf(`
<div style="display:grid;grid-template-columns:repeat(3,1fr);gap:16px;margin-bottom:24px;">
  %s %s %s
</div>`,
				statCard("📌", "Total", fmt.Sprintf("%d", total), "#e0e7ff"),
				statCard("✓", "Published", fmt.Sprintf("%d", pubCount), "#d1fae5"),
				statCard("✗", "Draft", fmt.Sprintf("%d", total-pubCount), "#fef3c7"),
			)
		}

		// Table rows
		rows := ""
		for i, item := range items {
			catName := "—"
			if item.Category.Name != "" {
				catName = item.Category.Name
			}
			publishedBadge := yesNo(item.IsPublished)
			bg := "#ffffff"
			if i%2 == 0 {
				bg = "#faf9fa"
			}
			rows += fmt.Sprintf(`
<tr style="background:%s;transition:background 0.15s;" onmouseover="this.style.background='#f5f4f5'" onmouseout="this.style.background='%s'">
  <td style="padding:12px 16px;font-size:13px;font-weight:500;">%s</td>
  <td style="padding:12px 16px;font-size:12px;color:#6b6a6b;">%s</td>
  <td style="padding:12px 16px;font-size:13px;">%s</td>
  <td style="padding:12px 16px;">%s</td>
  <td style="padding:12px 16px;display:flex;gap:8px;">
    <a href="/admin/stack-items/%s/edit" style="font-size:12px;font-weight:500;color:#1e1d1e;text-decoration:none;padding:4px 10px;background:#f0eff0;border-radius:4px;">Edit</a>
    <button onclick="confirmDelete('/admin/stack-items/%s', '%s')" style="font-size:12px;font-weight:500;color:#b91c1c;background:none;border:none;cursor:pointer;padding:4px 10px;background:#fee2e2;border-radius:4px;">Delete</button>
  </td>
</tr>`, bg, bg, item.Name, item.Slug, catName, publishedBadge, item.ID, item.ID, item.Name)
		}

		empty := ""
		if len(items) == 0 {
			empty = fmt.Sprintf(`<tr><td colspan="5" style="text-align:center;padding:40px;">%s</td></tr>`, emptyState("📌", "Belum ada stack item", "Tambahkan stack item pertama kamu", "New Stack Item", "/admin/stack-items/new"))
		}

		content := fmt.Sprintf(`
%s
<div style="text-align:right;margin-bottom:24px;">
  <a href="/admin/stack-items/new" style="display:inline-flex;align-items:center;gap:6px;padding:8px 18px;background:#1e1d1e;color:#f5f4f5;border-radius:6px;text-decoration:none;font-weight:500;font-size:13px;">
    + New Stack Item
  </a>
</div>
%s
<div style="background:#ffffff;border-radius:10px;border:1px solid #e5e3e4;overflow:hidden;">
  <table style="width:100%%;border-collapse:collapse;">
    <thead><tr style="background:#f8f7f8;">
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Name</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Slug</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Category</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Status</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Actions</th>
    </tr></thead>
    <tbody>%s %s</tbody>
  </table>
</div>

<!-- Delete Modal -->
<div id="deleteModal" style="display:none;position:fixed;top:0;left:0;width:100%%;height:100%%;background:rgba(0,0,0,0.4);z-index:1000;align-items:center;justify-content:center;">
  <div style="background:#fff;border-radius:16px;padding:24px;max-width:400px;width:90%%;box-shadow:0 20px 60px rgba(0,0,0,0.2);">
    <h3 style="font-size:16px;font-weight:600;color:#1e1d1e;">Hapus Stack Item</h3>
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
  fetch(_deleteUrl, { method: 'DELETE' }).then(r => { if (r.ok) window.location.href = '/admin/stack-items?flash=deleted'; });
});
</script>`, notifHTML, stats, rows, empty)

		return c.Type("html").SendString(layout("Stack Items", "stack_items", content))
	}
}

// ── Create / Edit Form ────────────────────────────────────────────────────────

func StackItemFormPage(db *gorm.DB, fallbackID string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id", fallbackID)
		isEdit := false
		var item model.StackItem

		if id != "" {
			if err := db.Preload("Category").First(&item, "id = ?", id).Error; err != nil {
				return c.Status(404).SendString("Stack item not found")
			}
			isEdit = true
		}

		// Category dropdown
		var categories []model.StackCategory
		db.Order("sort_order ASC, name ASC").Find(&categories)
		catOptions := `<option value="">— Pilih Category —</option>`
		for _, cat := range categories {
			sel := ""
			if item.CategoryID == cat.ID {
				sel = "selected"
			}
			catOptions += fmt.Sprintf(`<option value="%s" %s>%s</option>`, cat.ID, sel, cat.Name)
		}

		title := "Create Stack Item"
		formAction := "/admin/stack-items"
		submitLabel := "Create Stack Item"
		if isEdit {
			title = "Edit Stack Item"
			formAction = "/admin/stack-items/" + id
			submitLabel = "Save Changes"
		}

		content := fmt.Sprintf(`
<div style="display:flex; justify-content: flex-end; margin-bottom: 16px;">
  <a href="/admin/stack-items" style="display:inline-flex;align-items:center;gap:6px;font-size:13px;color:#4d4c4d;text-decoration:none;padding:6px 12px;border:1px solid #e5e3e4;border-radius:6px;">← Back to Stack Items</a>
</div>
<form method="POST" action="%s">
  <div style="background:#fff;border-radius:10px;border:1px solid #e5e3e4;padding:24px;margin-bottom:16px;">
    <h2 style="font-size:14px;font-weight:600;color:#1e1d1e;">Informasi Stack Item</h2>
    <div style="margin-top:16px;">
      <label style="font-size:12px;font-weight:600;color:#4d4c4d;">Category</label>
      <select name="category_id" style="width:100%%;margin-top:5px;border:1px solid #d1cfd0;border-radius:6px;padding:8px;font-size:13px;">%s</select>
    </div>
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
    <a href="/admin/stack-items" style="padding:8px 16px;border:1px solid #e5e3e4;border-radius:6px;color:#4d4c4d;text-decoration:none;">Batal</a>
    <button type="submit" style="padding:8px 20px;background:#1e1d1e;color:#f5f4f5;border:none;border-radius:6px;font-weight:500;">%s</button>
  </div>
</form>`,
			formAction,
			catOptions,
			inputField("Name", "name", item.Name, "text", true),
			inputField("Slug", "slug", item.Slug, "text", true),
			inputField("Icon URL", "icon_url", item.IconURL, "url", false),
			inputField("Sort Order", "sort_order", fmt.Sprintf("%d", item.SortOrder), "number", false),
			toggleSwitch("Published", "is_published", item.IsPublished, "Item bisa digunakan"),
			submitLabel,
		)

		return c.Type("html").SendString(layout(title, "stack_items", content))
	}
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func CreateStackItemHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		categoryID, err := uuid.Parse(c.FormValue("category_id"))
		if err != nil {
			return c.Status(400).SendString("Invalid category ID")
		}
		sortOrder := 0
		if so := c.FormValue("sort_order"); so != "" {
			fmt.Sscanf(so, "%d", &sortOrder)
		}
		item := model.StackItem{
			CategoryID:  categoryID,
			Name:        c.FormValue("name"),
			Slug:        c.FormValue("slug"),
			IconURL:     c.FormValue("icon_url"),
			SortOrder:   sortOrder,
			IsPublished: c.FormValue("is_published") == "true",
		}
		if err := db.Create(&item).Error; err != nil {
			return c.Status(500).SendString("Failed to create stack item")
		}
		return c.Redirect("/admin/stack-items?flash=created")
	}
}

func UpdateStackItemHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var item model.StackItem
		if err := db.First(&item, "id = ?", id).Error; err != nil {
			return c.Status(404).SendString("Stack item not found")
		}
		categoryID, err := uuid.Parse(c.FormValue("category_id"))
		if err != nil {
			return c.Status(400).SendString("Invalid category ID")
		}
		item.CategoryID = categoryID
		item.Name = c.FormValue("name")
		item.Slug = c.FormValue("slug")
		item.IconURL = c.FormValue("icon_url")
		if so := c.FormValue("sort_order"); so != "" {
			fmt.Sscanf(so, "%d", &item.SortOrder)
		}
		item.IsPublished = c.FormValue("is_published") == "true"
		db.Save(&item)
		return c.Redirect("/admin/stack-items?flash=updated")
	}
}

func DeleteStackItemHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		db.Delete(&model.StackItem{}, "id = ?", id)
		return c.SendStatus(200)
	}
}