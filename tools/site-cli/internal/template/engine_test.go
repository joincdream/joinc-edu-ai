package template

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/model"
)

func TestEngine(t *testing.T) {
	tempTheme := t.TempDir()

	// 1. messages.yaml 생성
	msgContent := `common:
  site_title: "Test Site"
nav:
  home: "Home Label"
`
	if err := os.WriteFile(filepath.Join(tempTheme, "messages.yaml"), []byte(msgContent), 0644); err != nil {
		t.Fatalf("failed to write messages.yaml: %v", err)
	}

	// 2. base.html 및 index.html 생성
	baseHTML := `{{ define "base" }}<!DOCTYPE html>
<html>
<head><title>{{ .SiteTitle }}</title></head>
<body>
  <nav>{{ .Messages.Nav.home }}</nav>
  {{ block "content" . }}{{ end }}
</body>
</html>{{ end }}`

	indexHTML := `{{ define "content" }}
  <main>Total Posts: {{ len .Posts }}</main>
{{ end }}`

	if err := os.WriteFile(filepath.Join(tempTheme, "base.html"), []byte(baseHTML), 0644); err != nil {
		t.Fatalf("failed to write base.html: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempTheme, "index.html"), []byte(indexHTML), 0644); err != nil {
		t.Fatalf("failed to write index.html: %v", err)
	}

	t.Run("Initialize engine and render page with messages", func(t *testing.T) {
		eng, err := NewEngine(tempTheme)
		if err != nil {
			t.Fatalf("failed to init engine: %v", err)
		}

		ctx := &model.TemplateContext{
			SiteTitle: "My Awesome Blog",
			Posts: []*model.Post{
				{Slug: "post-1"},
				{Slug: "post-2"},
			},
		}

		var buf bytes.Buffer
		if err := eng.RenderPage(&buf, "index.html", ctx); err != nil {
			t.Fatalf("failed to render page: %v", err)
		}

		rendered := buf.String()
		if !strings.Contains(rendered, "<title>My Awesome Blog</title>") {
			t.Errorf("title not rendered correctly: %s", rendered)
		}
		if !strings.Contains(rendered, "<nav>Home Label</nav>") {
			t.Errorf("messages.nav.home not injected: %s", rendered)
		}
		if !strings.Contains(rendered, "Total Posts: 2") {
			t.Errorf("posts length not rendered: %s", rendered)
		}
	})

	t.Run("Fail when theme directory not found", func(t *testing.T) {
		_, err := NewEngine("non-existent-theme-path")
		if err == nil {
			t.Fatal("expected error for non-existent theme dir, got nil")
		}
	})

	t.Run("Fail when messages.yaml is missing", func(t *testing.T) {
		emptyDir := t.TempDir()
		_, err := NewEngine(emptyDir)
		if err == nil {
			t.Fatal("expected error for missing messages.yaml, got nil")
		}
	})
}
