package admin

import (
	"fmt"
	"portfolio-cms/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ── List Running Activities ──────────────────────────────────────────────────

func RunningActivitiesPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var activities []model.RunningActivity
		db.Order("date DESC").Limit(100).Find(&activities)

		// Flash notification
		flash := c.Query("flash")
		notifHTML := ""
		if flash == "created" {
			notifHTML = notifBanner("success", "✓ Aktivitas berhasil dicatat.")
		} else if flash == "updated" {
			notifHTML = notifBanner("success", "✓ Aktivitas berhasil diperbarui.")
		} else if flash == "deleted" {
			notifHTML = notifBanner("info", "Aktivitas telah dihapus.")
		}

		// Statistik
		total := len(activities)
		pubCount := 0
		totalDistance := 0.0
		totalDuration := 0
		for _, a := range activities {
			if a.IsPublished {
				pubCount++
			}
			totalDistance += a.DistanceMeters
			totalDuration += a.DurationSec
		}
		stats := ""
		if total > 0 {
			stats = fmt.Sprintf(`
<div style="display:grid;grid-template-columns:repeat(4,1fr);gap:16px;margin-bottom:24px;">
  %s %s %s %s
</div>`,
				statCard("🏃", "Total Runs", fmt.Sprintf("%d", total), "#e0e7ff"),
				statCard("📏", "Total Distance", fmt.Sprintf("%.1f km", totalDistance/1000), "#d1fae5"),
				statCard("⏱️", "Total Time", fmt.Sprintf("%dh %dm", totalDuration/3600, (totalDuration%3600)/60), "#fef3c7"),
				statCard("✓", "Published", fmt.Sprintf("%d", pubCount), "#ede9fe"),
			)
		}

		// Table rows
		rows := ""
		for i, a := range activities {
			distanceKm := a.DistanceMeters / 1000
			durationMin := a.DurationSec / 60
			paceMin := a.AvgPaceSec / 60
			paceSec := a.AvgPaceSec % 60
			pace := fmt.Sprintf("%d:%02d /km", paceMin, paceSec)
			publishedBadge := yesNo(a.IsPublished)
			bg := "#ffffff"
			if i%2 == 0 {
				bg = "#faf9fa"
			}

			rows += fmt.Sprintf(`
<tr style="background:%s;transition:background 0.15s;" onmouseover="this.style.background='#f5f4f5'" onmouseout="this.style.background='%s'">
  <td style="padding:12px 16px;font-size:13px;">%s</td>
  <td style="padding:12px 16px;font-size:13px;font-weight:500;">%s</td>
  <td style="padding:12px 16px;font-size:13px;">%.2f km</td>
  <td style="padding:12px 16px;font-size:13px;">%d min</td>
  <td style="padding:12px 16px;font-size:13px;">%s</td>
  <td style="padding:12px 16px;font-size:13px;">%d bpm</td>
  <td style="padding:12px 16px;">%s</td>
  <td style="padding:12px 16px;display:flex;gap:8px;">
    <a href="/admin/running-activities/%s/edit" style="font-size:12px;font-weight:500;color:#1e1d1e;text-decoration:none;padding:4px 10px;background:#f0eff0;border-radius:4px;">Edit</a>
    <button onclick="confirmDelete('/admin/running-activities/%s', '%s')" style="font-size:12px;font-weight:500;color:#b91c1c;background:none;border:none;cursor:pointer;padding:4px 10px;background:#fee2e2;border-radius:4px;">Delete</button>
  </td>
</tr>`,
				bg, bg,
				a.Date.Format("2006-01-02"),
				a.Title,
				distanceKm,
				durationMin,
				pace,
				a.AvgHeartRate,
				publishedBadge,
				a.ID, a.ID, a.Date.Format("2006-01-02")+" "+a.Title,
			)
		}

		empty := ""
		if len(activities) == 0 {
			empty = fmt.Sprintf(`<tr><td colspan="8" style="text-align:center;padding:40px;">%s</td></tr>`, emptyState("🏃", "Belum ada aktivitas lari", "Catat lari pertama kamu", "Log New Run", "/admin/running-activities/new"))
		}

		content := fmt.Sprintf(`
%s
<div style="text-align:right;margin-bottom:24px;">
  <a href="/admin/running-activities/new" style="display:inline-flex;align-items:center;gap:6px;padding:8px 18px;background:#1e1d1e;color:#f5f4f5;border-radius:6px;text-decoration:none;font-weight:500;font-size:13px;">
    + New Activity
  </a>
</div>
%s
<div style="background:#ffffff;border-radius:10px;border:1px solid #e5e3e4;overflow:hidden;">
  <table style="width:100%%;border-collapse:collapse;">
    <thead><tr style="background:#f8f7f8;">
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Date</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Title</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Distance</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Duration</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Pace</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Avg HR</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Status</th>
      <th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Actions</th>
    </tr></thead>
    <tbody>%s %s</tbody>
  </table>
</div>

<!-- Delete Modal -->
<div id="deleteModal" style="display:none;position:fixed;top:0;left:0;width:100%%;height:100%%;background:rgba(0,0,0,0.4);z-index:1000;align-items:center;justify-content:center;">
  <div style="background:#fff;border-radius:16px;padding:24px;max-width:400px;width:90%%;box-shadow:0 20px 60px rgba(0,0,0,0.2);">
    <h3 style="font-size:16px;font-weight:600;color:#1e1d1e;">Hapus Aktivitas</h3>
    <p style="font-size:13px;color:#6b6a6b;margin-top:8px;">Kamu yakin ingin menghapus aktivitas <strong id="deleteItemName"></strong>?</p>
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
  fetch(_deleteUrl, { method: 'DELETE' }).then(r => { if (r.ok) window.location.href = '/admin/running-activities?flash=deleted'; });
});
</script>`, notifHTML, stats, rows, empty)

		return c.Type("html").SendString(layout("Running Activities", "running", content))
	}
}

// ── Create / Edit Form ────────────────────────────────────────────────────────

func RunningActivityFormPage(db *gorm.DB, fallbackID string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id", fallbackID)
		isEdit := false
		var activity model.RunningActivity

		if id != "" {
			if err := db.First(&activity, "id = ?", id).Error; err != nil {
				return c.Status(404).SendString("Activity not found")
			}
			isEdit = true
		}

		title := "Log New Run"
		formAction := "/admin/running-activities"
		submitLabel := "Save Run"
		if isEdit {
			title = "Edit Run"
			formAction = "/admin/running-activities/" + id
			submitLabel = "Save Changes"
		}

		content := fmt.Sprintf(`
<div style="display:flex; justify-content: flex-end; margin-bottom: 16px;">
  <a href="/admin/running-activities" style="display:inline-flex;align-items:center;gap:6px;font-size:13px;color:#4d4c4d;text-decoration:none;padding:6px 12px;border:1px solid #e5e3e4;border-radius:6px;">← Back to Activities</a>
</div>
<form method="POST" action="%s">
  <!-- Basic Info -->
  <div style="background:#fff;border-radius:10px;border:1px solid #e5e3e4;padding:24px;margin-bottom:16px;">
    <h2 style="font-size:14px;font-weight:600;color:#1e1d1e;">Informasi Dasar</h2>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:16px;">
      %s %s
    </div>
    <div style="margin-top:16px;">%s</div>
    <div style="margin-top:16px;">%s</div>
  </div>
  <!-- Distance & Time -->
  <div style="background:#fff;border-radius:10px;border:1px solid #e5e3e4;padding:24px;margin-bottom:16px;">
    <h2 style="font-size:14px;font-weight:600;color:#1e1d1e;">Distance & Time</h2>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:16px;">
      %s %s
    </div>
  </div>
  <!-- Calories -->
  <div style="background:#fff;border-radius:10px;border:1px solid #e5e3e4;padding:24px;margin-bottom:16px;">
    <h2 style="font-size:14px;font-weight:600;color:#1e1d1e;">Calories</h2>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:16px;">
      %s %s
    </div>
  </div>
  <!-- Performance Metrics -->
  <div style="background:#fff;border-radius:10px;border:1px solid #e5e3e4;padding:24px;margin-bottom:16px;">
    <h2 style="font-size:14px;font-weight:600;color:#1e1d1e;">Performance Metrics</h2>
    <div style="display:grid;grid-template-columns:1fr 1fr 1fr;gap:16px;margin-top:16px;">
      %s %s %s %s %s %s
    </div>
  </div>
  <!-- Settings -->
  <div style="background:#fff;border-radius:10px;border:1px solid #e5e3e4;padding:24px;margin-bottom:16px;">
    <h2 style="font-size:14px;font-weight:600;color:#1e1d1e;">Pengaturan</h2>
    <div style="display:grid;grid-template-columns:1fr;gap:16px;margin-top:16px;">
      %s
    </div>
  </div>
  <!-- Actions -->
  <div style="display:flex;justify-content:flex-end;gap:12px;margin-top:16px;padding:16px 0;">
    <a href="/admin/running-activities" style="padding:8px 16px;border:1px solid #e5e3e4;border-radius:6px;color:#4d4c4d;text-decoration:none;">Batal</a>
    <button type="submit" style="padding:8px 20px;background:#1e1d1e;color:#f5f4f5;border:none;border-radius:6px;font-weight:500;">%s</button>
  </div>
</form>`,
			formAction,
			inputField("Title", "title", activity.Title, "text", false),
			inputField("Date", "date", activity.Date.Format("2006-01-02"), "date", true),
			textareaField("Notes", "notes", activity.Notes, 2),
			inputField("Map Image URL", "map_image_url", activity.MapImageURL, "url", false),
			inputField("Duration (seconds)", "duration_sec", fmt.Sprintf("%d", activity.DurationSec), "number", true),
			inputField("Distance (meters)", "distance_meters", fmt.Sprintf("%.0f", activity.DistanceMeters), "number", true),
			inputField("Total Calories", "total_calories", fmt.Sprintf("%d", activity.TotalCalories), "number", false),
			inputField("Active Calories", "active_calories", fmt.Sprintf("%d", activity.ActiveCalories), "number", false),
			inputField("Avg Pace (sec/km)", "avg_pace_sec", fmt.Sprintf("%d", activity.AvgPaceSec), "number", false),
			inputField("Avg Speed (km/h)", "avg_speed_kph", fmt.Sprintf("%.1f", activity.AvgSpeedKph), "number", false),
			inputField("Avg Cadence", "avg_cadence", fmt.Sprintf("%d", activity.AvgCadence), "number", false),
			inputField("Avg Stride (m)", "avg_stride", fmt.Sprintf("%.2f", activity.AvgStride), "number", false),
			inputField("Steps", "steps", fmt.Sprintf("%d", activity.Steps), "number", false),
			inputField("Avg Heart Rate", "avg_heart_rate", fmt.Sprintf("%d", activity.AvgHeartRate), "number", false),
			toggleSwitch("Published", "is_published", activity.IsPublished, "Tampilkan di halaman publik"),
			submitLabel,
		)

		return c.Type("html").SendString(layout(title, "running", content))
	}
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func CreateRunningActivityHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		date, _ := time.Parse("2006-01-02", c.FormValue("date"))
		activity := model.RunningActivity{
			Date:           date,
			Title:          c.FormValue("title"),
			Notes:          c.FormValue("notes"),
			MapImageURL:    c.FormValue("map_image_url"),
			DurationSec:    formInt(c, "duration_sec"),
			DistanceMeters: formFloat(c, "distance_meters"),
			TotalCalories:  formInt(c, "total_calories"),
			ActiveCalories: formInt(c, "active_calories"),
			AvgPaceSec:     formInt(c, "avg_pace_sec"),
			AvgSpeedKph:    formFloat(c, "avg_speed_kph"),
			AvgCadence:     formInt(c, "avg_cadence"),
			AvgStride:      formFloat(c, "avg_stride"),
			Steps:          formInt(c, "steps"),
			AvgHeartRate:   formInt(c, "avg_heart_rate"),
			IsPublished:    c.FormValue("is_published") == "true",
		}
		if err := db.Create(&activity).Error; err != nil {
			return c.Status(500).SendString("Failed to create activity")
		}
		return c.Redirect("/admin/running-activities?flash=created")
	}
}

func UpdateRunningActivityHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var activity model.RunningActivity
		if err := db.First(&activity, "id = ?", id).Error; err != nil {
			return c.Status(404).SendString("Activity not found")
		}
		date, _ := time.Parse("2006-01-02", c.FormValue("date"))
		activity.Date = date
		activity.Title = c.FormValue("title")
		activity.Notes = c.FormValue("notes")
		activity.MapImageURL = c.FormValue("map_image_url")
		activity.DurationSec = formInt(c, "duration_sec")
		activity.DistanceMeters = formFloat(c, "distance_meters")
		activity.TotalCalories = formInt(c, "total_calories")
		activity.ActiveCalories = formInt(c, "active_calories")
		activity.AvgPaceSec = formInt(c, "avg_pace_sec")
		activity.AvgSpeedKph = formFloat(c, "avg_speed_kph")
		activity.AvgCadence = formInt(c, "avg_cadence")
		activity.AvgStride = formFloat(c, "avg_stride")
		activity.Steps = formInt(c, "steps")
		activity.AvgHeartRate = formInt(c, "avg_heart_rate")
		activity.IsPublished = c.FormValue("is_published") == "true"
		db.Save(&activity)
		return c.Redirect("/admin/running-activities?flash=updated")
	}
}

func DeleteRunningActivityHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		db.Delete(&model.RunningActivity{}, "id = ?", id)
		return c.SendStatus(200)
	}
}

// ── Helper kecil untuk parsing form ──
func formInt(c *fiber.Ctx, key string) int {
	val := 0
	fmt.Sscanf(c.FormValue(key), "%d", &val)
	return val
}

func formFloat(c *fiber.Ctx, key string) float64 {
	val := 0.0
	fmt.Sscanf(c.FormValue(key), "%f", &val)
	return val
}