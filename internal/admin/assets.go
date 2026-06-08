package admin

import (
	"fmt"
	"portfolio-cms/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ── List Assets ──────────────────────────────────────────────────────────────

func AssetsPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var assets []model.Asset
		// filter by owner_type and owner_id if provided
		query := db.Order("owner_type ASC, sort_order ASC, created_at DESC")
		ownerType := c.Query("owner_type")
		ownerID := c.Query("owner_id")
		if ownerType != "" {
			query = query.Where("owner_type = ?", ownerType)
		}
		if ownerID != "" {
			query = query.Where("owner_id = ?", ownerID)
		}
		query.Find(&assets)

		// Flash notification
		flash := c.Query("flash")
		notifHTML := ""
		if flash == "created" {
			notifHTML = notifBanner("success", "✓ Asset berhasil dibuat.")
		} else if flash == "updated" {
			notifHTML = notifBanner("success", "✓ Asset berhasil diperbarui.")
		} else if flash == "deleted" {
			notifHTML = notifBanner("info", "Asset telah dihapus.")
		}

		// Statistik
		total := len(assets)
		photoCount, videoCount, pdfCount, docCount := 0, 0, 0, 0
		for _, a := range assets {
			switch a.Type {
			case "photo":
				photoCount++
			case "video":
				videoCount++
			case "pdf":
				pdfCount++
			case "doc":
				docCount++
			}
		}
		stats := ""
		if total > 0 {
			stats = fmt.Sprintf(`
<div style="display:grid;grid-template-columns:repeat(4,1fr);gap:16px;margin-bottom:24px;">
  %s %s %s %s
</div>`,
				statCard("🖼️", "Photos", fmt.Sprintf("%d", photoCount), "#e0e7ff"),
				statCard("🎬", "Videos", fmt.Sprintf("%d", videoCount), "#d1fae5"),
				statCard("📄", "PDFs", fmt.Sprintf("%d", pdfCount), "#fef3c7"),
				statCard("📝", "Docs", fmt.Sprintf("%d", docCount), "#ede9fe"),
			)
		}

		// Filter form (mempertahankan nilai yang sudah dipilih)
		filterTypeSelected := map[string]string{"project": "", "experience": "", "": ""}
		if ownerType == "project" {
			filterTypeSelected["project"] = "selected"
		} else if ownerType == "experience" {
			filterTypeSelected["experience"] = "selected"
		} else {
			filterTypeSelected[""] = "selected"
		}

		// Table rows
		rows := ""
		for i, a := range assets {
			typeBadge := badge(string(a.Type), "blue")
			ownerLabel := fmt.Sprintf("%s/%s", a.OwnerType, a.OwnerID.String()[:8])
			bg := "#ffffff"
			if i%2 == 0 {
				bg = "#faf9fa"
			}
			titleDisplay := a.Title
			if titleDisplay == "" {
				titleDisplay = `<span style="color:#9a9899;">—</span>`
			}
			captionDisplay := a.Caption
			if captionDisplay == "" {
				captionDisplay = `<span style="color:#9a9899;">—</span>`
			}
			rows += fmt.Sprintf(`
<tr style="background:%s;transition:background 0.15s;" onmouseover="this.style.background='#f5f4f5'" onmouseout="this.style.background='%s'">
  <td style="padding:12px 16px;font-size:12px;">%s</td>
  <td style="padding:12px 16px;">%s</td>
  <td style="padding:12px 16px;max-width:200px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;font-size:12px;" title="%s">%s</td>
  <td style="padding:12px 16px;font-size:13px;">%s</td>
  <td style="padding:12px 16px;font-size:12px;">%s</td>
  <td style="padding:12px 16px;display:flex;gap:8px;">
    <a href="/admin/assets/%s/edit" style="font-size:12px;font-weight:500;color:#1e1d1e;text-decoration:none;padding:4px 10px;background:#f0eff0;border-radius:4px;">Edit</a>
    <button onclick="confirmDelete('/admin/assets/%s', 'Asset #%s')" style="font-size:12px;font-weight:500;color:#b91c1c;background:none;border:none;cursor:pointer;padding:4px 10px;background:#fee2e2;border-radius:4px;">Delete</button>
  </td>
</tr>`, bg, bg, ownerLabel, typeBadge, a.URL, a.URL, titleDisplay, captionDisplay, a.ID, a.ID, a.ID.String()[:8])
		}

		empty := ""
		if len(assets) == 0 {
			empty = fmt.Sprintf(`<tr><td colspan="6" style="text-align:center;padding:40px;">%s</td></tr>`, emptyState("🖼️", "Belum ada asset", "Tambahkan foto, video, atau dokumen", "New Asset", "/admin/assets/new"))
		}

		content := fmt.Sprintf(`
%s
<div style="display:flex; justify-content: space-between; align-items: center; margin-bottom: 24px;">
  <form method="GET" style="display:flex; gap:8px; align-items: center;">
    <select name="owner_type" style="border:1px solid #d1cfd0; border-radius:6px; padding:6px 10px; font-size:13px;">
      <option value="" %s>All types</option>
      <option value="project" %s>Project</option>
      <option value="experience" %s>Experience</option>
    </select>
    <button type="submit" style="padding:6px 12px; background:#f0eff0; border:1px solid #d1cfd0; border-radius:6px; cursor:pointer; font-size:13px;">Filter</button>
  </form>
  <a href="/admin/assets/new" style="display:inline-flex;align-items:center;gap:6px;padding:8px 18px;background:#1e1d1e;color:#f5f4f5;border-radius:6px;text-decoration:none;font-weight:500;font-size:13px;">
    + New Asset
  </a>
</div>
%s
<div style="background:#ffffff;border-radius:10px;border:1px solid #e5e3e4;overflow:hidden;">
  <table style="width:100%%;border-collapse:collapse;">
    <thead><tr style="background:#f8f7f8;">
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Owner</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Type</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">URL</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Title</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Caption</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Actions</th>
    </tr></thead>
    <tbody>%s %s</tbody>
  </table>
</div>

<!-- Delete Modal -->
<div id="deleteModal" style="display:none;position:fixed;top:0;left:0;width:100%%;height:100%%;background:rgba(0,0,0,0.4);z-index:1000;align-items:center;justify-content:center;">
  <div style="background:#fff;border-radius:16px;padding:24px;max-width:400px;width:90%%;box-shadow:0 20px 60px rgba(0,0,0,0.2);">
    <h3 style="font-size:16px;font-weight:600;color:#1e1d1e;">Hapus Asset</h3>
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
  fetch(_deleteUrl, { method: 'DELETE' }).then(r => { if (r.ok) window.location.href = '/admin/assets?flash=deleted'; });
});
</script>`,
			notifHTML,
			filterTypeSelected[""], filterTypeSelected["project"], filterTypeSelected["experience"],
			stats, rows, empty)

		return c.Type("html").SendString(layout("Assets", "assets", content))
	}
}

// ── Create / Edit Form ────────────────────────────────────────────────────────

func AssetFormPage(db *gorm.DB, fallbackID string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id", fallbackID)
		isEdit := false
		var asset model.Asset

		if id != "" {
			if err := db.First(&asset, "id = ?", id).Error; err != nil {
				return c.Status(404).SendString("Asset not found")
			}
			isEdit = true
		}

		// Owner type dropdown
		ownerTypeSelected := map[string]string{"project": "", "experience": ""}
		if asset.OwnerType == "project" {
			ownerTypeSelected["project"] = "selected"
		} else if asset.OwnerType == "experience" {
			ownerTypeSelected["experience"] = "selected"
		}

		title := "Create Asset"
		formAction := "/admin/assets"
		submitLabel := "Create Asset"
		if isEdit {
			title = "Edit Asset"
			formAction = "/admin/assets/" + id
			submitLabel = "Save Changes"
		}

		content := fmt.Sprintf(`
<div style="display:flex; justify-content: flex-end; margin-bottom: 16px;">
  <a href="/admin/assets" style="display:inline-flex;align-items:center;gap:6px;font-size:13px;color:#4d4c4d;text-decoration:none;padding:6px 12px;border:1px solid #e5e3e4;border-radius:6px;">← Back to Assets</a>
</div>
<form method="POST" action="%s">
  <div style="background:#fff;border-radius:10px;border:1px solid #e5e3e4;padding:24px;margin-bottom:16px;">
    <h2 style="font-size:14px;font-weight:600;color:#1e1d1e;">Informasi Asset</h2>
    <div style="margin-top:16px;">
      <label style="font-size:12px;font-weight:600;color:#4d4c4d;">Owner Type</label>
      <select name="owner_type" style="width:100%%;margin-top:5px;border:1px solid #d1cfd0;border-radius:6px;padding:8px;font-size:13px;">
        <option value="project" %s>Project</option>
        <option value="experience" %s>Experience</option>
      </select>
    </div>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:16px;">
      %s %s
    </div>
    <div style="margin-top:16px;">
      %s
    </div>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:16px;">
      %s %s
    </div>
    <div style="margin-top:16px;">
      %s
    </div>
  </div>
  <div style="display:flex;justify-content:flex-end;gap:12px;margin-top:16px;padding:16px 0;">
    <a href="/admin/assets" style="padding:8px 16px;border:1px solid #e5e3e4;border-radius:6px;color:#4d4c4d;text-decoration:none;">Batal</a>
    <button type="submit" style="padding:8px 20px;background:#1e1d1e;color:#f5f4f5;border:none;border-radius:6px;font-weight:500;">%s</button>
  </div>
</form>`,
			formAction,
			ownerTypeSelected["project"], ownerTypeSelected["experience"],
			inputField("Owner ID", "owner_id", asset.OwnerID.String(), "text", true),
			assetTypeSelect(asset.Type),
			inputField("URL", "url", asset.URL, "url", true),
			inputField("Title", "title", asset.Title, "text", false),
			inputField("Caption", "caption", asset.Caption, "text", false),
			inputField("Sort Order", "sort_order", fmt.Sprintf("%d", asset.SortOrder), "number", false),
			submitLabel,
		)

		return c.Type("html").SendString(layout(title, "assets", content))
	}
}

func assetTypeSelect(current model.AssetType) string {
	types := []string{"photo", "video", "pdf", "doc"}
	opts := ""
	for _, t := range types {
		sel := ""
		if string(current) == t {
			sel = "selected"
		}
		opts += fmt.Sprintf(`<option value="%s" %s>%s</option>`, t, sel, t)
	}
	return fmt.Sprintf(`
<div>
  <label style="font-size:12px;font-weight:600;color:#4d4c4d;">Type</label>
  <select name="type" style="width:100%%;margin-top:5px;border:1px solid #d1cfd0;border-radius:6px;padding:8px;font-size:13px;">%s</select>
</div>`, opts)
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func CreateAssetHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ownerID, err := uuid.Parse(c.FormValue("owner_id"))
		if err != nil {
			return c.Status(400).SendString("Invalid owner ID")
		}
		sortOrder := 0
		if so := c.FormValue("sort_order"); so != "" {
			fmt.Sscanf(so, "%d", &sortOrder)
		}
		asset := model.Asset{
			OwnerType: c.FormValue("owner_type"),
			OwnerID:   ownerID,
			Type:      model.AssetType(c.FormValue("type")),
			URL:       c.FormValue("url"),
			Title:     c.FormValue("title"),
			Caption:   c.FormValue("caption"),
			SortOrder: sortOrder,
		}
		if err := db.Create(&asset).Error; err != nil {
			return c.Status(500).SendString("Failed to create asset")
		}
		return c.Redirect("/admin/assets?flash=created")
	}
}

func UpdateAssetHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var asset model.Asset
		if err := db.First(&asset, "id = ?", id).Error; err != nil {
			return c.Status(404).SendString("Asset not found")
		}
		ownerID, err := uuid.Parse(c.FormValue("owner_id"))
		if err != nil {
			return c.Status(400).SendString("Invalid owner ID")
		}
		asset.OwnerType = c.FormValue("owner_type")
		asset.OwnerID = ownerID
		asset.Type = model.AssetType(c.FormValue("type"))
		asset.URL = c.FormValue("url")
		asset.Title = c.FormValue("title")
		asset.Caption = c.FormValue("caption")
		if so := c.FormValue("sort_order"); so != "" {
			fmt.Sscanf(so, "%d", &asset.SortOrder)
		}
		db.Save(&asset)
		return c.Redirect("/admin/assets?flash=updated")
	}
}

func DeleteAssetHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		db.Delete(&model.Asset{}, "id = ?", id)
		return c.SendStatus(200)
	}
}