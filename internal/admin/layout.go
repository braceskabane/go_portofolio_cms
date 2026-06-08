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

// ── Nav ──────────────────────────────────────────────────────────────────────

type navItem struct {
	Href    string
	Icon    string
	Label   string
	Section string // "overview" | "content" | "taxonomy" | "system"
	Badge   string // optional numeric badge
}

func buildNav(unreadContacts int) []navItem {
	unread := ""
	if unreadContacts > 0 {
		unread = fmt.Sprintf("%d", unreadContacts)
	}
	return []navItem{
		{"/admin", "◈", "Dashboard", "overview", ""},
		{"/admin/projects", "▤", "Projects", "content", ""},
		{"/admin/skills", "◎", "Skills", "content", ""},
		{"/admin/experiences", "◉", "Experiences", "content", ""},
		{"/admin/educations", "◫", "Educations", "content", ""},
		{"/admin/stack-items", "📌", "Stack Items", "content", ""},
		{"/admin/running-activities", "🏃", "Running", "content", ""},
		{"/admin/assets", "🖼️", "Assets", "content", ""},
		{"/admin/project-categories", "🏷️", "Project Cat.", "taxonomy", ""},
		{"/admin/experience-categories", "🏷️", "Exp. Cat.", "taxonomy", ""},
		{"/admin/stack-categories", "🗂️", "Stack Cat.", "taxonomy", ""},
		{"/admin/contacts", "✉️", "Contacts", "system", unread},
		{"/admin/profile", "◷", "Profile", "system", ""},
	}
}

// ── Layout ────────────────────────────────────────────────────────────────────

type layoutData struct {
	Title        string
	ActiveMenu   string
	Flash        string
	FlashIsError bool
	Nav          []navItem
	Content      template.HTML
	TopActions   template.HTML
}

type layoutCfg struct {
	flash          string
	topActions     string
	unreadContacts int
}

type layoutOpt func(*layoutCfg)

func WithFlash(msg string) layoutOpt {
	return func(c *layoutCfg) { c.flash = msg }
}

func WithTopActions(html string) layoutOpt {
	return func(c *layoutCfg) { c.topActions = html }
}

func WithUnreadContacts(n int) layoutOpt {
	return func(c *layoutCfg) { c.unreadContacts = n }
}

// layout renders the full admin shell.
func layout(title, activeMenu, content string, opts ...layoutOpt) string {
	cfg := &layoutCfg{}
	for _, o := range opts {
		o(cfg)
	}

	flash := cfg.flash
	flashIsError := false
	if strings.HasPrefix(flash, "ERR:") {
		flashIsError = true
		flash = strings.TrimPrefix(flash, "ERR:")
	}

	data := layoutData{
		Title:        title,
		ActiveMenu:   activeMenu,
		Flash:        flash,
		FlashIsError: flashIsError,
		Nav:          buildNav(cfg.unreadContacts),
		Content:      template.HTML(content), // #nosec
		TopActions:   template.HTML(cfg.topActions),
	}

	var buf bytes.Buffer
	if err := layoutTmpl.Execute(&buf, data); err != nil {
		return fmt.Sprintf(`<!DOCTYPE html><html><body><pre>template error: %s</pre></body></html>`, err.Error())
	}
	return buf.String()
}

// ── Page components ─────────────────────────────────────────────────────────

func pageHeader(title, subtitle string) string {
	sub := ""
	if subtitle != "" {
		sub = fmt.Sprintf(`<p style="font-size:13px;color:#6b6a6b;margin-top:3px;">%s</p>`, escapeHTML(subtitle))
	}
	return fmt.Sprintf(`
<div style="margin-bottom:24px;">
  <h1 style="font-size:20px;font-weight:700;color:#1e1d1e;letter-spacing:-0.02em;">%s</h1>
  %s
</div>`, escapeHTML(title), sub)
}

func sectionCard(title, content string) string {
	hdr := ""
	if title != "" {
		hdr = fmt.Sprintf(
			`<div style="padding:14px 20px;border-bottom:1px solid #f0eff0;"><h2 style="font-size:13.5px;font-weight:600;color:#1e1d1e;">%s</h2></div>`,
			escapeHTML(title),
		)
	}
	return fmt.Sprintf(`
<div style="background:#ffffff;border-radius:10px;border:1px solid #e5e3e4;margin-bottom:16px;overflow:hidden;">
  %s
  <div style="padding:20px;">%s</div>
</div>`, hdr, content)
}

func card(content string) string {
	return fmt.Sprintf(`<div style="background:#ffffff;border-radius:10px;border:1px solid #e5e3e4;padding:20px 24px;margin-bottom:16px;">%s</div>`, content)
}

func formGrid(fields ...string) string {
	return fmt.Sprintf(`<div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(240px,1fr));gap:16px;">%s</div>`, strings.Join(fields, ""))
}

func formRow(fields ...string) string {
	cols := len(fields)
	if cols == 0 {
		return ""
	}
	return fmt.Sprintf(`<div style="display:grid;grid-template-columns:repeat(%d,1fr);gap:16px;">%s</div>`, cols, strings.Join(fields, ""))
}

func divider(label string) string {
	if label == "" {
		return `<hr style="border:none;border-top:1px solid #f0eff0;margin:20px 0;" />`
	}
	return fmt.Sprintf(`
<div style="display:flex;align-items:center;gap:12px;margin:20px 0;">
  <hr style="flex:1;border:none;border-top:1px solid #f0eff0;" />
  <span style="font-size:11px;font-weight:600;text-transform:uppercase;letter-spacing:0.07em;color:#9a9899;">%s</span>
  <hr style="flex:1;border:none;border-top:1px solid #f0eff0;" />
</div>`, escapeHTML(label))
}

// ── Table ────────────────────────────────────────────────────────────────────

func tableWrapper(content string) string {
	return fmt.Sprintf(`<div style="background:#ffffff;border-radius:10px;border:1px solid #e5e3e4;overflow:hidden;"><div style="overflow-x:auto;">%s</div></div>`, content)
}

func tableHead(cols ...string) string {
	cells := ""
	for _, c := range cols {
		cells += fmt.Sprintf(`<th style="padding:10px 16px;text-align:left;font-size:11.5px;font-weight:600;text-transform:uppercase;letter-spacing:0.06em;color:#6b6a6b;white-space:nowrap;">%s</th>`, c)
	}
	return fmt.Sprintf(`<thead style="background:#f8f7f8;border-bottom:1px solid #e5e3e4;"><tr>%s</tr></thead>`, cells)
}

func tableRow(isEven bool, cells ...string) string {
	bg := "#ffffff"
	if isEven {
		bg = "#faf9fa"
	}
	content := ""
	for _, c := range cells {
		content += fmt.Sprintf(`<td style="padding:11px 16px;font-size:13px;color:#1e1d1e;border-bottom:1px solid #f0eff0;vertical-align:middle;">%s</td>`, c)
	}
	return fmt.Sprintf(`<tr style="background:%s;transition:background 0.1s;" onmouseover="this.style.background='#f5f4f5'" onmouseout="this.style.background='%s'">%s</tr>`, bg, bg, content)
}

// ── Buttons ──────────────────────────────────────────────────────────────────

func btnPrimary(label string) string {
	return fmt.Sprintf(`<button type="submit" style="display:inline-flex;align-items:center;gap:6px;background:#1e1d1e;color:#f5f4f5;border:none;border-radius:6px;padding:8px 18px;font-size:13px;font-weight:500;cursor:pointer;transition:background 0.12s;" onmouseover="this.style.background='#30353b'" onmouseout="this.style.background='#1e1d1e'">%s</button>`, label)
}

func btnSecondary(label, href string) string {
	return fmt.Sprintf(`<a href="%s" style="display:inline-flex;align-items:center;gap:6px;background:transparent;color:#4d4c4d;border:1px solid #e5e3e4;border-radius:6px;padding:7px 17px;font-size:13px;font-weight:500;cursor:pointer;transition:background 0.12s,color 0.12s;text-decoration:none;" onmouseover="this.style.background='#f5f4f5';this.style.color='#1e1d1e'" onmouseout="this.style.background='transparent';this.style.color='#4d4c4d'">%s</a>`, href, label)
}

func btnDanger(label, href string) string {
	return fmt.Sprintf(`<a href="%s" style="display:inline-flex;align-items:center;gap:4px;font-size:12px;color:#b91c1c;font-weight:500;text-decoration:none;transition:color 0.1s;" onmouseover="this.style.color='#7f1d1d'" onmouseout="this.style.color='#b91c1c'">%s</a>`, href, label)
}

func topbarNewBtn(label, href string) string {
	return fmt.Sprintf(`<a href="%s" style="display:inline-flex;align-items:center;gap:6px;background:#1e1d1e;color:#f5f4f5;border-radius:6px;padding:7px 14px;font-size:12.5px;font-weight:500;text-decoration:none;transition:background 0.12s;" onmouseover="this.style.background='#30353b'" onmouseout="this.style.background='#1e1d1e'">＋ %s</a>`, href, label)
}

func deleteBtn(url, confirmMsg string) string {
	return fmt.Sprintf(`<button hx-delete="%s" hx-confirm="%s" hx-target="closest tr" hx-swap="outerHTML swap:0.25s" style="display:inline-flex;align-items:center;font-size:12px;color:#b91c1c;font-weight:500;background:none;border:none;cursor:pointer;padding:0;transition:color 0.1s;" onmouseover="this.style.color='#7f1d1d'" onmouseout="this.style.color='#b91c1c'">Delete</button>`, url, confirmMsg)
}

// ── Form Fields ──────────────────────────────────────────────────────────────

func inputField(label, name, value, inputType string, required bool) string {
	req := ""
	if required {
		req = "required"
	}
	if inputType == "" {
		inputType = "text"
	}
	return fmt.Sprintf(`
<div>
  <label style="display:block;font-size:12px;font-weight:600;color:#4d4c4d;margin-bottom:5px;letter-spacing:0.01em;">%s</label>
  <input type="%s" name="%s" value="%s" %s
    style="width:100%%;border:1px solid #d1cfd0;border-radius:6px;padding:8px 12px;font-size:13px;color:#1e1d1e;background:#ffffff;box-sizing:border-box;outline:none;transition:border-color 0.12s,box-shadow 0.12s;"
    onfocus="this.style.borderColor='#1e1d1e';this.style.boxShadow='0 0 0 3px rgba(30,29,30,0.08)'"
    onblur="this.style.borderColor='#d1cfd0';this.style.boxShadow='none'" />
</div>`, label, inputType, name, value, req)
}

func textareaField(label, name, value string, rows int) string {
	if rows == 0 {
		rows = 4
	}
	return fmt.Sprintf(`
<div>
  <label style="display:block;font-size:12px;font-weight:600;color:#4d4c4d;margin-bottom:5px;letter-spacing:0.01em;">%s</label>
  <textarea name="%s" rows="%d"
    style="width:100%%;border:1px solid #d1cfd0;border-radius:6px;padding:8px 12px;font-size:13px;color:#1e1d1e;background:#ffffff;box-sizing:border-box;resize:vertical;font-family:inherit;outline:none;transition:border-color 0.12s,box-shadow 0.12s;"
    onfocus="this.style.borderColor='#1e1d1e';this.style.boxShadow='0 0 0 3px rgba(30,29,30,0.08)'"
    onblur="this.style.borderColor='#d1cfd0';this.style.boxShadow='none'">%s</textarea>
</div>`, label, name, rows, value)
}

func toggleField(label, name string, value bool) string {
	trueSelected, falseSelected := "", ""
	if value {
		trueSelected = "selected"
	} else {
		falseSelected = "selected"
	}
	return fmt.Sprintf(`
<div>
  <label style="display:block;font-size:12px;font-weight:600;color:#4d4c4d;margin-bottom:5px;letter-spacing:0.01em;">%s</label>
  <select name="%s"
    style="width:100%%;border:1px solid #d1cfd0;border-radius:6px;padding:8px 12px;font-size:13px;color:#1e1d1e;background:#ffffff;box-sizing:border-box;cursor:pointer;outline:none;transition:border-color 0.12s;"
    onfocus="this.style.borderColor='#1e1d1e'"
    onblur="this.style.borderColor='#d1cfd0'">
    <option value="true" %s>Yes</option>
    <option value="false" %s>No</option>
  </select>
</div>`, label, name, trueSelected, falseSelected)
}

func selectField(label, name, selected string, options map[string]string, required bool) string {
	req := ""
	if required {
		req = "required"
	}
	opts := ""
	for v, lbl := range options {
		sel := ""
		if v == selected {
			sel = "selected"
		}
		opts += fmt.Sprintf(`<option value="%s" %s>%s</option>`, v, sel, lbl)
	}
	return fmt.Sprintf(`
<div>
  <label style="display:block;font-size:12px;font-weight:600;color:#4d4c4d;margin-bottom:5px;letter-spacing:0.01em;">%s</label>
  <select name="%s" %s
    style="width:100%%;border:1px solid #d1cfd0;border-radius:6px;padding:8px 12px;font-size:13px;color:#1e1d1e;background:#ffffff;box-sizing:border-box;cursor:pointer;outline:none;transition:border-color 0.12s;"
    onfocus="this.style.borderColor='#1e1d1e'"
    onblur="this.style.borderColor='#d1cfd0'">%s</select>
</div>`, label, name, req, opts)
}

func formActions(submitLabel, cancelHref string) string {
	return fmt.Sprintf(`
<div style="display:flex;align-items:center;gap:10px;margin-top:8px;">
  %s
  %s
</div>`, btnPrimary(submitLabel), btnSecondary("Cancel", cancelHref))
}

// ── Checkbox Group (untuk pilihan many-to-many) ─────────────────────────────

func checkboxGroup(label, name string, items []struct{ ID, Text string; Checked bool }) string {
	checkboxes := ""
	for _, item := range items {
		checked := ""
		if item.Checked {
			checked = "checked"
		}
		checkboxes += fmt.Sprintf(`
<label style="display:inline-flex;align-items:center;gap:6px;margin-right:12px;margin-bottom:6px;font-size:13px;">
  <input type="checkbox" name="%s" value="%s" %s> %s
</label>`, name, item.ID, checked, item.Text)
	}
	return fmt.Sprintf(`
<div>
  <label style="display:block;font-size:12px;font-weight:600;color:#4d4c4d;margin-bottom:8px;letter-spacing:0.01em;">%s</label>
  <div>%s</div>
</div>`, label, checkboxes)
}

// ── Badges & Pills ───────────────────────────────────────────────────────────

var badgeStyles = map[string]string{
	"green":  "background:#f0fdf4;color:#166534;border:1px solid #bbf7d0;",
	"red":    "background:#fef2f2;color:#991b1b;border:1px solid #fecaca;",
	"gray":   "background:#f5f4f5;color:#4d4c4d;border:1px solid #e5e3e4;",
	"blue":   "background:#eff6ff;color:#1e40af;border:1px solid #bfdbfe;",
	"yellow": "background:#fefce8;color:#854d0e;border:1px solid #fef08a;",
	"purple": "background:#faf5ff;color:#6b21a8;border:1px solid #e9d5ff;",
}

func badge(text, color string) string {
	st, ok := badgeStyles[color]
	if !ok {
		st = badgeStyles["gray"]
	}
	return fmt.Sprintf(`<span style="display:inline-flex;align-items:center;padding:2px 8px;border-radius:4px;font-size:11.5px;font-weight:500;%s">%s</span>`, st, text)
}

func yesNo(v bool) string {
	if v {
		return badge("Yes", "green")
	}
	return badge("No", "red")
}

// ─── Stat Card (khusus dashboard) ───────────────────────────────────────────

func statCard(icon, label, value, accentBg string) string {
	return fmt.Sprintf(`
<div style="background:#ffffff; border-radius:12px; border:1px solid #e5e3e4; padding:20px; display:flex; align-items:center; gap:16px;">
  <div style="width:48px; height:48px; border-radius:12px; background:%s; display:flex; align-items:center; justify-content:center; font-size:24px; flex-shrink:0;">%s</div>
  <div>
    <p style="font-size:13px; color:#6b6a6b;">%s</p>
    <p style="font-size:28px; font-weight:600; color:#1e1d1e;">%s</p>
  </div>
</div>`, accentBg, icon, label, value)
}

// ── Empty State ──────────────────────────────────────────────────────────────

func emptyState(icon, title, description, ctaLabel, ctaHref string) string {
	cta := ""
	if ctaLabel != "" && ctaHref != "" {
		cta = fmt.Sprintf(`<div style="margin-top:16px;">%s</div>`, topbarNewBtn(ctaLabel, ctaHref))
	}
	return fmt.Sprintf(`
<div style="display:flex;flex-direction:column;align-items:center;justify-content:center;padding:64px 24px;text-align:center;">
  <div style="font-size:36px;margin-bottom:12px;opacity:0.5;">%s</div>
  <p style="font-size:15px;font-weight:600;color:#1e1d1e;margin-bottom:4px;">%s</p>
  <p style="font-size:13px;color:#6b6a6b;max-width:280px;">%s</p>
  %s
</div>`, icon, escapeHTML(title), escapeHTML(description), cta)
}

// ── Pagination ───────────────────────────────────────────────────────────────

func pagination(currentPage, totalPages int, baseURL string) string {
	if totalPages <= 1 {
		return ""
	}
	sep := "?"
	if strings.Contains(baseURL, "?") {
		sep = "&"
	}
	prev, next := "", ""
	if currentPage > 1 {
		prev = fmt.Sprintf(`<a href="%s%spage=%d" style="display:inline-flex;align-items:center;gap:4px;padding:6px 12px;border-radius:6px;border:1px solid #e5e3e4;font-size:12.5px;font-weight:500;color:#4d4c4d;text-decoration:none;transition:background 0.1s;" onmouseover="this.style.background='#f5f4f5'" onmouseout="this.style.background=''">← Prev</a>`, baseURL, sep, currentPage-1)
	}
	if currentPage < totalPages {
		next = fmt.Sprintf(`<a href="%s%spage=%d" style="display:inline-flex;align-items:center;gap:4px;padding:6px 12px;border-radius:6px;border:1px solid #e5e3e4;font-size:12.5px;font-weight:500;color:#4d4c4d;text-decoration:none;transition:background 0.1s;" onmouseover="this.style.background='#f5f4f5'" onmouseout="this.style.background=''">Next →</a>`, baseURL, sep, currentPage+1)
	}
	info := fmt.Sprintf(`<span style="font-size:12.5px;color:#6b6a6b;">Page %d of %d</span>`, currentPage, totalPages)
	return fmt.Sprintf(`<div style="display:flex;align-items:center;justify-content:space-between;padding:14px 0;margin-top:4px;"><div>%s</div>%s<div>%s</div></div>`, prev, info, next)
}

// ── Search Bar ───────────────────────────────────────────────────────────────

func searchBar(placeholder, name, value, action string) string {
	return fmt.Sprintf(`
<form method="GET" action="%s" style="display:flex;gap:8px;margin-bottom:16px;">
  <input type="text" name="%s" value="%s" placeholder="%s"
    style="flex:1;max-width:320px;border:1px solid #e5e3e4;border-radius:6px;padding:7px 12px;font-size:13px;color:#1e1d1e;background:#ffffff;outline:none;transition:border-color 0.12s;"
    onfocus="this.style.borderColor='#1e1d1e'"
    onblur="this.style.borderColor='#e5e3e4'" />
  <button type="submit"
    style="padding:7px 14px;border-radius:6px;background:#1e1d1e;color:#f5f4f5;border:none;font-size:13px;font-weight:500;cursor:pointer;transition:background 0.12s;"
    onmouseover="this.style.background='#30353b'"
    onmouseout="this.style.background='#1e1d1e'">Search</button>
</form>`, action, name, value, placeholder)
}

// ── Utilities ────────────────────────────────────────────────────────────────

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&#34;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "…"
}