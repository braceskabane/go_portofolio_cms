package admin

import (
	"fmt"
	"log"
	"portfolio-cms/internal/config"
	"portfolio-cms/internal/dto"
	"portfolio-cms/internal/service"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RunningAnalysisPage(db *gorm.DB, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		html := layout(
			"Running Analysis",
			"/admin/running-analysis",
			renderAnalysisShell("", cfg.Google.ClientID),
			WithTopActions(renderAnalysisTopActions()),
		)
		return c.Type("html").SendString(html)
	}
}

func GenerateAnalysisHandler(db *gorm.DB, cfg *config.Config) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // ✅ LOG REQUEST MASUK
        log.Println("===== GENERATE ANALYSIS HANDLER CALLED =====")
        log.Printf("FormValues: %+v", c.Request().PostArgs())

        req := &dto.GenerateRunningAnalysisRequest{
            GoalDescription:  c.FormValue("goal_description"),
            PreferredRunTime: c.FormValue("preferred_run_time", "06:00"),
        }

        // ✅ LOG OBJEK REQUEST YANG SUDAH DIBENTUK
        log.Printf("Parsed request: %+v", req)

        svc := service.NewRunningAnalysisService(db, cfg.Gemini.APIKey, cfg.Gemini.Model)
        result, err := svc.GenerateAnalysis(req)
        if err != nil {
            // ✅ LOG ERROR
            log.Printf("ERROR in GenerateAnalysis: %v", err)
            return c.Type("html").SendString(
                renderErrorBanner("Gagal generate analisis: " + err.Error()),
            )
        }

        // ✅ LOG SUKSES
        log.Println("SUCCESS: Analysis generated")

        return c.Type("html").SendString(renderAnalysisResult(result, cfg.Google.ClientID))
    }
}

// ─── Shell ────────────────────────────────────────────────────────────────────

func renderAnalysisShell(resultHTML, clientID string) string {
	return fmt.Sprintf(`
%s
<div id="analysis-form-wrap">%s</div>
<div id="analysis-result" style="margin-top:24px;">%s</div>

<script>
(function() {
  var STORAGE_KEY = 'running_analysis_v3';
  var resultEl    = document.getElementById('analysis-result');

  // Selalu restore calendar config dari sessionStorage dulu
  // (tidak bergantung pada script di dalam HTML result)
  var savedEvents   = sessionStorage.getItem('gcal_events');
  var savedClientID = sessionStorage.getItem('gcal_client_id');
  if (savedEvents)   window._gcalEvents   = JSON.parse(savedEvents);
  if (savedClientID) window._gcalClientID = savedClientID;
  if (savedEvents || savedClientID) {
    console.log('[Calendar] Config restored dari sessionStorage:',
      (window._gcalEvents || []).length, 'events,',
      window._gcalClientID ? 'clientID ada' : 'clientID kosong'
    );
  }

  // Restore HTML result (tanpa eksekusi script — tidak perlu lagi)
  if (!resultEl.innerHTML.trim()) {
    var saved = sessionStorage.getItem(STORAGE_KEY);
    if (saved) {
      // Hapus tag <script> dari HTML sebelum di-set ke innerHTML
      // karena data sudah di-restore dari sessionStorage di atas
      var cleaned = saved.replace(/<script[\s\S]*?<\/script>/gi, '');
      resultEl.innerHTML = cleaned;
      console.log('[Analysis] Restored dari sessionStorage');
    }
  }
})();
</script>`,
		analysisPageStyles(),
		renderAnalysisForm(),
		resultHTML,
	)
}

// ─── Form ─────────────────────────────────────────────────────────────────────

func renderAnalysisForm() string {
	return fmt.Sprintf(`
<div style="background:#fff;border:1px solid #e5e3e4;border-radius:10px;padding:24px;">
  <h2 style="font-size:14px;font-weight:600;color:#1e1d1e;margin-bottom:4px;">Generate Analisis Lari</h2>
  <p style="font-size:12.5px;color:#6b6a6b;margin-bottom:20px;">
    Gemini menganalisis seluruh sesi lari dan menghasilkan laporan pelatih + jadwal latihan.
  </p>
  <div style="display:grid;grid-template-columns:1fr 1fr;gap:12px;margin-bottom:14px;">
    %s
    %s
  </div>
  <div style="display:flex;align-items:center;gap:12px;">
    <button onclick="generateAnalysis()" id="gen-btn"
      style="display:inline-flex;align-items:center;gap:8px;background:#1e1d1e;color:#f5f4f5;border:none;border-radius:6px;padding:9px 20px;font-size:13px;font-weight:500;cursor:pointer;">
      🧠 Generate Analisis
    </button>
    <span id="gen-spinner" style="font-size:12.5px;color:#6b6a6b;display:none;align-items:center;gap:6px;">
      <svg style="width:14px;height:14px;animation:spin 1s linear infinite;" viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
      </svg>
      Sedang menganalisis… (bisa 10–30 detik)
    </span>
  </div>
</div>

<script>
async function generateAnalysis() {
  var btn     = document.getElementById('gen-btn');
  var spinner = document.getElementById('gen-spinner');
  var result  = document.getElementById('analysis-result');

  var goal    = document.querySelector('[name="goal_description"]');
  var runTime = document.querySelector('[name="preferred_run_time"]');

  btn.disabled          = true;
  btn.style.opacity     = '0.6';
  spinner.style.display = 'inline-flex';
  result.innerHTML      = '';

  // Bersihkan cache lama saat generate baru
  sessionStorage.removeItem('running_analysis_v3');
  sessionStorage.removeItem('gcal_events');
  sessionStorage.removeItem('gcal_client_id');

  try {
    var formData = new FormData();
    formData.append('goal_description', goal ? goal.value : '');
    formData.append('preferred_run_time', runTime ? runTime.value : '06:00');

    var resp = await fetch('/admin/running-analysis/generate', {
      method: 'POST',
      body:   formData,
    });

    var html = await resp.text();

    // 1. Eksekusi script dulu (agar sessionStorage gcal_events terisi)
    var tempDiv = document.createElement('div');
    tempDiv.innerHTML = html;
    tempDiv.querySelectorAll('script').forEach(function(oldScript) {
      var newScript = document.createElement('script');
      newScript.textContent = oldScript.textContent;
      document.head.appendChild(newScript);
      document.head.removeChild(newScript);
    });

    // 2. Set HTML tanpa script ke result
    var cleaned = html.replace(/<script[\s\S]*?<\/script>/gi, '');
    result.innerHTML = cleaned;

    // 3. Simpan HTML bersih ke sessionStorage
    sessionStorage.setItem('running_analysis_v3', cleaned);

    console.log('[Analysis] Selesai. Events:', (window._gcalEvents || []).length, 'clientID:', window._gcalClientID ? 'ada' : 'kosong');

  } catch(err) {
    result.innerHTML = '<div style="padding:14px;border-radius:8px;background:#fef2f2;color:#991b1b;">❌ Error: ' + err.message + '</div>';
  } finally {
    btn.disabled          = false;
    btn.style.opacity     = '1';
    spinner.style.display = 'none';
  }
}
</script>`,
		inputField("Tujuan Latihan (opsional)", "goal_description", "", "text", false),
		inputField("Jam Preferensi Lari", "preferred_run_time", "06:00", "time", false),
	)
}

// ─── Full result renderer ─────────────────────────────────────────────────────

func renderAnalysisResult(r *dto.RunningAnalysisResult, clientID string) string {
	if r == nil {
		return ""
	}

	sections := []string{
		renderCoachNarrative(r.CoachNarrative, r.GeneratedAt),
		renderFitnessAssessment(r.FitnessAssessment),
		renderPaceZones(r.PaceZones),
		renderWeeklyPlan(r.WeeklyPlan),
		renderWarnings(r.Warnings),
	}

	if len(r.CalendarEvents) > 0 {
		sections = append(sections, renderCalendarSection(r.CalendarEvents, clientID))
	}

	result := ""
	for _, s := range sections {
		result += s
	}
	return result
}

// ─── Misc ─────────────────────────────────────────────────────────────────────

func renderAnalysisTopActions() string {
	return `<span style="font-size:12px;color:#6b6a6b;">Powered by Gemini Flash · semua sesi lari</span>`
}

func renderErrorBanner(msg string) string {
	return fmt.Sprintf(`
<div style="padding:14px 16px;border-radius:8px;background:#fef2f2;border-left:3px solid #dc2626;font-size:13px;color:#991b1b;">
  ⚠️ %s
</div>`, escapeHTML(msg))
}

// ─── CSS & Animation ──────────────────────────────────────────────────────────

func analysisPageStyles() string {
	return `
<style>
  @keyframes spin { to { transform: rotate(360deg); } }
  @keyframes fadeIn { from { opacity:0; transform:translateY(6px); } to { opacity:1; transform:translateY(0); } }

  .htmx-indicator          { display:none !important; }
  .htmx-request .htmx-indicator { display:inline-flex !important; }

  .analysis-card {
    background:#fff;
    border:1px solid #e5e3e4;
    border-radius:10px;
    padding:24px;
    margin-bottom:16px;
    animation:fadeIn 0.25s ease;
  }
  .analysis-card-header {
    display:flex;
    align-items:center;
    gap:10px;
    margin-bottom:20px;
  }
  .analysis-card-title {
    font-size:14px;
    font-weight:600;
    color:#1e1d1e;
  }
  .analysis-card-meta {
    margin-left:auto;
    font-size:11px;
    color:#9a9899;
  }
</style>`
}

// ─── Google Calendar JS ───────────────────────────────────────────────────────

func calendarScript(events []dto.CalendarEventDTO, clientID string) string {
	eventsJSON := marshalCalendarEvents(events)
	return fmt.Sprintf(`
<script>
(function() {
  var events   = %s;
  var clientID = %q;

  // Simpan ke sessionStorage agar tersedia saat restore
  sessionStorage.setItem('gcal_events',    JSON.stringify(events));
  sessionStorage.setItem('gcal_client_id', clientID);

  // Daftarkan ke window langsung
  window._gcalEvents   = events;
  window._gcalClientID = clientID;

  console.log('[Calendar] Config terdaftar:', events.length, 'events, clientID:', clientID ? 'ada' : 'kosong');
})();
</script>`, eventsJSON, clientID)
}

// ─── Sub-component Renderers ──────────────────────────────────────────────────

func renderCoachNarrative(narrative, generatedAt string) string {
	return fmt.Sprintf(`
<div class="analysis-card">
  <div class="analysis-card-header">
    <span style="font-size:20px;">🏅</span>
    <h2 class="analysis-card-title">Laporan Pelatih</h2>
    <span class="analysis-card-meta">%s</span>
  </div>
  <p style="font-size:13.5px;color:#2d2c2d;line-height:1.8;white-space:pre-line;">%s</p>
</div>`, escapeHTML(generatedAt), escapeHTML(narrative))
}

func renderFitnessAssessment(fa dto.FitnessAssessment) string {
	trendColor := colorByKey(map[string]string{
		"improving": "green", "plateau": "yellow", "declining": "red",
	}, fa.Trend, "gray")

	levelColor := colorByKey(map[string]string{
		"beginner": "blue", "intermediate": "yellow", "advanced": "green",
	}, fa.Level, "gray")

	baseColor := colorByKey(map[string]string{
		"weak": "red", "building": "yellow", "solid": "green", "strong": "purple",
	}, fa.AerobicBase, "gray")

	return fmt.Sprintf(`
<div class="analysis-card">
  <div class="analysis-card-header">
    <span style="font-size:20px;">📊</span>
    <h2 class="analysis-card-title">Kondisi Kebugaran</h2>
  </div>

  <div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(160px,1fr));gap:10px;margin-bottom:18px;">
    %s %s %s
  </div>

  <div style="display:grid;grid-template-columns:repeat(3,1fr);gap:10px;margin-bottom:18px;">
    %s %s %s
  </div>

  <div>
    <p style="font-size:12px;font-weight:600;color:#4d4c4d;margin-bottom:6px;">
      Fatigue Score: %d / 100
    </p>
    %s
  </div>
</div>`,
		miniStatCard("Level", strings.Title(fa.Level), levelColor),
		miniStatCard("Aerobic Base", strings.Title(fa.AerobicBase), baseColor),
		miniStatCard("Trend", strings.Title(fa.Trend), trendColor),
		tsbCard("CTL", fa.CTL, "Kebugaran kronik (42 hari)"),
		tsbCard("ATL", fa.ATL, "Beban akut (7 hari)"),
		tsbCard("TSB", fa.TSB, "Keseimbangan latihan"),
		fa.FatigueScore,
		renderProgressBar(int(fa.FatigueScore), 100, fatigueBarColor(int(fa.FatigueScore))),
	)
}

func renderPaceZones(pz dto.PaceZones) string {
	return fmt.Sprintf(`
<div class="analysis-card">
  <div class="analysis-card-header">
    <span style="font-size:20px;">⚡</span>
    <h2 class="analysis-card-title">Zona Pace Personal</h2>
  </div>
  <div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(180px,1fr));gap:10px;">
    %s %s %s %s
  </div>
</div>`,
		paceZoneCard("Easy", pz.Easy.Label, "#f0fdf4", "#166534"),
		paceZoneCard("Aerobic", pz.Aerobic.Label, "#eff6ff", "#1e40af"),
		paceZoneCard("Tempo", pz.Tempo.Label, "#fefce8", "#854d0e"),
		paceZoneCard("Threshold", pz.Threshold.Label, "#fef2f2", "#991b1b"),
	)
}

func renderWeeklyPlan(plan []dto.DayPlan) string {
	rows := ""
	for _, day := range plan {
		rows += renderDayPlanRow(day)
	}
	return fmt.Sprintf(`
<div class="analysis-card">
  <div class="analysis-card-header">
    <span style="font-size:20px;">📅</span>
    <h2 class="analysis-card-title">Rencana Latihan Mingguan</h2>
  </div>
  <div style="display:flex;flex-direction:column;gap:8px;">%s</div>
</div>`, rows)
}

func renderWarnings(warnings []dto.RunningWarning) string {
	if len(warnings) == 0 {
		return ""
	}
	rows := ""
	for _, w := range warnings {
		rows += renderWarningBanner(w)
	}
	return fmt.Sprintf(`
<div class="analysis-card">
  <div class="analysis-card-header">
    <span style="font-size:20px;">⚠️</span>
    <h2 class="analysis-card-title">Peringatan & Saran</h2>
  </div>
  %s
</div>`, rows)
}

func renderCalendarSection(events []dto.CalendarEventDTO, clientID string) string {
	rows := ""
	for _, e := range events {
		rows += fmt.Sprintf(`
<div style="display:flex;align-items:center;gap:10px;padding:10px 14px;border-radius:6px;border:1px solid #f0eff0;background:#faf9fa;margin-bottom:6px;">
  %s
  <div style="flex:1;min-width:0;">
    <p style="font-size:13px;font-weight:500;color:#1e1d1e;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;">%s</p>
    <p style="font-size:11.5px;color:#6b6a6b;margin-top:1px;">%s · %s · %d menit</p>
  </div>
</div>`,
			calendarColorDot(e.ColorID),
			escapeHTML(e.Title),
			e.Date, e.StartTime, e.DurationMin,
		)
	}

	return fmt.Sprintf(`
<div class="analysis-card">
  <div class="analysis-card-header">
    <span style="font-size:20px;">🗓️</span>
    <h2 class="analysis-card-title">Jadwal Latihan (%d sesi)</h2>
    <div style="margin-left:auto;">
      <button id="gcal-btn" onclick="handleAddToCalendar()"
        style="display:inline-flex;align-items:center;gap:6px;background:#4285f4;color:#fff;border:none;border-radius:6px;padding:8px 16px;font-size:13px;font-weight:500;cursor:pointer;transition:background 0.12s;">
        <svg style="width:14px;height:14px;flex-shrink:0;" viewBox="0 0 24 24" fill="currentColor">
          <path d="M19 3h-1V1h-2v2H8V1H6v2H5a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2V5a2 2 0 00-2-2zm0 16H5V9h14v10zM5 7V5h14v2H5z"/>
        </svg>
        Add to Google Calendar
      </button>
    </div>
  </div>

  <div style="margin-bottom:16px;">%s</div>
  <div id="gcal-status" style="display:none;margin-top:8px;"></div>
</div>

%s`, len(events), rows, calendarScript(events, clientID))
}

// ─── Atom helpers (hanya dipakai di file ini) ─────────────────────────────────

func miniStatCard(label, value, color string) string {
	st, ok := badgeStyles[color]
	if !ok {
		st = badgeStyles["gray"]
	}
	return fmt.Sprintf(`
<div style="padding:14px 16px;border-radius:8px;border:1px solid #f0eff0;%s">
  <p style="font-size:11px;font-weight:600;text-transform:uppercase;letter-spacing:0.06em;opacity:0.7;margin-bottom:4px;">%s</p>
  <p style="font-size:15px;font-weight:600;">%s</p>
</div>`, st, label, value)
}

func tsbCard(label string, value float64, desc string) string {
	color := "#1e1d1e"
	if label == "TSB" {
		switch {
		case value > 10:
			color = "#166534"
		case value < -20:
			color = "#991b1b"
		default:
			color = "#854d0e"
		}
	}
	return fmt.Sprintf(`
<div style="padding:14px 16px;border-radius:8px;border:1px solid #f0eff0;background:#faf9fa;">
  <p style="font-size:11px;font-weight:600;text-transform:uppercase;letter-spacing:0.06em;color:#6b6a6b;margin-bottom:4px;">%s</p>
  <p style="font-size:18px;font-weight:600;color:%s;">%.1f</p>
  <p style="font-size:11px;color:#9a9899;margin-top:2px;">%s</p>
</div>`, label, color, value, desc)
}

func renderProgressBar(value, max int, barColor string) string {
	pct := 0
	if max > 0 {
		pct = value * 100 / max
		if pct > 100 {
			pct = 100
		}
	}
	return fmt.Sprintf(`
<div style="background:#f0eff0;border-radius:99px;height:8px;overflow:hidden;">
  <div style="width:%d%%;height:100%%;background:%s;border-radius:99px;transition:width 0.4s ease;"></div>
</div>`, pct, barColor)
}

func fatigueBarColor(score int) string {
	switch {
	case score < 30:
		return "#16a34a"
	case score < 60:
		return "#d97706"
	default:
		return "#dc2626"
	}
}

func paceZoneCard(label, pace, bg, textColor string) string {
	return fmt.Sprintf(`
<div style="padding:16px;border-radius:8px;background:%s;">
  <p style="font-size:11px;font-weight:700;text-transform:uppercase;letter-spacing:0.07em;color:%s;opacity:0.75;margin-bottom:6px;">%s</p>
  <p style="font-size:16px;font-weight:600;color:%s;">%s</p>
</div>`, bg, textColor, label, textColor, pace)
}

func renderDayPlanRow(day dto.DayPlan) string {
	type sessionStyle struct{ bg, text, dot string }
	styles := map[string]sessionStyle{
		"easy":      {"#f0fdf4", "#166534", "#16a34a"},
		"tempo":     {"#fefce8", "#854d0e", "#d97706"},
		"long":      {"#eff6ff", "#1e40af", "#3b82f6"},
		"threshold": {"#fef2f2", "#991b1b", "#dc2626"},
		"rest":      {"#f5f4f5", "#6b6a6b", "#9a9899"},
		"strength":  {"#faf5ff", "#6b21a8", "#9333ea"},
	}
	st, ok := styles[day.SessionType]
	if !ok {
		st = styles["rest"]
	}

	distStr := `<span style="font-size:12px;color:#9a9899;">—</span>`
	if day.TargetDistKm > 0 {
		distStr = fmt.Sprintf(
			`<span style="font-size:12px;font-weight:600;color:%s;">%.1f km</span>`,
			st.text, day.TargetDistKm,
		)
	}

	return fmt.Sprintf(`
<div style="display:grid;grid-template-columns:80px 110px 1fr auto;align-items:center;gap:12px;padding:12px 16px;border-radius:8px;background:%s;border:1px solid #f0eff0;">
  <span style="font-size:12.5px;font-weight:600;color:%s;">%s</span>
  <span style="display:inline-flex;align-items:center;gap:5px;font-size:11.5px;font-weight:500;color:%s;">
    <span style="width:6px;height:6px;border-radius:50%%;background:%s;flex-shrink:0;"></span>
    %s
  </span>
  <div style="min-width:0;">
    <p style="font-size:12.5px;color:#2d2c2d;font-weight:500;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;">%s</p>
    <p style="font-size:11.5px;color:#6b6a6b;margin-top:2px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;">%s</p>
  </div>
  %s
</div>`,
		st.bg, st.text, escapeHTML(day.DayName),
		st.text, st.dot, strings.Title(day.SessionType),
		escapeHTML(day.Focus), escapeHTML(day.Rationale),
		distStr,
	)
}

func renderWarningBanner(w dto.RunningWarning) string {
	type wStyle struct{ border, bg, text, icon string }
	styles := map[string]wStyle{
		"info":     {"#94a3b8", "#f8fafc", "#475569", "ℹ️"},
		"warning":  {"#d97706", "#fffbeb", "#92400e", "⚠️"},
		"critical": {"#dc2626", "#fef2f2", "#991b1b", "🚨"},
	}
	st, ok := styles[w.Severity]
	if !ok {
		st = styles["info"]
	}
	return fmt.Sprintf(`
<div style="display:flex;gap:10px;padding:12px 14px;border-radius:6px;border-left:3px solid %s;background:%s;margin-bottom:8px;">
  <span style="flex-shrink:0;">%s</span>
  <p style="font-size:13px;color:%s;line-height:1.6;margin:0;">%s</p>
</div>`, st.border, st.bg, st.icon, st.text, escapeHTML(w.Message))
}

func calendarColorDot(colorID string) string {
	colors := map[string]string{
		"2": "#0b8043", // sage   = easy
		"6": "#f4511e", // tangerine = tempo/threshold
		"9": "#3f51b5", // blueberry = long
	}
	c, ok := colors[colorID]
	if !ok {
		c = "#9a9899"
	}
	return fmt.Sprintf(
		`<span style="width:10px;height:10px;border-radius:50%%;background:%s;flex-shrink:0;display:inline-block;"></span>`,
		c,
	)
}

func marshalCalendarEvents(events []dto.CalendarEventDTO) string {
	parts := make([]string, 0, len(events))
	for _, e := range events {
		parts = append(parts, fmt.Sprintf(
			`{"title":%q,"date":%q,"start_time":%q,"duration_min":%d,"description":%q,"color_id":%q}`,
			e.Title, e.Date, e.StartTime, e.DurationMin, e.Description, e.ColorID,
		))
	}
	return "[" + strings.Join(parts, ",") + "]"
}

// ─── Generic helpers ──────────────────────────────────────────────────────────

func colorByKey(m map[string]string, key, fallback string) string {
	if v, ok := m[key]; ok {
		return v
	}
	return fallback
}