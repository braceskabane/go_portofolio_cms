package admin

import (
	"fmt"
	"io"
	"portfolio-cms/internal/config"
	"portfolio-cms/internal/dto"
	"portfolio-cms/internal/model"
	"portfolio-cms/internal/service"
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
<div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(180px,1fr));gap:16px;margin-bottom:24px;">
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
				<td style="padding:12px 16px;font-size:13px;">%d</td>
				<td style="padding:12px 16px;font-size:13px;">%s</td>
				<td style="padding:12px 16px;font-size:13px;font-weight:500;">%s</td>
				<td style="padding:12px 16px;font-size:13px;">%.2f</td>
				<td style="padding:12px 16px;font-size:13px;">%d</td>
				<td style="padding:12px 16px;font-size:13px;">%s</td>
				<td style="padding:12px 16px;font-size:13px;">%.1f</td>
				<td style="padding:12px 16px;font-size:13px;">%d</td>
				<td style="padding:12px 16px;font-size:13px;">%.0f</td>
				<td style="padding:12px 16px;font-size:13px;">%.0f</td>
				<td style="padding:12px 16px;font-size:13px;">%d</td>
				<td style="padding:12px 16px;font-size:13px;">%.2f</td>
				<td style="padding:12px 16px;font-size:13px;">%d</td>
				<td style="padding:12px 16px;">%s</td>
				<td style="padding:12px 16px;white-space:nowrap;">
  					<div style="display:flex;gap:8px;">
						<a href="/admin/running-activities/%s/edit" style="font-size:12px;font-weight:500;color:#1e1d1e;text-decoration:none;padding:4px 10px;background:#f0eff0;border-radius:4px;">Edit</a>
						<button onclick="confirmDelete('/admin/running-activities/%s', '%s')" style="font-size:12px;font-weight:500;color:#b91c1c;background:none;border:none;cursor:pointer;padding:4px 10px;background:#fee2e2;border-radius:4px;">Delete</button>
					</div>
				</td>
			</tr>`,
				bg, bg,
				i+1,                                // nomor
				a.Date.Format("2006-01-02"),        // tanggal
				a.Title,                            // title
				distanceKm,
    			durationMin,                      
				pace,                               // pace
				a.AvgSpeedKph,                      // speed (km/h)
				a.AvgHeartRate,                     // hr (bpm)
				a.TotalCalories,                    // kalori
				a.ActiveCalories,                   // kalori aktif
				a.AvgCadence,                       // cadence
				a.AvgStride,                        // stride (m)
				a.Steps,                            // steps
				publishedBadge,                     // status
				a.ID, a.ID, a.Date.Format("2006-01-02")+" "+a.Title, // untuk edit/delete
			)
		}

		empty := ""
		if len(activities) == 0 {
			empty = fmt.Sprintf(`<tr><td colspan="15" style="text-align:center;padding:40px;">%s</td></tr>`, emptyState("🏃", "Belum ada aktivitas lari", "Catat lari pertama kamu", "Log New Run", "/admin/running-activities/new"))
		}

		content := fmt.Sprintf(`
%s
<div style="display:flex;justify-content:flex-end;gap:10px;margin-bottom:24px;">
  <button onclick="openGeminiModal()" style="display:inline-flex;align-items:center;gap:6px;padding:8px 18px;background:#4f46e5;color:#fff;border:none;border-radius:6px;cursor:pointer;font-weight:500;font-size:13px;">
    ✦ Add via Gemini
  </button>
  <a href="/admin/running-activities/new" style="display:inline-flex;align-items:center;gap:6px;padding:8px 18px;background:#1e1d1e;color:#f5f4f5;border-radius:6px;text-decoration:none;font-weight:500;font-size:13px;">
    + New Activity
  </a>
</div>
%s
<div style="background:#ffffff;border-radius:10px;border:1px solid #e5e3e4;overflow:hidden;min-width:0;">
  <div style="overflow-x:auto;">
    <table style="width:100%%;border-collapse:collapse;min-width:1200px;">
		<thead><tr style="background:#f8f7f8;">
			<th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">#</th>
			<th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Tanggal</th>
			<th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Title</th>
			<th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Jarak (m)</th>
			<th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Durasi (s)</th>
			<th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Pace</th>
			<th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Speed (km/h)</th>
			<th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">HR (bpm)</th>
			<th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Kalori</th>
			<th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Kalori Aktif</th>
			<th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Cadence</th>
			<th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Stride (m)</th>
			<th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Steps</th>
			<th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Status</th>
			<th style="padding:12px 16px;text-align:left;font-size:11px;font-weight:600;text-transform:uppercase;color:#6b6a6b;">Actions</th>
		</tr></thead>
    	<tbody>%s %s</tbody>
  	</table>
  </div>
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

		return c.Type("html").SendString(layout("Running Activities", "running", content + geminiStyles + geminiModalHTML + geminiScript))
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
			TotalCalories:  formFloat(c, "total_calories"),
			ActiveCalories: formFloat(c, "active_calories"),
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
		activity.TotalCalories = formFloat(c, "total_calories")
		activity.ActiveCalories = formFloat(c, "active_calories")
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

// ── Helper untuk pointer dereference ──────────────────────────────────────────

func getInt(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

func getFloat64(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

func getString(v string) string {
	return v
}

// ── Batch Create Handler  ───────────────────────────────────────────────

func BatchCreateRunningActivitiesHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			Activities []dto.CreateRunningActivityRequest `json:"activities"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Invalid JSON body",
			})
		}

		var activities []model.RunningActivity
		for _, a := range req.Activities {
			// Konversi date string ke time.Time, fallback ke zero time jika gagal
			date := time.Time{}
			if a.Date != nil && *a.Date != "" {
				if parsed, err := time.Parse(time.RFC3339, *a.Date); err == nil {
					date = parsed
				} else if parsed, err := time.Parse("2006-01-02", *a.Date); err == nil {
					date = parsed
				}
			}

			act := model.RunningActivity{
				Date:           date,
				Title:          a.Title,
				Notes:          a.Notes,
				MapImageURL:    a.MapImageURL,
				DurationSec:    getInt(a.DurationSec),
				DistanceMeters: getFloat64(a.DistanceMeters),
				TotalCalories:  getFloat64(a.TotalCalories),
				ActiveCalories: getFloat64(a.ActiveCalories),
				AvgPaceSec:     getInt(a.AvgPaceSec),
				AvgSpeedKph:    getFloat64(a.AvgSpeedKph),
				AvgCadence:     getInt(a.AvgCadence),
				AvgStride:      getFloat64(a.AvgStride),
				Steps:          getInt(a.Steps),
				AvgHeartRate:   getInt(a.AvgHeartRate),
				IsPublished:    a.IsPublished,
			}
			activities = append(activities, act)
		}

		if len(activities) == 0 {
			return c.JSON(fiber.Map{"success": true, "count": 0})
		}

		if err := db.Create(&activities).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"count":   len(activities),
		})
	}
}

// ── Preview Screenshots via Gemini ────────────────────────────────────────────

func PreviewScreenshotsHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cfg := config.Cfg
		if cfg.Gemini.APIKey == "" {
			return c.Status(500).JSON(fiber.Map{
				"error": "GEMINI_API_KEY belum dikonfigurasi",
			})
		}

		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Gagal membaca form"})
		}

		files := form.File["screenshots"]
		if len(files) == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "Tidak ada file yang diupload"})
		}
		if len(files) > 10 {
			return c.Status(400).JSON(fiber.Map{"error": "Maksimal 10 foto sekaligus"})
		}

		geminiSvc := service.NewGeminiService(cfg.Gemini.APIKey, cfg.Gemini.Model)

		type Result struct {
			Index    int                              `json:"index"`
			Filename string                           `json:"filename"`
			Data     *dto.CreateRunningActivityRequest `json:"data"`
			Error    string                           `json:"error,omitempty"`
		}

		results := make([]Result, 0, len(files))

		for i, fileHeader := range files {
			mimeType := fileHeader.Header.Get("Content-Type")
			allowed := map[string]bool{
				"image/jpeg": true,
				"image/png":  true,
				"image/webp": true,
			}

			if !allowed[mimeType] {
				results = append(results, Result{
					Index:    i,
					Filename: fileHeader.Filename,
					Error:    "Format tidak didukung (hanya JPEG/PNG/WebP)",
				})
				continue
			}

			if fileHeader.Size > 10*1024*1024 {
				results = append(results, Result{
					Index:    i,
					Filename: fileHeader.Filename,
					Error:    "File terlalu besar (maks 10MB)",
				})
				continue
			}

			f, err := fileHeader.Open()
			if err != nil {
				results = append(results, Result{
					Index:    i,
					Filename: fileHeader.Filename,
					Error:    "Gagal membaca file",
				})
				continue
			}

			imageBytes, err := io.ReadAll(f)
			f.Close()
			if err != nil {
				results = append(results, Result{
					Index:    i,
					Filename: fileHeader.Filename,
					Error:    "Gagal membaca isi file",
				})
				continue
			}

			extracted, err := geminiSvc.ExtractRunningActivity(imageBytes, mimeType)
			if err != nil {
				results = append(results, Result{
					Index:    i,
					Filename: fileHeader.Filename,
					Error:    "Gagal ekstrak: " + err.Error(),
				})
				continue
			}

			results = append(results, Result{
				Index:    i,
				Filename: fileHeader.Filename,
				Data:     extracted,
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"total":   len(files),
			"results": results,
		})
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

// ─── Blok HTML/CSS/JS Modal Gemini (di luar Sprintf) ──────────────────────
const geminiStyles = `
<style>
@keyframes spin { to { transform: rotate(360deg); } }
#dropzone:hover { border-color: #4f46e5; }
#previewTable td[contenteditable="true"]:focus { outline: 2px solid #4f46e5; border-radius: 2px; background: #eef2ff; }
#previewTable td { padding: 8px 10px; border-bottom: 1px solid #f0eff0; }
#previewTable tr:hover td { background: #fafafa; }

/* Scrollbar tabel running */
.table-scroll::-webkit-scrollbar { height: 8px; }
.table-scroll::-webkit-scrollbar-track { background: #f0eff0; border-radius: 0 0 10px 10px; }
.table-scroll::-webkit-scrollbar-thumb { background: #9a9899; border-radius: 4px; }
.table-scroll::-webkit-scrollbar-thumb:hover { background: #6b6a6b; }
</style>
`

const geminiModalHTML = `
<!-- Gemini Upload Modal -->
<div id="geminiModal" style="display:none;position:fixed;top:0;left:0;width:100%;height:100%;background:rgba(0,0,0,0.5);z-index:2000;align-items:flex-start;justify-content:center;overflow-y:auto;padding:40px 16px;box-sizing:border-box;">
  <div style="background:#fff;border-radius:16px;padding:28px;max-width:900px;width:100%;box-shadow:0 20px 60px rgba(0,0,0,0.25);margin:auto;">

    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:20px;">
      <div>
        <h2 style="font-size:16px;font-weight:600;color:#1e1d1e;">✦ Import via Gemini AI</h2>
        <p style="font-size:12px;color:#6b6a6b;margin-top:4px;">Upload screenshot Huawei Health, data akan diekstrak otomatis.</p>
      </div>
      <button onclick="closeGeminiModal()" style="background:none;border:none;font-size:20px;cursor:pointer;color:#6b6a6b;line-height:1;">×</button>
    </div>

    <!-- Step 1: Upload -->
    <div id="geminiStep1">
      <!-- DIV bukan LABEL — hindari double trigger file dialog -->
      <div id="dropzone" onclick="document.getElementById('geminiFiles').click()"
        style="display:flex;flex-direction:column;align-items:center;justify-content:center;border:2px dashed #d1d5db;border-radius:10px;padding:40px;cursor:pointer;transition:border-color 0.2s,background 0.2s;text-align:center;">
        <span style="font-size:32px;">📷</span>
        <span style="font-size:14px;font-weight:500;color:#374151;margin-top:8px;">Klik atau drag foto ke sini</span>
        <span style="font-size:12px;color:#9ca3af;margin-top:4px;">JPEG, PNG, WebP — maks 10 foto, masing-masing 10MB</span>
      </div>
      <!-- Input TERPISAH dari dropzone, bukan di dalam -->
      <input type="file" id="geminiFiles" accept="image/jpeg,image/png,image/webp" multiple style="display:none;">

      <div id="fileList" style="margin-top:12px;display:none;">
        <p style="font-size:12px;color:#6b6a6b;" id="fileListLabel"></p>
      </div>
      <div style="display:flex;justify-content:flex-end;margin-top:16px;gap:10px;">
        <button onclick="closeGeminiModal()" style="padding:8px 16px;border:1px solid #e5e3e4;border-radius:6px;background:#fff;cursor:pointer;font-size:13px;">Batal</button>
        <button id="btnExtract" onclick="extractScreenshots()" disabled style="padding:8px 20px;background:#4f46e5;color:#fff;border:none;border-radius:6px;cursor:pointer;font-weight:500;font-size:13px;opacity:0.5;">
          Ekstrak Data →
        </button>
      </div>
    </div>

    <!-- Step 2: Loading -->
    <div id="geminiStep2" style="display:none;text-align:center;padding:40px 0;">
      <div style="font-size:32px;animation:spin 1s linear infinite;display:inline-block;">⏳</div>
      <p style="font-size:14px;color:#4b5563;margin-top:12px;" id="loadingText">Mengekstrak data dari foto...</p>
    </div>

    <!-- Step 3: Preview -->
    <div id="geminiStep3" style="display:none;">
      <p style="font-size:13px;color:#374151;margin-bottom:12px;">Periksa dan koreksi data sebelum disimpan. Klik sel untuk edit.</p>
      <div style="overflow-x:auto;">
        <table style="width:100%;border-collapse:collapse;font-size:12px;" id="previewTable">
          <thead>
            <tr style="background:#f8f7f8;">
              <th style="padding:8px 10px;text-align:left;white-space:nowrap;color:#6b6a6b;font-size:11px;">#</th>
              <th style="padding:8px 10px;text-align:left;white-space:nowrap;color:#6b6a6b;font-size:11px;">FILE</th>
              <th style="padding:8px 10px;text-align:left;white-space:nowrap;color:#6b6a6b;font-size:11px;">TANGGAL</th>
              <th style="padding:8px 10px;text-align:left;white-space:nowrap;color:#6b6a6b;font-size:11px;">TITLE</th>
              <th style="padding:8px 10px;text-align:left;white-space:nowrap;color:#6b6a6b;font-size:11px;">JARAK (m)</th>
              <th style="padding:8px 10px;text-align:left;white-space:nowrap;color:#6b6a6b;font-size:11px;">DURASI (s)</th>
              <th style="padding:8px 10px;text-align:left;white-space:nowrap;color:#6b6a6b;font-size:11px;">PACE (s/km)</th>
              <th style="padding:8px 10px;text-align:left;white-space:nowrap;color:#6b6a6b;font-size:11px;">SPEED (km/h)</th>
              <th style="padding:8px 10px;text-align:left;white-space:nowrap;color:#6b6a6b;font-size:11px;">HR (bpm)</th>
              <th style="padding:8px 10px;text-align:left;white-space:nowrap;color:#6b6a6b;font-size:11px;">KALORI</th>
              <th style="padding:8px 10px;text-align:left;white-space:nowrap;color:#6b6a6b;font-size:11px;">KALORI AKTIF</th>
              <th style="padding:8px 10px;text-align:left;white-space:nowrap;color:#6b6a6b;font-size:11px;">CADENCE</th>
              <th style="padding:8px 10px;text-align:left;white-space:nowrap;color:#6b6a6b;font-size:11px;">STRIDE (m)</th>
              <th style="padding:8px 10px;text-align:left;white-space:nowrap;color:#6b6a6b;font-size:11px;">STEPS</th>
              <th style="padding:8px 10px;text-align:left;white-space:nowrap;color:#6b6a6b;font-size:11px;">STATUS</th>
            </tr>
          </thead>
          <tbody id="previewBody"></tbody>
        </table>
      </div>
      <div id="errorList" style="margin-top:12px;display:none;">
        <p style="font-size:12px;color:#b91c1c;font-weight:500;">Foto yang gagal diproses:</p>
        <ul id="errorItems" style="font-size:12px;color:#b91c1c;margin-top:4px;padding-left:16px;"></ul>
      </div>
      <div style="display:flex;justify-content:space-between;align-items:center;margin-top:20px;">
        <button onclick="resetGeminiModal()" style="padding:8px 16px;border:1px solid #e5e3e4;border-radius:6px;background:#fff;cursor:pointer;font-size:13px;">← Upload Lagi</button>
        <button onclick="saveAllActivities()" id="btnSaveAll" style="padding:8px 24px;background:#16a34a;color:#fff;border:none;border-radius:6px;font-weight:500;font-size:13px;cursor:pointer;">
          Simpan Semua →
        </button>
      </div>
    </div>

  </div>
</div>
`

const geminiScript = `
<script>
let _previewData = [];
let _selectedFiles = null; // simpan FileList di sini, TIDAK di input.files

// ── Modal ──
function openGeminiModal() {
  document.getElementById('geminiModal').style.display = 'flex';
  resetGeminiModal();
}
function closeGeminiModal() {
  document.getElementById('geminiModal').style.display = 'none';
}
function resetGeminiModal() {
  document.getElementById('geminiStep1').style.display = 'block';
  document.getElementById('geminiStep2').style.display = 'none';
  document.getElementById('geminiStep3').style.display = 'none';
  document.getElementById('geminiFiles').value = '';
  document.getElementById('fileList').style.display = 'none';
  document.getElementById('btnExtract').disabled = true;
  document.getElementById('btnExtract').style.opacity = '0.5';
  _previewData = [];
  _selectedFiles = null;
}

// ── Input file change (klik biasa) ──
// Pasang listener sekali saja, BUKAN onchange di HTML
document.getElementById('geminiFiles').addEventListener('change', function() {
  if (this.files && this.files.length > 0) {
    _selectedFiles = this.files;
    showSelectedFiles(this.files);
  }
});

function showSelectedFiles(files) {
  const names = Array.from(files).map(f => f.name).join(', ');
  document.getElementById('fileListLabel').textContent = files.length + ' foto dipilih: ' + names;
  document.getElementById('fileList').style.display = 'block';
  document.getElementById('btnExtract').disabled = false;
  document.getElementById('btnExtract').style.opacity = '1';
}

// ── Drag & drop ──
const dz = document.getElementById('dropzone');
dz.addEventListener('dragover', function(e) {
  e.preventDefault();
  e.stopPropagation();
  this.style.borderColor = '#4f46e5';
  this.style.background = '#f5f3ff';
});
dz.addEventListener('dragleave', function(e) {
  e.stopPropagation();
  this.style.borderColor = '#d1d5db';
  this.style.background = '';
});
dz.addEventListener('drop', function(e) {
  e.preventDefault();
  e.stopPropagation();
  this.style.borderColor = '#d1d5db';
  this.style.background = '';
  const files = e.dataTransfer.files;
  if (files && files.length > 0) {
    // Simpan ke _selectedFiles, JANGAN assign ke input.files (read-only di kebanyakan browser)
    _selectedFiles = files;
    showSelectedFiles(files);
  }
});

// ── Ekstrak ke Gemini ──
async function extractScreenshots() {
  // Ambil dari _selectedFiles, bukan dari input.files
  const files = _selectedFiles;
  if (!files || files.length === 0) {
    alert('Pilih foto terlebih dahulu');
    return;
  }

  document.getElementById('geminiStep1').style.display = 'none';
  document.getElementById('geminiStep2').style.display = 'block';

  const formData = new FormData();
  // Loop semua file — ini yang handle multiple
  for (let i = 0; i < files.length; i++) {
    formData.append('screenshots', files[i]);
  }

  let dots = 0;
  const interval = setInterval(() => {
    dots = (dots + 1) % 4;
    document.getElementById('loadingText').textContent =
      'Mengekstrak data dari ' + files.length + ' foto' + '.'.repeat(dots);
  }, 500);

  try {
    const res = await fetch('/admin/running-activities/preview-screenshots', {
      method: 'POST',
      body: formData
    });
    const json = await res.json();
    clearInterval(interval);
    if (json.error) {
      alert('Error: ' + json.error);
      resetGeminiModal();
      return;
    }
    renderPreview(json);
  } catch (err) {
    clearInterval(interval);
    alert('Terjadi kesalahan: ' + err.message);
    resetGeminiModal();
  }
}

// ── Render preview — muncul setelah Gemini selesai proses ──
function renderPreview(json) {
  document.getElementById('geminiStep2').style.display = 'none';
  document.getElementById('geminiStep3').style.display = 'block';

  const tbody = document.getElementById('previewBody');
  tbody.innerHTML = '';
  _previewData = [];

  const errors = json.results.filter(r => r.error);
  const successes = json.results.filter(r => !r.error);

  if (successes.length === 0 && errors.length > 0) {
    alert('Semua foto gagal diproses. Coba lagi.');
    resetGeminiModal();
    return;
  }

  successes.forEach((r, idx) => {
    const d = r.data || {};
    _previewData.push({ ...d, _filename: r.filename });
    const dateVal = d.date ? d.date.substring(0, 10) : '';
    const tr = document.createElement('tr');
    tr.innerHTML =
      '<td style="color:#9ca3af;padding:8px 10px;">' + (idx + 1) + '</td>' +
      '<td style="color:#6b6a6b;max-width:120px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;padding:8px 10px;" title="' + r.filename + '">' + r.filename + '</td>' +
      editableCell(idx, 'date', dateVal) +
      editableCell(idx, 'title', d.title || '') +
      editableCell(idx, 'distance_meters', d.distance_meters != null ? d.distance_meters : '') +
      editableCell(idx, 'duration_sec', d.duration_sec != null ? d.duration_sec : '') +
      editableCell(idx, 'avg_pace_sec', d.avg_pace_sec != null ? d.avg_pace_sec : '') +
      editableCell(idx, 'avg_speed_kph', d.avg_speed_kph != null ? d.avg_speed_kph : '') +
      editableCell(idx, 'avg_heart_rate', d.avg_heart_rate != null ? d.avg_heart_rate : '') +
      editableCell(idx, 'total_calories', d.total_calories != null ? d.total_calories : '') +
      editableCell(idx, 'active_calories', d.active_calories != null ? d.active_calories : '') +
      editableCell(idx, 'avg_cadence', d.avg_cadence != null ? d.avg_cadence : '') +
      editableCell(idx, 'avg_stride', d.avg_stride != null ? d.avg_stride : '') +
      editableCell(idx, 'steps', d.steps != null ? d.steps : '') +
      '<td style="padding:8px 10px;"><select onchange="updateField(' + idx + ',\'is_published\',this.value===\'true\')" style="font-size:12px;border:1px solid #e5e3e4;border-radius:4px;padding:2px 6px;">' +
        '<option value="false"' + (!d.is_published ? ' selected' : '') + '>Draft</option>' +
        '<option value="true"' + (d.is_published ? ' selected' : '') + '>Published</option>' +
      '</select></td>';
    tbody.appendChild(tr);
  });

  if (errors.length > 0) {
    document.getElementById('errorList').style.display = 'block';
    const ul = document.getElementById('errorItems');
    ul.innerHTML = '';
    errors.forEach(r => {
      const li = document.createElement('li');
      li.textContent = r.filename + ': ' + r.error;
      ul.appendChild(li);
    });
  } else {
    document.getElementById('errorList').style.display = 'none';
  }
}

function editableCell(idx, field, value) {
  return '<td contenteditable="true" style="padding:8px 10px;border-bottom:1px solid #f0eff0;min-width:60px;" ' +
    'onblur="updateField(' + idx + ',\'' + field + '\',this.textContent.trim())">' +
    (value !== null && value !== undefined ? value : '') + '</td>';
}

function updateField(idx, field, value) {
  if (!_previewData[idx]) return;
  const intFields = ['duration_sec','avg_pace_sec','avg_heart_rate','total_calories','active_calories','avg_cadence','steps'];
  const floatFields = ['distance_meters','avg_speed_kph','avg_stride'];
  if (field === 'is_published') {
    _previewData[idx][field] = value === 'true' || value === true;
  } else if (floatFields.includes(field)) {
    _previewData[idx][field] = parseFloat(value) || null;
  } else if (intFields.includes(field)) {
    _previewData[idx][field] = parseInt(value) || null;
  } else {
    _previewData[idx][field] = value || null;
  }
}

// ── Simpan semua ──
async function saveAllActivities() {
  if (!_previewData.length) return;
  const btn = document.getElementById('btnSaveAll');
  btn.textContent = 'Menyimpan...';
  btn.disabled = true;

  const payload = _previewData.map(d => {
    const copy = { ...d };
    delete copy._filename;
    if (copy.date && copy.date.length === 10) {
      copy.date = copy.date + 'T00:00:00+07:00';
    }
    return copy;
  });

  try {
    const res = await fetch('/admin/running-activities/batch', {
      method: 'POST',
      credentials: 'same-origin',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ activities: payload })
    });
    const json = await res.json();
    if (json.success) {
      closeGeminiModal();
      window.location.href = '/admin/running-activities?flash=created';
    } else {
      alert('Gagal menyimpan: ' + (json.message || 'Unknown error'));
      btn.textContent = 'Simpan Semua →';
      btn.disabled = false;
    }
  } catch (err) {
    alert('Error: ' + err.message);
    btn.textContent = 'Simpan Semua →';
    btn.disabled = false;
  }
}
</script>
`