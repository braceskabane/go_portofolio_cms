package admin

import (
	"fmt"
	"portfolio-cms/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SkillsPage menampilkan daftar skill dan tombol + Add
func SkillsPage(db *gorm.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var skills []model.Skill
        db.Order("sort_order asc, category asc").Find(&skills)

        rows := ""
        for _, s := range skills {
            rows += fmt.Sprintf(`
<tr style="border-top:1px solid #e5e3e4;">
  <td style="padding:12px 16px;font-size:13px;font-weight:500;color:#1e1d1e;">%s</td>
  <td style="padding:12px 16px;font-size:13px;color:#4d4c4d;">%s</td>
  <td style="padding:12px 16px;font-size:13px;color:#4d4c4d;">%d%%</td>
  <td style="padding:12px 16px;font-size:13px;">%s</td>
  <td style="padding:12px 16px;font-size:13px;color:#9a9899;">%d</td>
  <td style="padding:12px 16px;font-size:13px;">
    <div style="display:flex; gap:12px;">
      <a href="/admin/skills/%s/edit" class="text-indigo-600 hover:text-indigo-800 font-medium transition">Edit</a>
      %s
    </div>
  </td>
</tr>`,
                escapeHTML(s.Name), escapeHTML(s.Category), s.Proficiency,
                yesNo(s.IsPublished), s.SortOrder,
                s.ID.String(),
                deleteBtn("/admin/skills/"+s.ID.String(), "Delete this skill?"),
            )
        }

        tableHTML := fmt.Sprintf(`
<table style="width:100%%; border-collapse:collapse;">
  <thead style="background:#f5f4f5; border-bottom:1px solid #e5e3e4;">
    <tr>
      <th style="padding:12px 16px; font-size:11.5px; font-weight:600; color:#4d4c4d; text-transform:uppercase; text-align:left;">Name</th>
      <th style="padding:12px 16px; font-size:11.5px; font-weight:600; color:#4d4c4d; text-transform:uppercase; text-align:left;">Category</th>
      <th style="padding:12px 16px; font-size:11.5px; font-weight:600; color:#4d4c4d; text-transform:uppercase; text-align:left;">Proficiency</th>
      <th style="padding:12px 16px; font-size:11.5px; font-weight:600; color:#4d4c4d; text-transform:uppercase; text-align:left;">Published</th>
      <th style="padding:12px 16px; font-size:11.5px; font-weight:600; color:#4d4c4d; text-transform:uppercase; text-align:left;">Order</th>
      <th style="padding:12px 16px; font-size:11.5px; font-weight:600; color:#4d4c4d; text-transform:uppercase; text-align:left;">Actions</th>
    </tr>
  </thead>
  <tbody>%s</tbody>
</table>`, rows)

        content := fmt.Sprintf(`
<div style="display:flex; align-items:center; justify-content:space-between; margin-bottom:24px;">
  <div>
    <h3 style="font-size:16px; font-weight:600; color:#1e1d1e;">All Skills</h3>
    <p style="font-size:13px; color:#6b6a6b;">%d skills total</p>
  </div>
  <a href="/admin/skills/new"
     style="background:#1e1d1e; color:#f5f4f5; text-decoration:none; font-size:13px; font-weight:500; padding:8px 16px; border-radius:6px; display:inline-flex; align-items:center; gap:6px; transition:background 0.12s;"
     onmouseover="this.style.background='#30353b'"
     onmouseout="this.style.background='#1e1d1e'">
     + Add Skill
  </a>
</div>
%s`, len(skills), tableWrapper(tableHTML))

        flash := c.Cookies("flash")
        c.ClearCookie("flash")
        return c.Type("html").SendString(layout("Skills", "Skills", content, flash))
    }
}

// SkillFormPage menampilkan form create/edit skill
func SkillFormPage(db *gorm.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        id := c.Params("id")
        var s model.Skill
        isEdit := id != "" && id != "new"

        if isEdit {
            if err := db.First(&s, "id = ?", id).Error; err != nil {
                return c.Redirect("/admin/skills")
            }
        }

        title := "New Skill"
        action := "/admin/skills"
        saveLabel := "Save Skill"
        if isEdit {
            title = "Edit Skill"
            action = "/admin/skills/" + id
            saveLabel = "Update Skill"
        }

        formContent := fmt.Sprintf(`
            <div style="margin-top:12px;">
                %s
                <form method="POST" action="%s" style="display:flex; flex-direction:column; gap:12px;">
                    %s %s %s %s %s
                    <div style="display:flex; align-items:center; gap:12px; margin-top:8px;">
                        %s
                        <a href="/admin/skills" style="font-size:13px; color:#6b6a6b; text-decoration:none; padding:8px 16px; border-radius:6px; border:1px solid #e5e3e4;">Cancel</a>
                    </div>
                </form>
            </div>`,
            card(`<p style="font-size:13px; color:#4d4c4d;">Fill skill details below.</p>`),
            action,
            inputField("Name *", "name", escapeHTML(s.Name), "text", true),
            inputField("Category", "category", escapeHTML(s.Category), "text", false),
            inputField("Proficiency (0-100)", "proficiency", fmt.Sprintf("%d", s.Proficiency), "number", false),
            inputField("Sort Order", "sort_order", fmt.Sprintf("%d", s.SortOrder), "number", false),
            toggleField("Published", "is_published", s.IsPublished),
            btnPrimary(saveLabel),
        )

        return c.Type("html").SendString(layout(title, "Skills", formContent))
    }
}

func CreateSkillHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		proficiency := 80
		fmt.Sscanf(c.FormValue("proficiency"), "%d", &proficiency)
		sortOrder := 0
		fmt.Sscanf(c.FormValue("sort_order"), "%d", &sortOrder)

		s := model.Skill{
			Name:        c.FormValue("name"),
			Category:    c.FormValue("category"),
			Proficiency: proficiency,
			SortOrder:   sortOrder,
			IsPublished: c.FormValue("is_published") == "true",
		}
		if err := db.Create(&s).Error; err != nil {
			setFlash(c, "ERR:"+err.Error())
		} else {
			setFlash(c, "Skill created!")
		}
		return c.Redirect("/admin/skills")
	}
}

func UpdateSkillHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		proficiency := 80
		fmt.Sscanf(c.FormValue("proficiency"), "%d", &proficiency)
		sortOrder := 0
		fmt.Sscanf(c.FormValue("sort_order"), "%d", &sortOrder)

		db.Model(&model.Skill{}).Where("id = ?", id).Updates(map[string]interface{}{
			"name":         c.FormValue("name"),
			"category":     c.FormValue("category"),
			"proficiency":  proficiency,
			"sort_order":   sortOrder,
			"is_published": c.FormValue("is_published") == "true",
			"updated_at":   time.Now(),
		})
		setFlash(c, "Skill updated!")
		return c.Redirect("/admin/skills")
	}
}

func DeleteSkillHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		db.Delete(&model.Skill{}, "id = ?", c.Params("id"))
		return c.SendString("")
	}
}