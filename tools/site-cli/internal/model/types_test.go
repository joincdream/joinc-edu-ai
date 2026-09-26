package model

import (
	"testing"
)

func TestFrontmatter_Validate(t *testing.T) {
	tests := []struct {
		name    string
		fm      Frontmatter
		wantErr bool
	}{
		{
			name: "Valid frontmatter",
			fm: Frontmatter{
				Title:       "Valid Title",
				Category:    "Tech",
				CreatedDate: "2026-09-26",
				Status:      "published",
			},
			wantErr: false,
		},
		{
			name: "Missing title",
			fm: Frontmatter{
				Category:    "Tech",
				CreatedDate: "2026-09-26",
			},
			wantErr: true,
		},
		{
			name: "Missing category falls back to General",
			fm: Frontmatter{
				Title:       "Valid Title",
				CreatedDate: "2026-09-26",
			},
			wantErr: false,
		},
		{
			name: "Missing created_date",
			fm: Frontmatter{
				Title:    "Valid Title",
				Category: "Tech",
			},
			wantErr: true,
		},
		{
			name: "Invalid created_date format",
			fm: Frontmatter{
				Title:       "Valid Title",
				Category:    "Tech",
				CreatedDate: "26-09-2026",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fm.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewTaxonomyIndex(t *testing.T) {
	posts := []*Post{
		{
			ID:   "p1",
			Slug: "post-1",
			Frontmatter: Frontmatter{
				Title:       "Oldest Post",
				Category:    "Go",
				CreatedDate: "2026-01-01",
			},
		},
		{
			ID:   "p2",
			Slug: "post-2",
			Frontmatter: Frontmatter{
				Title:       "Newest Post",
				Category:    "AI",
				CreatedDate: "2026-09-26",
			},
		},
		{
			ID:   "p3",
			Slug: "post-3",
			Frontmatter: Frontmatter{
				Title:       "Middle Post in AI",
				Category:    "AI",
				CreatedDate: "2026-05-15",
			},
		},
	}

	idx := NewTaxonomyIndex(posts)

	// 1. AllPosts 최신순 정렬 검증
	if len(idx.AllPosts) != 3 {
		t.Fatalf("expected 3 posts, got %d", len(idx.AllPosts))
	}
	if idx.AllPosts[0].ID != "p2" || idx.AllPosts[1].ID != "p3" || idx.AllPosts[2].ID != "p1" {
		t.Errorf("posts are not sorted by created_date DESC: got [%s, %s, %s]",
			idx.AllPosts[0].ID, idx.AllPosts[1].ID, idx.AllPosts[2].ID)
	}

	// 2. 카테고리 집계 검증
	if len(idx.Categories) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(idx.Categories))
	}

	aiCat := idx.CategoryMap["ai"]
	if aiCat == nil || aiCat.PostCount != 2 {
		t.Errorf("expected AI category with 2 posts, got %v", aiCat)
	}

	goCat := idx.CategoryMap["go"]
	if goCat == nil || goCat.PostCount != 1 {
		t.Errorf("expected Go category with 1 post, got %v", goCat)
	}
}

func TestSlugify(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"Agentic AI & LLMs", "agentic-ai-llms"},
		{"하네스 엔지니어링", "하네스-엔지니어링"},
		{"   Multiple---Spaces & Hyphens   ", "multiple-spaces-hyphens"},
		{"", "default"},
	}

	for _, c := range cases {
		got := Slugify(c.input)
		if got != c.expected {
			t.Errorf("Slugify(%q) = %q, want %q", c.input, got, c.expected)
		}
	}
}

func TestMessageBundle_Get(t *testing.T) {
	bundle := MessageBundle{
		Nav: map[string]string{
			"home": "Home",
		},
	}

	if val := bundle.Get("nav", "home", "Default"); val != "Home" {
		t.Errorf("expected 'Home', got %q", val)
	}
	if val := bundle.Get("nav", "missing_key", "Fallback"); val != "Fallback" {
		t.Errorf("expected 'Fallback', got %q", val)
	}
	if val := bundle.Get("unknown_section", "any", "Fallback"); val != "Fallback" {
		t.Errorf("expected 'Fallback', got %q", val)
	}
}
