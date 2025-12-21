package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewConfig_ReadsGOPACKAGE(t *testing.T) {
	t.Setenv("GOPACKAGE", "mypkg")
	cfg := NewConfig("sirish", "dev")
	assert.Equal(t, "mypkg", cfg.GoPackage)
}

func TestParse_Defaults(t *testing.T) {
	cfg := NewConfig("sirish", "dev")

	err := cfg.Parse([]string{}) // show banner default and format input is true
	require.NoError(t, err)

	assert.Equal(t, *cfg.ShowBanner, true)
	assert.Equal(t, *cfg.FormatImports, true)
	assert.Equal(t, len(*cfg.Types), 0)
	assert.Equal(t, *cfg.FilePath, "")
}

func TestParse_FilePathAndFlags(t *testing.T) {
	cfg := NewConfig("sirish", "dev")

	path := "visitors/test_samples/comment.go"
	err := cfg.Parse([]string{"-f", path, "-fmt=true", "-banner=false", "-tg=false"})
	require.NoError(t, err)

	assert.Equal(t, *cfg.ShowBanner, false)
	assert.Equal(t, *cfg.TraceGenerator, false)
	assert.Equal(t, *cfg.FormatImports, true)
	assert.Equal(t, len(*cfg.Types), 0)
	assert.Equal(t, *cfg.FilePath, path)
}

func TestParse_Types_CommaSeparated(t *testing.T) {
	cfg := NewConfig("sirish", "dev")

	// -banner registered as boolVar so -banner, false does not work for it
	err := cfg.Parse([]string{"-t", "Repo,ProfileStore", "-banner=false"})
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	assert.Equal(t, *cfg.ShowBanner, false)

	got := []string(*cfg.Types)
	want := []string{"Repo", "ProfileStore"}
	assert.ElementsMatch(t, got, want)
}

func TestParse_Types_Repeated(t *testing.T) {
	cfg := NewConfig("sirish", "dev")

	err := cfg.Parse([]string{"-t", "Repo", "-t", "ProfileStore"})
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	got := []string(*cfg.Types)
	want := []string{"Repo", "ProfileStore"}
	assert.ElementsMatch(t, got, want)
}

func TestParse_UnknownFlag_ReturnsError(t *testing.T) {
	cfg := NewConfig("sirish", "dev")

	err := cfg.Parse([]string{"-doesnotexist"}) // must have error
	assert.NotNil(t, err)
}
