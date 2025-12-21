package config

import (
	"flag"
	"github.com/pm1381/sirish/internal/dto"
	"os"
)

type Config struct {
	FormatImports  *bool
	TraceGenerator *bool
	Types          *dto.Types
	FilePath       *string // Relative Path
	GoPackage      string
	ShowBanner     *bool
	flagSet        *flag.FlagSet
}

func NewConfig(name string, env string) *Config {
	fg := flag.NewFlagSet(name, flag.ExitOnError)
	if env != "production" {
		fg = flag.NewFlagSet(name, flag.ContinueOnError)
	}
	cfg := Config{
		FormatImports:  new(bool),
		Types:          new(dto.Types),
		FilePath:       new(string),
		ShowBanner:     new(bool),
		TraceGenerator: new(bool),
		flagSet:        fg,
	}
	cfg.GoPackage = os.Getenv("GOPACKAGE")
	cfg.flagSet.BoolVar(cfg.ShowBanner, "banner", true, "Show program version")
	cfg.flagSet.BoolVar(cfg.FormatImports, "fmt", true, "format imports")
	cfg.flagSet.BoolVar(cfg.TraceGenerator, "tg", true, "if set to true it will init tracing where the context is not passed")
	cfg.flagSet.StringVar(cfg.FilePath, "f", "", "File path to parse. Can be overwritten with GOFILE")
	cfg.flagSet.Var(cfg.Types, "t", "list of interfaces which apm will be applied to."+
		"it can be either comma separated or repeated for example interface1, interface2 or -t interface1 -t interface2")

	return &cfg
}

func (c *Config) Parse(arguments []string) error {
	return c.flagSet.Parse(arguments)
}
