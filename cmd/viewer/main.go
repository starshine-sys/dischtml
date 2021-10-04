package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/dischtml"
)

// File ...
type File struct {
	Messages []discord.Message
	Channel  discord.Channel
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No file specified!")
		return
	}
	name := os.Args[1]

	if !filepath.IsAbs(name) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting working directory:", err)
			return
		}
		name = filepath.Join(cwd, name)
	}

	b, err := os.ReadFile(name)
	if err != nil {
		fmt.Printf("Couldn't read file %v: %v\n", name, err)
		return
	}

	var f File
	err = json.Unmarshal(b, &f)
	if err != nil {
		fmt.Println("Couldn't unmarshal file:", err)
		return
	}

	c := dischtml.Converter{
		Guild: discord.Guild{
			Name: "The Cove",
			ID:   739823533039026239,
			Icon: "5067667d382e53dabc4aba9095c45b49",
		},
	}

	html, err := c.ConvertHTML(f.Messages)
	if err != nil {
		fmt.Println("Couldn't convert file:", err)
		return
	}

	full, err := dischtml.Wrap(c.Guild, f.Channel, html)
	if err != nil {
		fmt.Println("Couldn't write full file:", err)
		return
	}

	err = os.WriteFile(name+".out.html", []byte(full), 0777)
	if err != nil {
		fmt.Println("Couldn't write file:", err)
	}
}
