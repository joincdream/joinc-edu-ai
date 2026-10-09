package builder

import (
	"encoding/json"
	"html/template"
	"os"
	"path/filepath"
	"testing"

	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/model"
)

func TestBuildDataBundleAndExport(t *testing.T) {
	tempDir := t.TempDir()
	jsonPath := filepath.Join(tempDir, "site-data.json")

	b := NewBuilder(Options{
		BaseURL: "/",
	})

	koPosts := []*model.Post{
		{
			ID:   "1",
			Slug: "test-post",
			Frontmatter: model.Frontmatter{
				Title:       "테스트 포스트",
				Category:    "AI",
				CreatedDate: "2026-10-09",
				Status:      "published",
			},
			HTMLContent: template.HTML("<p>안녕하세요</p>"),
			TOC: []model.TOCItem{
				{ID: "sec-1", Title: "섹션 1", Level: 2},
			},
		},
	}

	enPosts := []*model.Post{
		{
			ID:   "1",
			Slug: "test-post",
			Frontmatter: model.Frontmatter{
				Title:       "Test Post",
				Category:    "AI",
				CreatedDate: "2026-10-09",
				Status:      "published",
			},
			HTMLContent: template.HTML("<p>Hello</p>"),
		},
	}

	koTax := model.NewTaxonomyIndex(koPosts)
	enTax := model.NewTaxonomyIndex(enPosts)

	msgMap := map[string]model.MessageBundle{
		"ko": {
			Common: map[string]string{"site_title": "테스트 사이트"},
		},
		"en": {
			Common: map[string]string{"site_title": "Test Site"},
		},
	}

	pages := []PageExportData{
		{
			Slug:        "about",
			Title:       "소개",
			HTMLContent: "<p>소개글</p>",
			Lang:        "ko",
		},
		{
			Slug:        "about",
			Title:       "About",
			HTMLContent: "<p>About text</p>",
			Lang:        "en",
		},
	}

	bundle := b.BuildDataBundle(koPosts, enPosts, koTax, enTax, msgMap, pages)

	if bundle == nil {
		t.Fatal("expected non-nil bundle")
	}

	if len(bundle.Languages["ko"].Posts) != 1 {
		t.Fatalf("expected 1 ko post, got %d", len(bundle.Languages["ko"].Posts))
	}
	if bundle.Languages["ko"].Posts[0].AlternateURL != "/en/posts/test-post/" {
		t.Errorf("expected alternate URL /en/posts/test-post/, got %q", bundle.Languages["ko"].Posts[0].AlternateURL)
	}

	if len(bundle.Languages["en"].Posts) != 1 {
		t.Fatalf("expected 1 en post, got %d", len(bundle.Languages["en"].Posts))
	}
	if bundle.Languages["en"].Posts[0].AlternateURL != "/posts/test-post/" {
		t.Errorf("expected alternate URL /posts/test-post/, got %q", bundle.Languages["en"].Posts[0].AlternateURL)
	}

	// ExportSiteData 파일 쓰기 테스트
	if err := b.ExportSiteData(bundle, jsonPath); err != nil {
		t.Fatalf("ExportSiteData failed: %v", err)
	}

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("failed to read exported JSON: %v", err)
	}

	var parsed SiteDataBundle
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal exported JSON: %v", err)
	}

	if parsed.Languages["ko"].Messages.Common["site_title"] != "테스트 사이트" {
		t.Errorf("expected site_title '테스트 사이트', got %q", parsed.Languages["ko"].Messages.Common["site_title"])
	}
}
