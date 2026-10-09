package parser

import (
	"os"
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

	t.Run("ScanPostsByLang separates Korean and English posts", func(t *testing.T) {
		tempDir := t.TempDir()
		koPost := "---\ntitle: \"한글 포스트\"\ncreated_date: 2026-10-09\ncategory: \"General\"\n---\n내용"
		enPost := "---\ntitle: \"English Post\"\ncreated_date: 2026-10-09\ncategory: \"General\"\n---\nContent"
		if err := os.WriteFile(filepath.Join(tempDir, "2026-10-09-test.md"), []byte(koPost), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(tempDir, "2026-10-09-test.en.md"), []byte(enPost), 0644); err != nil {
			t.Fatal(err)
		}

		koPosts, err := ScanPostsByLang(tempDir, false, "ko")
		if err != nil || len(koPosts) != 1 || koPosts[0].Frontmatter.Title != "한글 포스트" {
			t.Fatalf("expected 1 Korean post, got %v (err: %v)", koPosts, err)
		}

		enPosts, err := ScanPostsByLang(tempDir, false, "en")
		if err != nil || len(enPosts) != 1 || enPosts[0].Frontmatter.Title != "English Post" {
			t.Fatalf("expected 1 English post, got %v (err: %v)", enPosts, err)
		}

		if koPosts[0].Slug != enPosts[0].Slug {
			t.Errorf("slug mismatch: ko=%q, en=%q", koPosts[0].Slug, enPosts[0].Slug)
		}
		if koPosts[0].Slug != "test" {
			t.Errorf("expected slug 'test', got %q", koPosts[0].Slug)
		}
	})
}
