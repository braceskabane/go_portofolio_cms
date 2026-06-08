package admin

import (
	"fmt"
	"portfolio-cms/internal/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ── List Contacts ────────────────────────────────────────────────────────────

func ContactsPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var contacts []model.Contact
		db.Order("created_at DESC").Find(&contacts)

		// Flash notification
		flash := c.Query("flash")
		notifHTML := ""
		if flash == "read" {
			notifHTML = notifBanner("success", "✓ Pesan ditandai sudah dibaca.")
		} else if flash == "deleted" {
			notifHTML = notifBanner("info", "Pesan telah dihapus.")
		}

		// Statistik
		total := len(contacts)
		unreadCount := 0
		for _, msg := range contacts {
			if !msg.IsRead {
				unreadCount++
			}
		}
		stats := ""
		if total > 0 {
			stats = fmt.Sprintf(`
<div style="display:grid;grid-template-columns:repeat(3,1fr);gap:16px;margin-bottom:24px;">
  %s %s %s
</div>`,
				statCard("✉️", "Total", fmt.Sprintf("%d", total), "#e0e7ff"),
				statCard("📬", "Unread", fmt.Sprintf("%d", unreadCount), "#fee2e2"),
				statCard("✓", "Read", fmt.Sprintf("%d", total-unreadCount), "#d1fae5"),
			)
		}

		// Table rows
		rows := ""
		for i, msg := range contacts {
			readBadge := badge("Unread", "red")
			if msg.IsRead {
				readBadge = badge("Read", "gray")
			}
			bg := "#ffffff"
			if i%2 == 0 {
				bg = "#faf9fa"
			}
			// Mark read button
			actionButton := ""
			if !msg.IsRead {
				actionButton = fmt.Sprintf(`<button onclick="markRead('/admin/contacts/%s/read')" style="font-size:12px;font-weight:500;color:#1e1d1e;background:#f0eff0;border:none;cursor:pointer;padding:4px 10px;border-radius:4px;">Mark Read</button>`, msg.ID)
			} else {
				actionButton = `<span style="font-size:12px;color:#9a9899;padding:4px 10px;">✓ Read</span>`
			}
			rows += fmt.Sprintf(`
<tr style="background:%s;transition:background 0.15s;" onmouseover="this.style.background='#f5f4f5'" onmouseout="this.style.background='%s'">
  <td style="padding:12px 16px;font-size:13px;font-weight:500;">%s</td>
  <td style="padding:12px 16px;font-size:13px;">%s</td>
  <td style="padding:12px 16px;font-size:13px;">%s</td>
  <td style="padding:12px 16px;">%s</td>
  <td style="padding:12px 16px;font-size:12px;">%s</td>
  <td style="padding:12px 16px;display:flex;gap:8px;">
    %s
    <button onclick="confirmDelete('/admin/contacts/%s', '%s')" style="font-size:12px;font-weight:500;color:#b91c1c;background:none;border:none;cursor:pointer;padding:4px 10px;background:#fee2e2;border-radius:4px;">Delete</button>
  </td>
</tr>`, bg, bg, msg.Name, msg.Email, msg.Subject, readBadge, msg.CreatedAt.Format("2006-01-02 15:04"), actionButton, msg.ID, msg.Name)
		}

		empty := ""
		if len(contacts) == 0 {
			empty = fmt.Sprintf(`<tr><td colspan="6" style="text-align:center;padding:40px;">%s</td></tr>`, emptyState("✉️", "Belum ada pesan", "Pesan dari pengunjung akan muncul di sini", "", ""))
		}

		content := fmt.Sprintf(`
%s
%s
<div style="background:#ffffff;border-radius:10px;border:1px solid #e5e3e4;overflow:hidden;">
  <table style="width:100%%;border-collapse:collapse;">
    <thead><tr style="background:#f8f7f8;">
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Name</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Email</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Subject</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Status</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Date</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Actions</th>
    </tr></thead>
    <tbody>%s %s</tbody>
  </table>
</div>

<!-- Delete Modal -->
<div id="deleteModal" style="display:none;position:fixed;top:0;left:0;width:100%%;height:100%%;background:rgba(0,0,0,0.4);z-index:1000;align-items:center;justify-content:center;">
  <div style="background:#fff;border-radius:16px;padding:24px;max-width:400px;width:90%%;box-shadow:0 20px 60px rgba(0,0,0,0.2);">
    <h3 style="font-size:16px;font-weight:600;color:#1e1d1e;">Hapus Pesan</h3>
    <p style="font-size:13px;color:#6b6a6b;margin-top:8px;">Kamu yakin ingin menghapus pesan dari <strong id="deleteItemName"></strong>?</p>
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
  fetch(_deleteUrl, { method: 'DELETE' }).then(r => { if (r.ok) window.location.href = '/admin/contacts?flash=deleted'; });
});
function markRead(url) {
  fetch(url, { method: 'POST' }).then(r => { if (r.ok) window.location.href = '/admin/contacts?flash=read'; });
}
</script>`, notifHTML, stats, rows, empty)

		return c.Type("html").SendString(layout("Contact Messages", "Contacts", content))
	}
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func MarkContactReadHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		db.Model(&model.Contact{}).Where("id = ?", id).Update("is_read", true)
		return c.Redirect("/admin/contacts?flash=read")
	}
}

func DeleteContactHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		db.Delete(&model.Contact{}, "id = ?", id)
		return c.SendStatus(200)
	}
}