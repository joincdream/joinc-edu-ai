package builder

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/markdown"
	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/model"
	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/parser"
	templateEngine "github.com/joincdream/joinc-ai.io/tools/site-cli/internal/template"
)

// Options 빌드 파이프라인 실행 옵션
type Options struct {
	SourceDir     string        // 마크다운 소스 디렉터리 (기본: "posts")
	PagesDir      string        // 독립 단일 페이지 디렉터리 (기본: "pages")
	ThemeDir      string        // 템플릿 테마 디렉터리 (기본: "templates/default")
	OutputDir     string        // 정적 산출물 출력 디렉터리 (기본: "dist")
	IncludeDrafts bool          // status: draft 포함 여부
	BaseURL       string        // 베이스 URL (기본: "/")
	Clean         bool          // 빌드 전 출력 디렉터리 정리 여부
	ExtraHead     template.HTML // Live Reload 스크립트 등
}

// Result 빌드 실행 결과 통계
type Result struct {
	TotalPosts      int
	TotalCategories int
	Duration        time.Duration
	OutputDir       string
}

// Builder 정적 사이트 일괄 컴파일 오케스트레이터
type Builder struct {
	opts Options
}

// NewBuilder 새 빌더 인스턴스를 생성합니다.
func NewBuilder(opts Options) *Builder {
	if opts.SourceDir == "" {
		opts.SourceDir = "posts"
	}
	if opts.PagesDir == "" {
		opts.PagesDir = "pages"
	}
	if opts.ThemeDir == "" {
		opts.ThemeDir = "templates/default"
	}
	if opts.OutputDir == "" {
		opts.OutputDir = "dist"
	}
	if opts.BaseURL == "" {
		opts.BaseURL = "/"
	}
	return &Builder{opts: opts}
}

// Build 전체 정적 사이트 컴파일 파이프라인을 실행합니다.
func (b *Builder) Build() (*Result, error) {
	startTime := time.Now()

	// 1. 출력 디렉터리 초기화
	if b.opts.Clean {
		if err := os.RemoveAll(b.opts.OutputDir); err != nil {
			return nil, fmt.Errorf("builder: failed to clean output dir %q: %w", b.opts.OutputDir, err)
		}
	}
	if err := os.MkdirAll(b.opts.OutputDir, 0755); err != nil {
		return nil, fmt.Errorf("builder: failed to create output dir %q: %w", b.opts.OutputDir, err)
	}

	// 2. 마크다운 포스트 스캔
	posts, err := parser.ScanPosts(b.opts.SourceDir, b.opts.IncludeDrafts)
	if err != nil {
		return nil, fmt.Errorf("builder: scanning posts failed: %w", err)
	}

	// 3. 마크다운 AST 변환 (TOC 추출 및 Mermaid 래핑)
	converter := markdown.NewConverter()
	for _, p := range posts {
		htmlContent, toc, convErr := converter.Convert([]byte(p.RawContent))
		if convErr != nil {
			return nil, fmt.Errorf("builder: failed to convert markdown for %s: %w", p.FilePath, convErr)
		}
		p.HTMLContent = htmlContent
		p.TOC = toc
	}

	// 4. 카테고리 색인 및 최신순 정렬
	taxonomy := model.NewTaxonomyIndex(posts)

	// 5. 외부 템플릿 엔진 초기화
	engine, err := templateEngine.NewEngine(b.opts.ThemeDir)
	if err != nil {
		return nil, fmt.Errorf("builder: failed to init template engine: %w", err)
	}
	messages := engine.GetMessages()

	baseCtx := model.TemplateContext{
		SiteTitle:    messages.Common["site_title"],
		SiteSubtitle: messages.Common["site_subtitle"],
		BaseURL:      b.opts.BaseURL,
		Messages:     messages,
		Categories:   taxonomy.Categories,
		ExtraHead:    b.opts.ExtraHead,
	}

	// 6. 메인 페이지 렌더링 (`dist/index.html`)
	indexCtx := baseCtx
	indexCtx.Posts = taxonomy.AllPosts
	indexCtx.CurrentPath = "/"
	if err := renderToFile(engine, "index.html", &indexCtx, filepath.Join(b.opts.OutputDir, "index.html")); err != nil {
		return nil, err
	}

	// 7. 카테고리별 페이지 렌더링 (`dist/category/{slug}/index.html`)
	for _, cat := range taxonomy.Categories {
		catDir := filepath.Join(b.opts.OutputDir, "category", cat.Slug)
		if err := os.MkdirAll(catDir, 0755); err != nil {
			return nil, err
		}

		catCtx := baseCtx
		catCtx.ActiveCategory = cat
		catCtx.Posts = cat.Posts
		catCtx.CurrentPath = fmt.Sprintf("/category/%s/", cat.Slug)

		if err := renderToFile(engine, "category.html", &catCtx, filepath.Join(catDir, "index.html")); err != nil {
			return nil, err
		}
	}

	// 8. 개별 포스트 상세 페이지 렌더링 (`dist/posts/{slug}/index.html`)
	for _, p := range taxonomy.AllPosts {
		postDir := filepath.Join(b.opts.OutputDir, "posts", p.Slug)
		if err := os.MkdirAll(postDir, 0755); err != nil {
			return nil, err
		}

		postCtx := baseCtx
		postCtx.Post = p
		postCtx.TOC = p.TOC
		postCtx.CurrentPath = fmt.Sprintf("/posts/%s/", p.Slug)

		if err := renderToFile(engine, "detail.html", &postCtx, filepath.Join(postDir, "index.html")); err != nil {
			return nil, err
		}
	}

	// 8.5. 독립 단일 페이지 렌더링
	renderedPages := make(map[string]bool)

	// (1) pages/ 디렉터리의 독립 단일 페이지 (*.md 및 *.html) 지원
	if b.opts.PagesDir != "" {
		if info, err := os.Stat(b.opts.PagesDir); err == nil && info.IsDir() {
			entries, _ := os.ReadDir(b.opts.PagesDir)
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				ext := strings.ToLower(filepath.Ext(entry.Name()))
				pageSlug := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
				pagePath := filepath.Join(b.opts.PagesDir, entry.Name())
				pageDir := filepath.Join(b.opts.OutputDir, pageSlug)

				if ext == ".md" {
					pagePost, parseErr := parser.ParseFile(pagePath)
					if parseErr != nil || pagePost == nil {
						continue
					}
					htmlContent, toc, convErr := converter.Convert([]byte(pagePost.RawContent))
					if convErr != nil {
						continue
					}
					pagePost.HTMLContent = htmlContent
					pagePost.TOC = toc

					if err := os.MkdirAll(pageDir, 0755); err != nil {
						return nil, err
					}

					pageCtx := baseCtx
					pageCtx.Post = pagePost
					pageCtx.TOC = toc
					pageCtx.CurrentPath = fmt.Sprintf("/%s/", pageSlug)

					tmplToUse := fmt.Sprintf("%s.html", pageSlug)
					if !engine.HasTemplate(tmplToUse) {
						tmplToUse = "page.html"
						if !engine.HasTemplate("page.html") {
							tmplToUse = "detail.html"
						}
					}

					if err := renderToFile(engine, tmplToUse, &pageCtx, filepath.Join(pageDir, "index.html")); err != nil {
						return nil, err
					}
					renderedPages[pageSlug] = true
				} else if ext == ".html" {
					if err := os.MkdirAll(pageDir, 0755); err != nil {
						return nil, err
					}

					pageCtx := baseCtx
					pageCtx.CurrentPath = fmt.Sprintf("/%s/", pageSlug)

					outPath := filepath.Join(pageDir, "index.html")
					outFile, err := os.Create(outPath)
					if err != nil {
						return nil, fmt.Errorf("builder: failed to create page %q: %w", outPath, err)
					}
					renderErr := engine.RenderCustomTemplate(outFile, pagePath, &pageCtx)
					outFile.Close()
					if renderErr != nil {
						return nil, renderErr
					}
					renderedPages[pageSlug] = true
				}
			}
		}
	}

	// (2) 템플릿 디렉터리의 독립 HTML 전용 페이지 자동 렌더링 (하위 호환)
	systemTemplates := map[string]bool{
		"base.html":     true,
		"index.html":    true,
		"category.html": true,
		"detail.html":   true,
		"page.html":     true,
	}

	themeEntries, _ := os.ReadDir(b.opts.ThemeDir)
	for _, entry := range themeEntries {
		name := entry.Name()
		if entry.IsDir() || strings.ToLower(filepath.Ext(name)) != ".html" || systemTemplates[name] {
			continue
		}
		pageSlug := strings.TrimSuffix(name, filepath.Ext(name))
		if renderedPages[pageSlug] {
			continue
		}

		pageDir := filepath.Join(b.opts.OutputDir, pageSlug)
		if err := os.MkdirAll(pageDir, 0755); err != nil {
			return nil, err
		}

		pageCtx := baseCtx
		pageCtx.CurrentPath = fmt.Sprintf("/%s/", pageSlug)
		if err := renderToFile(engine, name, &pageCtx, filepath.Join(pageDir, "index.html")); err != nil {
			return nil, err
		}
	}

	// 9. 에셋 복사 (templates/<theme>/assets/ -> dist/assets/)
	themeAssets := filepath.Join(b.opts.ThemeDir, "assets")
	distAssets := filepath.Join(b.opts.OutputDir, "assets")
	if info, err := os.Stat(themeAssets); err == nil && info.IsDir() {
		if err := copyDir(themeAssets, distAssets); err != nil {
			return nil, fmt.Errorf("builder: failed to copy theme assets: %w", err)
		}
	}

	// 10. 포스트 이미지 복사 (posts/assets/ -> dist/assets/images/)
	postsAssets := filepath.Join(b.opts.SourceDir, "assets")
	if info, err := os.Stat(postsAssets); err == nil && info.IsDir() {
		if err := copyDir(postsAssets, filepath.Join(distAssets, "images")); err != nil {
			return nil, fmt.Errorf("builder: failed to copy post assets: %w", err)
		}
	}

	// 10.5. 페이지 에셋 복사 (pages/assets/ -> dist/assets/pages/)
	pagesAssets := filepath.Join(b.opts.PagesDir, "assets")
	if info, err := os.Stat(pagesAssets); err == nil && info.IsDir() {
		if err := copyDir(pagesAssets, filepath.Join(distAssets, "pages")); err != nil {
			return nil, fmt.Errorf("builder: failed to copy page assets: %w", err)
		}
	}

	return &Result{
		TotalPosts:      len(taxonomy.AllPosts),
		TotalCategories: len(taxonomy.Categories),
		Duration:        time.Since(startTime),
		OutputDir:       b.opts.OutputDir,
	}, nil
}

// renderToFile 특정 템플릿과 컨텍스트를 지정된 목적지 파일로 렌더링하여 기록합니다.
func renderToFile(engine *templateEngine.Engine, tmplName string, ctx *model.TemplateContext, destPath string) error {
	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("builder: failed to create destination file %s: %w", destPath, err)
	}
	defer f.Close()

	if err := engine.RenderPage(f, tmplName, ctx); err != nil {
		return fmt.Errorf("builder: failed to render %s to %s: %w", tmplName, destPath, err)
	}

	return nil
}

// copyDir 디렉터리 내의 모든 파일을 목적지로 재귀 복사합니다.
func copyDir(src, dest string) error {
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, destPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile 단일 파일 복사
func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
