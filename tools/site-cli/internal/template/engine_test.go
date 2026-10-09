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

func TestEngine_DesignTokens(t *testing.T) {
	t.Run("Load DESIGN.md and inject into template context", func(t *testing.T) {
		tempTheme := t.TempDir()

		// messages.yaml
		msgContent := "common:\n  site_title: \"Test Site\"\n"
		if err := os.WriteFile(filepath.Join(tempTheme, "messages.yaml"), []byte(msgContent), 0644); err != nil {
			t.Fatalf("failed to write messages.yaml: %v", err)
		}

		// DESIGN.md
		designContent := `---
version: "alpha"
name: "Test Dark Theme"
colors:
  primary: "#1d4ed8"
  neutral: "#020617"
rounded:
  md: "8px"
---

# Test Theme
Markdown documentation follows here.
`
		if err := os.WriteFile(filepath.Join(tempTheme, "DESIGN.md"), []byte(designContent), 0644); err != nil {
			t.Fatalf("failed to write DESIGN.md: %v", err)
		}

		// base.html with design token CSS variables
		baseHTML := `{{ define "base" }}<!DOCTYPE html>
<html>
<head>
{{ if .Design }}
<style>
:root {
  --color-primary: {{ index .Design.Colors "primary" }};
  --rounded-md: {{ index .Design.Rounded "md" }};
}
</style>
{{ end }}
</head>
<body>{{ block "content" . }}{{ end }}</body>
</html>{{ end }}`

		indexHTML := `{{ define "content" }}<h1>{{ .Design.Name }}</h1>{{ end }}`

		if err := os.WriteFile(filepath.Join(tempTheme, "base.html"), []byte(baseHTML), 0644); err != nil {
			t.Fatalf("failed to write base.html: %v", err)
		}
		if err := os.WriteFile(filepath.Join(tempTheme, "index.html"), []byte(indexHTML), 0644); err != nil {
			t.Fatalf("failed to write index.html: %v", err)
		}

		eng, err := NewEngine(tempTheme)
		if err != nil {
			t.Fatalf("failed to init engine: %v", err)
		}

		if eng.GetDesign() == nil {
			t.Fatal("expected design tokens to be loaded, got nil")
		}
		if eng.GetDesign().Name != "Test Dark Theme" {
			t.Errorf("expected theme name 'Test Dark Theme', got %q", eng.GetDesign().Name)
		}
		if eng.GetDesign().Colors["primary"] != "#1d4ed8" {
			t.Errorf("expected primary color '#1d4ed8', got %q", eng.GetDesign().Colors["primary"])
		}

		var buf bytes.Buffer
		ctx := &model.TemplateContext{}
		if err := eng.RenderPage(&buf, "index.html", ctx); err != nil {
			t.Fatalf("failed to render page: %v", err)
		}

		rendered := buf.String()
		if !strings.Contains(rendered, "--color-primary: #1d4ed8;") {
			t.Errorf("rendered output missing CSS variable --color-primary: %s", rendered)
		}
		if !strings.Contains(rendered, "--rounded-md: 8px;") {
			t.Errorf("rendered output missing CSS variable --rounded-md: %s", rendered)
		}
		if !strings.Contains(rendered, "<h1>Test Dark Theme</h1>") {
			t.Errorf("rendered output missing design name in content: %s", rendered)
		}
	})

	t.Run("Graceful fallback when DESIGN.md is absent", func(t *testing.T) {
		tempTheme := t.TempDir()

		// messages.yaml only
		msgContent := "common:\n  site_title: \"Test Site\"\n"
		if err := os.WriteFile(filepath.Join(tempTheme, "messages.yaml"), []byte(msgContent), 0644); err != nil {
			t.Fatalf("failed to write messages.yaml: %v", err)
		}

		baseHTML := `{{ define "base" }}<html><body>{{ block "content" . }}{{ end }}</body></html>{{ end }}`
		indexHTML := `{{ define "content" }}<main>No Design</main>{{ end }}`

		if err := os.WriteFile(filepath.Join(tempTheme, "base.html"), []byte(baseHTML), 0644); err != nil {
			t.Fatalf("failed to write base.html: %v", err)
		}
		if err := os.WriteFile(filepath.Join(tempTheme, "index.html"), []byte(indexHTML), 0644); err != nil {
			t.Fatalf("failed to write index.html: %v", err)
		}

		eng, err := NewEngine(tempTheme)
		if err != nil {
			t.Fatalf("failed to init engine without DESIGN.md: %v", err)
		}

		if eng.GetDesign() != nil {
			t.Errorf("expected GetDesign() to be nil, got %v", eng.GetDesign())
		}

		var buf bytes.Buffer
		ctx := &model.TemplateContext{}
		if err := eng.RenderPage(&buf, "index.html", ctx); err != nil {
			t.Fatalf("failed to render page: %v", err)
		}

		if !strings.Contains(buf.String(), "<main>No Design</main>") {
			t.Errorf("unexpected render output: %s", buf.String())
		}
	})
}

func TestEngine_Multilingual(t *testing.T) {
	tempTheme := t.TempDir()

	// 1. messages.yaml (기본/한국어)
	koMsg := `common:
  site_title: "테스트 사이트"
nav:
  about: "소개"
`
	if err := os.WriteFile(filepath.Join(tempTheme, "messages.yaml"), []byte(koMsg), 0644); err != nil {
		t.Fatal(err)
	}

	// 2. messages_en.yaml (영문)
	enMsg := `common:
  site_title: "Test Site"
nav:
  about: "About"
`
	if err := os.WriteFile(filepath.Join(tempTheme, "messages_en.yaml"), []byte(enMsg), 0644); err != nil {
		t.Fatal(err)
	}

	baseHTML := `{{ define "base" }}<!DOCTYPE html><html lang="{{ .CurrentLang }}"><head><title>{{ .Messages.Common.site_title }}</title></head><body><nav>{{ .Messages.Nav.about }}</nav>{{ block "content" . }}{{ end }}</body></html>{{ end }}`
	indexHTML := `{{ define "content" }}<main>Content</main>{{ end }}`
	if err := os.WriteFile(filepath.Join(tempTheme, "base.html"), []byte(baseHTML), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tempTheme, "index.html"), []byte(indexHTML), 0644); err != nil {
		t.Fatal(err)
	}

	eng, err := NewEngine(tempTheme)
	if err != nil {
		t.Fatal(err)
	}

	// 한국어 렌더링 검증
	var koBuf bytes.Buffer
	koCtx := &model.TemplateContext{CurrentLang: "ko"}
	if err := eng.RenderPage(&koBuf, "index.html", koCtx); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(koBuf.String(), "테스트 사이트") || !strings.Contains(koBuf.String(), "소개") {
		t.Errorf("expected Korean messages in output: %s", koBuf.String())
	}

	// 영문 렌더링 검증
	var enBuf bytes.Buffer
	enCtx := &model.TemplateContext{CurrentLang: "en"}
	if err := eng.RenderPage(&enBuf, "index.html", enCtx); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(enBuf.String(), "Test Site") || !strings.Contains(enBuf.String(), "About") {
		t.Errorf("expected English messages in output: %s", enBuf.String())
	}
}

