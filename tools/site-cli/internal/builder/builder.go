package builder

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
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
	RedirectsFile string        // 하위 호환 리다이렉트 설정 파일 (기본: "redirect.yaml")
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
	if opts.RedirectsFile == "" {
		opts.RedirectsFile = "redirect.yaml"
	}
	return &Builder{opts: opts}
}

// Build 전체 정적 사이트 컴파일 파이프라인을 실행합니다.
// Build 전체 정적 사이트 컴파일 파이프라인을 실행합니다.
// 한국어 기본 사이트(dist/)와 영문 서브패스 사이트(dist/en/)를 순차 렌더링합니다.
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

	// 2. 마크다운 포스트 언어별 스캔 (ko vs en)
	koPosts, err := parser.ScanPostsByLang(b.opts.SourceDir, b.opts.IncludeDrafts, "ko")
	if err != nil {
		return nil, fmt.Errorf("builder: scanning Korean posts failed: %w", err)
	}
	enPosts, err := parser.ScanPostsByLang(b.opts.SourceDir, b.opts.IncludeDrafts, "en")
	if err != nil {
		return nil, fmt.Errorf("builder: scanning English posts failed: %w", err)
	}

	// 3. 마크다운 AST 변환 (TOC 추출 및 Mermaid 래핑)
	converter := markdown.NewConverter()
	for _, p := range koPosts {
		htmlContent, toc, convErr := converter.Convert([]byte(p.RawContent))
		if convErr != nil {
			return nil, fmt.Errorf("builder: failed to convert markdown for %s: %w", p.FilePath, convErr)
		}
		p.HTMLContent = htmlContent
		p.TOC = toc
	}
	for _, p := range enPosts {
		htmlContent, toc, convErr := converter.Convert([]byte(p.RawContent))
		if convErr != nil {
			return nil, fmt.Errorf("builder: failed to convert markdown for %s: %w", p.FilePath, convErr)
		}
		p.HTMLContent = htmlContent
		p.TOC = toc
	}

	// 4. 언어별 포스트 슬러그 맵 및 색인 생성
	koSlugMap := make(map[string]*model.Post)
	for _, p := range koPosts {
		koSlugMap[p.Slug] = p
	}
	enSlugMap := make(map[string]*model.Post)
	for _, p := range enPosts {
		enSlugMap[p.Slug] = p
	}

	koTaxonomy := model.NewTaxonomyIndex(koPosts)
	enTaxonomy := model.NewTaxonomyIndex(enPosts)

	// 5. pages/ 디렉터리 내 독립 페이지 언어별 맵 파악
	koPagesMap := make(map[string]bool)
	enPagesMap := make(map[string]bool)
	if b.opts.PagesDir != "" {
		if entries, readErr := os.ReadDir(b.opts.PagesDir); readErr == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				name := entry.Name()
				ext := strings.ToLower(filepath.Ext(name))
				if ext != ".html" && ext != ".md" {
					continue
				}
				nameNoExt := strings.TrimSuffix(name, ext)
				if strings.HasSuffix(nameNoExt, ".en") {
					slug := strings.TrimSuffix(nameNoExt, ".en")
					enPagesMap[slug] = true
				} else {
					koPagesMap[nameNoExt] = true
				}
			}
		}
	}

	// 6. 외부 템플릿 엔진 초기화
	engine, err := templateEngine.NewEngine(b.opts.ThemeDir)
	if err != nil {
		return nil, fmt.Errorf("builder: failed to init template engine: %w", err)
	}

	// 7. Svelte 하이브리드 SSG용 통합 데이터 번들 생성 및 내보내기 (최신순 정렬된 AllPosts 전달)
	pages := b.CollectPagesData(converter)
	bundle := b.BuildDataBundle(koTaxonomy.AllPosts, enTaxonomy.AllPosts, koTaxonomy, enTaxonomy, engine.GetMessagesMap(), pages)
	dataCachePath := filepath.Join(b.opts.OutputDir, ".cache", "site-data.json")
	if err := b.ExportSiteData(bundle, dataCachePath); err != nil {
		return nil, fmt.Errorf("builder: exporting site data failed: %w", err)
	}
	_ = b.ExportSiteData(bundle, filepath.Join(".cache", "site-data.json"))

	if b.isSvelteTheme() {
		// Svelte 5 SSR 기반 정적 HTML 일괄 렌더링
		if err := b.runSvelteSSG(dataCachePath, b.opts.OutputDir); err != nil {
			return nil, fmt.Errorf("builder: svelte SSG build failed: %w", err)
		}
	} else {
		// 기존 Go HTML 템플릿 기반 렌더링 (default-light 및 default 다크 테마 호환)
		// 8. 한국어 기본 사이트 렌더링 (dist/)
		if err := b.renderSite(engine, converter, &langRenderConfig{
			Lang:         "ko",
			BaseURL:      b.opts.BaseURL,
			OutputDir:    b.opts.OutputDir,
			Taxonomy:     koTaxonomy,
			PeerTaxonomy: enTaxonomy,
			ThisPostsMap: koSlugMap,
			PeerPostsMap: enSlugMap,
			ThisPagesMap: koPagesMap,
			PeerPagesMap: enPagesMap,
		}); err != nil {
			return nil, fmt.Errorf("builder: failed to render Korean site: %w", err)
		}

		// 9. 영문 서브패스 사이트 렌더링 (dist/en/)
		enOutputDir := filepath.Join(b.opts.OutputDir, "en")
		if err := os.MkdirAll(enOutputDir, 0755); err != nil {
			return nil, fmt.Errorf("builder: failed to create en output dir %q: %w", enOutputDir, err)
		}
		enBaseURL := "/en/"
		if b.opts.BaseURL != "/" {
			enBaseURL = strings.TrimSuffix(b.opts.BaseURL, "/") + "/en/"
		}

		if err := b.renderSite(engine, converter, &langRenderConfig{
			Lang:         "en",
			BaseURL:      enBaseURL,
			OutputDir:    enOutputDir,
			Taxonomy:     enTaxonomy,
			PeerTaxonomy: koTaxonomy,
			ThisPostsMap: enSlugMap,
			PeerPostsMap: koSlugMap,
			ThisPagesMap: enPagesMap,
			PeerPagesMap: koPagesMap,
		}); err != nil {
			return nil, fmt.Errorf("builder: failed to render English site: %w", err)
		}
	}

	// 10. 에셋 복사 (templates/<theme>/assets/ -> dist/assets/)
	themeAssets := filepath.Join(b.opts.ThemeDir, "assets")
	distAssets := filepath.Join(b.opts.OutputDir, "assets")
	if info, err := os.Stat(themeAssets); err == nil && info.IsDir() {
		if err := copyDir(themeAssets, distAssets); err != nil {
			return nil, fmt.Errorf("builder: failed to copy theme assets: %w", err)
		}
	}

	// 11. 포스트 이미지 복사 (posts/assets/ -> dist/assets/images/)
	postsAssets := filepath.Join(b.opts.SourceDir, "assets")
	if info, err := os.Stat(postsAssets); err == nil && info.IsDir() {
		if err := copyDir(postsAssets, filepath.Join(distAssets, "images")); err != nil {
			return nil, fmt.Errorf("builder: failed to copy post assets: %w", err)
		}
	}

	// 11.5. 페이지 에셋 복사 (pages/assets/ -> dist/assets/pages/)
	pagesAssets := filepath.Join(b.opts.PagesDir, "assets")
	if info, err := os.Stat(pagesAssets); err == nil && info.IsDir() {
		if err := copyDir(pagesAssets, filepath.Join(distAssets, "pages")); err != nil {
			return nil, fmt.Errorf("builder: failed to copy page assets: %w", err)
		}
	}

	// 12. 하위 호환 리다이렉트 페이지 생성
	if err := b.generateRedirects(); err != nil {
		return nil, fmt.Errorf("builder: generating redirects failed: %w", err)
	}

	totalPosts := len(koTaxonomy.AllPosts) + len(enTaxonomy.AllPosts)
	return &Result{
		TotalPosts:      totalPosts,
		TotalCategories: len(koTaxonomy.Categories),
		Duration:        time.Since(startTime),
		OutputDir:       b.opts.OutputDir,
	}, nil
}

// langRenderConfig 언어별 사이트 렌더링 매개변수
type langRenderConfig struct {
	Lang         string // "ko" 또는 "en"
	BaseURL      string // "/" 또는 "/en/"
	OutputDir    string // dist 또는 dist/en
	Taxonomy     *model.TaxonomyIndex
	PeerTaxonomy *model.TaxonomyIndex
	ThisPostsMap map[string]*model.Post
	PeerPostsMap map[string]*model.Post
	ThisPagesMap map[string]bool
	PeerPagesMap map[string]bool
}

// renderSite 특정 언어의 메인, 카테고리, 포스트 상세, 독립 페이지를 렌더링합니다.
func (b *Builder) renderSite(engine *templateEngine.Engine, converter *markdown.Converter, cfg *langRenderConfig) error {
	messages := engine.GetMessagesFor(cfg.Lang)

	baseCtx := model.TemplateContext{
		SiteTitle:    messages.Common["site_title"],
		SiteSubtitle: messages.Common["site_subtitle"],
		BaseURL:      cfg.BaseURL,
		CurrentLang:  cfg.Lang,
		Messages:     messages,
		Categories:   cfg.Taxonomy.Categories,
		ExtraHead:    b.opts.ExtraHead,
	}

	// 1. 메인 홈 페이지 렌더링 (`{OutputDir}/index.html`)
	indexCtx := baseCtx
	indexCtx.Posts = cfg.Taxonomy.AllPosts
	indexCtx.CurrentPath = cfg.BaseURL
	indexCtx.SwitchURL = map[string]template.URL{
		"ko": template.URL("/"),
		"en": template.URL("/en/"),
	}
	indexCtx.AlternateLangs = []model.AlternateLink{
		{Lang: "ko", URL: "/"},
		{Lang: "en", URL: "/en/"},
		{Lang: "x-default", URL: "/"},
	}
	if err := renderToFile(engine, "index.html", &indexCtx, filepath.Join(cfg.OutputDir, "index.html")); err != nil {
		return err
	}

	// 2. 카테고리별 페이지 렌더링 (`{OutputDir}/category/{slug}/index.html`)
	for _, cat := range cfg.Taxonomy.Categories {
		catDir := filepath.Join(cfg.OutputDir, "category", cat.Slug)
		if err := os.MkdirAll(catDir, 0755); err != nil {
			return err
		}

		catCtx := baseCtx
		catCtx.ActiveCategory = cat
		catCtx.Posts = cat.Posts
		catCtx.CurrentPath = fmt.Sprintf("%scategory/%s/", cfg.BaseURL, cat.Slug)

		peerCatURL := "/en/"
		if cfg.Lang == "en" {
			peerCatURL = "/"
			if cfg.PeerTaxonomy != nil && cfg.PeerTaxonomy.CategoryMap[cat.Slug] != nil {
				peerCatURL = fmt.Sprintf("/category/%s/", cat.Slug)
			}
		} else {
			if cfg.PeerTaxonomy != nil && cfg.PeerTaxonomy.CategoryMap[cat.Slug] != nil {
				peerCatURL = fmt.Sprintf("/en/category/%s/", cat.Slug)
			}
		}

		if cfg.Lang == "ko" {
			catCtx.SwitchURL = map[string]template.URL{
				"ko": template.URL(catCtx.CurrentPath),
				"en": template.URL(peerCatURL),
			}
		} else {
			catCtx.SwitchURL = map[string]template.URL{
				"ko": template.URL(peerCatURL),
				"en": template.URL(catCtx.CurrentPath),
			}
		}

		if err := renderToFile(engine, "category.html", &catCtx, filepath.Join(catDir, "index.html")); err != nil {
			return err
		}
	}

	// 3. 개별 포스트 상세 페이지 렌더링 (`{OutputDir}/posts/{slug}/index.html`)
	for _, p := range cfg.Taxonomy.AllPosts {
		postDir := filepath.Join(cfg.OutputDir, "posts", p.Slug)
		if err := os.MkdirAll(postDir, 0755); err != nil {
			return err
		}

		postCtx := baseCtx
		postCtx.Post = p
		postCtx.TOC = p.TOC
		postCtx.CurrentPath = fmt.Sprintf("%sposts/%s/", cfg.BaseURL, p.Slug)

		peerPostURL := "/en/"
		hasPeerPost := cfg.PeerPostsMap[p.Slug] != nil
		if cfg.Lang == "en" {
			peerPostURL = "/"
			if hasPeerPost {
				peerPostURL = fmt.Sprintf("/posts/%s/", p.Slug)
			}
			postCtx.SwitchURL = map[string]template.URL{
				"ko": template.URL(peerPostURL),
				"en": template.URL(postCtx.CurrentPath),
			}
		} else {
			if hasPeerPost {
				peerPostURL = fmt.Sprintf("/en/posts/%s/", p.Slug)
			}
			postCtx.SwitchURL = map[string]template.URL{
				"ko": template.URL(postCtx.CurrentPath),
				"en": template.URL(peerPostURL),
			}
		}

		if hasPeerPost {
			postCtx.AlternateLangs = []model.AlternateLink{
				{Lang: "ko", URL: fmt.Sprintf("/posts/%s/", p.Slug)},
				{Lang: "en", URL: fmt.Sprintf("/en/posts/%s/", p.Slug)},
				{Lang: "x-default", URL: fmt.Sprintf("/posts/%s/", p.Slug)},
			}
		}

		if err := renderToFile(engine, "detail.html", &postCtx, filepath.Join(postDir, "index.html")); err != nil {
			return err
		}
	}

	// 4. 독립 단일 페이지 렌더링 (pages/ 디렉터리 내 해당 언어 파일 대상)
	renderedPages := make(map[string]bool)
	if b.opts.PagesDir != "" {
		if info, err := os.Stat(b.opts.PagesDir); err == nil && info.IsDir() {
			entries, _ := os.ReadDir(b.opts.PagesDir)
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				name := entry.Name()
				ext := strings.ToLower(filepath.Ext(name))
				if ext != ".md" && ext != ".html" {
					continue
				}

				nameNoExt := strings.TrimSuffix(name, ext)
				isEnPage := strings.HasSuffix(nameNoExt, ".en")
				pageSlug := nameNoExt
				if isEnPage {
					pageSlug = strings.TrimSuffix(nameNoExt, ".en")
				}

				// 현재 렌더링 언어와 페이지 파일의 언어 매칭 필터링
				if cfg.Lang == "en" && !isEnPage {
					continue
				}
				if cfg.Lang == "ko" && isEnPage {
					continue
				}

				pagePath := filepath.Join(b.opts.PagesDir, name)
				pageDir := filepath.Join(cfg.OutputDir, pageSlug)

				peerPageURL := "/en/"
				if cfg.Lang == "en" {
					peerPageURL = "/"
					if cfg.PeerPagesMap[pageSlug] {
						peerPageURL = fmt.Sprintf("/%s/", pageSlug)
					}
				} else {
					if cfg.PeerPagesMap[pageSlug] {
						peerPageURL = fmt.Sprintf("/en/%s/", pageSlug)
					}
				}

				switchMap := map[string]template.URL{
					"ko": template.URL(fmt.Sprintf("/%s/", pageSlug)),
					"en": template.URL(peerPageURL),
				}
				if cfg.Lang == "en" {
					switchMap = map[string]template.URL{
						"ko": template.URL(peerPageURL),
						"en": template.URL(fmt.Sprintf("/en/%s/", pageSlug)),
					}
				}

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
						return err
					}

					pageCtx := baseCtx
					pageCtx.Post = pagePost
					pageCtx.TOC = toc
					pageCtx.CurrentPath = fmt.Sprintf("%s%s/", cfg.BaseURL, pageSlug)
					pageCtx.SwitchURL = switchMap

					tmplToUse := fmt.Sprintf("%s.html", pageSlug)
					if !engine.HasTemplate(tmplToUse) {
						tmplToUse = "page.html"
						if !engine.HasTemplate("page.html") {
							tmplToUse = "detail.html"
						}
					}

					if err := renderToFile(engine, tmplToUse, &pageCtx, filepath.Join(pageDir, "index.html")); err != nil {
						return err
					}
					renderedPages[pageSlug] = true
				} else if ext == ".html" {
					if err := os.MkdirAll(pageDir, 0755); err != nil {
						return err
					}

					pageCtx := baseCtx
					pageCtx.CurrentPath = fmt.Sprintf("%s%s/", cfg.BaseURL, pageSlug)
					pageCtx.SwitchURL = switchMap

					outPath := filepath.Join(pageDir, "index.html")
					outFile, err := os.Create(outPath)
					if err != nil {
						return fmt.Errorf("builder: failed to create page %q: %w", outPath, err)
					}
					renderErr := engine.RenderCustomTemplate(outFile, pagePath, &pageCtx)
					outFile.Close()
					if renderErr != nil {
						return renderErr
					}
					renderedPages[pageSlug] = true
				}
			}
		}
	}

	return nil
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

// isSvelteTheme 해당 테마 디렉터리가 Svelte 5 SSG 테마인지 확인합니다.
func (b *Builder) isSvelteTheme() bool {
	ssgScript := filepath.Join(b.opts.ThemeDir, "scripts", "build-ssg.ts")
	pkgJson := filepath.Join(b.opts.ThemeDir, "package.json")
	if _, err := os.Stat(ssgScript); err == nil {
		if _, err := os.Stat(pkgJson); err == nil {
			return true
		}
	}
	return false
}

// runSvelteSSG Svelte 5 SSR 기반 정적 HTML 일괄 렌더링 서브프로세스를 실행합니다.
func (b *Builder) runSvelteSSG(dataPath, outDir string) error {
	absData, err := filepath.Abs(dataPath)
	if err != nil {
		absData = dataPath
	}
	absOut, err := filepath.Abs(outDir)
	if err != nil {
		absOut = outDir
	}

	cmd := exec.Command("npm", "run", "build:ssg", "--", "--data", absData, "--out", absOut)
	cmd.Dir = b.opts.ThemeDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
