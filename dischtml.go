// Package dischtml contains functions to convert Discord messages to pretty HTML.
package dischtml

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/microcosm-cc/bluemonday"
)

// Converter is a struct holding all necessary info for converting messages.
// It is not safe for concurrent use.
type Converter struct {
	Guild    discord.Guild
	Channels []discord.Channel
	Roles    []discord.Role
	Users    []discord.User
	Members  []discord.Member

	tmpl *template.Template
}

// ConvertHTML converts messages to HTML.
func (c *Converter) ConvertHTML(msg []discord.Message) (template.HTML, error) {
	err := c.setTemplates()
	if err != nil {
		return "", err
	}

	sort.Slice(msg, func(i, j int) bool {
		return msg[i].ID < msg[j].ID
	})

	data := struct {
		Messages []discord.Message
	}{Messages: msg}

	var w bytes.Buffer

	err = c.tmpl.ExecuteTemplate(&w, "msgs.html", data)
	if err != nil {
		return "", err
	}
	return template.HTML(w.String()), err
}

// Wrap wraps the given string into a complete HTML page.
// For more control over this, use the exported HTMLPage and CSS variables.
func Wrap(guild discord.Guild, channel discord.Channel, s template.HTML, msgCount int) (string, error) {
	data := struct {
		HighlightCSS, CSS template.CSS
		HighlightJS       template.JS
		Guild             discord.Guild
		Channel           discord.Channel
		Content           template.HTML
		MsgCount          int
		Now               string
	}{CSS: template.CSS(CSS), HighlightCSS: template.CSS(HighlightCSS), HighlightJS: template.JS(HighlightJS), Guild: guild, Channel: channel, Content: s, MsgCount: msgCount, Now: time.Now().UTC().Format("2006-01-02 15:04")}

	var w bytes.Buffer

	err := htmlPageTmpl.Execute(&w, data)
	if err != nil {
		return "", err
	}
	return w.String(), err
}

func (c *Converter) setTemplates() (err error) {
	if c.tmpl != nil {
		return nil
	}

	c.tmpl, err = template.New("").Funcs(c.funcs()).ParseFS(fs, "embed/tmpls/*.html")
	return err
}

func (c *Converter) funcs() template.FuncMap {
	funcs := make(template.FuncMap, len(funcMap))
	for k, v := range funcMap {
		funcs[k] = v
	}

	funcs["parseMentions"] = c.parseMentions
	funcs["userColour"] = func(u discord.User) template.HTML {
		var clr discord.Color
		for _, m := range c.Members {
			if m.User.ID == u.ID {
				clr = userColour(c.Roles, m.RoleIDs)
				break
			}
		}

		if clr != 0 {
			s := fmt.Sprintf(`<span class="font-bold" style="color: #%06X;">%%s</span>`, clr)
			return template.HTML(fmt.Sprintf(s, bluemonday.UGCPolicy().Sanitize(u.Username)))
		}

		return template.HTML(
			fmt.Sprintf(`<span class="font-bold">%s</span>`,
				bluemonday.UGCPolicy().Sanitize(u.Username)),
		)
	}

	return funcs
}

func userColour(roles []discord.Role, roleIDs []discord.RoleID) discord.Color {
	sort.Slice(roles, func(i, j int) bool {
		return roles[i].Position < roles[j].Position
	})

	for _, r := range roles {
		for _, id := range roleIDs {
			if r.ID == id && r.Color != 0 {
				return r.Color
			}
		}
	}
	return 0
}
