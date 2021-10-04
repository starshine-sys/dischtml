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
// The messages should be in reverse order, i.e. newest first.
// They will be reversed in this function.
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
func Wrap(guild discord.Guild, channel discord.Channel, s template.HTML) (string, error) {
	data := struct {
		CSS     template.CSS
		Guild   discord.Guild
		Channel discord.Channel
		Content template.HTML
	}{CSS: template.CSS(CSS), Guild: guild, Channel: channel, Content: s}

	var w bytes.Buffer

	err := htmlPageTmpl.Execute(&w, data)
	if err != nil {
		return "", err
	}
	return w.String(), err
}
