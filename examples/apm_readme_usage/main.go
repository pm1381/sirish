package main

import (
	"github.com/pm1381/sirish/examples/apm_readme_usage/internal"
)

func main() {
	e := internal.StartEcho()
	h := internal.NewTestModuleSirishWrapperImpl("testHandler", internal.NewTestHandler(), "custom")
	internal.SetUpRoutes(e, h)
	e.Logger.Fatal(e.Start("localhost:8080"))
}
