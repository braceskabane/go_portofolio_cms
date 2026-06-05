package admin

import (
	"fmt"
	"portfolio-cms/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ExperiencesPage menampilkan daftar experience dan tombol + Add
func ExperiencesPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var items []model.Experience
		db.Order("sort_order asc, start_date desc").Find(&items)

		rows := ""
		for _, e := range items {
			endDate := "Present"
			if e.EndDate != nil {
				endDate = e.EndDate.Format("Jan 2006")
			}
			rows += fmt.Sprintf(`
<tr style="border-top:1px solid #e5e3e4;">
  <td style="padding:12px 16px;font-size:13px;font-weight:500;color:#1e1d1e;">%s</td>
  <td style="padding:12px 16px;font-size:13px;color:#4d4c4d;">%s</td>
  <td style="padding:12px 16px;font-size:13px;color:#9a9899;">%s – %s</td>
  <td style="padding:12px 16px;font-size:13px;">%s</td>
  <td style="padding:12px 16px;font-size:13px;">
    <div style="display:flex; gap:12px;">
      <a href="/admin/experiences/%s/edit" class="text-indigo-600 hover:text-indigo-800 font-medium transition">Edit</a>
      %s
    </div>
  </td>
</tr>`,
				escapeHTML(e.Company), escapeHTML(e.Position),
				e.StartDate.Format("Jan 2006"), endDate,
				yesNo(e.IsPublished),
				e.ID.String(),
				deleteBtn("/admin/experiences/"+e.ID.String(), "Delete this experience?"),
			)
		}

		tableHTML := fmt.Sprintf(`
<table style="width:100%%; border-collapse:collapse;">
  <thead style="background:#f5f4f5; border-bottom:1px solid #e5e3e4;">
    <tr>
      <th style="padding:12px 16px; font-size:11.5px; font-weight:600; color:#4d4c4d; text-transform:uppercase; text-align:left;">Company</th>
      <th style="padding:12px 16px; font-size:11.5px; font-weight:600; color:#4d4c4d; text-transform:uppercase; text-align:left;">Position</th>
      <th style="padding:12px 16px; font-size:11.5px; font-weight:600; color:#4d4c4d; text-transform:uppercase; text-align:left;">Period</th>
      <th style="padding:12px 16px; font-size:11.5px; font-weight:600; color:#4d4c4d; text-transform:uppercase; text-align:left;">Published</th>
      <th style="padding:12px 16px; font-size:11.5px; font-weight:600; color:#4d4c4d; text-transform:uppercase; text-align:left;">Actions</th>
    </tr>
  </thead>
  <tbody>%s</tbody>
</table>`, rows)

		content := fmt.Sprintf(`
<div style="display:flex; align-items:center; justify-content:space-between; margin-bottom:24px;">
  <div>
    <h3 style="font-size:16px; font-weight:600; color:#1e1d1e;">All Experiences</h3>
    <p style="font-size:13px; color:#6b6a6b;">%d experiences total</p>
  </div>
  <a href="/admin/experiences/new"
     style="background:#1e1d1e; color:#f5f4f5; text-decoration:none; font-size:13px; font-weight:500; padding:8px 16px; border-radius:6px; display:inline-flex; align-items:center; gap:6px; transition:background 0.12s;"
     onmouseover="this.style.background='#30353b'"
     onmouseout="this.style.background='#1e1d1e'">
     + Add Experience
  </a>
</div>
%s`, len(items), tableWrapper(tableHTML))

		flash := c.Cookies("flash")
		c.ClearCookie("flash")
		return c.Type("html").SendString(layout("Experiences", "Experiences", content, flash))
	}
}

// ExperienceFormPage menampilkan form create/edit experience
func ExperienceFormPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var e model.Experience
		isEdit := id != "" && id != "new"

		if isEdit {
			if err := db.First(&e, "id = ?", id).Error; err != nil {
				return c.Redirect("/admin/experiences")
			}
		}

		title := "New Experience"
		action := "/admin/experiences"
		saveLabel := "Save Experience"
		if isEdit {
			title = "Edit Experience"
			action = "/admin/experiences/" + id
			saveLabel = "Update Experience"
		}

		// start/end date bisa diisi jika edit
		startDate := ""
		endDate := ""
		if isEdit {
			startDate = e.StartDate.Format("2006-01-02")
			if e.EndDate != nil {
				endDate = e.EndDate.Format("2006-01-02")
			}
		}

		formContent := fmt.Sprintf(`
			<div style="margin-top:12px;">
				%s
				<form method="POST" action="%s" style="display:flex; flex-direction:column; gap:12px;">
					%s %s %s %s %s %s %s
					<div style="display:flex; align-items:center; gap:12px; margin-top:8px;">
						%s
						<a href="/admin/experiences" style="font-size:13px; color:#6b6a6b; text-decoration:none; padding:8px 16px; border-radius:6px; border:1px solid #e5e3e4;">Cancel</a>
					</div>
				</form>
			</div>`,
			card(`<p style="font-size:13px; color:#4d4c4d;">Fill experience details.</p>`),
			action,
			inputField("Company *", "company", escapeHTML(e.Company), "text", true),
			inputField("Position *", "position", escapeHTML(e.Position), "text", true),
			inputField("Location", "location", escapeHTML(e.Location), "text", false),
			inputField("Start Date *", "start_date", startDate, "date", true),
			inputField("End Date", "end_date", endDate, "date", false),
			toggleField("Currently Working Here", "is_current", e.IsCurrent),
			toggleField("Published", "is_published", e.IsPublished),
			btnPrimary(saveLabel),
		)

		return c.Type("html").SendString(layout(title, "Experiences", formContent))
	}
}

func CreateExperienceHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startDate, err := time.Parse("2006-01-02", c.FormValue("start_date"))
		if err != nil {
			setFlash(c, "ERR:Invalid start date format")
			return c.Redirect("/admin/experiences")
		}
		var endDate *time.Time
		if ed := c.FormValue("end_date"); ed != "" {
			t, err := time.Parse("2006-01-02", ed)
			if err == nil {
				endDate = &t
			}
		}
		e := model.Experience{
			Company:     c.FormValue("company"),
			Position:    c.FormValue("position"),
			Location:    c.FormValue("location"),
			StartDate:   startDate,
			EndDate:     endDate,
			IsCurrent:   c.FormValue("is_current") == "true",
			IsPublished: c.FormValue("is_published") == "true",
		}
		if err := db.Create(&e).Error; err != nil {
			setFlash(c, "ERR:"+err.Error())
		} else {
			setFlash(c, "Experience created!")
		}
		return c.Redirect("/admin/experiences")
	}
}

func UpdateExperienceHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		db.Model(&model.Experience{}).Where("id = ?", id).Updates(map[string]interface{}{
			"company":      c.FormValue("company"),
			"position":     c.FormValue("position"),
			"is_current":   c.FormValue("is_current") == "true",
			"is_published": c.FormValue("is_published") == "true",
			"updated_at":   time.Now(),
		})
		setFlash(c, "Experience updated!")
		return c.Redirect("/admin/experiences")
	}
}

func DeleteExperienceHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		db.Delete(&model.Experience{}, "id = ?", c.Params("id"))
		return c.SendString("")
	}
}