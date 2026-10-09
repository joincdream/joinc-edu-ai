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

func TestBuilder_Build_Multilingual(t *testing.T) {
	tempDist := t.TempDir()
	tempPosts := t.TempDir()
	tempPages := t.TempDir()
	themeDir := filepath.Join("..", "..", "..", "..", "templates", "default-light")

	// 1. 한국어 포스트 & 영문 포스트 생성
	koPost := "---\ntitle: \"쿠버네티스 심층 분석\"\ncreated_date: 2026-10-09\ncategory: \"Cloud\"\n---\n한국어 본문입니다."
	enPost := "---\ntitle: \"Kubernetes Deep Dive\"\ncreated_date: 2026-10-09\ncategory: \"Cloud\"\n---\nEnglish body content."
	if err := os.WriteFile(filepath.Join(tempPosts, "2026-10-09-k8s-dive.md"), []byte(koPost), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tempPosts, "2026-10-09-k8s-dive.en.md"), []byte(enPost), 0644); err != nil {
		t.Fatal(err)
	}

	// 2. 한국어 페이지 & 영문 페이지 생성
	koAbout := `{{ define "content" }}<h1>소개</h1>{{ end }}`
	enAbout := `{{ define "content" }}<h1>About Me</h1>{{ end }}`
	if err := os.WriteFile(filepath.Join(tempPages, "about.html"), []byte(koAbout), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tempPages, "about.en.html"), []byte(enAbout), 0644); err != nil {
		t.Fatal(err)
	}

	opts := Options{
		SourceDir: tempPosts,
		PagesDir:  tempPages,
		ThemeDir:  themeDir,
		OutputDir: tempDist,
		Clean:     true,
	}

	b := NewBuilder(opts)
	res, err := b.Build()
	if err != nil {
		t.Fatalf("Multilingual Build() failed: %v", err)
	}

	if res.TotalPosts != 2 {
		t.Errorf("expected 2 total posts (1 ko + 1 en), got %d", res.TotalPosts)
	}

	// 3. 한국어 메인 및 영문 메인 검증
	koIndex, err := os.ReadFile(filepath.Join(tempDist, "index.html"))
	if err != nil {
		t.Fatalf("missing dist/index.html: %v", err)
	}
	if !strings.Contains(string(koIndex), "쿠버네티스 심층 분석") {
		t.Errorf("missing Korean post title in dist/index.html")
	}

	enIndex, err := os.ReadFile(filepath.Join(tempDist, "en", "index.html"))
	if err != nil {
		t.Fatalf("missing dist/en/index.html: %v", err)
	}
	if !strings.Contains(string(enIndex), "Kubernetes Deep Dive") {
		t.Errorf("missing English post title in dist/en/index.html")
	}

	// 4. 포스트 상세 페이지 검증
	koPostHTML, err := os.ReadFile(filepath.Join(tempDist, "posts", "k8s-dive", "index.html"))
	if err != nil {
		t.Fatalf("missing dist/posts/k8s-dive/index.html: %v", err)
	}
	if !strings.Contains(string(koPostHTML), "한국어 본문입니다") {
		t.Errorf("missing Korean body in post detail")
	}

	enPostHTML, err := os.ReadFile(filepath.Join(tempDist, "en", "posts", "k8s-dive", "index.html"))
	if err != nil {
		t.Fatalf("missing dist/en/posts/k8s-dive/index.html: %v", err)
	}
	if !strings.Contains(string(enPostHTML), "English body content") {
		t.Errorf("missing English body in post detail")
	}

	// 5. 소개 페이지 검증
	koAboutHTML, err := os.ReadFile(filepath.Join(tempDist, "about", "index.html"))
	if err != nil {
		t.Fatalf("missing dist/about/index.html: %v", err)
	}
	if !strings.Contains(string(koAboutHTML), "소개") {
		t.Errorf("missing Korean about content")
	}

	enAboutHTML, err := os.ReadFile(filepath.Join(tempDist, "en", "about", "index.html"))
	if err != nil {
		t.Fatalf("missing dist/en/about/index.html: %v", err)
	}
	if !strings.Contains(string(enAboutHTML), "About Me") {
		t.Errorf("missing English about content")
	}
}
