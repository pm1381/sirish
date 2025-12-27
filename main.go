package main

import (
	"embed"
	"fmt"
	"github.com/pm1381/sirish/internal"
	"github.com/pm1381/sirish/internal/config"
	"github.com/pm1381/sirish/internal/view"
	"github.com/pm1381/sirish/internal/visitors"
	"github.com/pm1381/sirish/internal/wrapper"
	"log"
	"os"
	"path/filepath"
)

const (
	builtBy  = "golang"
	date     = "2025-12-20"
	progDesc = "sirish, a solution to create wrappers for your tracing needs"
)

//go:embed ascii.txt
var asciiArt string

//go:embed internal/templates/*.gotmpl
var templatesMemoryEmbed embed.FS

func main() {
	cfg := config.NewConfig(os.Args[0], "production")

	err := cfg.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	if *cfg.ShowBanner {
		banner := view.NewBanner(templatesMemoryEmbed, "internal/templates/banner.gotmpl", asciiArt, progDesc, "", date, builtBy)
		fmt.Print(banner.Show())
	}
	if *cfg.FilePath == "" {
		*cfg.FilePath = os.Getenv("GOFILE")
	}
	cwd, err := os.Getwd() // working directory
	if err != nil {
		log.Fatal(err)
	}
	abs := filepath.Join(cwd, *cfg.FilePath)
	*cfg.FilePath = abs
	fmt.Printf("sirish configs: filePath: %s goPackage: %s \n", *cfg.FilePath, cfg.GoPackage)

	// parse comments for //sirish:InterfaceName
	commentVisitor := visitors.NewCommentVisitor(*cfg.FilePath)
	commentVisitor.Traverse()
	// parse interfaces inside the file
	typeVisitor := visitors.NewTypeVisitor(*cfg.FilePath, internal.GenerateUniqueValues(*cfg.Types, commentVisitor.GetTargets()))
	err = typeVisitor.Traverse()
	if err != nil {
		log.Fatal(err)
	}

	generator := wrapper.NewApmWrapper("sirish", "internal/templates/wrapper.gotmpl",
		templatesMemoryEmbed, typeVisitor.GetWrappedInterfaces(), typeVisitor.GetImports())

	fmt.Println("sirish starts the firework...")

	err = generator.Generate(wrapper.APMTypeWrapperOptions{
		GeneralOptions: wrapper.GeneralOptions{
			Version:  version,
			Imports:  *cfg.FormatImports,
			CreateTx: *cfg.TraceGenerator,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
