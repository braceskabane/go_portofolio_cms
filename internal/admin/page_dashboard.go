package admin

import (
	"fmt"
	"portfolio-cms/internal/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DashboardPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var projectCount, skillCount, expCount, eduCount, contactCount, unreadCount int64
		db.Model(&model.Project{}).Count(&projectCount)
		db.Model(&model.Skill{}).Count(&skillCount)
		db.Model(&model.Experience{}).Count(&expCount)
		db.Model(&model.Education{}).Count(&eduCount) // ← tambahkan ini
		db.Model(&model.Contact{}).Count(&contactCount)
		db.Model(&model.Contact{}).Where("is_read = false").Count(&unreadCount)

		// Helper statCard dengan inline style Pure Charcoal
		statCard := func(icon, label, value, accentBg string) string {
			return fmt.Sprintf(`
<div style="background:#ffffff; border-radius:12px; border:1px solid #e5e3e4; padding:20px; display:flex; align-items:center; gap:16px;">
  <div style="width:48px; height:48px; border-radius:12px; background:%s; display:flex; align-items:center; justify-content:center; font-size:24px; flex-shrink:0;">%s</div>
  <div>
    <p style="font-size:13px; color:#6b6a6b;">%s</p>
    <p style="font-size:28px; font-weight:600; color:#1e1d1e;">%s</p>
  </div>
</div>`, accentBg, icon, label, value)
		}

		unreadBadge := ""
		if unreadCount > 0 {
			unreadBadge = fmt.Sprintf(`<span style="margin-left:auto; background:#b91c1c; color:#ffffff; font-size:11px; border-radius:999px; padding:2px 8px;">%d</span>`, unreadCount)
		}

		content := fmt.Sprintf(`
<div style="display:grid; grid-template-columns: repeat(auto-fill, minmax(180px,1fr)); gap:16px; margin-bottom:32px;">
  %s %s %s %s %s
</div>

<div style="display:grid; grid-template-columns: 1fr; gap:24px;">
  <div style="background:#ffffff; border-radius:12px; border:1px solid #e5e3e4; padding:24px;">
    <h3 style="font-size:16px; font-weight:600; color:#1e1d1e; margin-bottom:12px;">Quick Actions</h3>
    <div style="display:flex; flex-direction:column; gap:8px;">
      <a href="/admin/projects/new" style="display:flex; align-items:center; gap:12px; padding:12px; border-radius:8px; text-decoration:none; color:#1e1d1e; font-size:13px; transition:background 0.12s;" onmouseover="this.style.background='#f0eff0'" onmouseout="this.style.background='transparent'">
        <span style="width:32px; height:32px; background:#e5e3e4; border-radius:8px; display:flex; align-items:center; justify-content:center;">➕</span>
        Add New Project
      </a>
      <a href="/admin/skills" style="display:flex; align-items:center; gap:12px; padding:12px; border-radius:8px; text-decoration:none; color:#1e1d1e; font-size:13px;" onmouseover="this.style.background='#f0eff0'" onmouseout="this.style.background='transparent'">
        <span style="width:32px; height:32px; background:#fef3c7; border-radius:8px; display:flex; align-items:center; justify-content:center;">⚡</span>
        Manage Skills
      </a>
      <a href="/admin/profile" style="display:flex; align-items:center; gap:12px; padding:12px; border-radius:8px; text-decoration:none; color:#1e1d1e; font-size:13px;" onmouseover="this.style.background='#f0eff0'" onmouseout="this.style.background='transparent'">
        <span style="width:32px; height:32px; background:#d1fae5; border-radius:8px; display:flex; align-items:center; justify-content:center;">👤</span>
        Update Profile
      </a>
      <a href="/admin/contacts" style="display:flex; align-items:center; gap:12px; padding:12px; border-radius:8px; text-decoration:none; color:#1e1d1e; font-size:13px;" onmouseover="this.style.background='#f0eff0'" onmouseout="this.style.background='transparent'">
        <span style="width:32px; height:32px; background:#fee2e2; border-radius:8px; display:flex; align-items:center; justify-content:center;">✉️</span>
        View Messages %s
      </a>
    </div>
  </div>
  <div style="background:#ffffff; border-radius:12px; border:1px solid #e5e3e4; padding:24px;">
    <h3 style="font-size:16px; font-weight:600; color:#1e1d1e; margin-bottom:12px;">Public API Endpoints</h3>
    <div style="display:flex; flex-direction:column; gap:8px; font-size:13px; font-family:monospace;">
      <div style="background:#f5f4f5; border-radius:6px; padding:8px 12px; display:flex; align-items:center; gap:8px;"><span style="color:#16a34a; font-weight:500; font-family:sans-serif; font-size:11px;">GET</span> /api/v1/public/profile</div>
      <div style="background:#f5f4f5; border-radius:6px; padding:8px 12px; display:flex; align-items:center; gap:8px;"><span style="color:#16a34a; font-weight:500; font-family:sans-serif; font-size:11px;">GET</span> /api/v1/public/projects</div>
      <div style="background:#f5f4f5; border-radius:6px; padding:8px 12px; display:flex; align-items:center; gap:8px;"><span style="color:#16a34a; font-weight:500; font-family:sans-serif; font-size:11px;">GET</span> /api/v1/public/skills</div>
      <div style="background:#f5f4f5; border-radius:6px; padding:8px 12px; display:flex; align-items:center; gap:8px;"><span style="color:#16a34a; font-weight:500; font-family:sans-serif; font-size:11px;">GET</span> /api/v1/public/experiences</div>
      <div style="background:#f5f4f5; border-radius:6px; padding:8px 12px; display:flex; align-items:center; gap:8px;"><span style="color:#16a34a; font-weight:500; font-family:sans-serif; font-size:11px;">POST</span> /api/v1/public/contact</div>
    </div>
  </div>
</div>`,
			statCard("🗂", "Projects", fmt.Sprintf("%d", projectCount), "#e0e7ff"),    // indigo-50
			statCard("⚡", "Skills", fmt.Sprintf("%d", skillCount), "#fef3c7"),      // yellow-50
			statCard("💼", "Experiences", fmt.Sprintf("%d", expCount), "#dbeafe"),    // blue-50
			statCard("🎓", "Educations", fmt.Sprintf("%d", eduCount), "#f3e8ff"),    // purple-50  ← sekarang pakai eduCount
			statCard("✉️", "Messages", fmt.Sprintf("%d", contactCount), "#d1fae5"),  // green-50
			unreadBadge,
		)

		return c.Type("html").SendString(layout("Dashboard", "Dashboard", content))
	}
}