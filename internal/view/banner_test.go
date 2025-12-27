package view

import (
	"embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

//go:embed test_templates/banner.gotmpl
var testTemplates embed.FS

func TestBannerTemplate(t *testing.T) {
	bannerView := NewBanner(testTemplates, "test_templates/banner.gotmpl", "artTest", "test", "0.0.1", "today", "golang")
	res := bannerView.Show()

	fmt.Println(res)

	want := []string{"artTest", "today", "golang", "test"}
	for _, w := range want {
		assert.True(t, strings.Contains(res, w))
	}
}
