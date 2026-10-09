package parser

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	t.Run("Valid frontmatter and body", func(t *testing.T) {
		raw := `---
title: "Sample Post Title"
category: "Cloud"
tags:
  - AWS
  - GCP
created_date: 2026-09-26
status: published
summary: "Short summary"
---

# Heading 1
This is markdown content.`

		fm, body, err := Parse(strings.NewReader(raw))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fm.Title != "Sample Post Title" {
			t.Errorf("expected title 'Sample Post Title', got %q", fm.Title)
		}
		if fm.Category != "Cloud" {
			t.Errorf("expected category 'Cloud', got %q", fm.Category)
		}
		if len(fm.Tags) != 2 || fm.Tags[0] != "AWS" {
			t.Errorf("expected tags [AWS, GCP], got %v", fm.Tags)
		}
		if !strings.Contains(body, "This is markdown content.") {
			t.Errorf("body does not contain expected content: %q", body)
		}
	})

	t.Run("Missing required field returns error", func(t *testing.T) {
		raw := `---
category: "Cloud"
created_date: 2026-09-26
---
Body text`

		_, _, err := Parse(strings.NewReader(raw))
		if err == nil {
			t.Fatal("expected error for missing title, got nil")
		}
	})
}

func TestExtractSlug(t *testing.T) {
	cases := []struct {
		path     string
		expected string
	}{
		{"posts/deep-dive/2026-09-26-agentic-workflow.md", "agentic-workflow"},
		{"posts/deep-dive/2026-09-26-agentic-workflow.en.md", "agentic-workflow"},
		{"posts/market-trends/2026-08-01-korean-제목-테스트.md", "korean-제목-테스트"},
		{"posts/simple-post.md", "simple-post"},
		{"2026-01-01.md", "2026-01-01"},
	}

	for _, c := range cases {
		got := ExtractSlug(c.path)
		if got != c.expected {
			t.Errorf("ExtractSlug(%q) = %q, want %q", c.path, got, c.expected)
		}
	}
}

func TestCalculateReadingTime(t *testing.T) {
	// 500자 이상 본문
	longText := strings.Repeat("Hello world this is a test. ", 100)
	time := CalculateReadingTime(longText)
	if time < 1 {
		t.Errorf("expected at least 1 minute, got %d", time)
	}

	// 빈 본문
	emptyTime := CalculateReadingTime("")
	if emptyTime != 1 {
		t.Errorf("expected 1 minute for empty text, got %d", emptyTime)
	}
}
