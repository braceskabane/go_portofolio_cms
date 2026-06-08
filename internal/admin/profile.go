package admin

import (
	"fmt"
	"portfolio-cms/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ── Profile Page (View & Edit) ───────────────────────────────────────────────

func ProfilePage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var profile model.Profile
		db.First(&profile)

		isNew := profile.ID == uuid.Nil
		title := "Edit Profile"
		submitLabel := "Save Profile"
		if isNew {
			title = "Create Profile"
			submitLabel = "Create Profile"
		}

		flash := c.Query("flash")
		notifHTML := ""
		if flash == "updated" {
			notifHTML = notifBanner("success", "✓ Profile berhasil diperbarui.")
		} else if flash == "created" {
			notifHTML = notifBanner("success", "✓ Profile berhasil dibuat.")
		}

		// Personal Info
		personalFields := formGrid(
			inputField("Full Name", "full_name", profile.FullName, "text", true),
			inputField("Title", "title", profile.Title, "text", false),
			inputField("Email", "email", profile.Email, "email", false),
			inputField("Phone", "phone", profile.Phone, "text", false),
			inputField("Location", "location", profile.Location, "text", false),
			inputField("Website URL", "website_url", profile.WebsiteURL, "url", false),
			inputField("Avatar URL", "avatar_url", profile.AvatarURL, "url", false),
			inputField("Resume URL", "resume_url", profile.ResumeURL, "url", false),
		)

		// About & Bio
		aboutBio := formGrid(
			textareaField("About", "about", profile.About, 6),
			textareaField("Bio", "bio", profile.Bio, 4),
		)

		// Social Links
		socialFields := formGrid(
			inputField("GitHub URL", "github_url", profile.GithubURL, "url", false),
			inputField("LinkedIn URL", "linkedin_url", profile.LinkedinURL, "url", false),
			inputField("Twitter URL", "twitter_url", profile.TwitterURL, "url", false),
			inputField("Instagram URL", "instagram_url", profile.InstagramURL, "url", false),
			inputField("TikTok URL", "tiktok_url", profile.TiktokURL, "url", false),
			inputField("Strava URL", "strava_url", profile.StravaURL, "url", false),
			"",
			"",
		)

		// Settings
		settings := formGrid(
			toggleSwitch("Available for Hire", "available_for_hire", profile.AvailableForHire, "Tampilkan status tersedia"),
			inputField("Years Experience", "years_experience", fmt.Sprintf("%d", profile.YearsExperience), "number", false),
			toggleSwitch("Published", "is_published", profile.IsPublished, "Tampilkan profile ke publik"),
			"",
		)

		content := fmt.Sprintf(`
%s
<div style="display:flex; justify-content: flex-end; margin-bottom: 16px;">
  <a href="/admin" style="display:inline-flex;align-items:center;gap:6px;font-size:13px;color:#4d4c4d;text-decoration:none;padding:6px 12px;border:1px solid #e5e3e4;border-radius:6px;">← Back to Dashboard</a>
</div>
<form method="POST" action="/admin/profile">
  %s
  %s
  %s
  %s
  <div style="display:flex;justify-content:flex-end;gap:12px;margin-top:16px;padding:16px 0;">
    <a href="/admin" style="padding:8px 16px;border:1px solid #e5e3e4;border-radius:6px;color:#4d4c4d;text-decoration:none;">Batal</a>
    <button type="submit" style="padding:8px 20px;background:#1e1d1e;color:#f5f4f5;border:none;border-radius:6px;font-weight:500;">%s</button>
  </div>
</form>`,
			notifHTML,
			sectionCard("Personal Info", personalFields),
			sectionCard("About & Bio", aboutBio),
			sectionCard("Social Links", socialFields),
			sectionCard("Settings", settings),
			submitLabel,
		)

		return c.Type("html").SendString(layout(title, "Profile", content))
	}
}

// ── Handler for Create/Update (Upsert) ───────────────────────────────────────

func UpsertProfileHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var profile model.Profile
		db.First(&profile)

		profile.FullName = c.FormValue("full_name")
		profile.Title = c.FormValue("title")
		profile.About = c.FormValue("about")
		profile.Bio = c.FormValue("bio")
		profile.AvatarURL = c.FormValue("avatar_url")
		profile.Location = c.FormValue("location")
		profile.Email = c.FormValue("email")
		profile.Phone = c.FormValue("phone")
		profile.GithubURL = c.FormValue("github_url")
		profile.LinkedinURL = c.FormValue("linkedin_url")
		profile.TwitterURL = c.FormValue("twitter_url")
		profile.InstagramURL = c.FormValue("instagram_url")
		profile.TiktokURL = c.FormValue("tiktok_url")
		profile.StravaURL = c.FormValue("strava_url")
		profile.WebsiteURL = c.FormValue("website_url")
		profile.ResumeURL = c.FormValue("resume_url")
		profile.AvailableForHire = c.FormValue("available_for_hire") == "true"
		profile.YearsExperience, _ = c.ParamsInt("years_experience", 0)
		profile.IsPublished = c.FormValue("is_published") == "true"

		if profile.ID == uuid.Nil {
			if err := db.Create(&profile).Error; err != nil {
				return c.Status(500).SendString("Failed to create profile")
			}
			return c.Redirect("/admin/profile?flash=created")
		}
		db.Save(&profile)
		return c.Redirect("/admin/profile?flash=updated")
	}
}