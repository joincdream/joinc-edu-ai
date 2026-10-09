package builder

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/markdown"
	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/model"
	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/parser"
)

// PostExportData Svelte 프론트엔드로 전달되는 단일 포스트 DTO
type PostExportData struct {
	ID                 string            `json:"id"`
	Slug               string            `json:"slug"`
	FilePath           string            `json:"file_path"`
	Frontmatter        model.Frontmatter `json:"frontmatter"`
	HTMLContent        string            `json:"html_content"`
	TOC                []model.TOCItem   `json:"toc"`
	ReadingTimeMinutes int               `json:"reading_time_minutes"`
	Excerpt            string            `json:"excerpt"`
	Lang               string            `json:"lang"`
	AlternateURL       string            `json:"alternate_url,omitempty"`
}

// CategoryExportData Svelte 프론트엔드로 전달되는 카테고리 DTO
type CategoryExportData struct {
	Name      string   `json:"name"`
	Slug      string   `json:"slug"`
	PostCount int      `json:"post_count"`
	PostSlugs []string `json:"post_slugs"`
}

// PageExportData Svelte 프론트엔드로 전달되는 독립 단일 페이지 DTO (About 등)
type PageExportData struct {
	Slug         string `json:"slug"`
	Title        string `json:"title"`
	HTMLContent  string `json:"html_content"`
	Lang         string `json:"lang"`
	AlternateURL string `json:"alternate_url,omitempty"`
}

// LanguageDataBundle 언어별 데이터 세트
type LanguageDataBundle struct {
	Lang       string               `json:"lang"`
	BaseURL    string               `json:"base_url"`
	Messages   model.MessageBundle  `json:"messages"`
	Categories []CategoryExportData `json:"categories"`
	Posts      []PostExportData     `json:"posts"`
	Pages      []PageExportData     `json:"pages"`
}

// SiteDataBundle 전체 사이트 데이터 전송 계약 객체 (SSOT)
type SiteDataBundle struct {
	Version     string                         `json:"version"`
	GeneratedAt string                         `json:"generated_at"`
	Languages   map[string]*LanguageDataBundle `json:"languages"`
}

// BuildDataBundle 빌더 파이프라인 데이터로부터 SiteDataBundle을 조립합니다.
func (b *Builder) BuildDataBundle(
	koPosts []*model.Post,
	enPosts []*model.Post,
	koTaxonomy *model.TaxonomyIndex,
	enTaxonomy *model.TaxonomyIndex,
	messagesMap map[string]model.MessageBundle,
	pages []PageExportData,
) *SiteDataBundle {
	enBaseURL := "/en/"
	if b.opts.BaseURL != "/" {
		enBaseURL = fmt.Sprintf("%sen/", b.opts.BaseURL)
	}

	bundle := &SiteDataBundle{
		Version:     "1.0.0",
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Languages:   make(map[string]*LanguageDataBundle),
	}

	// 1. 한국어 데이터 번들 조립
	koBundle := &LanguageDataBundle{
		Lang:       "ko",
		BaseURL:    b.opts.BaseURL,
		Messages:   messagesMap["ko"],
		Categories: buildCategoryExportList(koTaxonomy),
		Posts:      buildPostExportList(koPosts, enPosts, "ko"),
		Pages:      filterPagesByLang(pages, "ko"),
	}
	bundle.Languages["ko"] = koBundle

	// 2. 영문 데이터 번들 조립
	enBundle := &LanguageDataBundle{
		Lang:       "en",
		BaseURL:    enBaseURL,
		Messages:   messagesMap["en"],
		Categories: buildCategoryExportList(enTaxonomy),
		Posts:      buildPostExportList(enPosts, koPosts, "en"),
		Pages:      filterPagesByLang(pages, "en"),
	}
	bundle.Languages["en"] = enBundle

	return bundle
}

// ExportSiteData 주어진 경로에 SiteDataBundle을 JSON 파일로 직렬화하여 저장합니다.
func (b *Builder) ExportSiteData(bundle *SiteDataBundle, outputPath string) error {
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("export: failed to create directory %q: %w", dir, err)
	}

	data, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return fmt.Errorf("export: failed to marshal site data bundle: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("export: failed to write site data to %q: %w", outputPath, err)
	}

	return nil
}

func buildCategoryExportList(tax *model.TaxonomyIndex) []CategoryExportData {
	if tax == nil {
		return nil
	}
	list := make([]CategoryExportData, 0, len(tax.Categories))
	for _, cat := range tax.Categories {
		var slugs []string
		for _, p := range cat.Posts {
			slugs = append(slugs, p.Slug)
		}
		list = append(list, CategoryExportData{
			Name:      cat.Name,
			Slug:      cat.Slug,
			PostCount: cat.PostCount,
			PostSlugs: slugs,
		})
	}
	return list
}

func buildPostExportList(thisPosts []*model.Post, peerPosts []*model.Post, lang string) []PostExportData {
	peerMap := make(map[string]bool, len(peerPosts))
	for _, p := range peerPosts {
		peerMap[p.Slug] = true
	}

	list := make([]PostExportData, 0, len(thisPosts))
	for _, p := range thisPosts {
		altURL := ""
		if peerMap[p.Slug] {
			if lang == "ko" {
				altURL = fmt.Sprintf("/en/posts/%s/", p.Slug)
			} else {
				altURL = fmt.Sprintf("/posts/%s/", p.Slug)
			}
		}

		list = append(list, PostExportData{
			ID:                 p.ID,
			Slug:               p.Slug,
			FilePath:           p.FilePath,
			Frontmatter:        p.Frontmatter,
			HTMLContent:        string(p.HTMLContent),
			TOC:                p.TOC,
			ReadingTimeMinutes: p.ReadingTimeMinutes,
			Excerpt:            p.Excerpt(),
			Lang:               lang,
			AlternateURL:       altURL,
		})
	}
	return list
}

func filterPagesByLang(pages []PageExportData, lang string) []PageExportData {
	var result []PageExportData
	for _, page := range pages {
		if page.Lang == lang {
			result = append(result, page)
		}
	}
	return result
}

// CollectPagesData pages/ 디렉터리 내 독립 단일 페이지를 수집하여 PageExportData 슬라이스로 반환합니다.
func (b *Builder) CollectPagesData(converter *markdown.Converter) []PageExportData {
	if b.opts.PagesDir == "" {
		return nil
	}
	info, err := os.Stat(b.opts.PagesDir)
	if err != nil || !info.IsDir() {
		return nil
	}

	entries, err := os.ReadDir(b.opts.PagesDir)
	if err != nil {
		return nil
	}

	var rawPages []PageExportData
	// 먼저 모든 페이지를 읽고 슬러그 맵을 구성합니다.
	slugLangMap := make(map[string]map[string]bool) // slug -> lang -> true

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
		lang := "ko"
		slug := nameNoExt
		if strings.HasSuffix(nameNoExt, ".en") {
			lang = "en"
			slug = strings.TrimSuffix(nameNoExt, ".en")
		}

		filePath := filepath.Join(b.opts.PagesDir, name)
		title := slug
		htmlContent := ""

		if ext == ".md" {
			pagePost, parseErr := parser.ParseFile(filePath)
			if parseErr == nil && pagePost != nil {
				if pagePost.Frontmatter.Title != "" {
					title = pagePost.Frontmatter.Title
				}
				if convHTML, _, convErr := converter.Convert([]byte(pagePost.RawContent)); convErr == nil {
					htmlContent = string(convHTML)
				}
			}
		} else if ext == ".html" {
			if contentBytes, readErr := os.ReadFile(filePath); readErr == nil {
				htmlContent = string(contentBytes)
			}
		}

		if slugLangMap[slug] == nil {
			slugLangMap[slug] = make(map[string]bool)
		}
		slugLangMap[slug][lang] = true

		rawPages = append(rawPages, PageExportData{
			Slug:        slug,
			Title:       title,
			HTMLContent: htmlContent,
			Lang:        lang,
		})
	}

	// 상호 대체 링크(AlternateURL) 설정
	for i := range rawPages {
		p := &rawPages[i]
		if p.Lang == "ko" && slugLangMap[p.Slug]["en"] {
			p.AlternateURL = fmt.Sprintf("/en/%s/", p.Slug)
		} else if p.Lang == "en" && slugLangMap[p.Slug]["ko"] {
			p.AlternateURL = fmt.Sprintf("/%s/", p.Slug)
		}
	}

	return rawPages
}
