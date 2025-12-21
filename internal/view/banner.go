package view

import (
	"bytes"
	"embed"
	"html/template"
	"log"
)

type banner struct {
	templatesMemoryEmbed embed.FS
	asciiArt             string
	programDesc          string
	version              string
	buildDate            string
	builtBy              string
	pattern              string
}

func (b *banner) Show() string {
	t := template.Must(template.ParseFS(b.templatesMemoryEmbed, b.pattern))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, map[string]string{
		"ASCII":     b.asciiArt,
		"ProgDesc":  b.programDesc,
		"ProgVer":   b.version,
		"BuildDate": b.buildDate,
		"BuiltBy":   b.builtBy,
	})
	if err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

func NewBanner(templatesMemoryEmbed embed.FS,
	pattern string,
	art string,
	programDesc string,
	version string,
	buildDate string,
	builtBy string) View {
	if pattern == "" {
		pattern = "internal/templates/banner.gotmpl"
	}
	return &banner{
		templatesMemoryEmbed: templatesMemoryEmbed,
		asciiArt:             art,
		programDesc:          programDesc,
		version:              version,
		buildDate:            buildDate,
		builtBy:              builtBy,
		pattern:              pattern,
	}
}
