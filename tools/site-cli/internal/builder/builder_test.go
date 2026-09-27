package builder

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuilder_Build(t *testing.T) {
	tempDist := t.TempDir()
	sourceDir := filepath.Join("..", "..", "testdata", "valid_posts")
	themeDir := filepath.Join("..", "..", "..", "..", "templates", "default")


	opts := Options{
		SourceDir:     sourceDir,
		ThemeDir:      themeDir,
		OutputDir:     tempDist,
		IncludeDrafts: false,
		Clean:         true,
	}

	b := NewBuilder(opts)
	res, err := b.Build()
	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	// 1. 결과 메트릭 검증
	if res.TotalPosts != 2 {
		t.Errorf("expected 2 published posts, got %d", res.TotalPosts)
	}
	if res.TotalCategories != 2 {
		t.Errorf("expected 2 categories (AI, Golang), got %d", res.TotalCategories)
	}
	if res.Duration <= 0 {
		t.Errorf("expected positive duration, got %v", res.Duration)
	}

	// 2. index.html 파일 검증
	indexPath := filepath.Join(tempDist, "index.html")
	indexBytes, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("missing index.html: %v", err)
	}
	indexContent := string(indexBytes)
	if !strings.Contains(indexContent, "Valid Post 1") || !strings.Contains(indexContent, "Valid Post 2") {
		t.Errorf("index.html missing post titles: %s", indexContent)
	}
	if !strings.Contains(indexContent, "AI Info") {
		t.Errorf("index.html missing site title: %s", indexContent)
	}

	// 3. category/ai/index.html 검증
	catPath := filepath.Join(tempDist, "category", "ai", "index.html")
	catBytes, err := os.ReadFile(catPath)
	if err != nil {
		t.Fatalf("missing category/ai/index.html: %v", err)
	}
	if !strings.Contains(string(catBytes), "Valid Post 1") {
		t.Errorf("category/ai/index.html missing Valid Post 1")
	}

	// 4. posts/valid-1/index.html 검증
	postPath := filepath.Join(tempDist, "posts", "valid-1", "index.html")
	postBytes, err := os.ReadFile(postPath)
	if err != nil {
		t.Fatalf("missing posts/valid-1/index.html: %v", err)
	}
	postContent := string(postBytes)
	if !strings.Contains(postContent, "Valid Post 1") {
		t.Errorf("posts/valid-1/index.html missing title")
	}
	if !strings.Contains(postContent, "Content 1") {
		t.Errorf("posts/valid-1/index.html missing body content")
	}

	// 5. assets/style.css 복사 검증
	cssPath := filepath.Join(tempDist, "assets", "style.css")
	if _, err := os.Stat(cssPath); os.IsNotExist(err) {
		t.Errorf("missing assets/style.css in output dir")
	}
}

func TestBuilder_Build_WithRedirects(t *testing.T) {
	tempDist := t.TempDir()
	sourceDir := filepath.Join("..", "..", "testdata", "valid_posts")
	themeDir := filepath.Join("..", "..", "..", "..", "templates", "default")

	// 임시 redirect.yaml 작성
	redirectYaml := `
query_redirects:
  - source_path: "/posts/detail"
    query_key: "id"
    default_target: "/"
    mappings:
      "19": "/posts/spec-driven-ai-harness/"
      "4": "/posts/weekly-ai-trend/"
path_redirects:
  - source_path: "/old-path"
    target_path: "/new-path/"
`
	redirectFile := filepath.Join(t.TempDir(), "redirect.yaml")
	if err := os.WriteFile(redirectFile, []byte(redirectYaml), 0644); err != nil {
		t.Fatalf("failed to write temp redirect file: %v", err)
	}

	opts := Options{
		SourceDir:     sourceDir,
		ThemeDir:      themeDir,
		OutputDir:     tempDist,
		IncludeDrafts: false,
		Clean:         true,
		RedirectsFile: redirectFile,
	}

	b := NewBuilder(opts)
	_, err := b.Build()
	if err != nil {
		t.Fatalf("Build() with redirects failed: %v", err)
	}

	// 1. query_redirects 검증: tempDist/posts/detail/index.html
	detailRedirectPath := filepath.Join(tempDist, "posts", "detail", "index.html")
	detailBytes, err := os.ReadFile(detailRedirectPath)
	if err != nil {
		t.Fatalf("expected redirect file at %s, but not found: %v", detailRedirectPath, err)
	}
	detailContent := string(detailBytes)
	if !strings.Contains(detailContent, "window.location.replace") {
		t.Errorf("expected window.location.replace in redirect HTML")
	}
	if !strings.Contains(detailContent, "spec-driven-ai-harness") {
		t.Errorf("expected mapping destination in redirect HTML")
	}
	if !strings.Contains(detailContent, "새로운 페이지로 안전하게 이동 중입니다...") {
		t.Errorf("expected user guide message in redirect HTML")
	}

	// 2. path_redirects 검증: tempDist/old-path/index.html
	oldPathRedirect := filepath.Join(tempDist, "old-path", "index.html")
	oldBytes, err := os.ReadFile(oldPathRedirect)
	if err != nil {
		t.Fatalf("expected path redirect file at %s, but not found: %v", oldPathRedirect, err)
	}
	oldContent := string(oldBytes)
	if !strings.Contains(oldContent, "/new-path/") {
		t.Errorf("expected /new-path/ in path redirect HTML")
	}
}
