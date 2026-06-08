package admin

import (
	"fmt"
	"portfolio-cms/internal/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DashboardPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ── Counts ────────────────────────────────────────────────────
		var (
			projectCount, skillCount, expCount, eduCount int64
			contactCount, unreadCount                    int64
			runCount, assetCount                         int64
			stackCatCount, projectCatCount, expCatCount  int64
		)
		db.Model(&model.Project{}).Count(&projectCount)
		db.Model(&model.Skill{}).Count(&skillCount)
		db.Model(&model.Experience{}).Count(&expCount)
		db.Model(&model.Education{}).Count(&eduCount)
		db.Model(&model.Contact{}).Count(&contactCount)
		db.Model(&model.Contact{}).Where("is_read = false").Count(&unreadCount)
		db.Model(&model.RunningActivity{}).Count(&runCount)
		db.Model(&model.Asset{}).Count(&assetCount)
		db.Model(&model.StackCategory{}).Count(&stackCatCount)
		db.Model(&model.ProjectCategory{}).Count(&projectCatCount)
		db.Model(&model.ExperienceCategory{}).Count(&expCatCount)

		// ── Stat card component ───────────────────────────────────────
		// icon: emoji, label: metric name, value: formatted count
		// accent: background color for icon container
		statCard := func(icon, label, value, accent, href string) string {
			change := ""
			if href != "" {
				change = fmt.Sprintf(
					`<a href="%s" style="display:inline-flex;align-items:center;gap:3px;font-size:11.5px;font-weight:500;color:#6b6a6b;text-decoration:none;margin-top:6px;transition:color 0.1s;" onmouseover="this.style.color='#1e1d1e'" onmouseout="this.style.color='#6b6a6b'">Manage →</a>`,
					href,
				)
			}
			return fmt.Sprintf(`
<div style="background:#ffffff;border-radius:10px;border:1px solid #e5e3e4;padding:18px;display:flex;flex-direction:column;gap:0;transition:box-shadow 0.15s,transform 0.15s;"
  onmouseover="this.style.boxShadow='0 4px 16px rgba(0,0,0,0.07)';this.style.transform='translateY(-1px)'"
  onmouseout="this.style.boxShadow='none';this.style.transform='none'">
  <div style="display:flex;align-items:flex-start;justify-content:space-between;margin-bottom:12px;">
    <div style="width:38px;height:38px;border-radius:9px;background:%s;display:flex;align-items:center;justify-content:center;font-size:18px;flex-shrink:0;">%s</div>
  </div>
  <p style="font-size:26px;font-weight:700;color:#1e1d1e;line-height:1;letter-spacing:-0.02em;">%s</p>
  <p style="font-size:12px;color:#6b6a6b;margin-top:4px;font-weight:500;">%s</p>
  %s
</div>`, accent, icon, value, label, change)
		}

		// ── Stat grid sections ────────────────────────────────────────

		// Group 1: Core content metrics
		contentStats := fmt.Sprintf(`
<div style="margin-bottom:8px;">
  <p style="font-size:11px;font-weight:700;text-transform:uppercase;letter-spacing:0.08em;color:#9a9899;margin-bottom:12px;">Content</p>
  <div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(160px,1fr));gap:12px;">
    %s %s %s %s
  </div>
</div>`,
			statCard("▤", "Projects", fmt.Sprintf("%d", projectCount), "#e0e7ff", "/admin/projects"),
			statCard("◎", "Skills", fmt.Sprintf("%d", skillCount), "#fef3c7", "/admin/skills"),
			statCard("◉", "Experiences", fmt.Sprintf("%d", expCount), "#dbeafe", "/admin/experiences"),
			statCard("◫", "Educations", fmt.Sprintf("%d", eduCount), "#f3e8ff", "/admin/educations"),
		)

		// Group 2: Activity & assets
		activityStats := fmt.Sprintf(`
<div style="margin-bottom:8px;">
  <p style="font-size:11px;font-weight:700;text-transform:uppercase;letter-spacing:0.08em;color:#9a9899;margin-bottom:12px;">Activity & Assets</p>
  <div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(160px,1fr));gap:12px;">
    %s %s
  </div>
</div>`,
			statCard("🏃", "Runs", fmt.Sprintf("%d", runCount), "#d1fae5", "/admin/running-activities"),
			statCard("🖼️", "Assets", fmt.Sprintf("%d", assetCount), "#fee2e2", "/admin/assets"),
		)

		// Group 3: Taxonomy counts
		taxonomyStats := fmt.Sprintf(`
<div style="margin-bottom:8px;">
  <p style="font-size:11px;font-weight:700;text-transform:uppercase;letter-spacing:0.08em;color:#9a9899;margin-bottom:12px;">Taxonomy</p>
  <div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(160px,1fr));gap:12px;">
    %s %s %s
  </div>
</div>`,
			statCard("🏷️", "Project Cats.", fmt.Sprintf("%d", projectCatCount), "#fce7f3", "/admin/project-categories"),
			statCard("🏷️", "Exp. Cats.", fmt.Sprintf("%d", expCatCount), "#cffafe", "/admin/experience-categories"),
			statCard("🗂️", "Stack Cats.", fmt.Sprintf("%d", stackCatCount), "#ede9fe", "/admin/stack-categories"),
		)

		// Group 4: Inbox — larger card
		inboxBadge := ""
		if unreadCount > 0 {
			inboxBadge = fmt.Sprintf(
				`<span style="margin-left:8px;background:#b91c1c;color:#fff;font-size:10px;font-weight:600;border-radius:999px;padding:1px 7px;">%d unread</span>`,
				unreadCount,
			)
		}
		inboxStats := fmt.Sprintf(`
<div style="margin-bottom:8px;">
  <p style="font-size:11px;font-weight:700;text-transform:uppercase;letter-spacing:0.08em;color:#9a9899;margin-bottom:12px;">Inbox</p>
  <div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(160px,1fr));gap:12px;">
    <div style="background:#ffffff;border-radius:10px;border:1px solid #e5e3e4;padding:18px;transition:box-shadow 0.15s,transform 0.15s;"
      onmouseover="this.style.boxShadow='0 4px 16px rgba(0,0,0,0.07)';this.style.transform='translateY(-1px)'"
      onmouseout="this.style.boxShadow='none';this.style.transform='none'">
      <div style="display:flex;align-items:flex-start;justify-content:space-between;margin-bottom:12px;">
        <div style="width:38px;height:38px;border-radius:9px;background:#fef9c3;display:flex;align-items:center;justify-content:center;font-size:18px;">✉️</div>
        %s
      </div>
      <p style="font-size:26px;font-weight:700;color:#1e1d1e;line-height:1;letter-spacing:-0.02em;">%d</p>
      <p style="font-size:12px;color:#6b6a6b;margin-top:4px;font-weight:500;">Messages</p>
      <a href="/admin/contacts" style="display:inline-flex;align-items:center;gap:3px;font-size:11.5px;font-weight:500;color:#6b6a6b;text-decoration:none;margin-top:6px;transition:color 0.1s;" onmouseover="this.style.color='#1e1d1e'" onmouseout="this.style.color='#6b6a6b'">View inbox →</a>
    </div>
  </div>
</div>`, inboxBadge, contactCount)

		statsBlock := fmt.Sprintf(`
<div style="display:flex;flex-direction:column;gap:24px;margin-bottom:28px;">
  %s %s %s %s
</div>`, contentStats, activityStats, taxonomyStats, inboxStats)

		// ── Quick actions ─────────────────────────────────────────────
		quickAction := func(icon, label, href, accent string) string {
			return fmt.Sprintf(`
<a href="%s" style="display:flex;align-items:center;gap:12px;padding:11px 14px;border-radius:8px;text-decoration:none;color:#1e1d1e;font-size:13px;font-weight:500;transition:background 0.1s;"
  onmouseover="this.style.background='#f5f4f5'"
  onmouseout="this.style.background='transparent'">
  <span style="width:32px;height:32px;flex-shrink:0;background:%s;border-radius:7px;display:flex;align-items:center;justify-content:center;font-size:15px;">%s</span>
  %s
</a>`, href, accent, icon, label)
		}

		quickActions := sectionCard("Quick Actions", fmt.Sprintf(`
<div style="display:flex;flex-direction:column;gap:2px;">
  %s %s %s %s %s %s
</div>`,
			quickAction("▤", "New Project", "/admin/projects/new", "#e0e7ff"),
			quickAction("◉", "New Experience", "/admin/experiences/new", "#dbeafe"),
			quickAction("◎", "Manage Skills", "/admin/skills", "#fef3c7"),
			quickAction("🏃", "Log a Run", "/admin/running-activities/new", "#d1fae5"),
			quickAction("🖼️", "Upload Asset", "/admin/assets/new", "#fee2e2"),
			quickAction("✉️", "View Messages", "/admin/contacts", "#fef9c3"),
		))

		// ── Public API reference ──────────────────────────────────────
		apiEndpoint := func(method, path string) string {
			methodColor := map[string]string{
				"GET":  "#16a34a",
				"POST": "#1d4ed8",
			}
			color, ok := methodColor[method]
			if !ok {
				color = "#6b6a6b"
			}
			return fmt.Sprintf(`
<div style="display:flex;align-items:center;gap:10px;background:#f8f7f8;border-radius:6px;padding:8px 12px;">
  <span style="font-size:10.5px;font-weight:700;font-family:monospace;color:%s;min-width:36px;">%s</span>
  <code style="font-size:12px;color:#4d4c4d;font-family:'Menlo','Consolas',monospace;">%s</code>
</div>`, color, method, path)
		}

		apiRef := sectionCard("Public API", fmt.Sprintf(`
<div style="display:flex;flex-direction:column;gap:6px;">
  %s %s %s %s %s %s %s
</div>`,
			apiEndpoint("GET", "/api/v1/public/profile"),
			apiEndpoint("GET", "/api/v1/public/projects"),
			apiEndpoint("GET", "/api/v1/public/skills"),
			apiEndpoint("GET", "/api/v1/public/experiences"),
			apiEndpoint("GET", "/api/v1/public/educations"),
			apiEndpoint("GET", "/api/v1/public/running-activities"),
			apiEndpoint("POST", "/api/v1/public/contact"),
		))

		// ── Bottom two-column layout ──────────────────────────────────
		bottomGrid := fmt.Sprintf(`
<div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;">
  %s %s
</div>`, quickActions, apiRef)

		content := statsBlock + bottomGrid

		return c.Type("html").SendString(
			layout("Dashboard", "Dashboard", content,
				WithUnreadContacts(int(unreadCount)),
			),
		)
	}
}