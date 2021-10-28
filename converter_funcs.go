package dischtml

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/microcosm-cc/bluemonday"
)

var (
	channelRe = regexp.MustCompile(`<#(\d{15,})>`)
	userRe    = regexp.MustCompile(`<@!?(\d{15,})>`)
	roleRe    = regexp.MustCompile(`<@&(\d{15,})>`)
)

func (c *Converter) parseMentions(content string) string {
	matches := channelRe.FindAllStringSubmatch(content, -1)
	for _, m := range matches {
		if len(m) != 2 {
			continue
		}

		var ch *discord.Channel

		for _, channel := range c.Channels {
			if channel.ID.String() == m[1] {
				ch = &channel
				break
			}
		}

		if ch == nil {
			content = strings.ReplaceAll(content, m[0],
				fmt.Sprintf(`<span class="rounded p-1 bg-mentionBlurple text-mentionText">&lt;#%s&gt;</span>`, m[1]))
			continue
		}

		content = strings.ReplaceAll(content, m[0],
			fmt.Sprintf(`<span class="rounded p-1 bg-mentionBlurple text-mentionText" title="ID: %d">#%s</span>`, ch.ID, bluemonday.UGCPolicy().Sanitize(ch.Name)))
	}

	matches = userRe.FindAllStringSubmatch(content, -1)
	for _, m := range matches {
		if len(m) != 2 {
			continue
		}

		var u *discord.User

		for _, user := range c.Users {
			if user.ID.String() == m[1] {
				u = &user
				break
			}
		}

		if u == nil {
			content = strings.ReplaceAll(content, m[0],
				fmt.Sprintf(`<span class="rounded p-1 bg-mentionBlurple text-mentionText">&lt;@!%s&gt;</span>`, m[1]))
			continue
		}

		content = strings.ReplaceAll(content, m[0],
			fmt.Sprintf(`<span class="rounded p-1 bg-mentionBlurple text-mentionText" title="ID: %d">@%s</span>`, u.ID, bluemonday.UGCPolicy().Sanitize(u.Username+"#"+u.Discriminator)))
	}

	matches = roleRe.FindAllStringSubmatch(content, -1)
	for _, m := range matches {
		if len(m) != 2 {
			continue
		}

		var r *discord.Role

		for _, role := range c.Roles {
			if role.ID.String() == m[1] {
				r = &role
				break
			}
		}

		if r == nil {
			content = strings.ReplaceAll(content, m[0],
				fmt.Sprintf(`<span class="rounded p-1 bg-mentionBlurple text-mentionText">&lt;@&%s&gt;</span>`, m[1]))
			continue
		}

		content = strings.ReplaceAll(content, m[0],
			fmt.Sprintf(`<span class="rounded p-1 bg-mentionBlurple text-mentionText" title="ID: %d">@%s</span>`, r.ID, bluemonday.UGCPolicy().Sanitize(r.Name)))
	}

	return content
}

func (c *Converter) extraInfo(id discord.MessageID) template.HTML {
	s, ok := c.ExtraUserInfo[id]
	if !ok {
		return ""
	}

	return template.HTML(fmt.Sprintf(`<span class="text-lighterGray">%s</span> `, bluemonday.UGCPolicy().Sanitize(s)))
}
