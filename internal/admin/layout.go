package admin

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"strings"
)

//go:embed templates/*.html
var templateFS embed.FS

var layoutTmpl = template.Must(
	template.ParseFS(templateFS, "templates/layout.html"),
)

// navItem defines a sidebar navigation entry.
type navItem struct {
	Href  string
	Icon  string
	Label string
}

var defaultNav = []navItem{
	{"/admin", "◈", "Dashboard"},
	{"/admin/projects", "▤", "Projects"},
	{"/admin/skills", "◎", "Skills"},
	{"/admin/experiences", "◉", "Experiences"},
	{"/admin/educations", "◫", "Educations"},
	{"/admin/profile", "◷", "Profile"},
	{"/admin/contacts", "◻", "Contacts"},
}

// layoutData is the template context passed to layout.html.
type layoutData struct {
	Title        string
	ActiveMenu   string
	Flash        string
	FlashIsError bool
	Nav          []navItem
	Content      template.HTML
}

// layout renders the full admin shell and returns the HTML string.
func layout(title, activeMenu, content string, flashMsg ...string) string {
	flash := ""
	flashIsError := false

	if len(flashMsg) > 0 && flashMsg[0] != "" {
		flash = flashMsg[0]
		if strings.HasPrefix(flash, "ERR:") {
			flashIsError = true
			flash = strings.TrimPrefix(flash, "ERR:")
		}
	}

	data := layoutData{
		Title:        title,
		ActiveMenu:   activeMenu,
		Flash:        flash,
		FlashIsError: flashIsError,
		Nav:          defaultNav,
		Content:      template.HTML(content), // #nosec — content generated internally
	}

	var buf bytes.Buffer
	if err := layoutTmpl.Execute(&buf, data); err != nil {
		return fmt.Sprintf(`<!DOCTYPE html><html><body><pre>template error: %s</pre></body></html>`, err.Error())
	}
	return buf.String()
}

// ── Component helpers ─────────────────────────────────────────────────────────
// Semua helper menggunakan gaya inline dengan palet Pure Charcoal:
//   #0f0e0f (950)  #161516 (900)  #1e1d1e (800)
//   #252425 (700)  #30353b (600)  #4d4c4d (500)
//   #6b6a6b (400)  #9a9899 (300)  #c5c3c4 (200)
//   #e8e6e7 (100)  #f5f4f5 (50)

// card wraps content in a clean white card.
func card(content string) string {
	return fmt.Sprintf(`<div style="background:#ffffff;border-radius:8px;border:1px solid #e5e3e4;padding:20px 24px;margin-bottom:16px;">%s</div>`, content)
}

// tableWrapper wraps a <table> in a styled overflow container.
func tableWrapper(content string) string {
	return fmt.Sprintf(
		`<div style="background:#ffffff;border-radius:8px;border:1px solid #e5e3e4;overflow:hidden;"><div style="overflow-x:auto;">%s</div></div>`,
		content,
	)
}

// deleteBtn renders an HTMX delete button targeting the closest <tr>.
func deleteBtn(url, confirmMsg string) string {
	return fmt.Sprintf(`<button
    hx-delete="%s"
    hx-confirm="%s"
    hx-target="closest tr"
    hx-swap="outerHTML swap:0.25s"
    style="font-size:12.5px;color:#b91c1c;font-weight:500;background:none;border:none;cursor:pointer;padding:0;transition:color 0.1s;"
    onmouseover="this.style.color='#7f1d1d'"
    onmouseout="this.style.color='#b91c1c'">Delete</button>`,
		url, confirmMsg)
}

// inputField renders a labeled <input>.
func inputField(label, name, value, inputType string, required bool) string {
	req := ""
	if required {
		req = "required"
	}
	if inputType == "" {
		inputType = "text"
	}
	return fmt.Sprintf(`
<div style="margin-bottom:0;">
  <label style="display:block;font-size:12.5px;font-weight:500;color:#4d4c4d;margin-bottom:5px;">%s</label>
  <input type="%s" name="%s" value="%s" %s
    style="width:100%%;border:1px solid #d1cfd0;border-radius:6px;padding:7px 11px;font-size:13px;color:#1e1d1e;background:#ffffff;box-sizing:border-box;transition:border-color 0.1s;"
    onfocus="this.style.borderColor='#4d4c4d';this.style.boxShadow='0 0 0 2px #9a9899'"
    onblur="this.style.borderColor='#d1cfd0';this.style.boxShadow='none'"/>
</div>`, label, inputType, name, value, req)
}

// textareaField renders a labeled <textarea>.
func textareaField(label, name, value string, rows int) string {
	if rows == 0 {
		rows = 4
	}
	return fmt.Sprintf(`
<div style="margin-bottom:0;">
  <label style="display:block;font-size:12.5px;font-weight:500;color:#4d4c4d;margin-bottom:5px;">%s</label>
  <textarea name="%s" rows="%d"
    style="width:100%%;border:1px solid #d1cfd0;border-radius:6px;padding:7px 11px;font-size:13px;color:#1e1d1e;background:#ffffff;box-sizing:border-box;resize:vertical;font-family:inherit;transition:border-color 0.1s;"
    onfocus="this.style.borderColor='#4d4c4d';this.style.boxShadow='0 0 0 2px #9a9899'"
    onblur="this.style.borderColor='#d1cfd0';this.style.boxShadow='none'">%s</textarea>
</div>`, label, name, rows, value)
}

// toggleField renders a Yes/No <select>.
func toggleField(label, name string, value bool) string {
	trueSelected, falseSelected := "", ""
	if value {
		trueSelected = "selected"
	} else {
		falseSelected = "selected"
	}
	return fmt.Sprintf(`
<div style="margin-bottom:0;">
  <label style="display:block;font-size:12.5px;font-weight:500;color:#4d4c4d;margin-bottom:5px;">%s</label>
  <select name="%s"
    style="width:100%%;border:1px solid #d1cfd0;border-radius:6px;padding:7px 11px;font-size:13px;color:#1e1d1e;background:#ffffff;box-sizing:border-box;cursor:pointer;transition:border-color 0.1s;"
    onfocus="this.style.borderColor='#4d4c4d'"
    onblur="this.style.borderColor='#d1cfd0'">
    <option value="true" %s>Yes</option>
    <option value="false" %s>No</option>
  </select>
</div>`, label, name, trueSelected, falseSelected)
}

// badge renders a small status pill.
func badge(text, color string) string {
	styles := map[string]string{
		"green": "background:#f0fdf4;color:#166534;border:1px solid #bbf7d0;",
		"red":   "background:#fef2f2;color:#991b1b;border:1px solid #fecaca;",
		"gray":  "background:#f5f4f5;color:#4d4c4d;border:1px solid #e5e3e4;",
		"blue":  "background:#eff6ff;color:#1e40af;border:1px solid #bfdbfe;",
	}
	st, ok := styles[color]
	if !ok {
		st = styles["gray"]
	}
	return fmt.Sprintf(
		`<span style="display:inline-flex;align-items:center;padding:2px 8px;border-radius:4px;font-size:11.5px;font-weight:500;%s">%s</span>`,
		st, text,
	)
}

// yesNo renders a green/red badge.
func yesNo(v bool) string {
	if v {
		return badge("Yes", "green")
	}
	return badge("No", "red")
}

// btnPrimary renders the charcoal-style primary action button.
func btnPrimary(label string) string {
	return fmt.Sprintf(`<button type="submit"
    style="background:#1e1d1e;color:#f5f4f5;border:none;border-radius:6px;padding:8px 20px;font-size:13px;font-weight:500;cursor:pointer;transition:background 0.12s;"
    onmouseover="this.style.background='#30353b'"
    onmouseout="this.style.background='#1e1d1e'">%s</button>`, label)
}

// escapeHTML escapes the five HTML-special characters.
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&#34;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}