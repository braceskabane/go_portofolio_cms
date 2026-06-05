package admin

import (
	"fmt"
	"portfolio-cms/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// EducationsPage menampilkan daftar education dan tombol + Add
func EducationsPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var items []model.Education
		db.Order("sort_order asc, start_date desc").Find(&items)

		rows := ""
		for _, e := range items {
			rows += fmt.Sprintf(`
<tr style="border-top:1px solid #e5e3e4;">
  <td style="padding:12px 16px;font-size:13px;font-weight:500;color:#1e1d1e;">%s</td>
  <td style="padding:12px 16px;font-size:13px;color:#4d4c4d;">%s</td>
  <td style="padding:12px 16px;font-size:13px;color:#9a9899;">%s</td>
  <td style="padding:12px 16px;font-size:13px;">%s</td>
  <td style="padding:12px 16px;font-size:13px;">
    <div class="flex items-center gap-3">
      <a href="/admin/educations/%s/edit" class="text-indigo-600 hover:text-indigo-800 font-medium transition">Edit</a>
      %s
    </div>
  </td>
</tr>`,
				escapeHTML(e.Institution), escapeHTML(e.Degree),
				e.StartDate.Format("Jan 2006"),
				yesNo(e.IsPublished),
				e.ID.String(),
				deleteBtn("/admin/educations/"+e.ID.String(), "Delete this education?"),
			)
		}

		tableHTML := fmt.Sprintf(`
<table style="width:100%%; border-collapse:collapse;">
  <thead style="background:#f5f4f5; border-bottom:1px solid #e5e3e4;">
    <tr>
      <th style="padding:12px 16px; font-size:11.5px; font-weight:600; color:#4d4c4d; text-transform:uppercase; text-align:left;">Institution</th>
      <th style="padding:12px 16px; font-size:11.5px; font-weight:600; color:#4d4c4d; text-transform:uppercase; text-align:left;">Degree</th>
      <th style="padding:12px 16px; font-size:11.5px; font-weight:600; color:#4d4c4d; text-transform:uppercase; text-align:left;">Start Date</th>
      <th style="padding:12px 16px; font-size:11.5px; font-weight:600; color:#4d4c4d; text-transform:uppercase; text-align:left;">Published</th>
      <th style="padding:12px 16px; font-size:11.5px; font-weight:600; color:#4d4c4d; text-transform:uppercase; text-align:left;">Actions</th>
    </tr>
  </thead>
  <tbody>%s</tbody>
</table>`, rows)

		content := fmt.Sprintf(`
<div style="display:flex; align-items:center; justify-content:space-between; margin-bottom:24px;">
  <div>
    <h3 style="font-size:16px; font-weight:600; color:#1e1d1e;">All Educations</h3>
    <p style="font-size:13px; color:#6b6a6b;">%d educations total</p>
  </div>
  <a href="/admin/educations/new"
     style="background:#1e1d1e; color:#f5f4f5; text-decoration:none; font-size:13px; font-weight:500; padding:8px 16px; border-radius:6px; display:inline-flex; align-items:center; gap:6px; transition:background 0.12s;"
     onmouseover="this.style.background='#30353b'"
     onmouseout="this.style.background='#1e1d1e'">
     + Add Education
  </a>
</div>
%s`, len(items), tableWrapper(tableHTML))

		flash := c.Cookies("flash")
		c.ClearCookie("flash")
		return c.Type("html").SendString(layout("Educations", "Educations", content, flash))
	}
}

// EducationFormPage menampilkan form create/edit education
func EducationFormPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var e model.Education
		isEdit := id != "" && id != "new"

		if isEdit {
			if err := db.First(&e, "id = ?", id).Error; err != nil {
				return c.Redirect("/admin/educations")
			}
		}

		title := "New Education"
		action := "/admin/educations"
		saveLabel := "Save Education"
		if isEdit {
			title = "Edit Education"
			action = "/admin/educations/" + id
			saveLabel = "Update Education"
		}

		startDate := ""
		if isEdit {
			startDate = e.StartDate.Format("2006-01-02")
		}

		formContent := fmt.Sprintf(`
			<div style="margin-top:12px;">
				%s
				<form method="POST" action="%s" style="display:flex; flex-direction:column; gap:12px;">
					%s %s %s %s %s %s
					<div style="display:flex; align-items:center; gap:12px; margin-top:8px;">
						%s
						<a href="/admin/educations" style="font-size:13px; color:#6b6a6b; text-decoration:none; padding:8px 16px; border-radius:6px; border:1px solid #e5e3e4;">Cancel</a>
					</div>
				</form>
			</div>`,
			card(`<p style="font-size:13px; color:#4d4c4d;">Fill education details.</p>`),
			action,
			inputField("Institution *", "institution", escapeHTML(e.Institution), "text", true),
			inputField("Degree", "degree", escapeHTML(e.Degree), "text", false),
			inputField("Field of Study", "field_of_study", escapeHTML(e.FieldOfStudy), "text", false),
			inputField("GPA", "gpa", e.GPA, "text", false),
			inputField("Start Date *", "start_date", startDate, "date", true),
			toggleField("Published", "is_published", e.IsPublished),
			btnPrimary(saveLabel),
		)

		return c.Type("html").SendString(layout(title, "Educations", formContent))
	}
}

func CreateEducationHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startDate, err := time.Parse("2006-01-02", c.FormValue("start_date"))
		if err != nil {
			setFlash(c, "ERR:Invalid start date")
			return c.Redirect("/admin/educations")
		}
		e := model.Education{
			Institution:  c.FormValue("institution"),
			Degree:       c.FormValue("degree"),
			FieldOfStudy: c.FormValue("field_of_study"),
			GPA:          c.FormValue("gpa"),
			StartDate:    startDate,
			IsPublished:  c.FormValue("is_published") == "true",
		}
		if err := db.Create(&e).Error; err != nil {
			setFlash(c, "ERR:"+err.Error())
		} else {
			setFlash(c, "Education created!")
		}
		return c.Redirect("/admin/educations")
	}
}

func UpdateEducationHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		db.Model(&model.Education{}).Where("id = ?", id).Updates(map[string]interface{}{
			"institution":    c.FormValue("institution"),
			"degree":         c.FormValue("degree"),
			"field_of_study": c.FormValue("field_of_study"),
			"is_published":   c.FormValue("is_published") == "true",
			"updated_at":     time.Now(),
		})
		setFlash(c, "Education updated!")
		return c.Redirect("/admin/educations")
	}
}

func DeleteEducationHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		db.Delete(&model.Education{}, "id = ?", c.Params("id"))
		return c.SendString("")
	}
}