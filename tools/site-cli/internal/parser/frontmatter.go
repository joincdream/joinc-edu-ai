package parser

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/model"
)

var datePrefixRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}-?`)

// ParseFrontmatter 메타데이터 구조체를 언마샬링하고 원시 본문 바이트를 반환합니다.
func ParseFrontmatter(r io.Reader, fm *model.Frontmatter) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	sanitized := sanitizeFrontmatter(data)
	bodyBytes, err := frontmatter.Parse(bytes.NewReader(sanitized), fm)
	if err != nil {
		return "", fmt.Errorf("failed to parse YAML frontmatter: %w", err)
	}
	return string(bodyBytes), nil
}

// sanitizeFrontmatter title에 콜론(: )이 포함되었으나 따옴표로 감싸지지 않은 YAML 문법 오류 자동 보정
func sanitizeFrontmatter(data []byte) []byte {
	lines := strings.Split(string(data), "\n")
	if len(lines) < 2 || !strings.HasPrefix(strings.TrimSpace(lines[0]), "---") {
		return data
	}

	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "---" {
			break
		}
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "title:") {
			val := strings.TrimSpace(strings.TrimPrefix(trimmed, "title:"))
			if !strings.HasPrefix(val, "\"") && !strings.HasPrefix(val, "'") && strings.Contains(val, ":") {
				valEscaped := strings.ReplaceAll(val, "\"", "\\\"")
				lines[i] = fmt.Sprintf("title: %q", valEscaped)
			}
		}
	}

	return []byte(strings.Join(lines, "\n"))
}

// Parse io.Reader로부터 YAML Frontmatter와 순수 마크다운 본문을 파싱하고 유효성을 검증합니다.
func Parse(r io.Reader) (*model.Frontmatter, string, error) {
	var fm model.Frontmatter
	rawContent, err := ParseFrontmatter(r, &fm)
	if err != nil {
		return nil, "", err
	}
	if err := fm.Validate(); err != nil {
		return nil, "", err
	}
	return &fm, rawContent, nil
}

// ParseFile 단일 마크다운 파일을 파싱하여 model.Post 객체로 반환합니다.
// Frontmatter가 없는 일반 마크다운 문서(README.md, dashboard.md 등)인 경우 nil, nil을 반환합니다.
func ParseFile(filePath string) (*model.Post, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read markdown file %s: %w", filePath, err)
	}

	trimmed := strings.TrimSpace(string(data))
	if !strings.HasPrefix(trimmed, "---") {
		// Frontmatter가 없는 문서는 포스트가 아니므로 스킵
		return nil, nil
	}

	var fm model.Frontmatter
	rawContent, err := ParseFrontmatter(bytes.NewReader(data), &fm)
	if err != nil {
		return nil, fmt.Errorf("file %s: %w", filePath, err)
	}

	// 1. 카테고리가 명시되지 않았거나 기본값인 경우 상위 디렉터리 명칭을 카테고리로 자동 설정
	if fm.Category == "General" || strings.TrimSpace(fm.Category) == "" {
		parentDir := filepath.Base(filepath.Dir(filePath))
		if parentDir != "" && parentDir != "." && parentDir != "posts" && parentDir != "drafts" {
			fm.Category = FormatCategoryName(parentDir)
		}
	}

	// 2. created_date가 누락된 경우 파일명 접두어 또는 파일 수정시간으로 fallback
	if strings.TrimSpace(fm.CreatedDate) == "" {
		baseName := filepath.Base(filePath)
		if datePrefixRegex.MatchString(baseName) {
			dateStr := datePrefixRegex.FindString(baseName)
			fm.CreatedDate = strings.TrimSuffix(dateStr, "-")
		} else if fileInfo, err := os.Stat(filePath); err == nil {
			fm.CreatedDate = fileInfo.ModTime().Format("2006-01-02")
		}
	}

	// [SE-05] 필수 필드 및 불변식 사전 검증
	if err := fm.Validate(); err != nil {
		return nil, fmt.Errorf("file %s: %w", filePath, err)
	}

	slug := ExtractSlug(filePath)
	readingTime := CalculateReadingTime(rawContent)

	return &model.Post{
		ID:                 slug,
		Slug:               slug,
		FilePath:           filePath,
		Frontmatter:        fm,
		RawContent:         rawContent,
		ReadingTimeMinutes: readingTime,
	}, nil
}

// ExtractSlug 파일 경로에서 파일명을 기반으로 URL-Safe 슬러그를 추출합니다.
// 예: "posts/deep-dive/2026-09-26-my-post.md" -> "my-post"
// 예: "posts/deep-dive/2026-09-26-my-post.en.md" -> "my-post"
func ExtractSlug(filePath string) string {
	base := filepath.Base(filePath)
	ext := filepath.Ext(base)
	nameWithoutExt := strings.TrimSuffix(base, ext)

	// 다국어 접미사(.en 등) 제거하여 원문과 동일 슬러그 유지
	if strings.HasSuffix(nameWithoutExt, ".en") {
		nameWithoutExt = strings.TrimSuffix(nameWithoutExt, ".en")
	}

	// 선행 날짜(YYYY-MM-DD-) 패턴 제거
	cleaned := datePrefixRegex.ReplaceAllString(nameWithoutExt, "")
	if cleaned == "" {
		cleaned = nameWithoutExt
	}

	return model.Slugify(cleaned)
}

// CalculateReadingTime 마크다운 본문 단어/글자 수를 기반으로 예상 읽기 시간(분)을 계산합니다.
// 한글 및 영문 혼합 기준 분당 300글자/단어 기준 (최소 1분)
func CalculateReadingTime(content string) int {
	clean := strings.TrimSpace(content)
	if clean == "" {
		return 1
	}

	// 공백으로 단어 분리
	words := len(strings.Fields(clean))
	chars := len([]rune(clean))

	// 가중치 계산 (단어 수와 글자 수 평균)
	est := chars / 500
	if words/150 > est {
		est = words / 150
	}
	if est < 1 {
		return 1
	}
	return est
}

// FormatCategoryName 디렉터리 슬러그를 읽기 좋은 카테고리 명칭으로 변환합니다.
// 예: "deep-dive" -> "Deep Dive"
func FormatCategoryName(name string) string {
	parts := strings.Split(name, "-")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}

