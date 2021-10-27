// Package dischtml contains functions to convert Discord messages to pretty HTML.
package dischtml

import (
	"bytes"
	"html/template"
	"sort"

	"github.com/diamondburned/arikawa/v3/discord"
)

// Converter is a struct holding all necessary info for converting messages.
type Converter struct {
	Guild    discord.Guild
	Channels []discord.Channel
	Roles    []discord.Role
	Users    []discord.User
}

// ConvertHTML converts messages to HTML.
func (c *Converter) ConvertHTML(msg []discord.Message) (template.HTML, error) {
	sort.Slice(msg, func(i, j int) bool {
		return msg[i].ID < msg[j].ID
	})

	data := struct {
		Messages []discord.Message
	}{Messages: msg}

	var w bytes.Buffer

	err := tmpls.ExecuteTemplate(&w, "msgs.html", data)
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
		Twemoji           template.JS
		Guild             discord.Guild
		Channel           discord.Channel
		Content           template.HTML
		MsgCount          int
	}{CSS: template.CSS(CSS), HighlightCSS: template.CSS(HighlightCSS), HighlightJS: template.JS(HighlightJS), Twemoji: template.JS(TwemojiJS), Guild: guild, Channel: channel, Content: s, MsgCount: msgCount}

	var w bytes.Buffer

	err := htmlPageTmpl.Execute(&w, data)
	if err != nil {
		return "", err
	}
	return w.String(), err
}
