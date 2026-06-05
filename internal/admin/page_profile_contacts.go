package admin

import (
	"fmt"
	"portfolio-cms/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ─── Profile ───────────────────────────────────────────────────────────────────

func ProfilePage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var p model.Profile
		db.First(&p) // may be empty — that's fine

		form := fmt.Sprintf(`
<div class="max-w-2xl">
  <form method="POST" action="/admin/profile" class="space-y-4">
    <div class="bg-white rounded-xl border border-gray-200 shadow-sm p-6 space-y-4">
      <h3 class="font-semibold text-gray-800 text-base">Personal Info</h3>
      %s %s %s %s %s %s
    </div>
    <div class="bg-white rounded-xl border border-gray-200 shadow-sm p-6 space-y-4">
      <h3 class="font-semibold text-gray-800 text-base">Social Links</h3>
      %s %s %s %s %s
    </div>
    <div class="bg-white rounded-xl border border-gray-200 shadow-sm p-6 space-y-4">
      <h3 class="font-semibold text-gray-800 text-base">Settings</h3>
      %s
    </div>
    <button type="submit" class="bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-medium px-6 py-2.5 rounded-lg transition">
      Save Profile
    </button>
  </form>
</div>`,
			inputField("Full Name *", "full_name", escapeHTML(p.FullName), "text", true),
			inputField("Title / Headline", "title", escapeHTML(p.Title), "text", false),
			textareaField("Bio", "bio", escapeHTML(p.Bio), 4),
			inputField("Avatar URL", "avatar_url", escapeHTML(p.AvatarURL), "url", false),
			inputField("Email", "email", escapeHTML(p.Email), "email", false),
			inputField("Phone", "phone", escapeHTML(p.Phone), "tel", false),
			inputField("GitHub URL", "github_url", escapeHTML(p.GithubURL), "url", false),
			inputField("LinkedIn URL", "linkedin_url", escapeHTML(p.LinkedinURL), "url", false),
			inputField("Twitter URL", "twitter_url", escapeHTML(p.TwitterURL), "url", false),
			inputField("Website URL", "website_url", escapeHTML(p.WebsiteURL), "url", false),
			inputField("Resume URL", "resume_url", escapeHTML(p.ResumeURL), "url", false),
			toggleField("Published", "is_published", p.IsPublished),
		)

		flash := c.Cookies("flash")
		c.ClearCookie("flash")
		return c.Type("html").SendString(layout("Profile", "Profile", form, flash))
	}
}

func UpsertProfileHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var p model.Profile
		isNew := db.First(&p).Error != nil

		data := map[string]interface{}{
			"full_name":    c.FormValue("full_name"),
			"title":        c.FormValue("title"),
			"bio":          c.FormValue("bio"),
			"avatar_url":   c.FormValue("avatar_url"),
			"email":        c.FormValue("email"),
			"phone":        c.FormValue("phone"),
			"github_url":   c.FormValue("github_url"),
			"linkedin_url": c.FormValue("linkedin_url"),
			"twitter_url":  c.FormValue("twitter_url"),
			"website_url":  c.FormValue("website_url"),
			"resume_url":   c.FormValue("resume_url"),
			"is_published": c.FormValue("is_published") == "true",
			"updated_at":   time.Now(),
		}

		if isNew {
			p = model.Profile{
				FullName: c.FormValue("full_name"), Title: c.FormValue("title"),
				Bio: c.FormValue("bio"), AvatarURL: c.FormValue("avatar_url"),
				Email: c.FormValue("email"), Phone: c.FormValue("phone"),
				GithubURL: c.FormValue("github_url"), LinkedinURL: c.FormValue("linkedin_url"),
				TwitterURL: c.FormValue("twitter_url"), WebsiteURL: c.FormValue("website_url"),
				ResumeURL: c.FormValue("resume_url"),
				IsPublished: c.FormValue("is_published") == "true",
			}
			db.Create(&p)
		} else {
			db.Model(&p).Updates(data)
		}

		setFlash(c, "Profile saved!")
		return c.Redirect("/admin/profile")
	}
}

// ─── Contacts ──────────────────────────────────────────────────────────────────

func ContactsPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var contacts []model.Contact
		db.Order("created_at desc").Find(&contacts)

		rows := ""
		for _, ct := range contacts {
			readBadge := badge("Unread", "blue")
			if ct.IsRead {
				readBadge = badge("Read", "gray")
			}
			rows += fmt.Sprintf(`
<tr id="contact-%s" class="border-t border-gray-100 hover:bg-gray-50 transition">
  <td class="px-4 py-3 text-sm font-medium text-gray-800">%s</td>
  <td class="px-4 py-3 text-sm text-gray-500">%s</td>
  <td class="px-4 py-3 text-sm text-gray-500">%s</td>
  <td class="px-4 py-3 text-sm text-gray-400 max-w-xs truncate">%s</td>
  <td class="px-4 py-3 text-sm">%s</td>
  <td class="px-4 py-3 text-sm text-gray-400">%s</td>
  <td class="px-4 py-3 text-sm">
    <div class="flex items-center gap-3">
      %s
      %s
    </div>
  </td>
</tr>`,
				ct.ID.String(),
				escapeHTML(ct.Name),
				escapeHTML(ct.Email),
				escapeHTML(ct.Subject),
				escapeHTML(ct.Message),
				readBadge,
				ct.CreatedAt.Format("02 Jan 2006 15:04"),
				func() string {
					if !ct.IsRead {
						return fmt.Sprintf(`<button
  hx-post="/admin/contacts/%s/read"
  hx-target="#contact-%s"
  hx-swap="outerHTML"
  class="text-blue-500 hover:text-blue-700 text-sm font-medium transition">Mark Read</button>`, ct.ID.String(), ct.ID.String())
					}
					return ""
				}(),
				deleteBtn("/admin/contacts/"+ct.ID.String(), "Delete this message?"),
			)
		}

		var unreadCount int64
		db.Model(&model.Contact{}).Where("is_read = false").Count(&unreadCount)

		content := fmt.Sprintf(`
<div class="flex items-center justify-between mb-6">
  <div>
    <h3 class="text-lg font-semibold text-gray-800">Contact Messages</h3>
    <p class="text-sm text-gray-500">%d total, %d unread</p>
  </div>
</div>
%s`,
			len(contacts), unreadCount,
			tableWrapper(fmt.Sprintf(`
<table class="w-full text-left">
  <thead class="bg-gray-50 border-b border-gray-200">
    <tr>
      <th class="px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Name</th>
      <th class="px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Email</th>
      <th class="px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Subject</th>
      <th class="px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Message</th>
      <th class="px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Status</th>
      <th class="px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Received</th>
      <th class="px-4 py-3 text-xs font-semibold text-gray-500 uppercase">Actions</th>
    </tr>
  </thead>
  <tbody>%s</tbody>
</table>`, rows)),
		)

		flash := c.Cookies("flash")
		c.ClearCookie("flash")
		return c.Type("html").SendString(layout("Contacts", "Contacts", content, flash))
	}
}

func MarkContactReadHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		db.Model(&model.Contact{}).Where("id = ?", id).Update("is_read", true)

		// Return the updated row for HTMX swap
		var ct model.Contact
		db.First(&ct, "id = ?", id)
		return c.Type("html").SendString(fmt.Sprintf(`
<tr id="contact-%s" class="border-t border-gray-100 hover:bg-gray-50 transition">
  <td class="px-4 py-3 text-sm font-medium text-gray-800">%s</td>
  <td class="px-4 py-3 text-sm text-gray-500">%s</td>
  <td class="px-4 py-3 text-sm text-gray-500">%s</td>
  <td class="px-4 py-3 text-sm text-gray-400 max-w-xs truncate">%s</td>
  <td class="px-4 py-3 text-sm">%s</td>
  <td class="px-4 py-3 text-sm text-gray-400">%s</td>
  <td class="px-4 py-3 text-sm">%s</td>
</tr>`,
			ct.ID.String(),
			escapeHTML(ct.Name), escapeHTML(ct.Email),
			escapeHTML(ct.Subject), escapeHTML(ct.Message),
			badge("Read", "gray"),
			ct.CreatedAt.Format("02 Jan 2006 15:04"),
			deleteBtn("/admin/contacts/"+ct.ID.String(), "Delete this message?"),
		))
	}
}

func DeleteContactHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		db.Delete(&model.Contact{}, "id = ?", c.Params("id"))
		return c.SendString("")
	}
}
