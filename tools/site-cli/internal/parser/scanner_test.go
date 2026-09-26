package parser

import (
	"path/filepath"
	"testing"
)

func TestScanPosts(t *testing.T) {
	validDir := filepath.Join("..", "..", "testdata", "valid_posts")
	invalidDir := filepath.Join("..", "..", "testdata", "invalid_posts")

	t.Run("Scan valid posts without drafts", func(t *testing.T) {
		posts, err := ScanPosts(validDir, false)
		if err != nil {
			t.Fatalf("unexpected scan error: %v", err)
		}
		if len(posts) != 2 {
			t.Errorf("expected 2 published posts, got %d", len(posts))
		}
		for _, p := range posts {
			if p.Frontmatter.IsDraft() {
				t.Errorf("found draft post in published scan: %s", p.Slug)
			}
		}
	})

	t.Run("Scan valid posts with drafts", func(t *testing.T) {
		posts, err := ScanPosts(validDir, true)
		if err != nil {
			t.Fatalf("unexpected scan error: %v", err)
		}
		if len(posts) != 3 {
			t.Errorf("expected 3 posts including drafts, got %d", len(posts))
		}
	})

	t.Run("Scan invalid posts triggers fail-fast error", func(t *testing.T) {
		_, err := ScanPosts(invalidDir, false)
		if err == nil {
			t.Fatal("expected fail-fast error on invalid frontmatter, got nil")
		}
	})

	t.Run("Scan non-existent directory returns error", func(t *testing.T) {
		_, err := ScanPosts("non-existent-dir-12345", false)
		if err == nil {
			t.Fatal("expected error for non-existent directory, got nil")
		}
	})
}
