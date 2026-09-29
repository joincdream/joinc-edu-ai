package markdown

import (
	"bytes"
	"fmt"
	stdhtml "html"
	"html/template"
	"regexp"
	"strings"

	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/model"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

var mermaidRegex = regexp.MustCompile(`(?s)<pre><code class="language-mermaid">(.*?)</code></pre>`)

// Converter 마크다운 본문을 HTML로 변환하고 TOC 및 Mermaid 블록을 처리하는 엔진
type Converter struct {
	md goldmark.Markdown
}

// NewConverter 새 마크다운 변환기 인스턴스를 생성합니다.
func NewConverter() *Converter {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			extension.Strikethrough,
			extension.TaskList,
			extension.Linkify,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			goldmarkhtml.WithHardWraps(),
			goldmarkhtml.WithXHTML(),
			goldmarkhtml.WithUnsafe(), // 사용자 HTML 태그 및 인라인 요소 보존
		),
	)

	return &Converter{md: md}
}

// Convert 마크다운 텍스트를 HTML과 목차 트리(TOC)로 변환합니다.
func (c *Converter) Convert(source []byte) (template.HTML, []model.TOCItem, error) {
	reader := text.NewReader(source)
	doc := c.md.Parser().Parse(reader)

	var toc []model.TOCItem
	usedIDs := make(map[string]int)

	// AST 노드 순회: H2, H3 헤딩 추출 및 슬러그 앵커 ID 정규화
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if heading, ok := n.(*ast.Heading); ok {
			if heading.Level == 2 || heading.Level == 3 {
				title := extractHeadingText(heading, source)
				if title != "" {
					slug := model.Slugify(title)
					// 중복 ID 방어 (slug, slug-1, slug-2 ...)
					if count, exists := usedIDs[slug]; exists {
						usedIDs[slug] = count + 1
						slug = fmt.Sprintf("%s-%d", slug, count)
					} else {
						usedIDs[slug] = 1
					}

					heading.SetAttributeString("id", []byte(slug))

					toc = append(toc, model.TOCItem{
						ID:    slug,
						Title: title,
						Level: heading.Level,
					})
				}
			}
		}

		// 이미지 상대 경로 자동 보정 (예: ../assets/img.png -> /assets/images/img.png)
		if img, ok := n.(*ast.Image); ok {
			dest := string(img.Destination)
			if newDest, modified := normalizeImagePath(dest); modified {
				img.Destination = []byte(newDest)
			}
		}

		return ast.WalkContinue, nil
	})

	if err != nil {
		return "", nil, fmt.Errorf("failed to process markdown AST: %w", err)
	}

	var buf bytes.Buffer
	if err := c.md.Renderer().Render(&buf, source, doc); err != nil {
		return "", nil, fmt.Errorf("failed to render markdown to HTML: %w", err)
	}

	renderedHTML := buf.String()

	// Mermaid 코드 블록을 클라이언트 렌더링 컨테이너로 치환
	renderedHTML = mermaidRegex.ReplaceAllStringFunc(renderedHTML, func(match string) string {
		submatches := mermaidRegex.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}
		rawCode := strings.TrimSpace(stdhtml.UnescapeString(submatches[1]))
		return fmt.Sprintf(
			`<div class="mermaid-container my-8 p-6 rounded-lg border border-[var(--color-border-subtle)] bg-white shadow-md flex flex-col items-center justify-center overflow-x-auto"><div class="mermaid-raw" style="display: none;">%s</div></div>`,
			stdhtml.EscapeString(rawCode),
		)
	})

	return template.HTML(renderedHTML), toc, nil
}

// extractHeadingText 헤딩 노드 하위의 모든 텍스트 요소를 안전하게 수집합니다.
func extractHeadingText(n ast.Node, source []byte) string {
	var buf bytes.Buffer
	_ = ast.Walk(n, func(child ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if textNode, ok := child.(*ast.Text); ok {
				buf.Write(textNode.Text(source))
			}
		}
		return ast.WalkContinue, nil
	})
	return strings.TrimSpace(buf.String())
}

// normalizeImagePath 마크다운 내 상대 에셋 경로를 정적 배포 절대 경로로 정규화합니다.
func normalizeImagePath(dest string) (string, bool) {
	// 1. 외부 URL (http://, https://, //) 또는 data URI는 보정 대상에서 제외
	if strings.Contains(dest, "://") || strings.HasPrefix(dest, "//") || strings.HasPrefix(dest, "data:") {
		return dest, false
	}

	// 2. 이미 배포 절대 경로(/assets/...)인 경우 보정 대상에서 제외
	if strings.HasPrefix(dest, "/assets/") {
		return dest, false
	}

	// 3. ../assets/, ./assets/, assets/, posts/assets/ 등 assets/ 디렉터리를 가리키는 상대 경로 치환
	if idx := strings.Index(dest, "assets/"); idx != -1 {
		subPath := dest[idx+len("assets/"):]
		return "/assets/images/" + strings.TrimPrefix(subPath, "/"), true
	}

	return dest, false
}
