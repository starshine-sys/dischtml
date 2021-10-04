package dischtml

import (
	"embed"
	"fmt"
	"html/template"

	"github.com/diamondburned/arikawa/v3/discord"
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

var htmlPageTmpl = template.Must(template.New("wrap").Parse(HTMLPage))

var tmpls = template.Must(template.New("").Funcs(funcMap).ParseFS(fs, "embed/tmpls/*.html"))

func init() {
	dir, _ := fs.ReadDir(".")

	for _, d := range dir {
		fmt.Println(d.Name(), d.Type())
	}
}

var funcMap = template.FuncMap{
	"msgMarkdown": func(s string) template.HTML {
		b := blackfriday.Run([]byte(s), blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.Autolink|blackfriday.Strikethrough|blackfriday.HardLineBreak))
		return template.HTML(b)
	},
	"msgIDTimestamp": func(id discord.MessageID) template.HTML {
		return template.HTML(id.Time().Format("2006-01-02 15:04"))
	},
}
