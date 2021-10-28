package dischtml

import (
	"embed"
	"fmt"
	"html/template"
	"regexp"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/dustin/go-humanize"
	"github.com/russross/blackfriday/v2"
)

// HTMLPage contains the HTML that converted messages should be wrapped around for stand-alone viewing.
//go:embed embed/wrap.html
var HTMLPage string

// CSS contains the CSS used to make the pages actually look good.
//go:embed embed/discord.css
var CSS string

//go:embed embed/tmpls/*
var fs embed.FS

// HighlightCSS contains the CSS used to make highlight.js look good.
//go:embed embed/highlightjs/theme.css
var HighlightCSS string

// HighlightJS is the highlight.js code.
//go:embed embed/highlightjs/highlight.min.js
var HighlightJS string

var htmlPageTmpl = template.Must(template.New("wrap").Parse(HTMLPage))

var funcMap = template.FuncMap{
	"msgMarkdown": func(s string) template.HTML {
		b := blackfriday.Run([]byte(s), blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.Autolink|blackfriday.Strikethrough|blackfriday.HardLineBreak))
		return template.HTML(b)
	},
	"msgIDTimestamp": func(id discord.MessageID) template.HTML {
		return template.HTML(id.Time().UTC().Format("2006-01-02 15:04"))
	},
	"msgIDTime": func(id discord.MessageID) template.HTML {
		return template.HTML(id.Time().UTC().Format("15:04"))
	},
	"colour": func(c discord.Color) template.CSS {
		return template.CSS(fmt.Sprintf("%06X", c))
	},
	"emojiToImgs": func(s string) string {
		big := ""
		if onlyEmoji.MatchString(s) {
			big = " big-emoji"
		}

		emojis := emojiMatch.FindAllString(s, -1)
		if emojis == nil {
			return s
		}

		for _, e := range emojis {
			ext := ".png"
			groups := emojiMatch.FindStringSubmatch(e)
			if groups[1] == "a" {
				ext = ".gif"
			}
			name := groups[2]
			url := emojiBaseURL + groups[3] + ext

			s = strings.NewReplacer(e, fmt.Sprintf(`<img class="emoji%s" src="%v" alt="%v" />`, big, url, name)).Replace(s)
		}

		return s
	},
	"isBot": func(m discord.Message) bool {
		return m.Author.Bot && !m.WebhookID.IsValid() && m.Author.Discriminator != "0000"
	},
	"isImage": func(a discord.Attachment) bool {
		switch {
		case strings.HasSuffix(a.Filename, ".jpg"), strings.HasSuffix(a.Filename, ".jpeg"):
			return true
		case strings.HasSuffix(a.Filename, ".gif"):
			return true
		case strings.HasSuffix(a.Filename, ".png"):
			return true
		case strings.HasSuffix(a.Filename, ".webp"):
			return true
		}
		return false
	},
	"byteSize": humanize.Bytes,
	"firstID": func(msgs []discord.Message) discord.MessageID {
		return msgs[0].ID
	},
	"largeMessageGap": func(prev, current discord.MessageID) bool {
		return current.Time().Sub(prev.Time()) > 5*time.Minute
	},
}

const emojiBaseURL = "https://cdn.discordapp.com/emojis/"

var (
	emojiMatch = regexp.MustCompile("<(?P<animated>a)?:(?P<name>\\w+):(?P<emoteID>\\d{15,})>")
	onlyEmoji  = regexp.MustCompile(`^(<(a)?:(\w+):(\d{15,})>|\s)*$`)
)
